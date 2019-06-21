# Detector

Detector is a golang application designed to detect suspicious account logins when the distance between the two login locations (gathered from login ip address) is too far to physically travel within a certain time/speed (500 mph). 

## Installation

#### With Docker (Preferred Method)
The preferred method to run detector is in a Docker container. So first make sure you have docker installed and configured on your machine. 
- [Install Docker](https://docs.docker.com/install/)

After the prerequisite installations, pull down the repo, and build the docker image

```bash
git clone git@github.com:abumere/detector.git
cd detector
docker build -t detector-docker .
```
#### Locally (Alternative Method)
If you wish to build and run the golang binaries locally then you need to ensure that golang is installed and configured on your machine
- [Install Go](https://golang.org/doc/install)

After the prerequisite installations, pull down the repo, and build the golang binaries

```bash
git clone git@github.com:abumere/detector.git
cd detector
go build
```

##
## Running 

#### With Docker (Preferred Method)
To run detector just spin up a docker container from the docker image built in the previous step

`````
docker run -d -p 8080:8080 --name detector detector-docker
`````
#### Locally (Alternative Method)
If you followed the alternative installation step then you can do the following

```bash
cd detector
./detector
```
##
## Usage

Detector accepts input through POST requests and returns outputs through http response

You can send data through the curl command: 
```bash
$ curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764800, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42", "ip_address": "206.81.252.6"}' http://localhost:8080/v1/
```
Alternatively you can use a tool like [Postman](https://www.getpostman.com/downloads/) to send POST requests


# POST Request Inputs
The POST must contain a properly formatted JSON object as a string. The fields are below

| Field             | Required?  | Format        |
| -------------     |:----------:| ------------: |
| username          | YES        | length > 0    |
| unix_timestamp    | YES        |   length > 0  |
| event_uuid        | YES        |    length > 0 |
| ip_address        | YES        |    IPv4       |


##
## Expected Results 
In response to the above POST request, your API should return a JSON document informing
the client about the geo information for the current IP access event as well as the nearest
previous and subsequent events (if they exist). For the preceding/subsequent events, it should
also include a field suspiciousTravel indicating whether travel to/from that geo is suspicious
or not, as well as the speed.
```bash
{  
   "currentGeo":{  
      "lat":39.1702,
      "lon":-76.8538,
      "radius":20
   },
   “travelToCurrentGeoSuspicious”:true,
   “travelFromCurrentGeoSuspicious”:false,
   "precedingIpAccess":{  
      "ip":"24.242.71.20",
      "speed":55,
      "lat":30.3764,
      "lon":-97.7078,
      "radius":5,
      "timestamp":1514764800
   },
   "subsequentIpAccess":{  
      "ip":"91.207.175.104",
      "speed":27600,
      "lat":34.0494,
      "lon":-118.2641,
      "radius":200,
      "timestamp":1514851200
   }
}
```

## 3rd Party Libraries & Resources 
- [MaxMind City Database Data](https://dev.maxmind.com/geoip/geoip2/geolite2/): Publically available city geolocation data 
- [geoip2-golang](https://github.com/oschwald/geoip2-golang): A MaxMind GeoIP2 Reader for Go
- [go-sqlite3](https://github.com/mattn/go-sqlite3): sqlite3 driver for go using database/sql
- [mux](https://github.com/gorilla/mux): A powerful URL router and dispatcher for golang
- [travel/travel.go](https://gist.github.com/cdipaolo/d3f8db3848278b49db68): Used to calculate distance using the Haversin Formula 

## Potential Bugs in Coding Challenge Discription
The following points are things in the coding challenge description that I believe to be bugs. I made certain assumptions around most points in order to finsih the project. 

- In the expected result the precedingIpAccess timestamp is the same as the timestamp from the POST request. If its a preceeding login then it should be an date that comes before the current login time. Based on the distance between the two points and the given speed (55 mph) I was able to deduce that the time for the preceeding login should be 1514677279
- In the expected result the subsequentIpAccess has a speed of 27600 which seems to be incorrect. I believe this is becuase the subsequentIpAccess has an IP of 91.207.175.104 which should make the distance from the currentGeo point about 2300 miles. The timestamp of the subsequentIpAccess is 1514851200. This means that there is only a 24 hour time difference between the currentGeo login and the subsequent one. Given the distance between the two points and the time elapsed the distance speed should be about 96 miles an hour. This means both of these logins should have been marked not suspicious 
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)