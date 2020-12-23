package model

var (
	packageTmpl = `// NOTE:generate by ormc, DO NOT EDIT!
	package {{.}}
	`
	dbImportTmpl = `
import(
	"context"

	"gorm.io/gorm"
)
`
	dbVarTmpl = `
var (
    _db      *gorm.DB
    _writeDB *gorm.DB
    _readDB  *gorm.DB
)
`
	dbContentTmpl = `
// SetDB set database  with write and read
func SetDB(db *gorm.DB) {
    _db = db
}

// GetDB get database with write and read
func GetDB() *gorm.DB {
    return _db
}

// SetWriteDB set database with write only
func SetWriteDB(db *gorm.DB) {
    _writeDB = db
}

// GetWriteDB get database with write only
func GetWriteDB() *gorm.DB {
	if _writeDB == nil {
		return _db
	}
    return _writeDB
}

// SetReadDB set database orm with read only
func SetReadDB(o *gorm.DB) {
    _readDB = o
}

// GetReadDB get database orm with read only
func GetReadDB() *gorm.DB {
	if _readDB == nil {
		return _db
	}
    return _readDB
}

// Table base database table
type Table struct {
    *gorm.DB
    ctx     *context.Context
    preLoad bool
}

// SetCtx set table context
func (t *Table) SetCtx(ctx *context.Context) {
    t.ctx = ctx
}

// GetCtx get table context
func (t *Table) GetCtx() *context.Context {
    return t.ctx
}

// SetDB set table db, default global _db or _writeDB or _readDB
func (t *Table) SetDB(db *gorm.DB) {
    t.DB = db
}

// GetDB get table db
func (t *Table) GetDB() *gorm.DB {
    return t.DB
}

// EnablePreload enable table preload
func (t *Table) EnablePreload() {
    t.preLoad = true
}

// DisablePreload disable table preload
func (t *Table) DisablePreload() {
    t.preLoad = false
}

type options struct {
	query map[string]interface{}
}

type Option func(o *options)
`
)
