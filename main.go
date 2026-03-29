package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /webhook", func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		if verify_token := values.Get("hub.verify_token"); verify_token == "" || verify_token != os.Getenv("VERIFY_TOKEN") {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintln(w, "Invalid token")
			return
		}

		challenge := values.Get("hub.challenge")
		if challenge == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintln(w, "Missing challenge")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, challenge)
	})

	mux.HandleFunc("POST /webhook", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		hash_payload := req.Header.Get("X-Hub-Signature-256")
		if hash_payload == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintln(w, "Missing X-Hub-Signature-256 header")
			return
		}
		hash_payload = strings.Replace(hash_payload, "sha256=", "", -1)

		rawBody, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Can't read req body, bad format")
			return
		}

		mac := hmac.New(sha256.New, []byte(os.Getenv("APP_SECRET")))
		mac.Write(rawBody)
		if expected_payload := hex.EncodeToString(mac.Sum(nil)); expected_payload != hash_payload {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintln(w, "JSON doesn't match hash payload")
			return
		}

		fmt.Println(string(rawBody))

		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Starting server at port :%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, mux))
}
