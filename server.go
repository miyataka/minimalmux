package main

import (
	"net/http"
)

func main() {
	srv := &http.Server{
		Addr: ":8080",
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
