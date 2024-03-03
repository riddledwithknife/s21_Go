package main

import (
	"context"
	"log"
	"net/http"

	"day06/server"

	"github.com/jackc/pgx/v5"
)

func main() {
	rateLimiter := server.NewRateLimiter()

	conn, err := pgx.Connect(context.Background(), server.GetDBCredentials())
	if err != nil {
		log.Fatalf("Error connection to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	mux := http.NewServeMux()

	fileServerCss := http.FileServer(http.Dir("./css"))
	mux.Handle("/css/", http.StripPrefix("/css/", fileServerCss))

	fileServerImgs := http.FileServer(http.Dir("./images"))
	mux.Handle("/images/", http.StripPrefix("/images/", fileServerImgs))

	mux.HandleFunc("/", rateLimiter.Middleware(server.HandleMainPage(*conn)))
	mux.HandleFunc("/login", rateLimiter.Middleware(server.HandleLoginPage()))
	mux.HandleFunc("/admin", rateLimiter.Middleware(server.HandleAdminPage(*conn)))
	mux.HandleFunc("/logout", rateLimiter.Middleware(server.HandleLogout()))

	log.Fatal(http.ListenAndServe(":8888", mux))
}
