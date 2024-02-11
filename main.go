package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

//go:embed assets/css
var cssFiles embed.FS

func main() {
	var store Store

	middleware := func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	}

	cssFS := fs.FS(cssFiles)
	http.Handle("/css", middleware(http.FileServer(http.FS(cssFS))))
	// http.Handle("/css/", middleware(http.FileServer(http.Dir("assets"))))

	http.Handle("/", middleware(handleGetIndex(&store)))
	http.Handle("/todos/new", middleware(handlePostToDo(&store)))
	http.Handle("/todos/edit/", middleware(handleGetEdit(&store)))
	http.Handle("/todos/done/", middleware(handlePatchDone(&store)))
	http.Handle("/todos/update/", middleware(handlePostUpdate(&store)))
	http.Handle("/todos/delete/", middleware(handlePostDelete(&store)))
	http.Handle("/todos/clear-completed", middleware(handlePostClearCompleted(&store)))

	log.Fatal(http.ListenAndServe(":38080", logRequest(http.DefaultServeMux)))
}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
