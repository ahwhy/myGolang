package templ

import (
	"html/template"
	"log"
	"os"
	"strings"
)

// Define a template.
const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`

// Prepare some data to insert into the template.
type Recipient struct {
	Name, Gift string
	Attended   bool
}

var recipients = []Recipient{
	{"Aunt Mildred", "bone china tea set", true},
	{"Uncle John", "moleskin pants", false},
	{"Cousin Rodney", "", false},
}

func Define() {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Parse(letter))
	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
}

// A simple template definition to test our function.
// We print the input text several ways:
// - the original
// - title-cased
// - title-cased and then printed with %q
// - printed with %q and then title-cased.
const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}
`

func Func() {
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
	}

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}

// func Glob() {
// 	// Here we create a temporary directory and populate it with our sample
// 	// template definition files; usually the template files would already
// 	// exist in some location known to the program.
// 	dir := createTestDir([]templateFile{
// 		// T0.tmpl is a plain template file that just invokes T1.
// 		{"T0.tmpl", `T0 invokes T1: ({{template "T1"}})`},
// 		// T1.tmpl defines a template, T1 that invokes T2.
// 		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
// 		// T2.tmpl defines a template T2.
// 		{"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
// 	})
// 	// Clean up after the test; another quirk of running as an example.
// 	defer os.RemoveAll(dir)
// 	// pattern is the glob pattern used to find all the template files.
// 	pattern := filepath.Join(dir, "*.tmpl")
// 	// Here starts the example proper.
// 	// T0.tmpl is the first name matched, so it becomes the starting template,
// 	// the value returned by ParseGlob.
// 	tmpl := template.Must(template.ParseGlob(pattern))
// 	err := tmpl.Execute(os.Stdout, nil)
// 	if err != nil {
// 		log.Fatalf("template execution: %s", err)
// 	}
// }

// func Helpers() {
// 	// Here we create a temporary directory and populate it with our sample
// 	// template definition files; usually the template files would already
// 	// exist in some location known to the program.
// 	dir := createTestDir([]templateFile{
// 		// T1.tmpl defines a template, T1 that invokes T2.
// 		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
// 		// T2.tmpl defines a template T2.
// 		{"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
// 	})
// 	// Clean up after the test; another quirk of running as an example.
// 	defer os.RemoveAll(dir)
// 	// pattern is the glob pattern used to find all the template files.
// 	pattern := filepath.Join(dir, "*.tmpl")
// 	// Here starts the example proper.
// 	// Load the helpers.
// 	templates := template.Must(template.ParseGlob(pattern))
// 	// Add one driver template to the bunch; we do this with an explicit template definition.
// 	_, err := templates.Parse("{{define `driver1`}}Driver 1 calls T1: ({{template `T1`}})\n{{end}}")
// 	if err != nil {
// 		log.Fatal("parsing driver1: ", err)
// 	}
// 	// Add another driver template.
// 	_, err = templates.Parse("{{define `driver2`}}Driver 2 calls T2: ({{template `T2`}})\n{{end}}")
// 	if err != nil {
// 		log.Fatal("parsing driver2: ", err)
// 	}
// 	// We load all the templates before execution. This package does not require
// 	// that behavior but html/template's escaping does, so it's a good habit.
// 	err = templates.ExecuteTemplate(os.Stdout, "driver1", nil)
// 	if err != nil {
// 		log.Fatalf("driver1 execution: %s", err)
// 	}
// 	err = templates.ExecuteTemplate(os.Stdout, "driver2", nil)
// 	if err != nil {
// 		log.Fatalf("driver2 execution: %s", err)
// 	}
// }

// func Share() {
// 	// Here we create a temporary directory and populate it with our sample
// 	// template definition files; usually the template files would already
// 	// exist in some location known to the program.
// 	dir := createTestDir([]templateFile{
// 		// T0.tmpl is a plain template file that just invokes T1.
// 		{"T0.tmpl", "T0 ({{.}} version) invokes T1: ({{template `T1`}})\n"},
// 		// T1.tmpl defines a template, T1 that invokes T2. Note T2 is not defined
// 		{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
// 	})
// 	// Clean up after the test; another quirk of running as an example.
// 	defer os.RemoveAll(dir)
// 	// pattern is the glob pattern used to find all the template files.
// 	pattern := filepath.Join(dir, "*.tmpl")
// 	// Here starts the example proper.
// 	// Load the drivers.
// 	drivers := template.Must(template.ParseGlob(pattern))
// 	// We must define an implementation of the T2 template. First we clone
// 	// the drivers, then add a definition of T2 to the template name space.
// 	// 1. Clone the helper set to create a new name space from which to run them.
// 	first, err := drivers.Clone()
// 	if err != nil {
// 		log.Fatal("cloning helpers: ", err)
// 	}a
// 	// 2. Define T2, version A, and parse it.
// 	_, err = first.Parse("{{define `T2`}}T2, version A{{end}}")
// 	if err != nil {
// 		log.Fatal("parsing T2: ", err)
// 	}
// 	// Now repeat the whole thing, using a different version of T2.
// 	// 1. Clone the drivers.
// 	second, err := drivers.Clone()
// 	if err != nil {
// 		log.Fatal("cloning drivers: ", err)
// 	}
// 	// 2. Define T2, version B, and parse it.
// 	_, err = second.Parse("{{define `T2`}}T2, version B{{end}}")
// 	if err != nil {
// 		log.Fatal("parsing T2: ", err)
// 	}
// 	// Execute the templates in the reverse order to verify the
// 	// first is unaffected by the second.
// 	err = second.ExecuteTemplate(os.Stdout, "T0.tmpl", "second")
// 	if err != nil {
// 		log.Fatalf("second execution: %s", err)
// 	}
// 	err = first.ExecuteTemplate(os.Stdout, "T0.tmpl", "first")
// 	if err != nil {
// 		log.Fatalf("first: execution: %s", err)
// 	}
// }
