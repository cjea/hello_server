package handle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"visits/pkg/storage"
)

/*
	All in a transaction:
		card = find_card(attempt.card_id) // NOTE: select for update
		spent = sum(previous transactions for card)
			select sum(amount) where card_id = attempt.card_id
		return false if spent + attempt.amount > card.limit
		return false if length(find_all_xactions_from_two_minutes_ago()) >= 3

		persist attempt
		return true
*/

type PurchaseRequest struct {
	CardNumber string `json:"card_number"`
	Cents      uint64 `json:"amount"`
}

func PurchasesHandler(db storage.PurchaseAllower) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var req PurchaseRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request body "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "bad request body "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("req = %#v\n", req)
		ok, err := db.IsAllowed(req.CardNumber, req.Cents)
		if err != nil {
			http.Error(w, "unexpected error"+err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "purchase not allowed", http.StatusForbidden)
			return
		}

		fmt.Fprintf(w, "success!")
	}
}
