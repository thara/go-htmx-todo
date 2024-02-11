package main

import (
	"embed"
	"log"
	"net/http"
)

//go:embed assets/css
var cssFiles embed.FS

func main() {
	var store Store

	// cssFS := fs.FS(cssFiles)
	// http.Handle("/css", http.FileServer(http.FS(cssFS)))
	http.Handle("/css/", http.FileServer(http.Dir("assets")))

	http.HandleFunc("/", handleGetIndex(&store))
	http.HandleFunc("/todos/new", handlePostToDo(&store))
	http.HandleFunc("/todos/edit/", handleGetEdit(&store))
	http.HandleFunc("/todos/done/", handlePatchDone(&store))
	http.HandleFunc("/todos/update/", handlePostUpdate(&store))
	http.HandleFunc("/todos/delete/", handlePostDelete(&store))
	http.HandleFunc("/todos/clear-completed", handlePostClearCompleted(&store))

	log.Fatal(http.ListenAndServe(":38080", logRequest(http.DefaultServeMux)))
}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
