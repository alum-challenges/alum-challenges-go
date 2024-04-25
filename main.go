package main

import (
    "fmt"
    "html/template"
    "io/fs"
    "os"
    "net/http"
    "path/filepath"
    "strings"
    "github.com/gorilla/mux"
    "github.com/russross/blackfriday"
)

type Stored_Html struct {
    Html_data template.HTML
}

type Problems_List struct {
    List []string
}

func main() {
    r := mux.NewRouter()

    tmpl := template.Must(template.ParseGlob("static/*.html"))

    file_names := file_names_slice()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

        Lista := Problems_List{
            List: file_names,
        }

        tmpl.ExecuteTemplate(w, "index.html", Lista)
    })

    r.HandleFunc("/python/week-{week}/{filename}", func (w http.ResponseWriter, r *http.Request) {
        // get the weeks name from url
        vars := mux.Vars(r)
        week := vars["week"]
        file_name := vars["filename"]

        a := fmt.Sprintf("python/week-%s/%s/%s.md", week, file_name, file_name)

        data, err := os.ReadFile(a)

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

    })

    http.ListenAndServe(":8000", r)
}

func file_names_slice() []string{
    temp := []string{}
    filepath.WalkDir("python",func (path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() { // skip all directories
            count := strings.Count(path, "/")
            if count >= 2{
                temp = append(temp, path)
            }
        }
        return nil
    })
    return temp
}
