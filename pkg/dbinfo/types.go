package dbinfo

import (
	"log"
	"regexp"
)

var Mysql2GoDic = map[string]string{
	"smallint":            "int16",
	"smallint unsigned":   "uint16",
	"int":                 "int",
	"int unsigned":        "uint",
	"bigint":              "int64",
	"bigint unsigned":     "uint64",
	"varchar":             "string",
	"char":                "string",
	"date":                "time.Time",
	"datetime":            "time.Time",
	"bit(1)":              "[]uint8",
	"tinyint":             "int8",
	"tinyint unsigned":    "uint8",
	"tinyint(1)":          "bool", // tinyint(1) 默认设置成bool
	"tinyint(1) unsigned": "bool", // tinyint(1) 默认设置成bool
	"json":                "string",
	"text":                "string",
	"timestamp":           "time.Time",
	"double":              "float64",
	"mediumtext":          "string",
	"longtext":            "string",
	"float":               "float32",
	"float unsigned":      "float32",
	"tinytext":            "string",
	"enum":                "string",
	"time":                "time.Time",
	"tinyblob":            "[]byte",
	"blob":                "[]byte",
	"mediumblob":          "[]byte",
	"longblob":            "[]byte",
}

var Mysql2GoRegexp = map[string]string{
	`^(tinyint)[(]\d+[)]`:            "int8",
	`^(tinyint)[(]\d+[)] unsigned`:   "uint8",
	`^(smallint)[(]\d+[)]`:           "int16",
	`^(int)[(]\d+[)]`:                "int",
	`^(bigint)[(]\d+[)]`:             "int64",
	`^(char)[(]\d+[)]`:               "string",
	`^(enum)[(](.)+[)]`:              "string",
	`^(varchar)[(]\d+[)]`:            "string",
	`^(varbinary)[(]\d+[)]`:          "[]byte",
	`^(binary)[(]\d+[)]`:             "[]byte",
	`^(decimal)[(]\d+,\d+[)]`:        "float64",
	`^(mediumint)[(]\d+[)]`:          "string",
	`^(double)[(]\d+,\d+[)]`:         "float64",
	`^(float)[(]\d+,\d+[)]`:          "float64",
	`^(float)[(]\d+,\d+[)] unsigned`: "float64",
	`^(datetime)[(]\d+[)]`:           "time.Time",
	`^(bit)[(]\d+[)]`:                "[]uint8",
}

func TypeMysql2Go(t string) string {
	if v, ok := Mysql2GoDic[t]; ok {
		return v
	}

	for p, v := range Mysql2GoRegexp {
		matched, err := regexp.MatchString(p, t)
		if err != nil {
			log.Fatalf("regexp match failed, pattern: %s, str: %s, err: %s", p, t, err)
		}

		if matched {
			return v
		}
	}

	return ""
}
