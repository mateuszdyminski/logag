// Package middlewares provides common middleware handlers.
package middlewares

import (
	"net/http"

	"time"

	"github.com/Sirupsen/logrus"
)

// Log is a middleware that logs details about the request/response.
func Log() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			t := time.Now()

			next.ServeHTTP(res, req)

			elapsed := time.Since(t)
			logrus.Printf("%s [%s] \"%s %s %s\" \"%s\" \"%s\" \"Took: %s\"", req.RemoteAddr,
				t.Format("02/Jan/2006:15:04:05 -0700"), req.Method, req.RequestURI, req.Proto, req.Referer(), req.UserAgent(), elapsed)
		})
	}
}
