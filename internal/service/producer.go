package service

import (
	"io"
	"net/http"
)

func NewDCProducerServer(addr string, bqueue chan byte) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/dc/produce", handleProducer(bqueue))
	return &http.Server{Handler: mux, Addr: addr}
}

func handleProducer(bqueue chan byte) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, b := range body {
			bqueue <- b
		}
	}
}
