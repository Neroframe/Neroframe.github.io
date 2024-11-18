package models

import (
    "database/sql"
)

type Specialize struct {
    ID    int
    Email string
}

func GetAllSpecializes(db *sql.DB) ([]Specialize, error) {
    rows, err := db.Query("SELECT id, email FROM Specialize")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var specializes []Specialize
    for rows.Next() {
        var s Specialize
        err := rows.Scan(&s.ID, &s.Email)
        if err != nil {
            return nil, err
        }
        specializes = append(specializes, s)
    }
    return specializes, nil
}

func GetSpecialize(db *sql.DB, id int, email string) (*Specialize, error) {
    var s Specialize
    err := db.QueryRow("SELECT id, email FROM Specialize WHERE id=$1 AND email=$2", id, email).
        Scan(&s.ID, &s.Email)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &s, nil
}

func CreateSpecialize(db *sql.DB, s *Specialize) error {
    _, err := db.Exec("INSERT INTO Specialize (id, email) VALUES ($1, $2)", s.ID, s.Email)
    return err
}

func DeleteSpecialize(db *sql.DB, id int, email string) error {
    _, err := db.Exec("DELETE FROM Specialize WHERE id=$1 AND email=$2", id, email)
    return err
}
