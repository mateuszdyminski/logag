### logag - Log Aggregator

Simple log aggregator based on elasticsearch. 

### Prerequisites

- Golang 1.5+
- Docker

### Run

To start elasticsearch on your local machine
```
./elastic.sh
```

To start log aggregator 
```
go run main.go
```

Check if everything is working: http://localhost:8001

### Test

Open http://localhost:8001 and click 'Real-Time' tab

Send sample logs to logag:
```
curl -i \
    -H "Accept: application/json" \
    -X POST -d '{ "user": "family_x", "logs": [ { "time": "2016-01-14T14:19:47.999999999+01:00", "level": "warning", "msg": "The groups number increased tremendously!" }, { "time": "2016-01-14T14:19:48.999999999+01:00", "level": "info", "msg": "The ice breaks!" }, { "time": "2016-01-14T14:19:49.999999999+01:00", "level": "fatal", "msg": "App died!" }, { "time": "2016-01-14T14:19:56.999999999+01:00", "level": "error", "msg": "Some error in app!" }, { "time": "2016-01-14T14:19:50.999999999+01:00", "level": "debug", "msg": "Everything ok!" } ] }' \
    http://localhost:8001/api/logs
```

Logs should appear in your web browser.