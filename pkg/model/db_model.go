package model

import (
	"path/filepath"

	"github.com/dinic/ormc/pkg/dbinfo"
	"github.com/dinic/ormc/pkg/utils"
)

type DBModel struct {
	db      *dbinfo.DB
	Package string
	file    string
}

func NewDbModel(db *dbinfo.DB) *DBModel {
	m := new(DBModel)

	m.db = db

	return m
}

func (m *DBModel) initPackage() {
	m.Package = m.db.Config.Package
}

func (m *DBModel) initFilename() {
	m.file = filepath.Join(m.db.Config.OutPutDir, utils.Camel2Hungarian(m.db.Config.DBName)+".go")
}

func (m *DBModel) init() {
	m.initPackage()
	m.initFilename()
}

func (m *DBModel) Render() []*Renderer {
	m.init()
	res := make([]*Renderer, 0, 4)

	r := NewRender(m.file)
	r.addItem(m.Package, packageTmpl)
	r.addItem(nil, dbImportTmpl)
	r.addItem(nil, dbVarTmpl)
	r.addItem(nil, dbContentTmpl)

	res = append(res, r)

	for _, table := range m.db.Tables {
		tm := NewTableModel(m.Package, m.db, table)
		r = tm.Render()
		res = append(res, r)
	}

	return res
}
