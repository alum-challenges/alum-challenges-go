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
    "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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

    r.HandleFunc("/courses/{course}/week-{week}/{exercise}/", func (w http.ResponseWriter, r *http.Request) {
        // get the weeks name from url
        vars := mux.Vars(r)
        course := filepath.Clean(vars["course"])
        week := filepath.Clean(vars["week"])
        exercise := filepath.Clean(vars["exercise"])

        a := filepath.Join("courses", course, "week-" + week,exercise, exercise + ".md")
        data, err := os.ReadFile(a)

        if err != nil {
            fmt.Println(err)
        } 
        // translate the file into markdown
        markdown := mdToHTML(data)
        
        // fmt.Println(markdown)
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
    filepath.WalkDir("courses",func (path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() { // skip all directories
            count := strings.Count(path, "/")
            if count == 3{
                temp = append(temp, path)
            }
        }
        return nil
    })
    return temp
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
