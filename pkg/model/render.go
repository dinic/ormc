package model

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type item struct {
	data interface{}
	tmpl string
}

type Renderer struct {
	file  string
	items []*item
}

func NewRender(file string) *Renderer {
	r := new(Renderer)
	r.file = file
	r.items = make([]*item, 0, 4)

	return r
}

func (r *Renderer) addItem(data interface{}, tmpl string) {
	i := &item{
		data: data,
		tmpl: tmpl,
	}

	if r.items == nil {
		r.items = make([]*item, 0, 4)
	}

	r.items = append(r.items, i)
}

func (r *Renderer) Render() {
	if len(r.items) == 0 {
		return
	}

	os.MkdirAll(filepath.Dir(r.file), 0755)

	f, err := os.OpenFile(r.file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("open %s failed, err: %s", r.file, err)
	}

	temp := template.New(r.file)
	defer f.Close()
	for _, i := range r.items {
		tpl, err := temp.Clone()
		if err != nil {
			log.Fatalf("template clone failed, err: %s", err)
		}
		t, err := tpl.Parse(i.tmpl)
		if err != nil {
			log.Fatal(err)
		}

		if err = t.Execute(f, i.data); err != nil {
			log.Fatal(err)
		}
	}

	if err := exec.Command("goimports", "-l", "-w", r.file).Run(); err != nil {
		log.Fatalf("auto import failed, file: %s, err: %s", r.file, err)
	}

	if err = exec.Command("gofmt", "-l", "-w", r.file).Run(); err != nil {
		log.Fatalf("fmt file: %s failed, err: %s", r.file, err)
	}
}
