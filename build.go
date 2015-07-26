package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func buildStylesheets() string {
	stylesheets, err := filepath.Glob("**/*.scss")
	check(err)

	var b bytes.Buffer
	for _, scssPath := range stylesheets {
		// only include 'full' scss files
		var base = filepath.Base(scssPath)
		if base[0] == '_' {
			continue
		}

		output, err := CompileSassFile(scssPath)
		check(err)

		var path = fmt.Sprintf("%s.css", base[:len(base)-5])

		err = ioutil.WriteFile(filepath.Join("public/css", path), []byte(output), 0644)
		check(err)

		b.Write([]byte(fmt.Sprintf("<link rel=\"stylesheet\" href=\"css/%s\" />\n", path)))
	}

	return b.String()
}

func buildOneFile(path string, funcMap template.FuncMap, files ...string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	check(err)
	defer f.Close()

	fmt.Printf("Created file %s\n", path)
	t := template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
	err = t.ExecuteTemplate(f, "main", nil)
	check(err)

	return nil
}

func createDirectories() {

	fmt.Printf("Clearing public directory\n")
	err := os.RemoveAll("./public")
	check(err)

	fmt.Printf("Creating public directory\n")
	err = os.MkdirAll("./public", 0770)
	check(err)

	fmt.Printf("Creating stylesheet directory\n")
	err = os.MkdirAll("./public/css", 0770)
	check(err)

	fmt.Printf("Creating javascript directory\n")
	err = os.MkdirAll("./public/js", 0770)
	check(err)
}

func BuildAll() error {
	createDirectories()

	layouts, err := filepath.Glob("./*.tmpl")
	check(err)

	htmlFiles, err := filepath.Glob("./*.html")
	check(err)

	var javascripts = buildJavascripts()
	var stylesheets = buildStylesheets()

	funcMap := template.FuncMap{
		"javascripts": func() string { return javascripts },
		"stylesheets": func() string { return stylesheets },
	}

	for _, html := range htmlFiles {
		var files []string
		files = append(files, html)
		files = append(files, layouts...)

		buildOneFile(filepath.Join("public", html), funcMap, files...)
	}

	return nil
}
