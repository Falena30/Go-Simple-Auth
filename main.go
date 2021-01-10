package main

import (
	"Go-Simple-Auth/handle"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	//r.Handle("/static/",
	//	http.StripPrefix("/static/",
	//		http.FileServer(http.Dir("asset"))))

	r.HandleFunc("/register", handle.PostRegister)
	//static route untuk memamngil local css atau js
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./asset/")))
	http.ListenAndServe(":8080", r)
}
