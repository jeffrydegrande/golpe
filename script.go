package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type graph map[string][]string
type inDegree map[string]int

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

		path := filepath.Join("javascripts/", "component_"+filepath.Base(js))
		ioutil.WriteFile(path, []byte(stdout.String()), 0660)
	}
}

// General purpose topological sort, not specific to the application of
// library dependencies.  Also adapted from Wikipedia pseudo code.
func topSortDFS(g graph) (order, cyclic []string) {
	L := make([]string, len(g))
	i := len(L)
	temp := map[string]bool{}
	perm := map[string]bool{}
	var cycleFound bool
	var cycleStart string
	var visit func(string)
	visit = func(n string) {
		switch {
		case temp[n]:
			cycleFound = true
			cycleStart = n
			return
		case perm[n]:
			return
		}
		temp[n] = true
		for _, m := range g[n] {
			visit(m)
			if cycleFound {
				if cycleStart > "" {
					cyclic = append(cyclic, n)
					if n == cycleStart {
						cycleStart = ""
					}
				}
				return
			}
		}
		delete(temp, n)
		perm[n] = true
		i--
		L[i] = n
	}
	for n := range g {
		if perm[n] {
			continue
		}
		visit(n)
		if cycleFound {
			return nil, cyclic
		}
	}
	return L, nil
}

// General purpose topological sort, not specific to the application of
// library dependencies.  Adapted from Wikipedia pseudo code, one main
// difference here is that this function does not consume the input graph.
// WP refers to incoming edges, but does not really need them fully represented.
// A count of incoming edges, or the in-degree of each node is enough.  Also,
// WP stops at cycle detection and doesn't output information about the cycle.
// A little extra code at the end of this function recovers the cyclic nodes.
func topSortKahn(g graph, in inDegree) (order, cyclic []string) {
	var L, S []string
	// rem for "remaining edges," this function makes a local copy of the
	// in-degrees and consumes that instead of consuming an input.
	rem := inDegree{}
	for n, d := range in {
		if d == 0 {
			// accumulate "set of all nodes with no incoming edges"
			S = append(S, n)
		} else {
			// initialize rem from in-degree
			rem[n] = d
		}
	}
	for len(S) > 0 {
		last := len(S) - 1 // "remove a node n from S"
		n := S[last]
		S = S[:last]
		L = append(L, n) // "add n to tail of L"
		for _, m := range g[n] {
			// WP pseudo code reads "for each node m..." but it means for each
			// node m *remaining in the graph.*  We consume rem rather than
			// the graph, so "remaining in the graph" for us means rem[m] > 0.
			if rem[m] > 0 {
				rem[m]--         // "remove edge from the graph"
				if rem[m] == 0 { // if "m has no other incoming edges"
					S = append(S, m) // "insert m into S"
				}
			}
		}
	}
	// "If graph has edges," for us means a value in rem is > 0.
	for c, in := range rem {
		if in > 0 {
			// recover cyclic nodes
			for _, nb := range g[c] {
				if rem[nb] > 0 {
					cyclic = append(cyclic, c)
					break
				}
			}
		}
	}
	if len(cyclic) > 0 {
		return nil, cyclic
	}
	return L, nil
}

func buildJavascripts() string {
	// this will give us compiled jsx
	runJsxCompiler()

	javascripts, err := filepath.Glob("javascripts/*.js")
	check(err)

	javascriptsDeps, err := filepath.Glob("javascripts/**/*.js")
	check(err)

	javascripts = append(javascripts, javascriptsDeps...)

	var fullPaths map[string]string
	fullPaths = make(map[string]string)

	g := graph{}
	in := inDegree{}

	// now we have to weight javascripts, by pulling out their dependencies
	re := regexp.MustCompile("//= require (.*)")
	for _, js := range javascripts {

		jsFile := filepath.Base(js)

		fullPaths[jsFile] = js
		g[jsFile] = g[jsFile]

		lines, err := readLines(js)
		check(err)
		for _, line := range lines {
			res := re.FindStringSubmatch(line)
			if res == nil {
				continue
			}

			jsDep := filepath.Base(res[1])

			if !strings.HasSuffix(jsDep, ".js") {
				jsDep = jsDep + ".js"
			}

			in[jsDep] = in[jsDep]
			successors := g[jsDep]
			for i := 0; ; i++ {
				if i == len(successors) {
					g[jsDep] = append(successors, jsFile)
					fmt.Println("Bumping ", jsFile)
					in[jsFile]++
					break
				}
				if jsDep == successors[i] {
					break
				}
			}
		}
	}

	order, cyclic := topSortKahn(g, in)
	if cyclic != nil {
		fmt.Println("Cyclic javascript dependencies:", cyclic)
		panic("can not continue")
	}

	var b bytes.Buffer
	for _, js := range order {
		CopyFile(fullPaths[js], filepath.Join("public/js", js))
		b.Write([]byte(fmt.Sprintf("<script src=\"js/%s\"></script>\n", js)))
	}

	return b.String()
}
