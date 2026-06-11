package model

import "database/sql"

type Recipient struct {
	ID          int
	Title       string
	Account     string
	Description string
	CreatedAt   string
}

func CreateRecipient(db *sql.DB, title, account, description string) error {
	_, err := db.Exec(
		"INSERT INTO recipients (title, account, description) VALUES (?, ?, ?)",
		title, account, description,
	)
	return err
}

func UpdateRecipient(db *sql.DB, id, title, account, description string) error {
	_, err := db.Exec(
		"UPDATE recipients SET title = ?, account = ?, description = ? WHERE id = ?",
		title, account, description, id,
	)
	return err
}

func DeleteRecipient(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM recipients WHERE id = ?", id)
	return err
}

func GetRecipientByID(db *sql.DB, id string) (Recipient, error) {
	var rcp Recipient
	err := db.QueryRow(
		"SELECT id, title, account, description, created_at FROM recipients WHERE id = ?", id,
	).Scan(&rcp.ID, &rcp.Title, &rcp.Account, &rcp.Description, &rcp.CreatedAt)
	return rcp, err
}

func ListRecipients(db *sql.DB) ([]Recipient, error) {
	rows, err := db.Query("SELECT id, title, account, description, created_at FROM recipients ORDER BY title")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recipients []Recipient
	for rows.Next() {
		var rcp Recipient
		rows.Scan(&rcp.ID, &rcp.Title, &rcp.Account, &rcp.Description, &rcp.CreatedAt)
		recipients = append(recipients, rcp)
	}
	return recipients, nil
}
