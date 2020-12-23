package model

var extraImports = map[string]string{
	"time.Time": "time",
}

func importPackage(t string) string {
	if v, ok := extraImports[t]; ok {
		return v
	}

	return ""
}
