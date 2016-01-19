package handlers

import (
	"github.com/mateuszdyminski/logag/model"
	"github.com/mateuszdyminski/logag/services"
	"github.com/mateuszdyminski/logag/ws"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type LogResources struct {
	srv *services.LogService
}

func ConfigureLogRest(r *mux.Router, srv *services.LogService) {
	rest := &LogResources{srv: srv}

	// add logs
	r.HandleFunc("/api/logs", rest.addLog).Methods("POST")

	// search logs
	r.HandleFunc("/api/logs", rest.search).Methods("GET")

	// register ws client
	r.HandleFunc("/wsapi/ws/{id}", rest.serveWs).Methods("GET")

	// register filter for real-time logs
	r.HandleFunc("/wsapi/filter/{id}", rest.registerFilter).Methods("POST")

	// unregister filter for real-time logs
	r.HandleFunc("/wsapi/filter/{id}", rest.unregisterFilter).Methods("DELETE")
}

func (l *LogResources) addLog(w http.ResponseWriter, r *http.Request) {
	// parse json
	args := new(struct {
		User string      `json:"user,omitempty"`
		Logs []model.Log `json:"logs,omitempty"`
	})

	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		WriteErr(w, err, http.StatusBadRequest)
		return
	}

	if err := l.srv.AddLog(args.User, args.Logs); err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (l *LogResources) search(w http.ResponseWriter, req *http.Request) {
	from := req.URL.Query().Get("from") + "+01:00"
	to := req.URL.Query().Get("to") + "+01:00"
	query := req.URL.Query().Get("query")
	level := req.URL.Query().Get("level")

	size, err := strconv.Atoi(req.URL.Query().Get("l"))
	if err != nil {
		size = 100
	}

	skip, err := strconv.Atoi(req.URL.Query().Get("s"))
	if err != nil {
		skip = 0
	}

	fromTime := time.Time{}
	if from != "" {
		if fromTime, err = time.Parse(time.RFC3339, from); err != nil {
			logrus.Warnf("Can't parse 'from' time: %s", from)
		}
	}

	toTime := time.Time{}
	if to != "" {
		if toTime, err = time.Parse(time.RFC3339, to); err != nil {
			logrus.Warnf("Can't parse 'to' time: %s", to)
		}
	}

	logs, err := l.srv.Search(query, level, fromTime, toTime, size, skip)
	if err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(logs)
	if err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

// serverWs handles websocket requests from the peer.
func (l *LogResources) serveWs(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "Id can't be empty", http.StatusBadGateway)
		return
	}

	logrus.Infof("Registering client with id: %s to WS", id)
	upg, err := ws.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.Errorf("Error %+v", err)
		return
	}
	c := &ws.Connection{ID: id, Send: make(chan *model.Log, 256), Ws: upg}
	l.srv.Ws.Register <- c
	go c.WritePump()
}

// registerFilter registers filter for real-time logs.
func (l *LogResources) registerFilter(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "Id can't be empty", http.StatusBadGateway)
		return
	}

	filter := model.Filter{}
	if err := json.NewDecoder(req.Body).Decode(&filter); err != nil {
		WriteErr(w, err, http.StatusBadRequest)
		return
	}

	logrus.Infof("Registering filter %v.", filter)
	if err := l.srv.RegisterFilter(filter); err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}
}

// unregisterFilter unregisters filter for real-time logs.
func (l *LogResources) unregisterFilter(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "Id can't be empty", http.StatusBadGateway)
		return
	}

	logrus.Infof("Unregistering filter for connection with id %v.", id)
	if err := l.srv.UnregisterFilter(id); err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}
}