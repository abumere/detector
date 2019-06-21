# Detector

Detector is a golang application designed to detect suspicious account logins when the distance between the two locations (gathered from login ip address) is too far to physically travel within a certain time/speed (500 mph). 

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

## Usage

Detector accepts input through POST requests and returns outputs through http response

You can send data through the curl command: 
```bash
import foobar

$ curl -X POST -d
'{"username": "bob",
"unix_timestamp": 1514764800,
"event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42",
"ip_address": "206.81.252.6"}' http://localhost:8080/v1/
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


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)