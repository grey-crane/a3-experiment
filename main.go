package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/julienschmidt/httprouter"
)

//go:embed site
var embedded embed.FS

func main() {
	// Log service version.
	fmt.Println("FASTLY_SERVICE_VERSION:", os.Getenv("FASTLY_SERVICE_VERSION"))

	// Re-root the embedded FS so "site/" is not part of request paths.
	site, err := fs.Sub(embedded, "site")
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	handler := staticHandler(site)
	router.GET("/*filepath", handler)
	router.HEAD("/*filepath", handler)

	fsthttp.Serve(fsthttp.Adapt(router))
}

func resolve(site fs.FS, name string) (string, fs.FileInfo, bool) {
	candidates := []string{"index.html"}
	if name != "" {
		candidates = []string{name, name + ".html", path.Join(name, "index.html")}
	}
	for _, c := range candidates {
		// returns first matching candidate
		if info, err := fs.Stat(site, c); err == nil && !info.IsDir() {
			return c, info, true
		}
	}
	return "", nil, false
}

func staticHandler(site fs.FS) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Rooting the path before cleaning keeps ".." from escaping the site dir.
		name := strings.TrimPrefix(path.Clean("/"+ps.ByName("filepath")), "/")

		file, info, ok := resolve(site, name)
		if !ok {
			http.NotFound(w, r)
			return
		}
		f, err := site.Open(file)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		ctype := mime.TypeByExtension(path.Ext(file)) // get mime type

		h := w.Header()
		h.Set("Content-Type", ctype)
		h.Set("Content-Length", strconv.FormatInt(info.Size(), 10))
		h.Set("Cache-Control", "no-cache")

		if r.Method == http.MethodHead {
			return
		}
		io.Copy(w, f)
	}
}
