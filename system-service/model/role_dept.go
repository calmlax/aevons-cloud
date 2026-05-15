/**
 * 角色和部门关联表 Model
 *
 * @author
 * @date 2026-04-09 01:49:40.928391294 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type RoleDept struct {

	//角色ID
	RoleId int64 `gorm:"column:role_id;primaryKey;type:bigint;comment:角色ID" json:"roleId,string"`

	//部门ID
	DeptId int64 `gorm:"column:dept_id;primaryKey;type:bigint;comment:部门ID" json:"deptId,string"`
}

// TableName 指定表名
func (RoleDept) TableName() string {
	return "sys_role_dept"
}
