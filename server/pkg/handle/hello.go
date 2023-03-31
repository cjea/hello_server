package handle

import (
	"fmt"
	"net/http"
	"visits/pkg/storage"
)

func HelloHandler(storage storage.VisitRecorderCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		if err := storage.RecordVisit(name); err != nil {
			http.Error(w, "unable to record visit", http.StatusInternalServerError)
			return
		}

		count, err := storage.CountVisits(name)
		if err != nil {
			http.Error(w, "unable to count visits", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Hello, %s! You have visited %d times.\n", name, count)
	}
}
