[![Go Report Card](https://goreportcard.com/badge/github.com/sauravhiremath/skeduler?style=for-the-badge)](https://goreportcard.com/report/github.com/sauravhiremath/skeduler)
[![Lines of code](https://img.shields.io/tokei/lines/github/sauravhiremath/skeduler?style=for-the-badge)](https://github.com/sauravhiremath/skeduler)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sauravhiremath/skeduler?style=for-the-badge)](https://github.com/golang/go)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/sauravhiremath/skeduler?style=for-the-badge)](https://github.com/sauravhiremath/skeduler)

![image](https://user-images.githubusercontent.com/28642011/96423730-5edcab80-1217-11eb-993d-46d4356d7cd1.png)

Barebones meeting scheduling API. Schedule Meetings with participants

# Documentation

**Postman Link**: https://documenter.getpostman.com/view/8269592/TVRrX5q4

## End-point: Get all Meetings within a given timeframe
### Description: endTime and startTime should be in UTC Format following RFC3399
Method: GET
>```
>http://localhost:8080/meetings?start=2006-01-02T15:04:05Z&end=2006-02-02T15:04:05Z
>```
### Query Params

|Param|value|
|---|---|
|start|2006-01-02T15:04:05Z|
|end|2006-02-02T15:04:05Z|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃


## End-point: Create a Meeting
### Description: 
Method: POST
>```
>http://localhost:8080/meetings
>```
### Body (**raw**)

```json
{
    "title": "Sample Meeting",
    "participants": [
        {
            "name": "abc",
            "email": "abc@def.com",
            "rsvp": "Yes"
        }
    ],
    "start_time": "2019-10-18T10:44:56Z",
    "end_time": "2019-10-18T12:51:56Z"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃


## End-point: Get Single Meeting
### Description: 
Method: GET
>```
>http://localhost:8080/meeting/3142D8E8-01A6-14E7-D3F8-0CE0AF0247EK
>```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃


## End-point: Get all Meetings for a participant 
### Description: 
Method: GET
>```
>http://localhost:8080/meetings?participant=abc@def.com
>```
### Query Params

|Param|value|
|---|---|
|participant|abc@def.com|



⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

# Development Setup

> Note: This setup uses go1.15.x, older versions may not use .mod files so be careful 

Installing the dependencies in vendor

```
go mod vendor
```

**Now, you can use the bash script to run all these at once**

```
./runner.sh
```


OR, run the commands individually,

Build the project to get the executable

```
go build -o build/ .
```

Run Tests for the project

```
go test ./controllers
```

Run the project directly without `build`

```
go run main.go
```


# License

MIT
