package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
        "path/filepath"

	"github.com/daviddengcn/go-colortext"
	"github.com/gorilla/handlers"
)


var (
        lastModificationTimestamp int64 = 0
        directoriesToWatch []string
)

func init() {
        say("Building list of directories")
        buildListOfDirectories()
}

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
		{"api", "generate APIs"},
	}
)

type Cli struct {
}

func NewCli() *Cli {
	return &Cli{}
}

func (cli *Cli) getMethod(args ...string) (func(...string) error, bool) {
	if len(args) == 0 {
		return nil, false
	}

	/*
		camelArgs := make([]string, len(args))

		for i, s := range args {
			if len(s) == 0 {
				return nil, false
			}
			camelArgs[i] = strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		}
	*/

	// methodName := "Cmd" + strings.Join(camelArgs, "")
	methodName := "Cmd" + strings.ToUpper(args[0][:1]) + strings.ToLower(args[0][1:])
	method := reflect.ValueOf(cli).MethodByName(methodName)

	if !method.IsValid() {
		fmt.Printf("Method is not valid\n")
		return nil, false
	}
	return method.Interface().(func(...string) error), true
}

func (cli *Cli) Cmd(args ...string) error {
	if len(args) > 1 {
		method, exists := cli.getMethod(args[:1]...)
		if exists {
			return method(args[1:]...)
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
	say("Create new project")
	return nil
}

func (cli *Cli) CmdBuild(args ...string) error {
	say("Starting build")
	BuildAll()
	return nil
}

func (cli *Cli) CmdScript(args ...string) error {
	result := buildJavascriptsFromAppJs()
	fmt.Printf("%s", result)
	return nil
}

func buildListOfDirectories() error {
        cwd, err := os.Getwd()

        err = filepath.Walk(cwd, func(path string, f os.FileInfo, err error) error {
                if f.IsDir() {
                        base := filepath.Base(path)
                        toSkip := []string{".git", "public", ".module-cache", "components"}

                        for _, skip := range toSkip { 
                                if (base == skip) {
                                  return filepath.SkipDir
                                }
                        }
                        directoriesToWatch = append(directoriesToWatch, path)
                }
                return nil
        })
        return err
} 

func shouldRebuild() bool {
        for _, path := range directoriesToWatch {

                fi, err := os.Stat(path)
                if err != nil {
                        log.Fatal("stat failed", path)
                        continue
                }

                var currentTimeStamp = fi.ModTime().Unix()

                if currentTimeStamp > lastModificationTimestamp {
                        say(fmt.Sprintf("%s has changed. Rebuilding", path))
                        lastModificationTimestamp = currentTimeStamp
                        return true
                }
        }

        return false
}

func middlewareRebuildOnChanges(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
              if strings.Contains(r.URL.Path, ".html") && shouldRebuild() {
                      BuildAll()
              }
              next.ServeHTTP(w, r)
        })
}

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
                log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	})
}


func (cli *Cli) CmdRun(args ...string) error {
	say("Listening...")

	http.Handle("/", middlewareLogging(
                         middlewareRebuildOnChanges(
                         http.FileServer(http.Dir("public")))))
	http.ListenAndServe(":3000", nil)
	return nil
}

func (cli *Cli) CmdNew(args ...string) error {
	say("New thingy")
	return nil
}

func (cli *Cli) CmdApi(args ...string) error {
	say("Building API")

	for _, arg := range args {
		api, err := ReadBlueprint(arg)
		check(err)
		api.RunService()
		// api.GenerateService()
	}
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
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

	start := time.Now()
	cli.Cmd(flag.Args()...)
	delta := time.Now().Sub(start)
	fmt.Printf("Took %0.3fs\n", delta.Seconds())

}
