package main

import (
	"html/template"
	"os"
)

type Data struct {
	URL string
}

func main() {
	tmpl := `<script>const SUPABASE_URL = {{.URL}};</script>`
	t := template.Must(template.New("").Parse(tmpl))
	t.Execute(os.Stdout, Data{URL: "https://example.com"})
}
