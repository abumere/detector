package main

import (
	"bytes"
	"database/sql"
	"detector/geo"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"gopkg.in/testfixtures.v2"
)

var db *sql.DB
var fixtures *testfixtures.Context
var env *Env

func TestMain(m *testing.M) {
	var err error
	// Open connection with the test database.
	// Existing data would be deleted
	db, err := sql.Open("sqlite3", "./dataTest.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS logins (id INTEGER PRIMARY KEY, username TEXT, tStamp TEXT, uuid TEXT, ipAddr TEXT, lat TEXT, lon TEXT, radius TEXT)")
	statement.Exec()

	if err != nil {
		log.Fatal(err)
	}

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFiles(db, &testfixtures.SQLite{}, "testdata/fixtures/fullNegative/logins.yml")
	if err != nil {
		log.Fatal(err)
	}
	geoDB, _ := geo.NewGeo("./geo/GeoLite2-City.mmdb")
	env = &Env{ loginDB: db, geoDB: geoDB	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func clearTestDatabase() {
	statement, err := env.loginDB.Prepare("delete from logins")
	statement.Exec()

	if err != nil {
		log.Fatal(err)
	}
}

func TestSingleResult(t *testing.T) {
	clearTestDatabase()

	jsonBody := []byte(`{"username": "bob", "unix_timestamp": 1514851200, "event_uuid": "15ad929a-db03-4bf4-9541-8f728fa12e42", "ip_address": "91.207.175.104"}`)
	req, err := http.NewRequest("POST", "/v1/", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.HandlePost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"currentGeo":{"lat":34.0549,"lon":-118.2578,"radius":200}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFullNegativeResult(t *testing.T) {
	prepareTestDatabase()

	jsonBody := []byte(`{"username": "bob", "unix_timestamp": 1514764800, "event_uuid": "35ad929a-db03-4bf4-9541-8f728fa12e42", "ip_address": "206.81.252.6"}`)
	req, err := http.NewRequest("POST", "/v1/", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.HandlePost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"currentGeo":{"lat":39.2293,"lon":-76.6907,"radius":10},"precedingIpAccess":{"ip":"24.242.71.20","speed":55,"lat":30.3773,"lon":-97.71,"radius":5,"unix_timestamp":1514677279},"subsequentIpAccess":{"ip":"91.207.175.104","speed":8330887,"lat":34.0549,"lon":-118.2578,"radius":200,"unix_timestamp":1514764801},"travelFromCurrentGeoSuspicious":true,"travelToCurrentGeoSuspicious":false}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}