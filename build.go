package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func buildStylesheets() string {
	stylesheets, err := filepath.Glob("**/*.scss")
	check(err)

	for _, css := range stylesheets {
		CompileSassFile(css)
	}

	fmt.Printf("%v\n", stylesheets)
	return ""
}

func buildJavascripts() string {
	javascripts, err := filepath.Glob("**/*.js")
	check(err)

	fmt.Printf("%v\n", javascripts)

	var b bytes.Buffer

	for _, js := range javascripts {
		b.Write([]byte(fmt.Sprintf("<script src=\"%s\"></script>", js)))
	}

	return b.String()
}

func buildOneFile(path string, files ...string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	check(err)
	defer f.Close()

	fmt.Printf("Created file %s\n", path)

	funcMap := template.FuncMap{
		"javascripts": buildJavascripts,
		"stylesheets": buildStylesheets,
	}

	t := template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
	err = t.ExecuteTemplate(f, "main", nil)
	check(err)

	return nil
}

func BuildAll() error {
	layouts, err := filepath.Glob("./*.tmpl")
	check(err)

	htmlFiles, err := filepath.Glob("./*.html")
	check(err)

	if len(htmlFiles) == 0 {
		return nil
	}

	fmt.Printf("Creating public directory\n")
	err = os.MkdirAll("./public", 0770)
	check(err)

	for _, html := range htmlFiles {
		var files []string
		files = append(files, html)
		files = append(files, layouts...)

		buildOneFile(filepath.Join("public", html), files...)
	}
	return nil
}
