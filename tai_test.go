package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

const hostname = "http://localhost:3000"

// Helper functions
func postHelper(url string, body string, expectedStatus int, expectedBody string) (bool, string) {
	// the first element returned indicates if there was an error with a true value
	contentType := "application/json"
	bodyBytes := []byte(body)
	response, err := http.Post(hostname+url, contentType, bytes.NewBuffer(bodyBytes))
	if err != nil {
		message := fmt.Sprintf("Error in http post request to %s\n", url) +
			fmt.Sprint(err)
		return true, message
	}
	defer response.Body.Close()
	if expectedStatus != response.StatusCode {
		message := "post status=" + strconv.Itoa(response.StatusCode) +
			", expected=" + strconv.Itoa(expectedStatus) +
			" url=" + url
		return true, message
	}

	if expectedBody == "" { // no body testing required
		return false, "n/a"
	}
	rawrespbody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		message := "http post, unable to read body from " + url
		return true, message
	}
	respbody := string(rawrespbody)
	if expectedBody == respbody {
		return false, "n/a"
	}

	message := fmt.Sprintf("http post body: %s\ndid not match:%s\nurl = %s\n",
		respbody, expectedBody, url)
	return true, message
}

func deleteHelper(url string, body string, expectedStatus int) (bool, string) {

	reader := strings.NewReader(body)
	request, err := http.NewRequest("DELETE", hostname+url, reader)
	if err != nil {
		message := "http fail to create new request"
		return true, message
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		message := fmt.Sprintf("Error in http delete request to %s\n", url)
		return true, message
	}
	if expectedStatus != response.StatusCode {
		message := "delete status=" + strconv.Itoa(response.StatusCode) +
			", expected=" + strconv.Itoa(expectedStatus) +
			" url=" + url
		return true, message
	}
	return false, "n/a"
}

func getHelper(url string, expectedStatus int, expectedBody string) (bool, string) {

	response, err := http.Get(hostname + url)
	if err != nil {
		message := fmt.Sprintf("Error in http get request to %s\n", url)
		return true, message
	}

	defer response.Body.Close()
	if expectedStatus != response.StatusCode {
		message := "get status=" + strconv.Itoa(response.StatusCode) +
			", expected=" + strconv.Itoa(expectedStatus) +
			" url=" + url
		return true, message
	}

	rawrespbody, err := ioutil.ReadAll(response.Body)
	respbody := strings.TrimRight(string(rawrespbody), "\n\r")

	if err != nil {
		message := "http get, unable to read body from " + url
		return true, message
	}

	if expectedBody == respbody {
		return false, "n/a"
	}

	message := fmt.Sprintf("http get body: \n%s\ndid not match:\n%s\nurl = %s\n",
		respbody, expectedBody, url)
	return true, message
}

// create new member
func TestCreateMember1(t *testing.T) {
	// request: POST /create-member, body (json): {"name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
	// response: 201  body (json) {"member_id": "1"}
	url := "/create-member"
	requestBody := `{ "name": "Kevin Buchs", "thetype": "employee", ` +
		`"data": "Software Engineer" }`
	expectedStatus := 201
	expectedBody := `{ "member_id": "1" }`
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// read created member back
func TestReadMember1(t *testing.T) {
	// request: GET /read-member/Kevin Buchs
	// response: 200, json = [ { "member_id": "1", "name": "Kevin Buchs",
	//   "thetype": "employee", "data": "Software Engineer" } ]
	name := url.QueryEscape(`Kevin Buchs`)
	url := `/read-member/` + name
	expectedStatus := 200
	expectedBody := `[ { "member_id": "1", "name": "Kevin Buchs", ` +
		`"thetype": "employee", ` +
		`"data": "Software Engineer" } ]`
	err, message := getHelper(url, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// add tag to member
func TestTagMember1(t *testing.T) {
	// request: POST /create-tag, body (json): {"member_id": "1", "tag": "C#"}
	// response: 201
	url := "/create-tag"
	requestBody := `{"member_id": "1", "tag": "C#"}`
	expectedStatus := 201
	err, message := postHelper(url, requestBody, expectedStatus, "")
	if err {
		t.Errorf(message)
	}
}

// add second member
func TestCreateMember2(t *testing.T) {
	// request: POST /create-member, body (json): {"name": "Adam Buchs", "thetype": "contractor", "data": "6"}
	// response: 201  body (json) {"member_id": "2"}
	url := "/create-member"
	requestBody := `{ "name": "Adam Buchs", "thetype": "contractor", ` +
		`"data": "6" }`
	expectedStatus := 201
	expectedBody := `{ "member_id": "2" }`
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// add tag to second member
func TestTagMember2(t *testing.T) {
	// request: POST /create-tag, body (json): {"member_id": "2", "tag": "Angular"}
	// response: 201
	url := "/create-tag"
	requestBody := `{"member_id": "2", "tag": "Angular"}`
	expectedStatus := 201
	err, message := postHelper(url, requestBody, expectedStatus, "")
	if err {
		t.Errorf(message)
	}
}

// get all members
func TestReadAll1(t *testing.T) {
	// request: GET /read-all-members
	// response: 200, json = [ {"member_id": "1", "name": "Kevin Buchs", "thetype":
	//   "employee", "data": "Software Engineer"}, {"member_id": "2", "name": "Adam
	//   Buchs", "thetype": "contractor", "data": "6"} ]
	url := "/read-all-members"
	expectedBody := `[ { "member_id": "1", "name": "Kevin Buchs", ` +
		`"thetype": "employee", "data": "Software Engineer" }, { "member_id": "2", ` +
		`"name": "Adam Buchs", "thetype": "contractor", "data": "6" } ]`
	expectedStatus := 200
	err, message := getHelper(url, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// find members with tag "C#"
func TestFindTag1(t *testing.T) {
	// request: POST /find-member-tags  body (json): {"tags": ["C#"]}
	// response: 200, json = [
	//	{"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
	//	]
	url := "/find-member-tags"
	requestBody := `{ "tags": ["C#"] }`
	expectedBody := `[ { "member_id": "1", "name": "Kevin Buchs", ` +
		`"thetype": "employee", "data": "Software Engineer" } ]`
	expectedStatus := 200
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// find members with tag "Angular"
func TestFindTag2(t *testing.T) {
	// request: POST /find-member-tags  body (json): {"tags": ["Angular"]}
	// response: 200, json = [
	//	{ "member_id": "2", "name": "Adam Buchs", "thetype": "contractor", "data": "6" }
	//	]
	url := "/find-member-tags"
	requestBody := `{ "tags": ["Angular"] }`
	expectedBody := `[ { "member_id": "2", "name": "Adam Buchs", ` +
		`"thetype": "contractor", "data": "6" } ]`
	expectedStatus := 200
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// find members with tag "Jugular"
func TestFindTag3(t *testing.T) {
	// request: POST /find-member-tags  body (json): {"tags": ["Jugular"]}
	// response: 200, json = {[]}
	url := "/find-member-tags"
	requestBody := `{ "tags": ["Jugular"] }`
	expectedBody := `[  ]`
	expectedStatus := 200
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// update member type of Adam to employee and data to Software Engineer
func TestUpdateType(t *testing.T) {
	// request: POST /update-type  body (json): {"member_id": "2", "thetype": "employee", "Data": "Software Engineer"}
	// response: 201
	url := "/update-type"
	requestBody := `{ "member_id": "2", "thetype": "employee", "data": "Software Engineer" }`
	expectedStatus := 201
	err, message := postHelper(url, requestBody, expectedStatus, "")
	if err {
		t.Errorf(message)
	}
}

// update member name Adam to Kevin - creating name duplicates
func TestUpdateName(t *testing.T) {
	// request: POST /update-name  body: { "member_id": "2", "name": "Kevin Buchs" }
	// response: 201
	url := "/update-name"
	requestBody := `{ "member_id": "2", "name": "Kevin Buchs" }`
	expectedStatus := 201
	err, message := postHelper(url, requestBody, expectedStatus, "")
	if err {
		t.Errorf(message)
	}
}

// get member "Kevin Buchs"
func TestReadMember2(t *testing.T) {
	// request: GET /read-member/Kevin Buchs
	// response: 200, json = [
	//  { "member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer" },
	//  { "member_id": "2", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer" }
	//   ]
	name := url.QueryEscape(`Kevin Buchs`)
	url := `/read-member/` + name
	expectedStatus := 200
	expectedBody := `[ { "member_id": "1", "name": "Kevin Buchs", ` +
		`"thetype": "employee", "data": "Software Engineer" }, ` +
		`{ "member_id": "2", "name": "Kevin Buchs", ` +
		`"thetype": "employee", "data": "Software Engineer" } ]`
	err, message := getHelper(url, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// get all members
func TestReadAll2(t *testing.T) {
	// request: GET /read-all-members
	// response: 200, json = {[
	//	{"member_id": "1", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"},
	//	{"member_id": "2", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
	//	]}
	url := "/read-all-members"
	expectedBody := `[ { "member_id": "1", "name": "Kevin Buchs", ` +
		`"thetype": "employee", "data": "Software Engineer" }, { "member_id": "2", ` +
		`"name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer" } ]`
	expectedStatus := 200
	err, message := getHelper(url, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// delete tag from member Kevin2 "Angular"
func TestDeleteTag(t *testing.T) {
	// request: DELETE /delete-member-tag
	//         body(json): { "member_id": "2", "tag": "Angular" }
	// response: 200
	url := `/delete-member-tag`
	body := `{ "member_id": "2", "tag": "Angular" }`
	expectedStatus := 200
	err, message := deleteHelper(url, body, expectedStatus)
	if err {
		t.Errorf(message)
	}
}

// find member on tags "Angular"
func TestFindTag4(t *testing.T) {
	// request: POST /find-member-tags  body (json): {"tags": ["Angular"]}
	// response: 200, json = [  ]
	url := "/find-member-tags"
	requestBody := `{ "tags": ["Angular"] }`
	expectedBody := `[  ]`
	expectedStatus := 200
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// delete member Kevin 1
func TestDeleteMember(t *testing.T) {
	// request: DELETE /delete-member/1
	// response: 200
	url := "/delete-member/1"
	fake_body := `"nothing"`
	expectedStatus := 200
	err, message := deleteHelper(url, fake_body, expectedStatus)
	if err {
		t.Errorf(message)
	}
}

// find member on tags "C#"
func TestFindTag5(t *testing.T) {
	// request: POST /find-member-tags  body (json): {"tags": ["C#"]}
	// response: 200, json = {[]}
	url := "/find-member-tags"
	requestBody := `{ "tags": ["C#"] }`
	expectedBody := `[  ]`
	expectedStatus := 200
	err, message := postHelper(url, requestBody, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// // get all members
func TestReadAll3(t *testing.T) {
	// request: GET /read-all-members
	// response: 200, json = {[
	//	{"member_id": "2", "name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer"}
	//	]}
	url := "/read-all-members"
	expectedBody := `[ { "member_id": "2", ` +
		`"name": "Kevin Buchs", "thetype": "employee", "data": "Software Engineer" } ]`
	expectedStatus := 200
	err, message := getHelper(url, expectedStatus, expectedBody)
	if err {
		t.Errorf(message)
	}
}

// One final step to empty the database (so test suite is idempotent)
// delete member Kevin 2 - BUT, due to the fragility of member number,
// Kevin 2 is now in rowid 1 - this is BROKEN!
func TestDeleteMember2(t *testing.T) {
	// request: DELETE /delete-member/1
	// response: 200
	url := "/delete-member/1"
	fake_body := `"nothing"`
	expectedStatus := 200
	err, message := deleteHelper(url, fake_body, expectedStatus)
	if err {
		t.Errorf(message)
	}
}
