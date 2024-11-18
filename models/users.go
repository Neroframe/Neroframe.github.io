package models

import (
    "database/sql"
)

type User struct {
    Email   string
    Name    string
    Surname string
    Salary  sql.NullInt64
    Phone   sql.NullString
    CName   string
}

func GetAllUsers(db *sql.DB) ([]User, error) {
    rows, err := db.Query("SELECT email, name, surname, salary, phone, cname FROM Users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.Email, &user.Name, &user.Surname, &user.Salary, &user.Phone, &user.CName)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func GetUser(db *sql.DB, email string) (*User, error) {
    var user User
    err := db.QueryRow("SELECT email, name, surname, salary, phone, cname FROM Users WHERE email=$1", email).
        Scan(&user.Email, &user.Name, &user.Surname, &user.Salary, &user.Phone, &user.CName)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func CreateUser(db *sql.DB, user *User) error {
    _, err := db.Exec("INSERT INTO Users (email, name, surname, salary, phone, cname) VALUES ($1, $2, $3, $4, $5, $6)",
        user.Email, user.Name, user.Surname, user.Salary, user.Phone, user.CName)
    return err
}

func UpdateUser(db *sql.DB, user *User) error {
    _, err := db.Exec("UPDATE Users SET name=$1, surname=$2, salary=$3, phone=$4, cname=$5 WHERE email=$6",
        user.Name, user.Surname, user.Salary, user.Phone, user.CName, user.Email)
    return err
}

func DeleteUser(db *sql.DB, email string) error {
    _, err := db.Exec("DELETE FROM Users WHERE email=$1", email)
    return err
}
