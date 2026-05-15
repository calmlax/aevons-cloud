/**
 * 部门表 Repository
 *
 * @author
 * @date 2026-04-09 01:02:34.478942247 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"context"
	"fmt"
	"strconv"
	"system-service/dto"
	"system-service/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

type DeptRepository interface {
	base.BaseRepository[model.Dept]
	CreateDept(ctx context.Context, d dto.CreateDeptDTO) (*model.Dept, error)
	UpdateDept(ctx context.Context, id int64, d dto.UpdateDeptDTO) error
	ListTree(query dto.DeptQuery) ([]model.Dept, error)
	HasChildren(id int64) (bool, error)
}

type deptRepository struct {
	base.BaseRepository[model.Dept]
	db *gorm.DB
}

func NewDeptRepository(db *gorm.DB) DeptRepository {
	return &deptRepository{
		BaseRepository: base.NewBaseRepository[model.Dept](db),
		db:             db,
	}
}

// getAncestors 根据 parentId 计算祖级路径
func (r *deptRepository) getAncestors(parentId int64) (string, error) {
	if parentId == 0 {
		return "0", nil
	}
	var parent model.Dept
	if err := r.db.Select("id, ancestors").Where("id = ?", parentId).First(&parent).Error; err != nil {
		return "", fmt.Errorf("父级部门不存在")
	}
	return parent.Ancestors + "," + strconv.FormatInt(parentId, 10), nil
}

func (r *deptRepository) CreateDept(ctx context.Context, d dto.CreateDeptDTO) (*model.Dept, error) {
	ancestors, err := r.getAncestors(d.ParentId)
	if err != nil {
		return nil, err
	}
	dept := model.Dept{
		ParentId:  d.ParentId,
		Ancestors: ancestors,
		DeptType:  d.DeptType,
		DeptName:  d.DeptName,
		Sort:      d.Sort,
		Status:    d.Status,
		Remark:    d.Remark,
	}
	if err := r.db.WithContext(ctx).Create(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *deptRepository) UpdateDept(ctx context.Context, id int64, d dto.UpdateDeptDTO) error {
	ancestors, err := r.getAncestors(*d.ParentId)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Dept{}).Where("id = ?", id).Updates(map[string]any{
		"parent_id": d.ParentId,
		"ancestors": ancestors,
		"dept_type": d.DeptType,
		"dept_name": d.DeptName,
		"sort":      d.Sort,
		"status":    d.Status,
		"remark":    d.Remark,
	}).Error
}

func (r *deptRepository) ListTree(query dto.DeptQuery) ([]model.Dept, error) {
	db := r.db.Model(&model.Dept{})
	if query.DeptName != nil && *query.DeptName != "" {
		db = db.Where("dept_name LIKE ?", "%"+*query.DeptName+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.DeptType != nil {
		db = db.Where("dept_type = ?", *query.DeptType)
	}
	var list []model.Dept
	err := db.Order("parent_id, sort").Find(&list).Error
	return list, err
}

func (r *deptRepository) HasChildren(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Dept{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}
