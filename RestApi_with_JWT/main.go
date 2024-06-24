package main

import (
	"bufio"
	"context"
	"ex04/httpHand"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/places", httpHand.HandleRequest)
	mux.HandleFunc("/api/recommendd", httpHand.HandleRequestGeo)
	mux.HandleFunc("/api/recommend", httpHand.VerifyJWT(httpHand.HandleRequestGeo))
	mux.HandleFunc("/api/get_token", httpHand.HandleRequestToken)

	srv := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	go func() {
		log.Println("Podnimayus")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server ne podnyalsya: %v", err)
		}
	}()

	done := make(chan struct{})
	go shutDown(srv, done)
	<-done
}

func shutDown(srv *http.Server, done chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "s" {
			break
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Cannot fall: %s", err)
	}
	log.Println("Prileg")
	close(done)
}
