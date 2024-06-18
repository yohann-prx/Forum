package server

import "html/template"

var templates = template.Must(template.ParseGlob("./web/templates/*.html"))
