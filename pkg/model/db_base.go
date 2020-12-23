package model

/*
import (
	"context"

	"gorm.io/gorm"
)

var (
	_db      *gorm.DB
	_writeDB *gorm.DB
	_readDB  *gorm.DB
	_preLoad bool
)

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
	return _writeDB
}

// SetReadDB set database orm with read only
func SetReadDB(o *gorm.DB) {
	_readDB = o
}

// GetReadDB get database orm with read only
func GetReadDB(o *gorm.DB) *gorm.DB {
	return _readDB
}

// EnablePreload global enable mult table preload, default disabled
func EnablePreload() {
	_preLoad = true
}

// DisablePreload global disable mul table preload, default disabled
func DisablePreload() {
	_preLoad = false
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

type Option interface {
	Apply(map[string]interface{})
}
*/
