package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dinic/ormc/pkg/utils"

	"gopkg.in/yaml.v3"
)

type MysqlConfig struct {
	DBName    string                  `yaml:"-"`
	Host      string                  `yaml:"host"`
	Port      string                  `yaml:"port"`
	User      string                  `yaml:"user"`
	Password  string                  `yaml:"password"`
	Package   string                  `yaml:"package"`
	OutPutDir string                  `yaml:"out_dir"`
	Tables    map[string]*TableConfig `yaml:"tables"`
}

func (c *MysqlConfig) initDefault() {
	for k, v := range c.Tables {
		v.TableName = k
		v.initDefault()
	}
}

func (c *MysqlConfig) GetTableConfig(table string) *TableConfig {
	if v, ok := c.Tables[table]; ok {
		return v
	}

	return nil
}

func (c *MysqlConfig) GetColConfig(table, col string) *ColumnConfig {
	tc := c.GetTableConfig(table)
	if tc == nil {
		return nil
	}

	return tc.GetColConfig(col)
}

type TableConfig struct {
	TableName  string                   `yaml:"-"`
	StructName string                   `yaml:"struct_name"`
	Columns    map[string]*ColumnConfig `yaml:"columns"`
}

func (c *TableConfig) initDefault() {
	for k, v := range c.Columns {
		v.ColumnName = k
	}
}

func (c *TableConfig) GetColConfig(col string) *ColumnConfig {
	if v, ok := c.Columns[col]; ok {
		return v
	}

	return nil
}

type ColumnConfig struct {
	ColumnName  string `yaml:"-"`
	StructName  string `yaml:"struct_name"`
	DisableJson bool   `yaml:"disable_json"`
	ForeignTag  string `yaml:"foreign_tag"`
	Tag         string `yaml:"tag"`
}

func (c *ColumnConfig) parseTag(tag string) map[string]string {
	if tag == "" {
		return nil
	}

	res := make(map[string]string)
	str := strings.Split(tag, ",")
	for _, s := range str {
		ss := strings.SplitN(s, ":", 2)
		key := strings.TrimSpace(ss[0])
		value := ""
		if len(ss) == 2 {
			value = strings.TrimSpace(ss[1])
		}
		res[key] = value
	}

	return res
}

func (c *ColumnConfig) GetForeignTag() map[string]string {
	return c.parseTag(c.ForeignTag)
}

func (c *ColumnConfig) GetTag() map[string]string {
	return c.parseTag(c.Tag)
}

type TagConfig struct {
	ColumnName string `yaml:"column_name"`
	Tag        string `yaml:"tag"`
}

type Config struct {
	OutPutDir string                  `yaml:"out_dir"`
	Database  map[string]*MysqlConfig `yaml:"database"`
}

var (
	globalConfig *Config
)

func (c *Config) initDefault() {
	if c.OutPutDir == "" {
		c.OutPutDir = "gen_by_ormc"
	}

	if !filepath.IsAbs(c.OutPutDir) {
		dir, err := filepath.Abs(c.OutPutDir)
		if err != nil {
			log.Fatalf("get work dir failed, err: %s", err)
		}
		c.OutPutDir = dir
	}

	for k, mc := range c.Database {
		if mc.Package == "" {
			mc.Package = utils.Camel2Hungarian(mc.DBName)
		}

		if mc.OutPutDir == "" {
			mc.OutPutDir = filepath.Join(c.OutPutDir, mc.Package)
		}

		mc.DBName = k
		mc.initDefault()
	}
}

func (c *Config) GetDBConfig(db string) *MysqlConfig {
	if v, ok := c.Database[db]; ok {
		return v
	}

	log.Fatalf("db config not found, db: %s", db)
	return nil
}

func (c *Config) GetTableConfig(db, table string) *TableConfig {
	dc := c.GetDBConfig(db)
	if dc == nil {
		log.Fatalf("db config not found, db: %s", db)
	}

	return dc.GetTableConfig(table)
}

func (c *Config) GetColConfig(db, table, col string) *ColumnConfig {
	tc := c.GetTableConfig(db, table)
	if tc == nil {
		log.Fatalf("table config not found, db: %s, table: %s", db, table)
	}

	return tc.GetColConfig(col)
}

func Parse(path string) *Config {
	if globalConfig != nil {
		return globalConfig
	}

	c := new(Config)
	defer c.initDefault()

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s failed, err: %s", path, err)
	}

	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		log.Fatalf("config decode failed, err: %s", err)
	}

	globalConfig = c

	return globalConfig
}

func Get() *Config {
	return globalConfig
}
