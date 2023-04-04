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
	vdb, err := postgres.NewPurchaseDB(mustEnv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer vdb.DB.Close()

	mux := http.NewServeMux()
	mux.Handle("/authorize", handle.PurchasesHandler(vdb))

	fmt.Printf("Hello, world!")
	log.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:8080"), mux))
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Missing environment variable " + key)
	}
	return val
}
