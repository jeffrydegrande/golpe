package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

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

		var path = filepath.Join("public", html)

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		check(err)
		defer f.Close()

		fmt.Printf("Created file %s\n", path)

		t := template.Must(template.ParseFiles(files...))
		err = t.ExecuteTemplate(f, "main", nil)
		check(err)
	}
	return nil
}
