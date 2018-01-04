package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Server is a server
type Server struct {
	db *sql.DB
}

// TodoData is tododata
type TodoData struct {
	Todos []Todo
}

// Todo is todo
type Todo struct {
	ID, CreatedOn, DueOn, Status, Description string
}

const (
	dbUser = "postgres"
	dbPass = "123456"
	dbName = "postgres"
)

func (server Server) showList(w http.ResponseWriter, r *http.Request) {
	query := fmt.Sprintf("SELECT * FROM todolist;")
	rows, err := server.db.Query(query)
	checkError(err)

	t, err := template.ParseFiles("index.html")
	checkError(err)

	var todolist TodoData

	for rows.Next() {
		todo := Todo{}

		err := rows.Scan(&todo.ID, &todo.CreatedOn, &todo.DueOn, &todo.Status, &todo.Description)
		createdDate, _ := time.Parse(time.RFC3339, todo.CreatedOn)
		todo.CreatedOn = createdDate.Format("2006-01-02")
		dueDate, _ := time.Parse(time.RFC3339, todo.DueOn)
		todo.DueOn = dueDate.Format("2006-01-02")
		checkError(err)

		todolist.Todos = append(todolist.Todos, todo)
	}

	t.Execute(w, todolist)
}

func (server Server) addItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	createdDate := time.Now().Format("2006/01/02")
	dueDate := r.Form["due"][0]
	description := r.Form["todo"][0]

	query := fmt.Sprintf("INSERT INTO todolist(created_on, due_on, status, description) VALUES('%s', '%s', 'TODO', '%s');",
		createdDate, dueDate, description)

	_, err := server.db.Exec(query)
	checkError(err)

	server.showList(w, r)
}

func (server Server) completeItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ID := r.Form["submit"][0]

	query := fmt.Sprintf("UPDATE todolist SET \"status\" = 'DONE' WHERE id = %s", ID)

	_, err := server.db.Exec(query)
	checkError(err)

	server.showList(w, r)
}

func (server Server) archiveItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ID := r.Form["submit"][0]

	query := fmt.Sprintf("DELETE FROM todolist WHERE id = %s", ID)

	_, err := server.db.Exec(query)
	checkError(err)

	server.showList(w, r)
}

func main() {
	var err error
	server := Server{}
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)
	server.db, err = sql.Open("postgres", dbInfo)

	if err != nil {
		fmt.Println("Error opening database")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", server.showList).Methods("GET")
	router.HandleFunc("/", server.addItem).Methods("POST")
	router.HandleFunc("/done", server.completeItem).Methods("POST")
	router.HandleFunc("/archive", server.archiveItem).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":8000", router))

	server.db.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
