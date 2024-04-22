//main.go
package main

import (
    "fmt"
    "html/template"
    "net/http"
    "github.com/gorilla/mux"
)


func main() {
    // Parse templates
    tmpl := template.Must(template.ParseGlob("templates/*.html"))

    // Router
    router := mux.NewRouter()

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "index.html", nil)
    })
    http.ListenAndServe(":8000", router)
}

func handler (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}
