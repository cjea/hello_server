package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"visits/pkg/handle"
	"visits/pkg/storage/postgres"

	_ "github.com/lib/pq"
)

func withLogs(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("➡️ \t%s\n%v\n\n", r.URL.String(), time.Now())
	}
}

func main() {
	port := mustEnv("PORT")
	vdb, err := postgres.NewVisitsDB(mustEnv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to create HelloRequestDB: %v", err)
	}
	defer vdb.DB.Close()

	mux := http.NewServeMux()
	mux.Handle("/hello", withLogs(handle.HelloHandler(vdb)))

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), mux))
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Missing environment variable " + key)
	}
	return val
}
