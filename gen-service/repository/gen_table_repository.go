/**
 * 代码生成表 Repository
 *
 * @author
 * @date 2026-04-08 03:51:33.212851159 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"fmt"
	"gen-service/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type GenTableRepository interface {
	base.BaseRepository[model.GenTable]

	// DBTables 列出当前数据库中所有表（排除已导入的）
	DBTables(dbName string, excludePatterns []string, alreadyImported []string) ([]model.GenTable, error)

	// DBColumns 列出指定表的所有字段
	DBColumns(dbName, tableName string) ([]model.GenTableColumn, error)

	// GetDBName 获取当前连接的数据库名
	GetDBName() (string, error)

	// DBTablesByTableNames 列出当前数据库中所有表（排除已导入的）
	DBTablesByTableNames(dbName string, tableNames []string, excludePatterns []string) ([]model.GenTable, error)
}

// genTableRepository 结构体
type genTableRepository struct {
	base.BaseRepository[model.GenTable]
	db      *gorm.DB
	isMySQL bool
	isPg    bool
}

// 创建 GenTableRepository 实例
func NewGenTableRepository(db *gorm.DB) GenTableRepository {
	name := db.Dialector.Name()
	return &genTableRepository{
		BaseRepository: base.NewBaseRepository[model.GenTable](db),
		db:             db,
		isMySQL:        name == "mysql",
		isPg:           name == "postgres",
	}
}

func (r *genTableRepository) GetDBName() (string, error) {
	var dbName string
	if r.isMySQL {
		err := r.db.Raw("SELECT DATABASE()").Scan(&dbName).Error
		return dbName, err
	}
	if r.isPg {
		err := r.db.Raw("SELECT current_database()").Scan(&dbName).Error
		return dbName, err
	}
	return "", fmt.Errorf("unsupported database type")
}

func (r *genTableRepository) DBTables(dbName string, excludePatterns []string, alreadyImported []string) ([]model.GenTable, error) {
	var rows []struct {
		TableName    string
		TableComment string
	}

	if r.isMySQL {
		err := r.db.Raw(`
			SELECT table_name AS table_name, IFNULL(table_comment,'') AS table_comment
			FROM information_schema.tables
			WHERE table_schema = ? AND table_type = 'BASE TABLE'
			ORDER BY table_name`, dbName).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	if r.isPg {
		err := r.db.Raw(`
			SELECT c.relname AS table_name,
			       COALESCE(d.description, '') AS table_comment
			FROM pg_catalog.pg_class c
			LEFT JOIN pg_catalog.pg_description d ON d.objoid = c.oid AND d.objsubid = 0
			WHERE c.relkind IN ('r', 'p')
			  AND c.relnamespace = (SELECT oid FROM pg_catalog.pg_namespace WHERE nspname = 'public')
			ORDER BY c.relname`,
		).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	imported := make(map[string]bool, len(alreadyImported))
	for _, n := range alreadyImported {
		imported[n] = true
	}

	var result []model.GenTable
	for _, row := range rows {
		if imported[row.TableName] {
			continue
		}
		if r.matchesAnyPattern(row.TableName, excludePatterns) {
			continue
		}
		result = append(result, model.GenTable{
			DbTableName:    row.TableName,
			DbTableComment: row.TableComment,
		})
	}
	return result, nil
}

func (r *genTableRepository) DBTablesByTableNames(dbName string, tableNames []string, excludePatterns []string) ([]model.GenTable, error) {
	var rows []struct {
		TableName    string
		TableComment string
	}

	// 如果没有传入表名，就查全部
	if len(tableNames) == 0 {
		return nil, fmt.Errorf("未指定需要查询的表名")
	}

	if r.isMySQL {
		// MySQL 多表 IN 查询
		err := r.db.Raw(`
			SELECT table_name AS table_name, IFNULL(table_comment,'') AS table_comment
			FROM information_schema.tables
			WHERE table_schema = ? 
			AND table_type = 'BASE TABLE' 
			AND table_name IN (?)
			ORDER BY table_name`,
			dbName, tableNames).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	if r.isPg {
		// PostgreSQL 多表 IN 查询
		err := r.db.Raw(`
			SELECT c.relname AS table_name,
			       COALESCE(d.description, '') AS table_comment
			FROM pg_catalog.pg_class c
			LEFT JOIN pg_catalog.pg_description d ON d.objoid = c.oid AND d.objsubid = 0
			WHERE c.relkind IN ('r', 'p')
			  AND c.relnamespace = (SELECT oid FROM pg_catalog.pg_namespace WHERE nspname = 'public')
			  AND c.relname IN (?)
			ORDER BY c.relname`,
			tableNames).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	var result []model.GenTable
	for _, row := range rows {
		if r.matchesAnyPattern(row.TableName, excludePatterns) {
			continue
		}
		result = append(result, model.GenTable{
			DbTableName:    row.TableName,
			DbTableComment: row.TableComment,
		})
	}
	return result, nil
}

func (r *genTableRepository) DBColumns(dbName, tableName string) ([]model.GenTableColumn, error) {
	var rows []struct {
		ColumnName      string
		ColumnComment   string
		ColumnType      string
		ColumnKey       string
		Extra           string
		IsNullable      string
		ColumnDefault   *string
		OrdinalPosition int
	}

	// MySQL 查询
	if r.isMySQL {
		err := r.db.Raw(`
			SELECT
				column_name,
				IFNULL(column_comment,'') AS column_comment,
				column_type,
				column_key,
				extra,
				is_nullable,
				column_default,
				ordinal_position
			FROM information_schema.columns
			WHERE table_schema = ? AND table_name = ?
			ORDER BY ordinal_position`, dbName, tableName).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	// PostgreSQL 查询
	if r.isPg {
		err := r.db.Raw(`
			SELECT
				c.column_name,
				COALESCE(pd.description, '') AS column_comment,
				c.data_type AS column_type,
				c.is_nullable,
				c.column_default,
				c.ordinal_position
			FROM information_schema.columns c
			LEFT JOIN pg_catalog.pg_description pd
				ON pd.objoid = c.table_catalog::regclass::oid
				AND pd.objsubid = c.ordinal_position
			WHERE c.table_schema = 'public'
			AND c.table_name = ?
			ORDER BY c.ordinal_position;`, tableName).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	result := make([]model.GenTableColumn, 0, len(rows))
	for _, row := range rows {
		dataType := utils.MysqlOrPgTypeToDataType(row.ColumnType)
		json := ""
		if dataType == "int64" || dataType == "uint64" {
			//解决 JS 丢失精度问题：int64/uint64 → 返回 string，其他原样返回
			json = utils.ToLowerCamel(row.ColumnName) + ",string"
		} else {
			json = utils.ToLowerCamel(row.ColumnName)
		}
		col := model.GenTableColumn{
			ColumnName:      row.ColumnName,
			ColumnComment:   row.ColumnComment,
			ColumnType:      row.ColumnType,
			IsPrimaryKey:    row.ColumnKey == "PRI",
			IsAutoIncrement: strings.Contains(row.Extra, "auto_increment"),
			IsRequired:      row.IsNullable != "YES",
			Sort:            row.OrdinalPosition,
			Json:            json,
			DataType:        dataType,
			FieldName:       utils.ToUpperCamel(row.ColumnName),
			IsInsert:        row.ColumnKey != "PRI",
			IsEdit:          true,
			IsList:          true,
			Sortable:        false,
			Filterable:      false,
		}
		col.Condition = r.SmartCondition(col)
		length, scale := r.ColumnTypeLengthScale(col.ColumnType)
		col.DataLength = length
		col.DataPrecision = scale
		if row.ColumnDefault != nil {
			col.DefaultValue = *row.ColumnDefault
		}
		result = append(result, col)
	}
	return result, nil
}

var typeRegex = regexp.MustCompile(`(\w+)(?:\((\d+)(?:,(\d+))?\))?`)

// 解析数据库字段类型
func (r *genTableRepository) ColumnTypeLengthScale(dataType string) (int, int8) {
	dataType = strings.ToLower(dataType)

	match := typeRegex.FindStringSubmatch(dataType)
	if len(match) == 0 {
		return 0, 0
	}

	length := 0
	scale := 0
	// 解析长度
	if match[2] != "" {
		length, _ = strconv.Atoi(match[2])
	}

	// 解析小数位
	if match[3] != "" {
		scale, _ = strconv.Atoi(match[3])
	}
	return length, int8(scale)
}

// SmartCondition 智能匹配查询条件
func (r *genTableRepository) SmartCondition(col model.GenTableColumn) string {
	columnName := strings.ToLower(col.ColumnName)
	// dataType := strings.ToLower(col.DataType)
	var condition string
	// 主键 / ID → eq
	if strings.HasSuffix(columnName, "_id") ||
		strings.Contains(columnName, "status") ||
		strings.Contains(columnName, "type") ||
		strings.Contains(columnName, "state") {
		condition = "eq"
	}

	// 字典 / 状态 / 类型 → in
	if col.DictType != "" ||
		strings.Contains(columnName, "status") ||
		strings.Contains(columnName, "type") ||
		strings.Contains(columnName, "state") {
		condition = "in"
	}

	// 时间 / 日期 → between 区间
	if strings.Contains(columnName, "_time") ||
		strings.Contains(columnName, "_date") ||
		strings.Contains(columnName, "created_at") {
		condition = "between"
	}

	// 名称、标题、内容 → like
	if strings.Contains(columnName, "name") ||
		strings.Contains(columnName, "title") ||
		strings.Contains(columnName, "content") ||
		strings.Contains(columnName, "desc") {
		condition = "like"
	}

	// 手机、邮箱、账号 → left_like
	if strings.Contains(columnName, "phone") ||
		strings.Contains(columnName, "mobile") ||
		strings.Contains(columnName, "email") ||
		strings.Contains(columnName, "_no") {
		condition = "like_l"
	}

	return condition
}

// 数字类型
func (r *genTableRepository) isNumber(dt string) bool {
	return strings.Contains(dt, "int") ||
		strings.Contains(dt, "float") ||
		strings.Contains(dt, "double") ||
		strings.Contains(dt, "decimal")
}

// 字符串类型
func (r *genTableRepository) isString(dt string) bool {
	return strings.Contains(dt, "char") || strings.Contains(dt, "text")
}

// matchesAnyPattern 判断 name 是否匹配任意 pattern（支持 % 通配符前缀/后缀）
func (r *genTableRepository) matchesAnyPattern(name string, patterns []string) bool {
	for _, p := range patterns {
		if p == "" {
			continue
		}
		if r.matchPattern(name, p) {
			return true
		}
	}
	return false
}

func (r *genTableRepository) matchPattern(name, pattern string) bool {
	if len(pattern) == 0 {
		return false
	}
	if pattern[len(pattern)-1] == '%' {
		prefix := pattern[:len(pattern)-1]
		return len(name) >= len(prefix) && name[:len(prefix)] == prefix
	}
	if pattern[0] == '%' {
		suffix := pattern[1:]
		return len(name) >= len(suffix) && name[len(name)-len(suffix):] == suffix
	}
	return name == pattern
}
