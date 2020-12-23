package dbinfo

import (
	"strings"

	"github.com/dinic/ormc/pkg/utils"
)

var (
	defaultColName map[string]string = map[string]string{
		"id": "ID",
	}
)

type Column struct {
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
	IsNullAble    string `gorm:"column:IS_NULLABLE"`
	DataType      string `gorm:"column:COLUMN_TYPE"`
	ColumnKey     string `gorm:"column:COLUMN_KEY"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT"`
	Extra         string `gorm:"column:EXTRA"`
	colStructName string `gorm:"-"`
	table         *Table `gorm:"-"`
}

func (c *Column) StructName() string {
	if c.colStructName != "" {
		return c.colStructName
	}

	cf := c.table.db.Config.GetColConfig(c.table.TableName, c.ColumnName)

	if cf != nil && cf.StructName != "" {
		c.colStructName = cf.StructName
		return c.colStructName
	}

	if v, ok := defaultColName[c.ColumnName]; ok && v != "" {
		c.colStructName = v
		return c.colStructName
	}

	c.colStructName = utils.Hungarian2Camel(c.ColumnName)
	return c.colStructName
}

func (c *Column) IsPrimaryKey() bool {
	if strings.Contains(strings.ToLower(c.ColumnKey), "pri") {
		return true
	}
	return false
}

func (c *Column) IsUniqueKey() bool {
	if strings.Contains(strings.ToLower(c.ColumnKey), "uni") {
		return true
	}
	return false
}

func (c *Column) IsNotNull() bool {
	if strings.Contains(strings.ToLower(c.IsNullAble), "yes") {
		return true
	}

	return false
}

func (c *Column) IsAutoIncrement() bool {
	if strings.Contains(strings.ToLower(c.Extra), "auto_increment") {
		return true
	}
	return false
}

type ReferenceCol struct {
	DBname           string `gorm:"column:REFERENCED_TABLE_SCHEMA"`
	TableName        string `gorm:"column:REFERENCED_TABLE_NAME"`
	ColumnName       string `gorm:"column:REFERENCED_COLUMN_NAME"`
	ReferencedDB     string `gorm:"column:REFERENCED_TABLE_SCHEMA"`
	ReferencedTable  string `gorm:"column:REFERENCED_TABLE_NAME"`
	ReferencedColumn string `gorm:"column:REFERENCED_TABLE_NAME"`
}
