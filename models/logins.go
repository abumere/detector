package models

import (
	"database/sql"
	"sort"
)

type Login struct {
	Id			  int    `json:id`
	Username      string `json:"username"`
	UnixTimestamp int64  `json:"unix_timestamp"`
	EventUUID     string `json:"event_uuid"`
	IPAddr        string `json:"ip_address"`
	Lat 		  float64 `json:"lat"`
	Lon 		  float64 `json:"lon"`
	Radius		  uint16 `json:"radius"`
}

func AllLogins(db *sql.DB) ([]*Login, error) {
	rows, err := db.Query("SELECT * FROM logins")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logins := make([]*Login, 0)
	for rows.Next() {

		//Grab each login and add it to slice
		login := new(Login)
		err := rows.Scan(&login.Id, &login.Username, &login.UnixTimestamp, &login.EventUUID, &login.IPAddr, &login.Lat, &login.Lon, &login.Radius)

		if err != nil {
			return nil, err
		}
		logins = append(logins, login)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Sort logins in descending (most recent to oldest) time
	sort.Slice(logins, func(i, j int) bool { return logins[i].UnixTimestamp > logins[j].UnixTimestamp})
	return logins, nil
}

func LoginsByUsername(db *sql.DB, username string) ([]*Login, error) {
	rows, err := db.Query("SELECT * FROM logins WHERE username=?", username)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logins := make([]*Login, 0)
	for rows.Next() {

		//Grab each login and add it to slice
		login := new(Login)
		err := rows.Scan(&login.Id, &login.Username, &login.UnixTimestamp, &login.EventUUID, &login.IPAddr, &login.Lat, &login.Lon, &login.Radius)

		if err != nil {
			return nil, err
		}
		logins = append(logins, login)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Sort logins in descending (most recent to oldest) time
	sort.Slice(logins, func(i, j int) bool { return logins[i].UnixTimestamp < logins[j].UnixTimestamp})
	return logins, nil
}


func InsertLogin(db *sql.DB, row Login) {
	statement, err := db.Prepare("INSERT INTO logins (username,tStamp,uuid,ipAddr,lat,lon,radius) VALUES (?,?,?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(row.Username, row.UnixTimestamp, row.EventUUID, row.IPAddr, row.Lat, row.Lon, row.Radius)

	if err != nil {
		panic(err)
	}
}

func GetAdjacentLogins(allLogins []*Login, cLogin Login) (Login, Login) {
	// Find the login entries before and after the current one so we can calcucate
	var prevIndx, postIndx = -1, -1
	var prevLogin, postLogin Login

	for i, login := range allLogins {
		if cLogin.EventUUID == login.EventUUID {
			if i > 0 {
				prevIndx = i -1
			}

			if i < len(allLogins) - 1 {
				postIndx = i + 1
			}
		}
	}
	if prevIndx != -1 {
		prevLogin = *allLogins[prevIndx]
	}
	if postIndx != -1 {
		postLogin = *allLogins[postIndx]
	}
	return prevLogin, postLogin
}