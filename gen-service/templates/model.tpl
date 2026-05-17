package model

/**
 * {{.Comment}} Model
 *
 * @author {{.Author}}
 * @date {{.Date}}
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
{{if ne .BaseClass ""}}
import (
	"aevo/pkg/base"
)
{{end}}

type {{.ClassName}} struct {
  {{ range .Fields }}
  {{- if not (contains $.ExcludeFields .ColumnName) }}
  //{{.ColumnComment}}
  {{.FieldName}} {{.DataType}} `gorm:"column:{{.ColumnName}};
  {{- if .IsPrimaryKey -}}primaryKey;{{- end -}}
  {{- if .IsAutoIncrement -}}autoIncrement;{{- end -}}type:{{.ColumnType}};{{- if ne .ColumnComment "" -}} comment:{{.ColumnComment}}{{- end -}}" json:"{{.Json}}"`
  {{- end -}}
  {{ end }}
  {{.BaseClass}}
}

// TableName 指定表名
func ({{.ClassName}}) TableName() string {
	return "{{.TableName}}"
}