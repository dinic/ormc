package model

var (
	tableImportTmpl = `
import (
	{{range $obj :=  .}}{{if not $obj}}{{$obj}}{{else}}"{{$obj}}"{{end}}
	{{end}}
)
`

	tableDefineTmpl = `
type {{.TableStructName}} struct {
	{{range $col := .Cols}}{{$col.ColStructName}} {{$col.ColStructType}} {{$col.Extra}}
	{{end}}
}
`
	tableColNameFuncTmpl = `
{{$TableStructName :=.TableStructName}}
func (t *{{$TableStructName}}) TableName() string {
	return "{{.TableName}}"
}

{{range $col := .Cols}}
func (t *{{$TableStructName}})  {{$col.ColStructName}}ColName() string {
	return "{{$col.ColName}}"
}
{{end}}
`
	tableFuncTmpl = `
{{$TableStructName :=.TableStructName}}
// Get{{$TableStructName}}DB() return {{.TableName}} DB with write and read
func Get{{$TableStructName}}DB() *gorm.DB {
	return GetDB().Table("{{.TableName}}")
}

// Get{{$TableStructName}}ReadDB() return {{.TableName}} DB with  read
func Get{{$TableStructName}}ReadDB() *gorm.DB {
	return GetReadDB().Table("{{.TableName}}")
}

// Get{{$TableStructName}}WriteDB() return {{.TableName}} DB with write
func Get{{$TableStructName}}WriteDB() *gorm.DB {
	return GetWriteDB().Table("{{.TableName}}")
}

{{range $col := .Cols}}
// GetBy{{$col.ColStructName}} Get {{$TableStructName}} by column {{$col.ColName}}
func (t *{{$TableStructName}}) GetBy{{$col.ColStructName}}(value {{$col.ColStructType}})([]*{{$TableStructName}},  error) {
	var res []*{{$TableStructName}}
	err := GetReadDB().Table(t.TableName()).Where("{{$col.ColName}} = ?", value).Scan(&res).Error

	return res, err
}
{{end}}

{{range $col := .Cols}}
// BatchGetBy{{$col.ColStructName}} bacth get {{$TableStructName}} by column {{$col.ColName}}
func (t *{{$TableStructName}}) BatchGetBy{{$col.ColStructName}}(value ...{{$col.ColStructType}})([]*{{$TableStructName}},  error) {
	var res []*{{$TableStructName}}
	err := GetReadDB().Table(t.TableName()).Where("{{$col.ColName}} IN ?", value).Scan(&res).Error

	return res, err
}
{{end}}

{{range $col := .Cols}}
// With{{$col.ColStructName}} set Option with {{$col.ColName}} = value
func (t *{{$TableStructName}}) With{{$col.ColStructName}}(value {{$col.ColStructType}}) Option {
	return func(o *options) {
		var vs []{{$col.ColStructType}}
		v, ok := o.query["{{$col.ColName}}"]
		if !ok {
			vs = make([]{{$col.ColStructType}}, 0, 2)
		} else {
			vs = v.([]{{$col.ColStructType}})
		}
		vs = append(vs, value)
		o.query["{{$col.ColName}}"] = vs
	}
}
{{end}}

// Get get by options
func (t *{{$TableStructName}}) Get(opts ...Option) ([]*{{$TableStructName}}, error){
	var res []*{{$TableStructName}}

	o := options {
		query:make(map[string]interface{}, len(opts)),
	}

	for _, opt := range opts {
		opt(&o)
	}

	tx := GetReadDB().Table(t.TableName())

	for key, value := range o.query {
		tx = tx.Where(key + " IN ?", value)
	}

	err := tx.Scan(&res).Error
	return res, err
}

// Update update by cond, set value to newValues
func (t *{{$TableStructName}}) Update(cond []Option, values ...Option) error{
	cc := options {
		query:make(map[string]interface{}, len(cond)),
	}

	val := options {
		query: make(map[string]interface{}, len(values)),
	}

	for _, c := range cond {
		c(&cc)
	}

	for _, v :=range values {
		v(&val)
	}
	
	return GetWriteDB().Table(t.TableName()).Where(cc.query).Updates(val.query).Error
}
`
)
