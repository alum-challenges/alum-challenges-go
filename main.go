package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	// "os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed courses
var coursesFS embed.FS

type Stored_Html struct {
    Html_data template.HTML
}

type Problems_List struct {
    List []string
}

func main() {
    r := http.NewServeMux()

    tmpl := template.Must(template.ParseGlob("static/*.html"))

    file_names := file_names_slice()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        Lista := Problems_List{
            List: file_names,
        }
        tmpl.ExecuteTemplate(w, "index.html", Lista)
    })

    r.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusGone) })
    r.Handle("GET /courses/{course}/{week}/{exercise}/", http.FileServerFS(coursesFS))

    r.HandleFunc("GET /courses/{course}/{week}/{exercise}/{$}", func(w http.ResponseWriter, r *http.Request) {
        course := r.PathValue("course")
        week := r.PathValue("week")
        exercise := r.PathValue("exercise")

        a := filepath.Join("courses", course, week, exercise, exercise+".md")
        data, err := coursesFS.ReadFile(a)

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
