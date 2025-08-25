package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServeTLS(":8080", "./cert/server.crt", "./cert/server.key", nil)) 
}

func setupAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.LoginHandler)
	log.Println("API server started on :8080\nReach the website via https://localhost:8080")
}
