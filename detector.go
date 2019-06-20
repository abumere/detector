package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oschwald/geoip2-golang"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"superman-detector/detector/geo"
	"superman-detector/detector/models"
	"superman-detector/detector/travel"
)


type loginRecord struct {
	Username      string `json:"username"`
	UnixTimestamp int64  `json:"unix_timestamp"`
	EventUUID     string `json:"event_uuid"`
	IPAddr        string `json:"ip_address"`
}

type currentGeo struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Radius  uint16  `json:"radius"`
}

type ipAccess struct {
	IP    	       string   `json:"ip"`
	Speed          int      `json:"speed"`
	Lat            float64  `json:"lat"`
	Lon            float64  `json:"lon"`
	Radius         uint16   `json:"radius"`
	Timestamp      int64    `json:"unix_timestamp"`
}

type Env struct {
	loginDB *sql.DB
	geoDB *geoip2.Reader
}


func isValidIP(ip string) bool {
	matched, err := regexp.Match(`^(?:(?:^|\.)(?:2(?:5[0-5]|[0-4]\d)|1?\d?\d)){4}$`, []byte(ip))
	if err != nil {
		panic(err)
	}
	return matched
}

func validateInputs(lr loginRecord) bool {
	isValidIP := isValidIP(lr.IPAddr)
	return isValidIP && len(lr.Username) > 0 && len(lr.EventUUID) > 0
}

func parsePostBody(reqBody io.ReadCloser) loginRecord {
	var lr loginRecord
	decoder := json.NewDecoder(reqBody)
	err := decoder.Decode(&lr)

	if err != nil {
		fmt.Println("Could not parse post request. Please check json formatting: ")
		panic(err)
	}
	if !validateInputs(lr) {
		panic("Invalid inputs. Please check format of post request and try again")
	}
	return lr
}

func getTravelSpeed(postLogin, prevLogin  models.Login) int {
	//Calc distance between prev login and current
	dist := travel.Distance(postLogin.Lat,postLogin.Lon,prevLogin.Lat,prevLogin.Lon)
	speed := travel.Speed(dist, prevLogin.UnixTimestamp, postLogin.UnixTimestamp)
	return speed
}

func (env *Env) handlePost(rw http.ResponseWriter, request *http.Request) {
	//Parse and validate post query input values
	var lr = parsePostBody(request.Body)

	ip := net.ParseIP(lr.IPAddr)
	record, err := env.geoDB.City(ip)
	if err != nil {
		panic(err)
	}

	cg := currentGeo{
		Lat:    record.Location.Latitude,
		Lon:    record.Location.Longitude,
		Radius: record.Location.AccuracyRadius,
	}

	loginRow := models.Login{
		Username: lr.Username,
		UnixTimestamp: lr.UnixTimestamp,
		EventUUID: lr.EventUUID,
		IPAddr: lr.IPAddr,
		Lat: cg.Lat,
		Lon: cg.Lon,
		Radius: cg.Radius,
	}

	// Add this login entry to the datastore
	models.InsertLogin(env.loginDB, loginRow)
	allLogins, err := models.LoginsByUsername(env.loginDB, loginRow.Username)

	if err != nil {
		http.Error(rw, http.StatusText(500), 500)
		panic(err)
	}

	//Get preceding and subsequent logins if applicable
	prevLogin, postLogin := models.GetAdjacentLogins(allLogins, loginRow)

	//fmt.Printf("Pre: %+v\n", prevLogin)
	//fmt.Printf("Current: %+v\n", loginRow)
	//fmt.Printf("Post: %+v\n", postLogin)
	//fmt.Println("\n\n")


	repOutput := map[string]interface{}{
		"currentGeo": cg,

	}

	if len(prevLogin.Username) != 0 {
		speed := getTravelSpeed(prevLogin, loginRow)
		if speed > 500 {
			repOutput["travelToCurrentGeoSuspicious"] = true
		} else {
			repOutput["travelToCurrentGeoSuspicious"] = false
		}

		repOutput["precedingIpAccess"] = ipAccess{
			IP: prevLogin.IPAddr,
			Speed: speed,
			Lat: prevLogin.Lat,
			Lon: prevLogin.Lon,
			Radius: prevLogin.Radius,
			Timestamp: prevLogin.UnixTimestamp,
		}
	}

	if len(postLogin.Username) != 0 {
		speed := getTravelSpeed(postLogin, loginRow)
		if speed > 500 {
			repOutput["travelFromCurrentGeoSuspicious"] = true
		} else {
			repOutput["travelFromCurrentGeoSuspicious"] = false
		}

		repOutput["subsequentIpAccess"] = ipAccess{
			IP: postLogin.IPAddr,
			Speed: speed,
			Lat: postLogin.Lat,
			Lon: postLogin.Lon,
			Radius: postLogin.Radius,
			Timestamp: postLogin.UnixTimestamp,
		}
	}


	jsonOutput, err := json.Marshal(repOutput)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonOutput)
}

func main() {
	loginDB, err1 := models.NewDB("./data.db")
	geoDB, err2 := geo.NewGeo("./geo/GeoLite2-City.mmdb")

	if err1 != nil {
		log.Panic(err1)
	}
	if err2 != nil {
		log.Panic(err2)
	}

	env := &Env{ loginDB: loginDB, geoDB: geoDB	}

	http.HandleFunc("/v1/", env.handlePost)
	http.ListenAndServe(":8080", nil)
}

