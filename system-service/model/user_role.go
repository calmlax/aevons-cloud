/**
 * 用户和角色关联表 Model
 *
 * @author
 * @date 2026-04-09 02:01:11.850369127 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type UserRole struct {

	//用户ID
	UserId int64 `gorm:"column:user_id;primaryKey;type:bigint;comment:用户ID" json:"userId,string"`

	//角色ID
	RoleId int64 `gorm:"column:role_id;primaryKey;type:bigint;comment:角色ID" json:"roleId,string"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "sys_user_role"
}
