package model

import "database/sql"

type User struct {
    ID    int
    Name  string
    Email string
    Role  string
}

func GetUserByID(db *sql.DB, ID string) (User, error) {
    var u User
    err := db.QueryRow(
        "SELECT id, name, email, role FROM users WHERE ID = ?", ID,
    ).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
    return u, err
}

func GetUserByEmail(db *sql.DB, email string) (User, error) {
    var u User
    err := db.QueryRow(
        "SELECT id, name, email, role FROM users WHERE email = ?", email,
    ).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
    return u, err
}

func GetUserPassword(db *sql.DB, email string) (string, error) {
    var hash string
    err := db.QueryRow(
        "SELECT password FROM users WHERE email = ?", email,
    ).Scan(&hash)
    return hash, err
}

func CreateUser(db *sql.DB, name, email, password, role string) (User, error) {
    res, err := db.Exec(
        "INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
        name, email, password, role,
    )
    if err != nil {
        return User{}, err
    }
    id, _ := res.LastInsertId()
    return User{ID: int(id), Name: name, Email: email, Role: role}, nil
}

func ListUsers(db *sql.DB) ([]User, error) {
    rows, err := db.Query("SELECT id, name, email, role FROM users ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role)
        users = append(users, u)
    }
    return users, nil
}
