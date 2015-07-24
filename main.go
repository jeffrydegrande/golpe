package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/daviddengcn/go-colortext"
)

func check(e error) {
	if e != nil {
		ct.ChangeColor(ct.Red, true, ct.None, false)
		panic(e)
		ct.ResetColor()
	}
}

func say(msg string) {
	ct.ChangeColor(ct.Green, true, ct.None, true)
	log.Println(msg)
	ct.ResetColor()
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
		{"watch", "watch for changes in the current directory and build."},
	}
)

type cli struct {
}

func newCli() *cli {
	return &cli{}
}

func (cli *cli) getMethod(args ...string) (func(...string) error, bool) {
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

func (cli *cli) cmd(args ...string) error {
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
	return cli.cmdHelp()
}

func (cli *cli) cmdHelp(args ...string) error {
	flag.Usage()
	return nil
}

func (cli *cli) cmdCreate(args ...string) error {
	say("Create new project")
	return nil
}

func (cli *cli) cmdBuild(args ...string) error {
	say("Building")
	buildAll()
	return nil
}

func (cli *cli) cmdRun(args ...string) error {
	say("Listening...")
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":3000", nil)
	return nil
}

func (cli *cli) cmdNew(args ...string) error {
	say("New thingy")
	return nil
}

func (cli *cli) cmdWatch(args ...string) error {
	say("Watching current directory")
	watch()
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
	cli := newCli()
	cli.cmd(flag.Args()...)
}
