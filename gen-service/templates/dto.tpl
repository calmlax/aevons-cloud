package dto

/**
 * {{.Comment}} DTO
 *
 * @author {{.Author}}
 * @date {{.Date}}
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"aevo/pkg/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type {{.ClassName}}Query struct {
	{{range .Fields -}}
	{{ if ne .Condition ""}}
	// {{.ColumnComment}}
	{{.FieldName}}  *{{- if or (eq .Condition "in") (eq .Condition "not_in") (eq .Condition "between") -}}[]
	{{- end -}}{{.DataType}} `form:"{{.FieldName | toLowerCamel}}" gorm:"column:{{.ColumnName}}" q:"{{.Condition}}"`{{end}}{{end}}
	base.BaseQuery
}

type Create{{.ClassName}}DTO struct {
	{{range .Fields -}}{{- if .IsInsert -}}
	//{{.ColumnComment}}
	{{.FieldName}} {{.DataType}} `json:"{{.Json}}" 
	{{- if or (and .IsRequired (ne .DataType "bool")) (and (gt .DataLength 0) (eq .DataType "string")) }} binding:"
	 {{- if .IsRequired -}}
	 required
	 {{- end -}}
	 {{- if and .IsRequired  (and (gt .DataLength 0) (eq .DataType "string")) -}},{{- end -}}
	 {{- if and (gt .DataLength 0) (eq .DataType "string") -}}max={{.DataLength}}{{- end -}}"{{- end -}}`
	 {{- end }}
	 {{end}}
}

type Update{{.ClassName}}DTO struct {
	{{range .Fields -}}{{- if .IsEdit -}}
	//{{.ColumnComment}}
	{{.FieldName}} *{{.DataType}} `json:"{{.Json}}" 
	{{- if or (and .IsRequired (ne .DataType "bool")) (and (gt .DataLength 0) (eq .DataType "string")) }} binding:"
	 {{- if .IsRequired -}}
	 required
	 {{- end -}}
	 {{- if and .IsRequired  (and (gt .DataLength 0) (eq .DataType "string")) -}},{{- end -}}
	 {{- if and (gt .DataLength 0) (eq .DataType "string") -}}max={{.DataLength}}{{- end -}}"{{- end -}}`
	 {{- end }}
	 {{end}}
}


type {{.ClassName}}DTO struct {
	{{range $i, $f := .Fields -}}
	//{{$f.ColumnComment}}
	{{$f.FieldName}} {{$f.DataType}} `excel:"column:{{$f.ColumnComment}};index:{{$i}};dict:{{$f.DictType}}" json:"{{$f.Json}}"
	{{- if or (and $f.IsRequired (ne $f.DataType "bool")) (and (gt $f.DataLength 0) (eq $f.DataType "string")) }} binding:"
	{{- if $f.IsRequired -}}
	required
	{{- end -}}
	{{- if and $f.IsRequired (and (gt $f.DataLength 0) (eq $f.DataType "string")) -}},{{- end -}}
	{{- if and (gt $f.DataLength 0) (eq $f.DataType "string") -}}max={{$f.DataLength}}{{- end -}}"{{- end -}}`
	{{end}}
}
