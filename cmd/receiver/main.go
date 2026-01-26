package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var v any
		if err := json.Unmarshal(body, &v); err == nil {
			b, _ := json.MarshalIndent(v, "", "  ")
			log.Printf("SIGNAL RECEIVED:\n%s\n", string(b))
		} else {
			log.Printf("SIGNAL RECEIVED (raw): %s\n", string(body))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Println("receiver listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
