package main

import (
	"database/sql"
	"fmt"
    _ "github.com/mattn/go-sqlite3"
	"os"
	"encoding/json"	
	"net/http"
	//"errors"
	//"io"
)

var dashboard string = "./index.html"
var home string = os.Getenv("HOME")
var database string = home + "/.local/inventory.sqlite3"
type jsonData struct {
	ID int
	Quantity int
	Name string
	Category string
}

func main(){
	db, err := sql.Open("sqlite3", database)
	if err != nil{
		fmt.Println(err)
	}
	defer db.Close()
    	createtable := `
		create table if not exists entries (
	    	ID integer primary key autoincrement,
	    	Name text,
	    	Quantity integer,
			Category text
		);
	`
	_, err = db.Exec(createtable)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("program is running")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", root)
	http.HandleFunc("/dashboard", dash)
	http.HandleFunc("/getinventory", func(w http.ResponseWriter, r *http.Request) {getInventory(w, r, db)})
	http.HandleFunc("/additem", func(w http.ResponseWriter, r *http.Request) {addItem(w, r, db)})
	http.HandleFunc("/edititem", func(w http.ResponseWriter, r *http.Request) {editItem(w, r, db)})
	http.HandleFunc("/deleteitem", func(w http.ResponseWriter, r *http.Request) {deleteItem(w, r, db)})

	err = http.ListenAndServe(":8000", nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	dash(w, r)
}

func dash(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		http.ServeFile(w, r, dashboard)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getInventory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "GET") {
		rows, err := db.Query("SELECT * FROM entries;")
		if err != nil {
			panic(err)
		}
		defer rows.Close()


		var entries []jsonData
		
		for rows.Next(){
			var entry jsonData
			err := rows.Scan(&entry.ID, &entry.Name, &entry.Quantity, &entry.Category)
			if err != nil {
				panic(err)
			}
			entries = append(entries, entry)
		}
		if err = rows.Err(); err != nil {
        		panic(err)
    		}
		
		Json, err := json.Marshal(entries)
		if err != nil {
			panic(err)
		}

		w.Write(Json)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func addItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "POST") {
		var Data jsonData

		err := json.NewDecoder(r.Body).Decode(&Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stmt, err := db.Prepare("INSERT INTO entries (ID, Quantity, Name, Category) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		
		_, err = stmt.Exec(Data.ID, Data.Quantity, Data.Name, Data.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}		
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func editItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "POST") {
		var Data jsonData

		err := json.NewDecoder(r.Body).Decode(&Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stmt, err := db.Prepare("UPDATE entries SET Quantity=?, Name=?, Category=? WHERE ID=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(Data.Quantity, Data.Name, Data.Category, Data.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}		
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
func deleteItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "POST") {
		var Data jsonData

		err := json.NewDecoder(r.Body).Decode(&Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stmt, err := db.Prepare("DELETE FROM  entries WHERE ID=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(Data.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}		
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

