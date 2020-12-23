package model

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/dinic/ormc/pkg/dbinfo"
	"github.com/dinic/ormc/pkg/utils"
)

type TableModel struct {
	db      *dbinfo.DB
	table   *dbinfo.Table
	Package string
}

func NewTableModel(pkg string, db *dbinfo.DB, table *dbinfo.Table) *TableModel {
	tm := new(TableModel)

	tm.Package = pkg
	tm.db = db
	tm.table = table

	return tm
}

func (tm *TableModel) fileName() string {
	return filepath.Join(tm.db.Config.OutPutDir, utils.Camel2Hungarian(tm.table.TableName)+"_table.go")
}

func (tm *TableModel) renderImport(r *Renderer) {
	imports := []string{"gorm.io/gorm"}

	for _, col := range tm.table.Cols {
		t := dbinfo.TypeMysql2Go(col.DataType)
		if t == "" {
			log.Fatalf("not found mysql type: %s", col.DataType)
		}

		p := importPackage(t)
		if p == "" {
			continue
		}

		found := false
		for _, s := range imports {
			if s == p {
				found = true
				break
			}
		}

		if !found {
			imports = append(imports, p)
		}
	}

	r.addItem(imports, tableImportTmpl)
}

type tableDefine struct {
	TableName       string
	TableStructName string
	Cols            []*colDefine
}

type colDefine struct {
	ColName       string
	ColStructName string
	ColStructType string
	Extra         string
}

func (tm *TableModel) colGormTag(col *dbinfo.Column) string {
	tags := make([]string, 0, 4)

	tags = append(tags, fmt.Sprintf("column:%s", col.ColumnName))
	tags = append(tags, fmt.Sprintf("type:%s", col.DataType))
	if col.IsNotNull() {
		tags = append(tags, "not null")
	}

	if col.IsUniqueKey() {
		tags = append(tags, "unique")
	}

	if col.IsPrimaryKey() {
		tags = append(tags, "primaryKey")
	}

	if col.IsAutoIncrement() {
		tags = append(tags, "autoIncrement")
	}
	return fmt.Sprintf("gorm:\"%s\"", strings.Join(tags, ";"))
}

func (tm *TableModel) colJsonTag(col *dbinfo.Column) string {
	cf := tm.db.Config.GetColConfig(tm.table.TableName, col.ColumnName)
	if cf != nil && cf.DisableJson == true {
		return ""
	}

	return fmt.Sprintf("json:\"%s\"", utils.Camel2Hungarian(col.StructName()))
}

func (tm *TableModel) extra(col *dbinfo.Column) string {
	gtag := tm.colGormTag(col)
	jtag := tm.colJsonTag(col)
	tag := gtag + " " + jtag
	tag = strings.TrimSpace(tag)
	tag = "`" + tag + "`"
	if len(col.ColumnComment) == 0 {
		return tag
	}

	return tag + " // " + col.ColumnComment
}

func (tm *TableModel) renderTableDefine(r *Renderer) {
	td := new(tableDefine)
	td.Cols = make([]*colDefine, 0, len(tm.table.Cols))

	td.TableName = tm.table.TableName
	td.TableStructName = tm.table.StructName()

	for _, col := range tm.table.Cols {
		cd := new(colDefine)
		cd.ColName = col.ColumnName
		cd.ColStructName = col.StructName()
		cd.ColStructType = dbinfo.TypeMysql2Go(col.DataType)
		cd.Extra = tm.extra(col)
		td.Cols = append(td.Cols, cd)
	}
	// TODO surrport reference column

	r.addItem(td, tableDefineTmpl)
	r.addItem(td, tableColNameFuncTmpl)
	r.addItem(td, tableFuncTmpl)
}

func (tm *TableModel) Render() *Renderer {
	r := NewRender(tm.fileName())
	r.addItem(tm.Package, packageTmpl)
	tm.renderImport(r)
	tm.renderTableDefine(r)

	return r
}
