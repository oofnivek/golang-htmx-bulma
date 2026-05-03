package view

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin/render"
)

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer(rootDir string) *Renderer {
	r := &Renderer{
		templates: make(map[string]*template.Template),
	}

	// Load layouts and partials
	var layoutFiles []string
	filepath.Walk(rootDir+"/layouts", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".html" {
			layoutFiles = append(layoutFiles, path)
		}
		return nil
	})

	var partialFiles []string
	filepath.Walk(rootDir+"/partials", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".html" {
			partialFiles = append(partialFiles, path)
		}
		return nil
	})

	// Load all files (pages and partials) and create dedicated template sets
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".html" {
			rel, _ := filepath.Rel(rootDir, path)
			
			// Skip files in layouts/ for individual sets (they are included as support)
			if filepath.HasPrefix(rel, "layouts/") {
				return nil
			}

			files := append([]string{path}, layoutFiles...)
			files = append(files, partialFiles...)
			
			tmpl := template.New(filepath.Base(path)).Funcs(template.FuncMap{
				"in_timezone": func(t time.Time, tz string) string {
					loc, err := time.LoadLocation(tz)
					if err != nil {
						loc = time.UTC
					}
					return t.In(loc).Format("2006-01-02 03:04:05 PM")
				},
				"dict": func(values ...interface{}) (map[string]interface{}, error) {
					if len(values)%2 != 0 {
						return nil, fmt.Errorf("invalid dict call")
					}
					dict := make(map[string]interface{}, len(values)/2)
					for i := 0; i < len(values); i += 2 {
						key, ok := values[i].(string)
						if !ok {
							return nil, fmt.Errorf("dict keys must be strings")
						}
						dict[key] = values[i+1]
					}
					return dict, nil
				},
				"add": func(a, b int) int { return a + b },
				"sub": func(a, b int) int { return a - b },
				"seq": func(start, end int) []int {
					var s []int
					for i := start; i <= end; i++ {
						s = append(s, i)
					}
					return s
				},
				"safe_id": func(s string) string {
					res := ""
					for _, r := range s {
						if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
							res += string(r)
						} else {
							res += "-"
						}
					}
					return res
				},
			})
			r.templates[rel] = template.Must(tmpl.ParseFiles(files...))
		}
		return nil
	})

	return r
}

func (r *Renderer) Instance(name string, data interface{}) render.Render {
	tmpl, ok := r.templates[name]
	if !ok {
		// Fallback if not found, though template.Must should have caught most
		return render.HTML{
			Template: nil,
			Name:     name,
			Data:     data,
		}
	}

	// Determine entry point
	entryPoint := "layout"
	// If the template name is a partial or doesn't have a 'layout' defined, use the name itself
	// Actually, we can check if the requested template is in pages/ or not.
	if !filepath.HasPrefix(name, "pages/") {
		entryPoint = name
	}

	return render.HTML{
		Template: tmpl,
		Name:     entryPoint,
		Data:     data,
	}
}

func SafeID(s string) string {
	res := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			res += string(r)
		} else {
			res += "-"
		}
	}
	return res
}
