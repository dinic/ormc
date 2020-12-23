package dbinfo

import (
	"fmt"
	"log"

	"github.com/dinic/ormc/pkg/utils"
)

type Table struct {
	db              *DB
	TableName       string
	CreateSQL       string
	Comment         string
	Cols            map[string]*Column
	Refers          []*ReferenceCol
	tableStructName string
}

func newTable(tbname, comment string) *Table {
	t := new(Table)

	t.TableName = tbname
	t.Comment = comment
	t.Cols = make(map[string]*Column)
	return t
}

func (t *Table) load() {
	sqlfmt := "SELECT COLUMN_NAME,COLUMN_DEFAULT,IS_NULLABLE,COLUMN_TYPE,COLUMN_KEY,COLUMN_COMMENT,EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'"

	var cols []*Column

	err := t.db.o.Raw(fmt.Sprintf(sqlfmt, t.db.Config.DBName, t.TableName)).Scan(&cols).Error
	if err != nil {
		log.Fatalf("%s.%s load failed, err: %s", t.db.Config.DBName, t.TableName, err)
	}

	for _, col := range cols {
		col.table = t
		col.table.db = t.db
		t.Cols[col.ColumnName] = col
	}

	t.loadForeignKey()
}

func (t *Table) loadForeignKey() {
	sqlfmt := "select TABLE_SCHEMA ,TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_SCHEMA,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where TABLE_SCHEMA = '%s' AND REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_NAME='%s'"
	var refers []*ReferenceCol

	err := t.db.o.Raw(fmt.Sprintf(sqlfmt, t.db.Config.DBName, t.TableName)).Scan(&refers).Error
	if err != nil {
		log.Fatalf("%s.%s load foreign key failed, err: %s", t.db.Config.DBName, t.TableName, err)
	}

	t.Refers = refers
}

func (t *Table) StructName() string {
	if t.tableStructName != "" {
		return t.tableStructName
	}

	cf := t.db.Config.GetTableConfig(t.TableName)

	if cf != nil && cf.StructName != "" {
		t.tableStructName = cf.StructName
		return t.tableStructName
	}

	t.tableStructName = utils.Hungarian2Camel(t.TableName)
	return t.tableStructName
}

func (t *Table) GetColumn(col string) *Column {
	if v, ok := t.Cols[col]; ok {
		return v
	}

	return nil
}

func (t *Table) IsForeignCol(col string) bool {
	for _, r := range t.Refers {
		if r.ColumnName == col {
			return true
		}
	}

	return false
}
