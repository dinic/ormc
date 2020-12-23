package dbinfo

import (
	"fmt"
	"log"

	"github.com/dinic/ormc/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	o      *gorm.DB
	Config *config.MysqlConfig
	Tables map[string]*Table
}

func NewDB(cf *config.MysqlConfig) *DB {
	db := new(DB)

	db.Config = cf
	db.Tables = make(map[string]*Table)

	return db
}

func (db *DB) open() {
	c := mysql.Config{}
	c.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		db.Config.User, db.Config.Password, db.Config.Host, db.Config.Port, db.Config.DBName)
	c.DefaultStringSize = 256
	c.DisableDatetimePrecision = false
	c.DontSupportRenameIndex = false
	c.DontSupportRenameColumn = false
	c.SkipInitializeWithVersion = false

	d, err := gorm.Open(mysql.New(c), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %s failed, err: %s", db.Config.DBName, err)
	}

	db.o = d
}

func (db *DB) Load() {
	db.open()
	db.readAllTables()

	for n, t := range db.Tables {
		log.Printf("load %s.%s table info", db.Config.DBName, n)
		t.load()
	}
}

func (db *DB) readAllTables() {
	sqlfmt := "select TABLE_SCHEMA as table_schema,TABLE_NAME as table_name,TABLE_COMMENT as comment  from information_schema.TABLES  where TABLE_SCHEMA = '%s'"

	var tables []struct {
		DBname    string `gorm:"column:table_schema"`
		TableName string `gorm:"column:table_name"`
		Comment   string `gorm:"column:comment"`
	}

	if err := db.o.Raw(fmt.Sprintf(sqlfmt, db.Config.DBName)).Scan(&tables).Error; err != nil {
		log.Fatalf("read all tables faield, err: %s", err)
	}

	for _, t := range tables {
		table := newTable(t.TableName, t.Comment)
		table.db = db
		db.Tables[table.TableName] = table
	}
}

func (db *DB) GetTable(table string) *Table {
	if v, ok := db.Tables[table]; ok {
		return v
	}

	return nil
}

func (db *DB) GetColumn(table, col string) *Column {
	t := db.GetTable(table)
	if t == nil {
		return nil
	}

	return t.GetColumn(col)
}
