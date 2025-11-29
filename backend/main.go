package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello: %s\n", r.URL.Path)
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true) // for debugging
		if err != nil {
			log.Println("Error dumping request:", err)
		} else {
			log.Println("========================")
			log.Println(string(dump))
			log.Println("========================")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"users": ["Rohan"]}`)
	})

	log.Println("Backend server starting on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
