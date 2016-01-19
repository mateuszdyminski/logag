package model

import "time"

type Log struct {
	User  string    `json:"user,omitempty"`
	Time  time.Time `json:"time,omitempty"`
	Level string    `json:"level,omitempty"`
	Msg   string    `json:"msg,omitempty"`
	Score *float64  `json:"score,omitempty"`
}

const LogMapping = `
        {
            "mappings" : {
                "log" : {
                    "properties" : {
                        "user" : { "type" : "string" },
                        "time" : { "type" : "date" },
                        "level" : { "type" : "string" },
                        "msg" : { "type" : "string" }
                    }
                }
            }
        }`
