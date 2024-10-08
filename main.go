package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	addr   = flag.String("addr", "127.0.0.1:8080", "address to bind to")
	action = flag.String("action", "false", "command to invoke on webhook")
)

func run() {
	log.Printf("running action in goroutine: %s\n", *action)
	args := strings.Split(*action, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Printf("webook received, invoking action\n")
	go run()
	w.WriteHeader(200)
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("POST /update", http.HandlerFunc(update))

	server := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	log.Printf("listening at http://%s", *addr)
	log.Fatal(server.ListenAndServe())
}
