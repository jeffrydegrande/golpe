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

func buildJavascripts() string {

	javascripts, err := filepath.Glob("javascripts/vendor/*.js")
	check(err)

	lib, err := filepath.Glob("javascripts/lib/*.js")
	check(err)
	javascripts = append(javascripts, lib...)

	components, err := filepath.Glob("javascripts/components/*.js")
	check(err)
	javascripts = append(javascripts, components...)

	toplevel, err := filepath.Glob("javascripts/*.js")
	check(err)
	javascripts = append(javascripts, toplevel...)

	var b bytes.Buffer
	for _, path := range javascripts {
		js := filepath.Base(path)
		pathTo := filepath.Join("public/js/", js)

		CopyFile(path, pathTo)
		b.Write([]byte(fmt.Sprintf("<script src=\"js/%s\"></script>\n", js)))
	}

	return b.String()
}
