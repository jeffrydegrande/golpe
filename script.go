package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
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

func buildJavascripts(concat bool) string {
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

	return b.String()
}
