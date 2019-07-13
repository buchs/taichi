package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/go-chi/chi.v4"
	_ "gopkg.in/mattn/go-sqlite3.v1"
)

// globals
var db *sql.DB

func setup_db() (bool, *sql.DB) {

	var sql_statements string
	var sqliteFile string
	var createTables bool

	db = nil
	environ := os.Getenv("TAI_ENVIRONMENT")
	if environ == "test" {
		sqliteFile = "tai_db_dev.db"
		os.Remove(sqliteFile) // wipe file before we start
	} else {
		sqliteFile = "tai_db.db"
	}

	_, err := os.Stat(sqliteFile)
	if os.IsNotExist(err) {
		createTables = true
	} else {
		createTables = false
	}

	db, err = sql.Open("sqlite3", sqliteFile)
	if err != nil {
		fmt.Println("Failed to open sqliteFile: " + sqliteFile)
		return true, nil
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	if createTables {
		sql_statements =
			`CREATE TABLE IF NOT EXISTS members
			  ( name TEXT NOT NULL,
			    thetype TEXT NOT NULL,
			    data TEXT NOT NULL);
			CREATE TABLE IF NOT EXISTS tags
			  ( member_id INTEGER NOT NULL,
		        tag TEXT NOT NULL);
		`
		_, err = db.Exec(sql_statements)
		if err != nil {
			fmt.Printf("Failed to execute sql statements\n%q\n%s\n", err, sql_statements)
			return true, nil
		}
	}
	return false, db
}

func main() {

	err, _ := setup_db()
	if err {
		fmt.Println("failed to setup database")
		return
	}

	r := chi.NewRouter()
	r.Post("/create-member", route_create_member)
	r.Post("/create-tag", route_create_tag)
	r.Get("/read-member/{name}", route_read_member)
	r.Get("/read-all-members", route_read_all_members)
	r.Post("/find-member-tags", route_find_tags)
	r.Post("/update-type", route_update_type)
	r.Post("/update-name", route_update_name)
	r.Delete("/delete-member-tag", route_delete_tag)
	r.Delete("/delete-member/{memberid}", route_delete_member)
	http.ListenAndServe(":3000", r)

	// fmt.Println("and I reached here")
	// db.Close()

}
