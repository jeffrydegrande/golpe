package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type API struct {
	Version        string          `json:"_version"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Metadata       []Metadata      `json:"metadata"`
	ResourceGroups []ResourceGroup `json:"resourceGroups"`
	Content        []Content       `json:"content"`
}

type Host struct {
	Value string `json:"value"`
}

type Format struct {
	Value string `json:"value"`
}

type Metadata struct {
	Format Format `json:"FORMAT"`
	Host   Host   `json:"HOST"`
}

type ResourceGroup struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Resources   []Resource `json:"resources"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Model struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Headers     []Header `json:"headers"`
	Body        string   `json:"body"`
	Schema      string   `json:"schema"`
}

type Parameter struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Default     string   `json:"default"`
	Example     string   `json:"example"`
	Values      []string `json:"values"`
}

type Example struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Requests    []Request  `json:"requests"`
	Responses   []Response `json:"responses"`
}

type Request struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Headers     []Header `json:"headers"`
	Body        string   `json:"body"`
	Schema      string   `json:"schema"`
}

type Response struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Headers     []Header `json:"headers"`
	Body        string   `json:"body"`
	Schema      string   `json:"schema"`
}

type Action struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Method      string      `json:"method"`
	Parameters  []Parameter `json:"parameters"`
	Headers     []Header    `json:"headers"`
	Examples    []Example   `json:"examples"`
}

type Resource struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	UriTemplate string      `json:"uriTemplate"`
	Model       Model       `json:"model"`
	Parameters  []Parameter `json:"parameters"`
	Headers     []Header    `json:"headers"`
	Actions     []Action    `json:"actions"`
}

type VariableName struct {
	Literal  string `json:"name"`
	Variable bool   `json:"variable"`
}

type TypeSpecification struct {
	Name        string         `json:"name"`
	NestedTypes []VariableName `json:"nestedTypes"`
}

type TypeDefinition struct {
}

type Content struct {
	Element string       `json:"element"`
	Name    VariableName `json:"name"`
}

func parseJSON(r io.Reader) (*API, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	api := new(API)
	err = json.Unmarshal(b, &api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

func parseMarkdown(r io.Reader) ([]byte, error) {
	path, err := exec.LookPath("drafter")
	if err != nil {
		return nil, errors.New("Couldn't find drafter. Please install it first https://github.com/apiaryio/drafter")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	echo := exec.Command("echo", string(b))
	out, err := echo.StdoutPipe()
	if err != nil {
		return nil, err
	}

	echo.Start()

	cmd := exec.Command(path, "--format", "json")
	cmd.Stdin = out
	return cmd.Output()
}

func ReadBlueprint(filename string) (*API, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	json, err := parseMarkdown(f)

	if err != nil {
		return nil, err
	}

	api, err := parseJSON(bytes.NewBuffer(json))
	return api, nil
}

func (api *API) RunService() {
	for _, group := range api.ResourceGroups {
		fmt.Printf("%s %s\n", group.Name, group.Description)
		for _, resource := range group.Resources {
			for _, action := range resource.Actions {

				if action.Method == "GET" {
					for _, example := range action.Examples {
						for _, response := range example.Responses {

							handler := func(w http.ResponseWriter, r *http.Request) {
								fmt.Println(r.RequestURI)
								httpStatus, _ := strconv.Atoi(response.Name)
								w.WriteHeader(httpStatus)

								for _, header := range response.Headers {
									w.Header().Set(header.Name, header.Value)
								}
								io.WriteString(w, response.Body)
							}

							log.Println("Adding route for", resource.UriTemplate)

							http.HandleFunc(resource.UriTemplate, handler)
							http.HandleFunc(fmt.Sprintf("%s/", resource.UriTemplate), handler)
						}
					}
				}
			}
		}
	}
}

func (api *API) GenerateService() {
	for _, c := range api.Content {
		fmt.Printf("%v", c)
	}
	/*
		for _, group := range api.ResourceGroups {
			fmt.Printf("%s %s\n", group.Name, group.Description)
			for _, resource := range group.Resources {
				fmt.Printf("url: %s\n", resource.UriTemplate)
				fmt.Printf("Parameters:\n")
				for _, parameter := range resource.Parameters {

					fmt.Printf("{name: %s, default: %s, type: %s, example: %s, description: %s, values: %s, required: %t}\n",
						parameter.Name, parameter.Default, parameter.Type, parameter.Example, parameter.Description, parameter.Values, parameter.Required)
				}

				fmt.Printf("Actions:\n")
				for _, action := range resource.Actions {
					fmt.Printf("%v\n", action)
				}
			}
		}
	*/

}
