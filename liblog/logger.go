package liblog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
)

// LogConfig is Logrus configuration.
type Logger struct {
	file *os.File
}

// NewLogger creates and initialized new log configuration.
func NewLogger(logsPath, logsLevel string) (*Logger, error) {
	log.SetFlags(0)
	log.SetOutput(&logruswriter{})
	logrus.SetFormatter(&formatter{})
	logFile, err := getLogFile(logsPath)
	if err != nil {
		return nil, fmt.Errorf("can't init log file: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
	level, err := logrus.ParseLevel(logsLevel)
	if err != nil {
		return nil, fmt.Errorf("can't set logging level: %v", err)
	}
	logrus.SetLevel(level)
	return &Logger{
		file: logFile,
	}, nil
}

func getLogFile(logsPath string) (*os.File, error) {
	appName := path.Base(os.Args[0])
	logFile, err := os.OpenFile(path.Join(logsPath, appName+".logag.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	logrus.Printf("Writing logs to %v", logFile.Name())
	return logFile, nil
}

// Close closes log file.
func (l *Logger) Close() error {
	return l.file.Close()
}

type logruswriter struct{}

func (l *logruswriter) Write(b []byte) (int, error) {
	logrus.Print(string(b[:len(b)-1]))
	return len(b), nil
}

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s [%s] %s\n",
		entry.Time.Format("2006/01/02 15:04:05.000"),
		entry.Level.String(),
		entry.Message)), nil
}
