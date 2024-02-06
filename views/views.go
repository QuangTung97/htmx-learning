package views

import (
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func findAndParseTemplates(rootDir string) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	prefix := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.WalkDir(cleanRoot, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		name := path[prefix:]
		_, err = root.New(name).Parse(string(data))
		return err
	})
	if err != nil {
		return nil, err
	}

	return root, nil
}

var initTemplates sync.Once
var globalTemplate *template.Template

func getTemplates() *template.Template {
	initTemplates.Do(func() {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("Can not get caller line")
		}

		templDir := path.Join(path.Dir(filename), "templ")
		globalTemplate = template.Must(findAndParseTemplates(templDir))
	})

	return globalTemplate
}

func Execute(w io.Writer, templateName string, data any) error {
	return getTemplates().ExecuteTemplate(w, templateName, data)
}
