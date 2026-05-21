/**
 * 代码生成表 Service
 *
 * @author
 * @date 2026-04-08 03:51:33.212851159 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"gen-service/internal/dto"
	"gen-service/internal/model"
	"gen-service/internal/repository"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"
	"github.com/calmlax/aevons-framework/core/scope"

	"gorm.io/gorm"
)

// 继承BaseService
type GenTableService interface {
	base.BaseService[model.GenTable, *dto.GenTableQuery]

	DBTables() ([]model.GenTable, error)

	ImportDbTables(ctx context.Context, tableNames []string) error

	SynchDbTable(ctx context.Context, tableId int64) error

	CodePreview(tableId int64) ([]dto.CodeDTO, error)

	CommandGenerateCode(tableNames []string, moduleName string) (bool, error)

	DownloadCodeZip(tableIds []int64) ([]byte, error)
}

type genTableService struct {
	base.BaseService[model.GenTable, *dto.GenTableQuery]
	repo       repository.GenTableRepository
	columnRepo repository.GenTableColumnRepository
	cfg        config.GenConfig
}

func NewGenTableService(repo repository.GenTableRepository, columnRepo repository.GenTableColumnRepository, cfg config.GenConfig) GenTableService {
	baseSrv := base.NewBaseService[model.GenTable, *dto.GenTableQuery](repo)
	return &genTableService{
		BaseService: baseSrv,
		repo:        repo,
		columnRepo:  columnRepo,
		cfg:         cfg,
	}
}

// Delete 批量删除数据表（级联删除对应字段）
func (s *genTableService) BatchDelete(ids []int64, scopes ...scope.DBScope) error {
	// 先删除关联字段
	if err := s.columnRepo.DeleteByTableIds(ids); err != nil {
		return err
	}

	return s.BaseService.BatchDelete(ids, scopes...)
}

func (s *genTableService) DBTables() ([]model.GenTable, error) {
	dbName, err := s.repo.GetDBName()
	if err != nil {
		return nil, err
	}

	// 已导入的表名
	existing, err := s.repo.List(&dto.GenTableQuery{})
	if err != nil {
		return nil, err
	}
	imported := make([]string, 0, len(existing))
	for _, t := range existing {
		imported = append(imported, t.DbTableName)
	}
	return s.repo.DBTables(dbName, s.cfg.ExcludeTables, imported)
}
func (s *genTableService) ImportDbTables(ctx context.Context, tableNames []string) error {
	dbName, err := s.repo.GetDBName()
	if err != nil {
		return err
	}

	// 获取所有表的基础信息（用于后续获取注释）
	dbTables, err := s.repo.DBTables(dbName, nil, nil)
	if err != nil {
		return err
	}

	// 开启事务
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		for _, tableName := range tableNames {
			// 校验表是否已经存在
			var count int64
			// 注意：在事务中查询必须使用 tx
			if err := tx.Model(&model.GenTable{}).Where("table_name = ?", tableName).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return fmt.Errorf("表[%s]已存在，请勿重复导入", tableName)
			}
			// 获取数据库字段信息
			cols, err := s.repo.DBColumns(dbName, tableName)
			if err != nil {
				return fmt.Errorf("获取表[%s]字段失败: %w", tableName, err)
			}
			// 计算元数据

			var dbTable model.GenTable
			for _, t := range dbTables {
				if t.DbTableName == tableName {
					dbTable = t
					break
				}
			}
			cleanName := s.RemoveTablePrefix(dbTable.DbTableName, s.cfg.CleanPrefixs)
			dbTable.ClassName = utils.ToUpperCamel(cleanName)
			if idx := strings.Index(dbTable.DbTableName, "_"); idx > -1 {
				dbTable.ModuleName = dbTable.DbTableName[:idx]
			}
			dbTable.Router = strings.ReplaceAll(cleanName, "_", "/")
			dbTable.Permission = strings.ReplaceAll(tableName, "_", ":")
			dbTable.MenuId = 0
			// 关键点：使用 tx 进行创建
			if err := tx.Create(&dbTable).Error; err != nil {
				return fmt.Errorf("保存表[%s]记录失败: %w", tableName, err)
			}

			if len(cols) > 0 {
				for i := range cols {
					cols[i].TableId = dbTable.Id
				}
				if err := tx.Create(&cols).Error; err != nil {
					return fmt.Errorf("保存表[%s]字段失败: %w", tableName, err)
				}
			}
		}
		return nil
	})
}

// SynchDbTable 同步数据库表结构
func (s *genTableService) SynchDbTable(ctx context.Context, tableId int64) error {
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		// 获取本地表信息
		var table model.GenTable
		if err := tx.First(&table, tableId).Error; err != nil {
			return errors.New("表不存在")
		}

		// 获取数据库实际表结构 (模拟 baseMapper.dbTableByTableName)
		dbName, err := s.repo.GetDBName()
		dbTables, err := s.repo.DBTablesByTableNames(dbName, []string{table.DbTableName}, s.cfg.ExcludeTables)
		if err != nil || len(dbTables) == 0 {
			return fmt.Errorf("数据库表[%s]不可用或不存在", table.TableName)
		}
		dbTable := dbTables[0]
		dbColumns, err := s.repo.DBColumns(dbName, table.DbTableName)
		if err != nil {
			return err
		}

		// 获取本地已有的字段列表
		var localFields []model.GenTableColumn
		tx.Where("table_id = ?", tableId).Find(&localFields)

		// 建立 Map 方便快速查找 (FieldName -> GenTableField)
		localFieldMap := make(map[string]*model.GenTableColumn)
		for i := range localFields {
			localFieldMap[localFields[i].FieldName] = &localFields[i]
		}

		dbFieldMap := make(map[string]*model.GenTableColumn)
		for i := range dbColumns {
			dbFieldMap[dbColumns[i].FieldName] = &dbColumns[i]
		}

		// --- 开始比对逻辑 ---

		var removeFieldIds []int64
		var updateFields []model.GenTableColumn
		var addFields []model.GenTableColumn

		// 找出要删除的 (本地有，数据库没了)
		for _, local := range localFields {
			if _, ok := dbFieldMap[local.FieldName]; !ok {
				removeFieldIds = append(removeFieldIds, local.Id)
			}
		}

		// 找出要更新或新增的
		for _, dbF := range dbColumns {
			if local, ok := localFieldMap[dbF.FieldName]; ok {
				// 更新字段属性
				local.ColumnComment = dbF.ColumnComment
				local.IsRequired = dbF.IsRequired
				local.IsPrimaryKey = dbF.IsPrimaryKey
				local.IsAutoIncrement = dbF.IsAutoIncrement
				local.DataType = dbF.DataType
				local.DefaultValue = dbF.DefaultValue
				local.DataLength = dbF.DataLength
				local.DataPrecision = dbF.DataPrecision
				updateFields = append(updateFields, *local)
			} else {
				// 新增字段
				dbF.TableId = tableId // 关联 ID
				addFields = append(addFields, dbF)
			}
		}

		// --- 数据库操作 ---

		cleanName := s.RemoveTablePrefix(dbTable.DbTableName, s.cfg.CleanPrefixs)
		dbTable.ClassName = utils.ToUpperCamel(cleanName)
		// 更新表注释
		if err := tx.Model(&table).Updates(map[string]any{
			"table_comment": dbTable.DbTableComment,
			"class_name":    dbTable.ClassName,
		}).Error; err != nil {
			return err
		}

		// 批量删除
		if len(removeFieldIds) > 0 {
			if err := tx.Delete(&model.GenTableColumn{}, removeFieldIds).Error; err != nil {
				return err
			}
		}

		// 批量更新 (GORM 批量更新通常需要循环或使用特定的 Save)
		for _, f := range updateFields {
			if err := tx.Save(&f).Error; err != nil {
				return err
			}
		}

		// 批量新增
		if len(addFields) > 0 {
			if err := tx.Create(&addFields).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *genTableService) GetGens(isFrontend bool) []dto.CodeDTO {

	var gens []dto.CodeDTO

	gens = append(gens, dto.CodeDTO{
		Name:         "Model",
		FileName:     "model.tpl",
		SaveFileName: "%s-service/internal/model/%s.go",
		FileType:     "go",
	})
	gens = append(gens, dto.CodeDTO{
		Name:         "DTO",
		FileName:     "dto.tpl",
		SaveFileName: "%s-service/internal/dto/%s_dto.go",
		FileType:     "go",
	})
	gens = append(gens, dto.CodeDTO{
		Name:         "Service",
		FileName:     "service.tpl",
		SaveFileName: "%s-service/internal/service/%s_service.go",
		FileType:     "go",
	})
	gens = append(gens, dto.CodeDTO{
		Name:         "Repository",
		FileName:     "repository.tpl",
		SaveFileName: "%s-service/internal/repository/%s_repository.go",
		FileType:     "go",
	})
	gens = append(gens, dto.CodeDTO{
		Name:         "Handler",
		FileName:     "handler.tpl",
		SaveFileName: "%s-service/internal/handler/%s_handler.go",
		FileType:     "go",
	})
	gens = append(gens, dto.CodeDTO{
		Name:         "Router",
		FileName:     "router.tpl",
		SaveFileName: "%s-service/internal/router/%s_router.go",
		FileType:     "go",
	})
	if !isFrontend {
		gens = append(gens, dto.CodeDTO{
			Name:         "SQL",
			FileName:     "sql.tpl",
			SaveFileName: "%s-service/gen_code/sql/%s.sql",
			FileType:     "sql",
		})
		gens = append(gens, dto.CodeDTO{
			Name:         "API",
			FileName:     "api.tpl",
			SaveFileName: "%s-service/gen_code/api/%s.ts",
			FileType:     "ts",
		})
		gens = append(gens, dto.CodeDTO{
			Name:         "Vue",
			FileName:     "vue.tpl",
			SaveFileName: "%s-service/gen_code/vue/%s.vue",
			FileType:     "vue",
		})
	} else {
		gens = append(gens, dto.CodeDTO{
			Name:         "sql",
			FileName:     "sql.tpl",
			SaveFileName: "%s-service/gen_code/sql/%s.menu.sql",
			FileType:     "sql",
		})
		gens = append(gens, dto.CodeDTO{
			Name:         "api",
			FileName:     "api.tpl",
			SaveFileName: "%s-service/gen_code/api/%s.ts",
			FileType:     "ts",
		})
		gens = append(gens, dto.CodeDTO{
			Name:         "vue",
			FileName:     "vue.tpl",
			SaveFileName: "%s-service/gen_code/vue/%s.vue",
			FileType:     "vue",
		})
	}
	return gens
}

func (s *genTableService) GenerateCode(table *model.GenTable, columns []model.GenTableColumn, isFrontend bool) ([]dto.CodeDTO, error) {

	baseClass, excludeFields := s.getBaseClass(columns)
	cleanName := s.RemoveTablePrefix(table.DbTableName, s.cfg.CleanPrefixs)
	// fields := s.FilterColumns(columns, excludeFields)
	// if fields == nil {
	// 	fields = []model.GenTableColumn{}
	// }
	data := map[string]any{
		"TableName":     table.DbTableName,
		"Comment":       table.DbTableComment,
		"CleanName":     cleanName,
		"ClassName":     table.ClassName,
		"ModuleName":    table.ModuleName,
		"Author":        table.Author,
		"Router":        table.Router,
		"BaseClass":     baseClass,
		"MenuId":        table.MenuId,
		"Permission":    table.Permission,
		"Remark":        table.Remark,
		"ApiPrefix":     s.cfg.ApiPrefix,
		"Date":          time.Now().UTC().Format("2006-01-02"),
		"ExcludeFields": excludeFields,
		"Fields":        columns,
	}
	gens := s.GetGens(isFrontend)
	for i := range gens {
		gen := &gens[i]
		// 渲染模板
		code, err := utils.ParseTemplateFile(gen.FileName, data)
		if err != nil {
			return nil, err
		}
		gen.Code = code
	}
	return gens, nil
}

func (s *genTableService) CodePreview(tableId int64) ([]dto.CodeDTO, error) {
	table, err := s.repo.GetById(tableId)
	if err != nil {
		return nil, err
	}
	if utils.IsEmpty(table) {
		return nil, errors.New("表不存在")
	}
	columns, err := s.columnRepo.ListByTableId(table.Id)
	if err != nil {
		return nil, err
	}
	return s.GenerateCode(table, columns, false)
}

// 生成代码
func (s *genTableService) CommandGenerateCode(tableNames []string, moduleName string) (bool, error) {
	fmt.Println("-------------------GenerateCode----------------------")
	fmt.Printf("tableName=%v\n", tableNames)
	fmt.Printf("moduleName=%v\n", moduleName)
	dbName, err := s.repo.GetDBName()
	if err != nil {
		return false, err
	}
	fmt.Printf("dbName=%v\n", dbName)
	tables, err := s.repo.DBTablesByTableNames(dbName, tableNames, s.cfg.ExcludeTables)
	if err != nil {
		return false, err
	}
	if utils.IsEmpty(tables) {
		return false, errors.New("表不存在或在'generator.exclude_tables'排除规则列表中")
	}

	for _, table := range tables {
		cleanName := s.RemoveTablePrefix(table.DbTableName, s.cfg.CleanPrefixs)
		table.ClassName = utils.ToUpperCamel(cleanName)
		table.Router = strings.ReplaceAll(cleanName, "_", "/")
		table.Permission = strings.ReplaceAll(table.DbTableName, "_", ":")
		table.MenuId = 0
		table.ModuleName = moduleName
		table.Author = ""
		columns, err := s.repo.DBColumns(dbName, table.DbTableName)
		if err != nil {
			return false, err
		}
		gens, err := s.GenerateCode(&table, columns, false)
		if err != nil {
			return false, err
		}
		for _, gen := range gens {

			// 拼接最终保存路径（自动替换 moduleName 和 className）
			savePath := fmt.Sprintf(gen.SaveFileName, "../"+moduleName, cleanName)
			dir := filepath.Dir(savePath)

			// 自动创建目录
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return false, err
			}

			// 4. 写入文件
			err = os.WriteFile(savePath, []byte(gen.Code), 0644)
			if err != nil {
				return false, err
			}

			fmt.Printf("✅ 已生成：%s\n", savePath)
		}
	}

	fmt.Println("\n🎉 所有表代码生成完成！")
	return true, nil
}

// 下载
func (s *genTableService) DownloadCodeZip(tableIds []int64) ([]byte, error) {

	tables, err := s.repo.ListByIds(tableIds)
	if err != nil {
		return nil, err
	}
	if utils.IsEmpty(tables) {
		return nil, errors.New("表不存在或在'generator.exclude_tables'排除规则列表中")
	}

	// 创建缓冲区
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// 遍历写入 Zip
	for _, table := range tables {
		columns, err := s.columnRepo.ListByTableId(table.Id)
		if err != nil {
			return nil, err
		}
		gens, err := s.GenerateCode(&table, columns, true)
		cleanName := s.RemoveTablePrefix(table.DbTableName, s.cfg.CleanPrefixs)
		for _, model := range gens {
			f, err := zipWriter.Create(fmt.Sprintf(model.SaveFileName, table.ModuleName, cleanName))
			if err != nil {
				return nil, err
			}
			_, err = f.Write([]byte(model.Code))
			if err != nil {
				return nil, err
			}
		}

	}

	// 4. 必须先关闭 zipWriter 才能完成数据的刷盘（写入 buf）
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 清理表名前缀 例如 sys_user → user (prefix: sys_)
func (s *genTableService) RemoveTablePrefix(tableName string, prefixs []string) string {
	for _, prefix := range prefixs {
		if strings.HasPrefix(tableName, prefix) {
			return strings.TrimPrefix(tableName, prefix)
		}
	}
	return tableName
}

func (s *genTableService) getBaseClass(columns []model.GenTableColumn) (string, []string) {

	fieldMap := make(map[string]bool)
	for _, col := range columns {
		fieldMap[col.ColumnName] = true
	}
	if utils.ContainsAllMap(fieldMap, []string{"is_deleted", "created_by", "created_at", "updated_by", "updated_at"}) {
		return "base.CommonModel", []string{"is_deleted", "created_by", "created_at", "updated_by", "updated_at"}
	} else if utils.ContainsAllMap(fieldMap, []string{"created_by", "created_at", "updated_by", "updated_at"}) {
		return "base.DefaultModel", []string{"created_by", "created_at", "updated_by", "updated_at"}
	} else if utils.ContainsAllMap(fieldMap, []string{"created_at", "updated_at"}) {
		return "base.BaseModel", []string{"created_at", "updated_at"}
	} else if utils.ContainsAllMap(fieldMap, []string{"is_deleted"}) {
		return "base.LogicDeleteModel", []string{"is_deleted"}
	}
	return "", []string{}
}

// 过滤掉 columns 中在 excludeFields 里的字段（按 ColumnName 匹配）
func (s *genTableService) FilterColumns(columns []model.GenTableColumn, excludeFields []string) []model.GenTableColumn {
	if len(excludeFields) == 0 {
		return columns
	}
	// 把排除列表转成 map 加速查找
	excludeMap := make(map[string]bool)
	for _, f := range excludeFields {
		excludeMap[f] = true
	}

	// 只保留不在排除列表中的字段
	var filtered []model.GenTableColumn
	for _, col := range columns {
		if !excludeMap[col.ColumnName] {
			filtered = append(filtered, col)
		}
	}
	return filtered
}
