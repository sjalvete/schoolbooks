package model

import (
	"database/sql"
	"fmt"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Date        string
	Price       string
	Skippable   bool
	RecipientID int
}

func CreateEvent(db *sql.DB, title, description, price, date string, skippable bool, recipientID *int) error {
	_, err := db.Exec(
		"INSERT INTO events (title, description, price, date, skippable, recipient_id) VALUES (?, ?, ?, ?, ?, ?)",
		title, description, price, date, skippable, recipientID,
	)
	return err
}

func UpdateEvent(db *sql.DB, id, title, description, price, date string, skippable bool, recipientID *int) error {
	_, err := db.Exec(
		"UPDATE events SET title = ?, description = ?, price = ?, date = ?, skippable = ?, recipient_id = ? WHERE id = ?",
		title, description, price, date, skippable, recipientID, id,
	)
	return err
}

func DeleteEvent(db *sql.DB, id string) error {
	_, err := db.Exec(
		"DELETE FROM events WHERE id = ?",
		id,
	)
	return err
}

func GetEventByID(db *sql.DB, id string) (Event, error) {
	var e Event
	var recipientID sql.NullInt64
	err := db.QueryRow(
		"SELECT id, title, description, price, date, skippable, recipient_id FROM events WHERE id = ?", id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Price, &e.Date, &e.Skippable, &recipientID)
	if recipientID.Valid {
		e.RecipientID = int(recipientID.Int64)
	}

	return e, err
}

func ListEvents(db *sql.DB) ([]Event, error) {
	rows, err := db.Query("SELECT id, title, description, price, date FROM events ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Price, &e.Date)
		events = append(events, e)
	}
	return events, nil
}

func ListFutureEvents(db *sql.DB) ([]Event, error) {
	rows, err := db.Query("SELECT id, title, description, price, date FROM events WHERE date >= date('now') ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Price, &e.Date)
		events = append(events, e)
	}
	return events, nil
}

func ListEventsByMonth(db *sql.DB, year, month int) ([]Event, error) {
	rows, err := db.Query(
		"SELECT id, title, description, price, date FROM events WHERE strftime('%Y', date) = ? AND strftime('%m', date) = ? ORDER BY date",
		fmt.Sprintf("%04d", year),
		fmt.Sprintf("%02d", month),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.ID, &e.Title, &e.Description, &e.Price, &e.Date)
		events = append(events, e)
	}
	return events, nil
}
