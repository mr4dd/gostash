package main

import (
	"database/sql"
	"fmt"
    _ "github.com/go-sql-driver/mysql"
	"os"
	"encoding/json"	
	"net/http"
	//"errors"
	//"io"
)

var dashboard string = "./index.html"
	// var home string = os.Getenv("HOME")
var database string = os.Getenv("DB")
var dsn string = os.Getenv("DSN")
func init() {
	if database == "" {
		database = "root@tcp(localhost:3306)/inventory"
	}
	if dsn == "" {
		dsn = "mysql"
	}
}

type jsonData struct {
	ID int
	Quantity int
	Name string
	Category string
	Description string
}

func main(){
	db, err := sql.Open(dsn, database)
	if err != nil{
		fmt.Println(err)
	}
	defer db.Close()
	
	fmt.Println("program is running")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", root)
	http.HandleFunc("/dashboard", dash)
	http.HandleFunc("/getinventory", func(w http.ResponseWriter, r *http.Request) {getInventory(w, r, db)})
	http.HandleFunc("/gettags", func(w http.ResponseWriter, r *http.Request) {getTags(w, r, db)})
	http.HandleFunc("/additem", func(w http.ResponseWriter, r *http.Request) {addItem(w, r, db)})
	http.HandleFunc("/edititem", func(w http.ResponseWriter, r *http.Request) {editItem(w, r, db)})
	http.HandleFunc("/deleteitem", func(w http.ResponseWriter, r *http.Request) {deleteItem(w, r, db)})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {search(w, r, db)})

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

func getTags(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "GET") {
		rows, err := db.Query("SELECT DISTINCT category FROM categories;")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var entries []map[string]string
		
		for rows.Next(){
			var entry_string string
			entry := make(map[string]string)
			err := rows.Scan(&entry_string)
			if err != nil {
				panic(err)
			}
			entry["category"] = entry_string
			entries = append(entries, entry)
		}
		if err = rows.Err(); err != nil {
        		panic(err)
    		}
		
		Json, err := json.Marshal(entries)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Cache-Control", "max-age=604800")
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

		descID, err := getOrCreateDescriptionID(db, Data.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		catID, err := getOrCreateCategory(db, Data.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare("INSERT INTO entries (quantity, name, category, description) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		
		_, err = stmt.Exec(Data.Quantity, Data.Name, catID, descID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}		
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getDescription(id int, db *sql.DB) (string) {
	var description string
	err := db.QueryRow("SELECT description FROM descriptions WHERE id=?", id).Scan(&description)
	if err != nil {
		return ""
	}
	return description
}

func getInventory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "GET") {
		rows, err := db.Query(`
			SELECT e.id, e.name, e.quantity, c.category, d.description
			FROM entries e
			JOIN categories c ON e.category = c.id
			JOIN descriptions d ON e.description = d.id
		`)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var entries []jsonData
		
		for rows.Next(){
			var entry jsonData
			err := rows.Scan(&entry.ID, &entry.Name, &entry.Quantity, &entry.Category, &entry.Description)
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

func editItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "POST") {
		var Data jsonData

		err := json.NewDecoder(r.Body).Decode(&Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		catID, err := getOrCreateCategory(db, Data.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare("UPDATE entries SET quantity=?, name=?, category=? WHERE id=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		stmt.Exec(Data.Quantity, Data.Name, catID, Data.ID)

		descID, err := getDescriptionIDByText(db, Data.Description)
		if descID == 0 || err != nil {
			descID, _ = insertDescription(db, Data.Description)
		}

		stmt, err = db.Prepare("UPDATE descriptions SET description=? WHERE id=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		stmt.Exec(Data.Description, descID)
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

		stmt, _ := db.Prepare("DELETE FROM entries WHERE id=?")
		defer stmt.Close()
		stmt.Exec(Data.ID)

		stmt, _ = db.Prepare("DELETE FROM descriptions WHERE id=?")
		defer stmt.Close()
		stmt.Exec(Data.ID)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func search(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (r.Method == "POST") {
		
		var Data jsonData
		err := json.NewDecoder(r.Body).Decode(&Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT e.id, e.name, e.quantity, c.category, d.description
			FROM entries e
			JOIN categories c ON e.category = c.id
			JOIN descriptions d ON e.description = d.id
			WHERE e.name COLLATE NOCASE LIKE '%' || ? || '%'
			OR c.category COLLATE NOCASE LIKE '%' || ? || '%'
		`, Data.Name, Data.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var entries []jsonData
		for rows.Next(){
			var entry jsonData
			rows.Scan(&entry.ID, &entry.Name, &entry.Quantity, &entry.Category, &entry.Description)
			entries = append(entries, entry)
		}
		
		Json, err := json.Marshal(entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(Json)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getOrCreateDescriptionID(db *sql.DB, desc string) (int, error) {
	id, _ := getDescriptionIDByText(db, desc)
	if id != 0 {
		return id, nil
	}
	return insertDescription(db, desc)
}

func getDescriptionIDByText(db *sql.DB, desc string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM descriptions WHERE description=?", desc).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertDescription(db *sql.DB, desc string) (int, error) {
	res, err := db.Exec("INSERT INTO descriptions(description) VALUES(?)", desc)
	if err != nil {
		return 0, err
	}
	lastID, _ := res.LastInsertId()
	return int(lastID), nil
}

func getOrCreateCategory(db *sql.DB, category string) (int, error) {
	var id int

	err := db.QueryRow("SELECT id FROM categories WHERE category = ?", category).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec("INSERT INTO categories (category) VALUES (?)", category)
		if err != nil {
			return 0, err
		}
		newID, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		id = int(newID)
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
