package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("invalid args count: %d", len(os.Args)-1)
	}

	pkg, types, out := os.Args[1], strings.Split(os.Args[2], ","), os.Args[3]
	if err := run(pkg, types, out); err != nil {
		log.Fatal(err)
	}

	p, _ := os.Getwd()
	fmt.Printf("%v generated\n", filepath.Join(p, out)) //nolint:forbidigo
}

func run(pkg string, types []string, outFile string) error {
	file, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(file)
	t, err := template.New("genTypes").Parse(codeForTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %v", err)
	}

	err = t.Execute(writer, &struct {
		Package   string
		TypeNames []string
		LastIndex int
	}{
		Package:   pkg,
		TypeNames: types,
		LastIndex: len(types) - 1,
	})
	if err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("flush file: %v", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("close file %v", err)
	}
	return nil
}

var codeForTemplate = `// Code generated by cmd/gen-types; DO NOT EDIT.

package {{ .Package }}

import (
	"fmt"

	"github.com/google/uuid"
)

{{range $typeName := .TypeNames }}
type {{$typeName}} struct {
	uuid.UUID
}

var {{$typeName}}Nil = {{$typeName}}{uuid.Nil}

func New{{$typeName}}() {{$typeName}} {
	return {{$typeName}}{
		UUID: uuid.New(),
	}
}

func (r {{$typeName}}) Validate() error {
	if r.UUID == uuid.Nil {
		return fmt.Errorf("validate error")
	}
	return nil
}

func (r {{$typeName}}) Matches(x any) bool {
	_, ok := x.({{$typeName}})
	return ok
}

func (r {{$typeName}}) IsZero() bool {
	return r.UUID == uuid.Nil
}

{{end -}}

type Types interface {
	{{range $idx, $name := .TypeNames -}}
	{{$name}}{{if lt $idx $.LastIndex}} | {{end}} 
	{{- end}}
}

func Parse[T Types](s string) (T, error) {
	var t T
	u, err := uuid.Parse(s)
	if err != nil {
		return t, err
	}
	switch any(t).(type) {
	{{range $name := .TypeNames -}}
	case {{$name}}:
		return T({{$name}}{u}), nil
	{{end -}}
	default:
		return t, fmt.Errorf("wrong type")
	}
}

func MustParse[T Types](s string) T {
	u := uuid.MustParse(s)

	var t T
	switch any(t).(type) {
	{{range $name := .TypeNames -}}
	case {{$name}}:
		return T({{$name}}{u})
	{{end -}}
	default:
		panic("wrong type")
	}
}

`
