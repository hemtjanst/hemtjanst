package web

import (
	"github.com/hemtjanst/hemtjanst/device"
	"html/template"
	"net/http"
)

func indexHandler(d *device.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/index.tmpl")
		t.Execute(w, d.GetAll())
	}
}
