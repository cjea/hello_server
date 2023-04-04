package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PurchaseDB struct {
	DB *sql.DB
}

func NewPurchaseDB(dbUrl string) (*PurchaseDB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	/*
			`cards` table (number primary key, limit)
		`purchases` (p_id primary key, card_number references cards(number), amount, at timestamp with timezone)
			possible indexes:
				- FK index on card_id for `purchases`
				- index on "at" column in `purchases`
	*/
	_, err = db.Exec(`
	create table if not exists cards(
		card_number text primary key not null,
		limit_cents numeric not null
	);
`)
	if err != nil {
		return nil, err
	}
	// TODO(cjea): build the indexes too
	_, err = db.Exec(`
	create table if not exists purchases(
		p_id text primary key not null,
		card_number text not null references cards(card_number),
		amount_cents numeric not null,
		at timestamp with time zone not null
	);
`)
	if err != nil {
		return nil, err
	}
	return &PurchaseDB{DB: db}, nil
}

func (db *PurchaseDB) IsAllowed(cardNumber string, cents uint64) (bool, error) {
	tx, err := db.DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return false, err
	}
	var exceedsLimit bool
	row := tx.QueryRow(
		`select for update coalesce(sum(amount_cents), 0) + $1 > cards.limit_cents
		from purchases, cards
		where cards.card_number = purchases.card_number and purchases.card_number = $2
		group by cards.limit_cents`,
		cents, cardNumber,
	)
	if err = row.Scan(&exceedsLimit); err != nil {
		return false, err
	}
	if exceedsLimit {
		return false, nil
	}

	var numPurchasesInWindow int
	row = tx.QueryRow(`
		select count(*) from purchases
		where card_number = $1 and at >= now()-'2 minutes'::interval
	`, cardNumber)
	if err = row.Scan(&numPurchasesInWindow); err != nil {
		return false, err
	}
	if numPurchasesInWindow >= 3 {
		return false, nil
	}
	id := uuid.New()
	_, err = db.DB.Exec(`
	insert into purchases(p_id, card_number, amount_cents, at)
	values ($1, $2, $3, now());
	`, id, cardNumber, cents)
	if err != nil {
		return false, err
	}

	return true, nil
}
