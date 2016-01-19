package main

import (
	"github.com/mateuszdyminski/logag/liblog"
	"github.com/mateuszdyminski/logag/libcfg"
	"github.com/mateuszdyminski/logag/handlers"
	"github.com/mateuszdyminski/logag/services"
	"github.com/mateuszdyminski/logag/middlewares"
	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/tylerb/graceful"
	"gopkg.in/olivere/elastic.v3"
	"net/http"
	"time"
	"github.com/mateuszdyminski/logag/ws"
	"github.com/boltdb/bolt"
)

// Application is the application object that runs HTTP server.
type Application struct {
	logSrv *services.LogService
}

// TODO: use different middleware library ?
func (app *Application) middlewares(cfg *libcfg.Cfg) (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.Log())
	middle.UseHandler(app.mux(cfg))

	return middle, nil
}

func (app *Application) mux(cfg *libcfg.Cfg) *mux.Router {
	router := mux.NewRouter()

	// initialize handlers per group of endpoints
	handlers.ConfigureLogRest(router, app.logSrv)

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.StaticsPath)))

	return router
}

func NewApp(cfg *libcfg.Cfg, db *bolt.DB) (*Application, error) {
	// Create a elastic client
	client, err := elastic.NewClient(elastic.SetURL(cfg.Elastics...))
	if err != nil {
		return nil, err
	}

	// create websocket server
	ws := ws.NewHub()

	// start ws in separate goroutine
	go ws.Run()

	// create all services
	logSrv := services.NewLogService(cfg, client, ws, db)

	// create indexes in ES
	if err := logSrv.CreateIndex(); err != nil {
		return nil, err
	}

	return &Application{logSrv: logSrv}, nil
}

func main() {
	cfg, err := libcfg.LoadCfg()
	if err != nil {
		logrus.Fatalf("Can't laod configuration! Err: %v", err)
	}

	logCfg, err := liblog.NewLogger("/tmp", "info")
	if err != nil {
		logrus.Fatalf("Can't create logger! Err: %v", err.Error())
	}
	defer logCfg.Close()

	db, err := bolt.Open(cfg.BoltDDPath, 0600, nil)
	if err != nil {
		logrus.Fatal("Can't open boltDB! Err: %v", err.Error())
	}
	defer db.Close()

	app, err := NewApp(cfg, db)
	if err != nil {
		logrus.Fatalf("Can't create application! Err: %v", err.Error())
	}

	middle, err := app.middlewares(cfg)
	if err != nil {
		logrus.Fatalf("Can't create http middlewares! Err: %v", err.Error())
	}

	drainInterval, err := time.ParseDuration(cfg.HttpDrainInterval)
	if err != nil {
		logrus.Fatalf("Can't parse drain interval! Err: %v", err.Error())
	}

	// TODO: check facebook-go as graceful server
	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: cfg.Host, Handler: middle},
	}

	logrus.Infoln("Running HTTP server on " + cfg.Host)

	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatal(err.Error())
	}
}
