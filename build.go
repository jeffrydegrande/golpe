package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

		var newPath = filepath.Join("public/css", path)

		err = ioutil.WriteFile(newPath, []byte(output), 0644)
		check(err)

		fmt.Println(newPath, "created")

		b.Write([]byte(fmt.Sprintf("<link rel=\"stylesheet\" href=\"css/%s\" />\n", path)))
	}

	return b.String()
}

func buildOneFile(path string, funcMap template.FuncMap, files ...string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	check(err)
	defer f.Close()

	fmt.Println(path, "created")
	t := template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
	err = t.ExecuteTemplate(f, "main", nil)
	check(err)

	return nil
}

func createDirectories() {

	err := os.MkdirAll("./public/css", 0770)
	check(err)

	err = os.MkdirAll("./public/js", 0770)
	check(err)
}

func copyImages() {
	cmd := exec.Command("cp", "-r", "images/", "public/images")
	err := cmd.Run()
	check(err)
}

func copyFonts() {
	cmd := exec.Command("cp", "-r", "fonts/", "public/fonts")
	err := cmd.Run()
	check(err)
}

func BuildAll() error {
	say("Creating directories")
	createDirectories()

	say("Compiling jsx")
	runJsxCompiler()

	say("Compile javascripts")
	var javascripts = buildJavascripts()

	say("Compiling stylesheets")
	var stylesheets = buildStylesheets()

	say("Copying images")
	copyImages()

	say("Coping fonts")
	copyFonts()

	funcMap := template.FuncMap{
		"javascripts": func() string { return javascripts },
		"stylesheets": func() string { return stylesheets },
	}

	say("Compiling templates")

	layouts, err := filepath.Glob("./*.tmpl")
	check(err)

	htmlFiles, err := filepath.Glob("./*.html")
	check(err)

	for _, html := range htmlFiles {
		var files []string
		files = append(files, html)
		files = append(files, layouts...)

		buildOneFile(filepath.Join("public", html), funcMap, files...)
	}

	return nil
}
