package main

import "net/http"

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				errRes(w, map[string]string{"error": err.(string)}, nil, 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
