package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

const (
	MockUserName = "username"
	MockPassword = "password"

	RequestIDHeader = "request-id"
)

func genRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestIDHeader, uuid.NewString())
}

func Middleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(genRequestID(r.Context()))

			ctxTimeOut, cancel := context.WithTimeout(r.Context(), time.Second*3)
			defer cancel()

			r = r.WithContext(ctxTimeOut)

			username, password, ok := r.BasicAuth()

			if ok {
				if username != MockUserName || password != MockPassword {
					http.Error(w, "", http.StatusForbidden)
					return
				}
			} else {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
