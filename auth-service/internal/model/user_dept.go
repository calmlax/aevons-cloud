/**
 * 用户部门 Model
 *
 * @author
 * @date 2026-04-09 02:01:11.848928919 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type UserDept struct {

	//用户编号
	UserId int64 `gorm:"column:user_id;primaryKey;type:bigint;comment:用户编号" json:"userId,string"`

	//部门编号
	DeptId int64 `gorm:"column:dept_id;primaryKey;type:bigint;comment:部门编号" json:"deptId,string"`

	//岗位编号
	PostId int64 `gorm:"column:post_id;type:bigint;comment:岗位编号" json:"postId,string"`
}

// TableName 指定表名
func (UserDept) TableName() string {
	return "sys_user_dept"
}
