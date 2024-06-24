package main

import (
	"bufio"
	"context"
	"ex01/DB"
	"ex01/httpHand"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHand.HandleRequest)

	srv := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	done := make(chan struct{})

	go func() {
		log.Println("Podnimayus")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server ne podnyalsya: %v", err)
		}
	}()

	go func() {
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
	}()

	<-done
}

type Store interface {
	GetPlaces(limit int, offset int) ([]DB.Place, int, error)
}
