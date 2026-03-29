package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /webhook", func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		if verify_token := values.Get("hub.verify_token"); verify_token == "" || verify_token != os.Getenv("VERIFY_TOKEN") {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Invalid token")
			return
		}

		challenge := values.Get("hub.challenge")
		if challenge == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Missing challenge")
			return
		}

		fmt.Fprint(w, challenge)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Starting server at port :%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, mux))
}
