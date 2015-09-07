package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func runJsxCompiler() {
	cmd := exec.Command("jsx", "jsx", "javascripts/components")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		panic(stderr.String())
	}
}

func jsFiles() []string {
	javascripts, err := filepath.Glob("javascripts/*.js")
	check(err)
	return javascripts
}

func jsFilesFrom(group string) []string {
	javascripts, err := filepath.Glob(fmt.Sprintf("javascripts/%s/*.js", group))
	check(err)
	return javascripts
}

func buildJavascripts() string {
	javascripts := jsFilesFrom("vendor")
	javascripts = append(javascripts, jsFilesFrom("lib")...)
	javascripts = append(javascripts, jsFilesFrom("components")...)
	javascripts = append(javascripts, jsFiles()...) // javascripts/*.js that is

	var b bytes.Buffer
	for _, path := range javascripts {
		js := filepath.Base(path)

		CopyFile(path, filepath.Join("public/js/", js))
		b.Write([]byte(fmt.Sprintf("<script src=\"js/%s\"></script>\n", js)))
	}

	CopyFile("javascripts/app.js", "public/js/app.js")
	return b.String()
}

func buildJavascriptsFromAppJs() string {
	// build a list of all the javascripts files available
	// and turn them into a map from base -> full path
	fullPaths := makeJavascriptPathMapping()

	re := regexp.MustCompile("//= require (.*)")
	lines, err := readLines("javascripts/app.js")
	check(err)

	var b bytes.Buffer

	for _, line := range lines {
		res := re.FindStringSubmatch(line)
		if res == nil {
			continue
		}

		// resolve the full path of this dependency
		path := filepath.Base(res[1])
		if !strings.HasSuffix(path, ".js") {
			path = path + ".js"
		}

		fullPath := fullPaths[path]

		pathTo := filepath.Join("public/js/", path)
		CopyFile(fullPath, pathTo)
		fmt.Println(pathTo, "created")
		b.Write([]byte(fmt.Sprintf("<script src=\"/js/%s\"></script>\n", path)))
	}

	CopyFile("javascripts/app.js", "public/js/app.js")
	// b.Write([]byte("<script src=\"/js/app.js\"></script>\n"))

	return b.String()
}

func makeJavascriptPathMapping() map[string]string {
	javascripts, err := filepath.Glob("javascripts/*.js")
	check(err)

	extraJavascripts, err := filepath.Glob("javascripts/**/*.js")
	check(err)

	javascripts = append(javascripts, extraJavascripts...)

	var fullPath map[string]string
	fullPath = make(map[string]string)

	for _, js := range javascripts {
		base := filepath.Base(js)
		fullPath[base] = js

	}
	return fullPath
}

