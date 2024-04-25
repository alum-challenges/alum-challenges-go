package main

import (
	"fmt"
	"net/http"
    "html/template"
    "github.com/gorilla/mux"
	"io/ioutil"
	"github.com/russross/blackfriday"
)

type Stored_Html struct {
	Html_data template.HTML
}

func main() {
	r := mux.NewRouter()

	tmpl := template.Must(template.ParseGlob("static/*.html"))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})

	r.HandleFunc("/python/{week}/{filename}", func (w http.ResponseWriter, r *http.Request) {
		// get the weeks name from url
		vars := mux.Vars(r)
		week := vars["week"]
		file_name := vars["filename"]
		
		a := fmt.Sprintf("python/week-%s/%s.md", week, file_name)
	
		data, err := ioutil.ReadFile(a)
	
		if err != nil {
			fmt.Println("Couldnt open the file")
			//! add to quit or something
		} 
			// translate the file into markdown
		markdown := blackfriday.MarkdownCommon(data)
	
		x := template.HTML(string(markdown))
	
		// add the data to the slice
		Datas := Stored_Html {
			Html_data: x,
		}
	
		// render the template
		tmpl.ExecuteTemplate(w, "layout.html", Datas)
		// tmpl.Execute(w, Datas)	
	})

    http.ListenAndServe(":8000", r)
}

