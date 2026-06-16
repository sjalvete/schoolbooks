package model

import (
	"database/sql"
	"time"
)

type AttendanceRow struct {
	ChildID    int
	ChildName  string
	ParentName string
	Going      *bool
	Paid       bool
	AmountPaid int
}

type EventAttendance struct {
	EventID   int
	Title     string
	Date      string
	Price     int
	Recipient *Recipient
	Skippable bool
	Rows      []AttendanceRow
}

// ListEventAttendance returns, for every event, the attendance/payment status
// of every relevant child. If userID is 0 all children are included (admin
// view), otherwise only the children belonging to that user are included.
func ListEventAttendance(db *sql.DB, userID int) ([]EventAttendance, error) {
	query := `
		SELECT e.id, e.title, e.date, e.price, e.skippable,
		       c.id, c.name, u.name,
		       a.going, a.paid, a.amount_paid,
			   r.title, r.account, r.details
		FROM events e
		CROSS JOIN children c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN attendance a ON a.child_id = c.id AND a.event_id = e.id
		LEFT JOIN recipients r ON r.id = e.recipient_id`

	args := []any{}
	if userID != 0 {
		query += " WHERE c.user_id = ?"
		args = append(args, userID)
	}
	query += " ORDER BY e.date, e.id, c.name"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EventAttendance
	for rows.Next() {
		var (
			ea         EventAttendance
			row        AttendanceRow
			going      sql.NullInt64
			paid       sql.NullInt64
			amountPaid sql.NullInt64
			rTitle     sql.NullString
			rAccount   sql.NullString
			rDetails   sql.NullString
		)
		if err := rows.Scan(
			&ea.EventID, &ea.Title, &ea.Date, &ea.Price, &ea.Skippable,
			&row.ChildID, &row.ChildName, &row.ParentName,
			&going, &paid, &amountPaid,
			&rTitle, &rAccount, &rDetails,
		); err != nil {
			return nil, err
		}

		if rTitle.Valid {
			ea.Recipient = &Recipient{
				Title:   rTitle.String,
				Account: rAccount.String,
				Details: rDetails.String,
			}
		}

		if going.Valid {
			v := going.Int64 != 0
			row.Going = &v
		}
		row.Paid = paid.Int64 != 0
		row.AmountPaid = int(amountPaid.Int64)

		if len(events) == 0 || events[len(events)-1].EventID != ea.EventID {
			events = append(events, ea)
		}
		last := &events[len(events)-1]
		last.Rows = append(last.Rows, row)
	}
	return events, rows.Err()
}

// Bucket classifies the event into one of three groups used by the payments
// view:
//
//	0 - at least one child's parent flagged the event as paid and is
//	    waiting for the admin to confirm it
//	1 - the event still has unsettled payments and hasn't expired yet
//	2 - everything else (settled, declined or expired)
func (ea EventAttendance) Bucket(today string) int {
	flagged := false
	unpaid := false
	for _, row := range ea.Rows {
		if row.Paid {
			flagged = true
		}
		declined := row.Going != nil && !*row.Going
		settled := row.AmountPaid >= ea.Price
		if !declined && !settled {
			unpaid = true
		}
	}

	expired := ea.Date < today

	switch {
	case flagged:
		return 0
	case unpaid && !expired:
		return 1
	default:
		return 2
	}
}

// BucketEvents groups events into the three payment buckets, preserving the
// date ordering produced by ListEventAttendance.
func BucketEvents(events []EventAttendance) [3][]EventAttendance {
	var buckets [3][]EventAttendance
	today := time.Now().Format("2006-01-02")
	for _, ea := range events {
		b := ea.Bucket(today)
		buckets[b] = append(buckets[b], ea)
	}
	return buckets
}

func ensureAttendanceRow(db *sql.DB, childID, eventID int) (int, error) {
	var id int
	err := db.QueryRow(
		"SELECT id FROM attendance WHERE child_id = ? AND event_id = ?", childID, eventID,
	).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec(
			"INSERT INTO attendance (child_id, event_id, amount_paid) VALUES (?, ?, 0)", childID, eventID,
		)
		if err != nil {
			return 0, err
		}
		lastID, err := res.LastInsertId()
		return int(lastID), err
	}
	return id, err
}

// SetGoing records the attendance decision for a child/event pair.
func SetGoing(db *sql.DB, childID, eventID int, going bool) error {
	id, err := ensureAttendanceRow(db, childID, eventID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE attendance SET going = ? WHERE id = ?", going, id)
	return err
}

// TogglePaid flips the "paid" communication flag for a child/event pair.
func TogglePaid(db *sql.DB, childID, eventID int) error {
	id, err := ensureAttendanceRow(db, childID, eventID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE attendance SET paid = CASE WHEN paid = 1 THEN 0 ELSE 1 END WHERE id = ?", id)
	return err
}

// SetAmountPaid records how much has actually been paid for a child/event
// pair and clears the "paid" communication flag.
func SetAmountPaid(db *sql.DB, childID, eventID, amount int) error {
	id, err := ensureAttendanceRow(db, childID, eventID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE attendance SET amount_paid = ?, paid = 0 WHERE id = ?", amount, id)
	return err
}

// ChildBelongsToUser reports whether the given child belongs to the user,
// used to authorize self-service attendance/payment updates.
func ChildBelongsToUser(db *sql.DB, childID, userID int) (bool, error) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM children WHERE id = ? AND user_id = ?", childID, userID,
	).Scan(&count)
	return count > 0, err
}
