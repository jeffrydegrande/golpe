package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type command struct {
	name        string
	description string
}

var (
	commands = []command{
		{"run", "runs project in current directoy"},
		{"create", "sets up new project"},
		{"build", "builds the project"},
		{"new", "add a new file"},
	}
)

type Cli struct {
}

func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) getMethod(args ...string) (func(...string) error, bool) {
	camelArgs := make([]string, len(args))
	for i, s := range args {
		if len(s) == 0 {
			return nil, false
		}
		camelArgs[i] = strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	methodName := "Cmd" + strings.Join(camelArgs, "")
	method := reflect.ValueOf(cli).MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("Method is not valid\n")
		return nil, false
	}

	return method.Interface().(func(...string) error), true
}

func (cli *Cli) Cmd(args ...string) error {
	if len(args) > 1 {
		method, exists := cli.getMethod(args[:2]...)
		if exists {
			return method(args[2:]...)
		}
	}
	if len(args) > 0 {
		method, exists := cli.getMethod(args[0])
		if !exists {
			return fmt.Errorf("%s is not a command\n", args[0])
		}
		return method(args[1:]...)
	}
	return cli.CmdHelp()
}

func (cli *Cli) CmdHelp(args ...string) error {
	flag.Usage()
	return nil
}

func (cli *Cli) CmdCreate(args ...string) error {
	fmt.Printf("Create new project")
	return nil
}

func (cli *Cli) CmdBuild(args ...string) error {
	fmt.Printf("Building\n")
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

func (cli *Cli) CmdRun(args ...string) error {
	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Println("Listening on port 3000.")
	http.ListenAndServe(":3000", nil)
	return nil
}

func (cli *Cli) CmdNew(args ...string) error {
	fmt.Printf("New thingy")
	return nil
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stdout, "Usage: sitegen [OPTIONS] COMMAND\n\nA static web project tool")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()

		help := "\nCommands:\n"

		for _, cmd := range commands {
			help += fmt.Sprintf("	%-10.10s%s\n", cmd.name, cmd.description)
		}

		help += "\nRun 'sitegen COMMAND --help' for mor information on a command."
		fmt.Fprintf(os.Stdout, "%s\n", help)
	}
}

func main() {
	flag.Parse()
	cli := NewCli()
	cli.Cmd(flag.Args()...)
}
