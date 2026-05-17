/**
 * 代码生成表字段 Model
 *
 * @author
 * @date 2026-04-08 03:51:33.21431843 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type GenTableColumn struct {

	//编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:编号" json:"id,string"`

	//所属表编号
	TableId int64 `gorm:"column:table_id;type:bigint;comment:所属表编号" json:"tableId,string"`

	//DB字段名称
	ColumnName string `gorm:"column:column_name;type:varchar(50);comment:DB字段名称" json:"columnName"`

	//DB字段解释
	ColumnComment string `gorm:"column:column_comment;type:varchar(50);comment:DB字段解释" json:"columnComment"`

	//DB字段数据类型
	ColumnType string `gorm:"column:column_type;type:varchar(20);comment:DB字段数据类型" json:"columnType"`

	//数据类型
	DataType string `gorm:"column:data_type;type:varchar(20);comment:数据类型" json:"dataType"`

	//字段名称
	FieldName string `gorm:"column:field_name;type:varchar(50);comment:字段名称" json:"fieldName"`

	//JSON
	Json string `gorm:"column:json;type:varchar(100);comment:JSON" json:"json"`

	//是否主键（0否，1是）
	IsPrimaryKey bool `gorm:"column:is_primary_key;type:tinyint(1);comment:是否主键（0否，1是）" json:"isPrimaryKey"`

	//是否自增（0否，1是）
	IsAutoIncrement bool `gorm:"column:is_auto_increment;type:tinyint(1);comment:是否自增（0否，1是）" json:"isAutoIncrement"`

	//是否必填（0否，1是）
	IsRequired bool `gorm:"column:is_required;type:tinyint(1);comment:是否必填（0否，1是）" json:"isRequired"`

	//是否插入字段（0否，1是）
	IsInsert bool `gorm:"column:is_insert;type:tinyint(1);comment:是否插入字段（0否，1是）" json:"isInsert"`

	//是否编辑字段（0否，1是）
	IsEdit bool `gorm:"column:is_edit;type:tinyint(1);comment:是否编辑字段（0否，1是）" json:"isEdit"`

	//是否列表字段（0否，1是）
	IsList bool `gorm:"column:is_list;type:tinyint(1);comment:是否列表字段（0否，1是）" json:"isList"`

	//列表排序字段（0否，1是）
	Sortable bool `gorm:"column:sortable;type:tinyint(1);comment:是否列表字段（0否，1是）" json:"sortable"`

	//列表筛选字段（0否，1是）
	Filterable bool `gorm:"column:filterable;type:tinyint(1);comment:是否列表字段（0否，1是）" json:"filterable"`

	//查询条件
	Condition string `gorm:"column:condition;type:varchar(10);comment:查询条件" json:"condition"`

	//字典类型
	DictType string `gorm:"column:dict_type;type:varchar(64);comment:字典类型" json:"dictType"`

	//排序
	Sort int `gorm:"column:sort;type:int;comment:排序" json:"sort"`

	//UI组件
	Component string `gorm:"column:component;type:varchar(20);comment:UI组件" json:"component"`

	//默认值
	DefaultValue string `gorm:"column:default_value;type:varchar(50);comment:默认值" json:"defaultValue"`

	//数据长款
	DataLength int `gorm:"column:data_length;type:int;comment:数据长款" json:"dataLength"`

	//数据精度
	DataPrecision int8 `gorm:"column:data_precision;type:tinyint(1);comment:数据精度" json:"dataPrecision"`
}

// TableName 指定表名
func (GenTableColumn) TableName() string {
	return "gen_table_column"
}
