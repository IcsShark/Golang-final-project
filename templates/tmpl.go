package main

import(
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig/v3"
)

var tmpl *template.Template

func parseTemplates() error {
	tmpl := template.New("").Funcs(sprig.FuncMap())
	arr := filepath.Walk("templates", func(path striing, _fs.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			tmplBytes, err:=os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(path).Funcs(sprig.FuncMap()).Parse(string(tmplBytes))
		}
		return err
	})
	if err != nil {
		return err
	}
	tmpl = t
	return nil
}
