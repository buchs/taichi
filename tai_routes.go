package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/antonholmquist/jason.v1"
	"gopkg.in/go-chi/chi.v4"
	_ "gopkg.in/mattn/go-sqlite3.v1"
)

// Helper Functions

// test membership in array of ints
func membership(members []int, member_id string) bool {
	member_id_num, err := strconv.Atoi(member_id)
	if err != nil {
		panic(err)
	}
	for _, a := range members {
		if a == member_id_num {
			return true
		}
	}
	members = append(members, member_id_num)
	return false
}

// the routing destination functions

func route_create_member(w http.ResponseWriter, r *http.Request) {

	bodyJson, err := jason.NewObjectFromReader(r.Body)
	if err != nil {
		panic(err)
	}
	sql_statement_template := `INSERT INTO members (name, thetype, data)
	     VALUES ('%s', '%s', '%s');`
	name, err := bodyJson.GetString("name")
	if err != nil {
		panic(err)
	}
	thetype, err := bodyJson.GetString("thetype")
	if err != nil {
		panic(err)
	}
	data, err := bodyJson.GetString("data")
	if err != nil {
		panic(err)
	}
	sql_statements := fmt.Sprintf(sql_statement_template,
		name, thetype, data)
	results, err := db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		return
	}
	results1, _ := results.LastInsertId()
	response := fmt.Sprintf(`{ "member_id": "%d" }`, results1)
	w.WriteHeader(201)
	w.Write([]byte(response))
}

func route_create_tag(w http.ResponseWriter, r *http.Request) {

	bodyJson, err := jason.NewObjectFromReader(r.Body)
	if err != nil {
		panic(err)
	}

	sql_statement_template := `INSERT INTO tags (member_id, tag)
		 VALUES (%d, '%s');`
	raw_member, err := bodyJson.GetString("member_id")
	if err != nil {
		panic(err)
	}
	member_id, err := strconv.Atoi(raw_member)
	if err != nil {
		panic(err)
	}
	tag, err := bodyJson.GetString("tag")
	if err != nil {
		panic(err)
	}
	sql_statements := fmt.Sprintf(sql_statement_template,
		member_id, tag)
	// results, err := db.Exec(sql_statements)
	_, err = db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		return
	}
	// results1, _ := results.LastInsertId()
	response := `{ "response": "n/a" }`
	w.WriteHeader(201)
	w.Write([]byte(response))
}

func route_read_member(w http.ResponseWriter, r *http.Request) {

	// request: GET /read-member/Kevin+Buchs
	// response: 200, json = [ { "member_id": "1", "name": "Kevin Buchs",
	//    "thetype": "employee", "data": "Software Engineer"} ]
	name, err := url.QueryUnescape(chi.URLParam(r, "name"))
	if err != nil {
		panic(err)
	}
	sprintf_template :=
		`SELECT rowid,name,thetype,data from members where name = '%s';`
	sql_statements := fmt.Sprintf(sprintf_template, name)
	rows, err := db.Query(sql_statements)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var rmemberid, rname, rthetype, rdata string
	termination := ""
	rerow := "[ "
	for rows.Next() {
		rerow += termination
		termination = ", "
		err := rows.Scan(&rmemberid, &rname, &rthetype, &rdata)
		if err != nil {
			panic(err)
		}
		rerow += fmt.Sprintf(`{ "member_id": "%s", "name": "%s", `+
			`"thetype": "%s", "data": "%s" }`,
			rmemberid, rname, rthetype, rdata)
	}
	rerow += " ]"
	w.Write([]byte(rerow))
}

func route_read_all_members(w http.ResponseWriter, r *http.Request) {

	// request: GET /read-all-members
	// response: 200, json = [ {"member_id": "1", "name": "Kevin Buchs", "thetype":
	//   "employee", "data": "Software Engineer"}, {"member_id": "2", "name": "Adam
	//   Buchs", "thetype": "contractor", "data": "6"} ]
	sql_statements := "SELECT rowid,name,thetype,data from members;"
	rows, err := db.Query(sql_statements)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var rmemberid, rname, rthetype, rdata, rerow, termination string
	termination = ""
	rerow = "[ "
	for rows.Next() {
		rerow += termination
		termination = ", "
		err := rows.Scan(&rmemberid, &rname, &rthetype, &rdata)
		if err != nil {
			panic(err)
		}
		rerow += fmt.Sprintf(`{ "member_id": "%s", "name": "%s", "thetype": "%s", "data": "%s" }`,
			rmemberid, rname, rthetype, rdata)
	}
	rerow += " ]"
	w.Write([]byte(rerow))
}

func route_find_tags(w http.ResponseWriter, r *http.Request) {

	sql_statement_template :=
		`SELECT rowid,name,thetype,data from members
		 WHERE rowid in 
		   (SELECT member_id from tags 
			  WHERE tags.tag = '%s');`

	bodyJson, err := jason.NewObjectFromReader(r.Body)
	if err != nil {
		panic(err)
	}

	tags_array, err := bodyJson.GetStringArray("tags")
	if err != nil {
		panic(err)
	}

	var rmemberid, rname, rthetype, rdata string
	members_in := make([]int, 100)
	termination := ""
	rerow := "[ "

	for _, tag := range tags_array {
		sql_statements := fmt.Sprintf(sql_statement_template, tag)
		rows, err := db.Query(sql_statements)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			rerow += termination
			termination = ", "
			err := rows.Scan(&rmemberid, &rname, &rthetype, &rdata)
			if err != nil {
				panic(err)
			}
			if membership(members_in, rmemberid) {
				continue // already have this one in the results
			}

			rerow += fmt.Sprintf(`{ "member_id": "%s", "name": "%s", "thetype": "%s", "data": "%s" }`,
				rmemberid, rname, rthetype, rdata)
		}
	}
	rerow += " ]"
	w.Write([]byte(rerow))
}

func route_update_type(w http.ResponseWriter, r *http.Request) {

	bodyJson, err := jason.NewObjectFromReader(r.Body)
	if err != nil {
		panic(err)
	}

	member_id, err := bodyJson.GetString("member_id")
	if err != nil {
		panic(err)
	}
	thetype, err := bodyJson.GetString("thetype")
	if err != nil {
		panic(err)
	}
	data, err := bodyJson.GetString("data")
	if err != nil {
		panic(err)
	}

	sql_statement_template := `UPDATE members SET thetype = '%s',
	    data = '%s'  WHERE rowid = '%s';`
	sql_statements := fmt.Sprintf(sql_statement_template,
		thetype, data, member_id)

	_, err = db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		panic(err)
	}
	w.WriteHeader(201)
	w.Write([]byte("updated"))
}

func route_update_name(w http.ResponseWriter, r *http.Request) {

	bodyJson, err := jason.NewObjectFromReader(r.Body)
	if err != nil {
		panic(err)
	}

	member_id, err := bodyJson.GetString("member_id")
	if err != nil {
		panic(err)
	}
	name, err := bodyJson.GetString("name")
	if err != nil {
		panic(err)
	}

	sql_statement_template := `UPDATE members SET name = '%s'
	    WHERE rowid = '%s';`
	sql_statements := fmt.Sprintf(sql_statement_template,
		name, member_id)

	_, err = db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		panic(err)
	}
	w.WriteHeader(201)
	w.Write([]byte("updated"))
}

func route_delete_tag(w http.ResponseWriter, r *http.Request) {

	bodyJson, err := jason.NewObjectFromReader(r.Body)
	if err != nil {
		panic(err)
	}

	member_id, err := bodyJson.GetString("member_id")
	if err != nil {
		panic(err)
	}
	tag, err := bodyJson.GetString("tag")
	if err != nil {
		panic(err)
	}

	sql_statement_template := `DELETE FROM tags WHERE member_id = '%s'
	    AND tag = '%s';`
	sql_statements := fmt.Sprintf(sql_statement_template,
		member_id, tag)

	_, err = db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		panic(err)
	}
	w.WriteHeader(200)
	w.Write([]byte("deleted"))
}

func route_delete_member(w http.ResponseWriter, r *http.Request) {

	memberid := chi.URLParam(r, "memberid")

	sql_statement_template := `DELETE FROM tags WHERE member_id = '%s';`
	sql_statements := fmt.Sprintf(sql_statement_template, memberid)

	_, err := db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		panic(err)
	}

	sql_statement_template = `DELETE FROM members WHERE rowid = '%s';`
	sql_statements = fmt.Sprintf(sql_statement_template, memberid)

	_, err = db.Exec(sql_statements)
	if err != nil {
		fmt.Printf("Failed to execute sql statements\n%q\n%s\n",
			err, sql_statements)
		panic(err)
	}

	w.WriteHeader(200)
	w.Write([]byte("deleted"))
}
