package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"
)

//go:embed templates
var templateFiles embed.FS

type indexPage struct {
	ToDoList  []Task
	Filter    string
	ItemsLeft int
}

var indexT = must(template.ParseFS(templateFiles, "templates/index.html.tmpl"))

var itemT = must(template.ParseFS(templateFiles, "templates/includes/todo-item.tmpl"))

type itemCountView struct {
	ItemsLeft int
}

var itemCountT = must(template.ParseFS(templateFiles, "templates/includes/item-count.tmpl"))

var editItemT = must(template.ParseFS(templateFiles, "templates/includes/edit-item.tmpl"))

type todoListView struct {
	ToDoList []Task
}

var todoListT = must(template.ParseFS(templateFiles, "templates/includes/todo-list.tmpl"))

func handleGetIndex(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		filter := r.URL.Query().Get("filter")

		var filterFunc func(n Task) bool
		switch filter {
		case "active":
			filterFunc = func(n Task) bool { return n.Done }
		case "completed":
			filterFunc = func(n Task) bool { return !n.Done }
		case "all":
			fallthrough
		default:
			filterFunc = func(n Task) bool { return false }
		}
		filtered := slices.DeleteFunc(store.GetTodoList(), filterFunc)

		if err := indexT.Execute(w, indexPage{
			ToDoList:  filtered,
			Filter:    filter,
			ItemsLeft: store.GetItemLeft(),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlePostToDo(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo := r.FormValue("todo")
		newToDo := store.AddTask(todo)

		if err := itemT.Execute(w, newToDo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := itemCountT.Execute(w, itemCountView{ItemsLeft: store.GetItemLeft()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGetEdit(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/todos/edit/")

		todo, ok := store.FindToDoByID(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if err := editItemT.Execute(w, todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlePatchDone(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/todos/done/")

		if ok := store.UpdateStatus(id); !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		task, ok := store.FindToDoByID(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if err := itemT.Execute(w, task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := itemCountT.Execute(w, itemCountView{ItemsLeft: store.GetItemLeft()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlePostUpdate(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/todos/update/")

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		note := r.FormValue("todo")

		if ok := store.Update(id, note); !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		task, ok := store.FindToDoByID(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if err := itemT.Execute(w, task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := itemCountT.Execute(w, itemCountView{ItemsLeft: store.GetItemLeft()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlePostDelete(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/todos/delete/")

		if ok := store.Delete(id); !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if err := itemCountT.Execute(w, itemCountView{ItemsLeft: store.GetItemLeft()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlePostClearCompleted(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		store.ClearCompleted()

		if err := todoListT.Execute(w, todoListView{ToDoList: store.GetTodoList()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := itemCountT.Execute(w, itemCountView{ItemsLeft: store.GetItemLeft()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}
