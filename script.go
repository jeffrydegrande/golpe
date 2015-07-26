package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func runJsxCompiler() {
	jsx, err := filepath.Glob("jsx/*.js")
	check(err)

	os.MkdirAll("public/js/components", 0770)

	for _, js := range jsx {
		fmt.Printf("Compiling %s\n", js)
		cmd := exec.Command("jsx", js)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()

		if err != nil {
			fmt.Printf("%s\n", stderr.String())
			os.Exit(1)
		}

		path := filepath.Join("public/js/components", filepath.Base(js))
		ioutil.WriteFile(path, []byte(stdout.String()), 0660)
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func CopyFile(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	check(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, []byte(data), 0644)
	check(err)
}

func buildJavascripts() string {
	// this will give us compiled jsx
	runJsxCompiler()

	javascripts, err := filepath.Glob("javascripts/*.js")
	check(err)

	// javascriptsComponents, err := filepath.Glob("javascripts/components/*.js")
	// check(err)

	// javascripts = append(javascripts, javascriptsComponents...)

	var m map[string]int
	m = make(map[string]int)

	// now we have to weight javascripts, by pulling out their dependencies
	re := regexp.MustCompile("//= require (.*)")
	for _, js := range javascripts {
		// current file has default priority
		m[filepath.Base(js)] = 1

		lines, err := readLines(js)
		check(err)
		for _, line := range lines {
			res := re.FindStringSubmatch(line)
			if res == nil {
				continue
			}

			jsFile := res[1]

			if !strings.HasSuffix(jsFile, ".js") {
				jsFile = jsFile + ".js"
			}

			_, err := os.Stat(filepath.Join("javascripts", jsFile))
			check(err)

			prio := 0
			switch filepath.Dir(jsFile) {
			case "components":
				prio = 100
				break
			case "vendor":
				prio = 1000
				break
			}

			m[jsFile] += prio
		}
	}

	err = os.MkdirAll("public/js/components", 0770)
	check(err)

	err = os.MkdirAll("public/js/vendor", 0770)
	check(err)

	keys := sortedKeys(m)
	var b bytes.Buffer
	for _, js := range keys {
		fmt.Printf("%s\n", js)

		CopyFile(filepath.Join("javascripts", js), filepath.Join("public/js/", js))

		b.Write([]byte(fmt.Sprintf("<script src=\"js/%s\"></script>\n", js)))
	}

	return b.String()
}
