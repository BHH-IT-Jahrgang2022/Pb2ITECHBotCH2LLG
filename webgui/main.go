package main

import (
	"embed"
	"flag"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var (
	// Flags
	port = flag.String("port", "8080", "Port to listen on")
)

//go:embed templates
var tmplEmbed embed.FS

//go:embed static
var staticEmbedFS embed.FS

type staticFS struct {
	fs fs.FS
}

func (sfs *staticFS) Open(name string) (fs.File, error) {
	return sfs.fs.Open(filepath.Join("static", name))
}

var staticEmbed = &staticFS{staticEmbedFS}

func main() {
	flag.Parse()

	// Create a new Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", filepath.Join("webgui", "static"))

	// Load templates
	tmpl, err := fs.Sub(tmplEmbed, "templates")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(gin.HTMLDebugMode(gin.DefaultHTMLRender, tmpl))

	// Index page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run("0.0.0.0:" + *port)
}
