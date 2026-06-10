package model

import (
	"database/sql"
	"fmt"
)

type Attendee struct {
	ID         int
	Name       string
	Attending  bool
	Paid       bool
	AmountPaid bool
}

type Event struct {
	ID          int
	Title       string
	Description string
	Date        string
	Price       string
	Attendees   []Attendee
}

func CreateEvent(db *sql.DB, title, description, price, date string) error {
	_, err := db.Exec(
		"INSERT INTO events (title, description, price, date) VALUES (?, ?, ?, ?)",
		title, description, price, date,
	)
	return err
}

func UpdateEvent(db *sql.DB, id, title, description, price, date string) error {
	_, err := db.Exec(
		"UPDATE events SET title = ?, description = ?, price = ?, date = ? WHERE id = ?",
		title, description, price, date, id,
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
	err := db.QueryRow(
		"SELECT id, title, description, price, date FROM events WHERE id = ?", id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Price, &e.Date)

	return e, err
}

func GetAttendees(db *sql.DB, user, event string) ([]Attendee, error) {
	var attendees []Attendee
	rows, err := db.Query(`
		SELECT c.id, c.name, a.attending, a.paid, a.amount_paid FROM children c
		JOIN attendance a ON c.id = a.child_id
		WHERE user_id = ? AND event_id = ?`,
		user, event,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Attendee
		rows.Scan(&a.ID, &a.Name, &a.Attending, &a.Paid, &a.AmountPaid)
		attendees = append(attendees, a)
	}
	return attendees, nil
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
