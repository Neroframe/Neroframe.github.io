package models

import (
    "database/sql"
)

type Doctor struct {
    Email  string
    Degree string
}

func GetAllDoctors(db *sql.DB) ([]Doctor, error) {
    rows, err := db.Query("SELECT email, degree FROM Doctor")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var doctors []Doctor
    for rows.Next() {
        var d Doctor
        err := rows.Scan(&d.Email, &d.Degree)
        if err != nil {
            return nil, err
        }
        doctors = append(doctors, d)
    }
    return doctors, nil
}

func GetDoctor(db *sql.DB, email string) (*Doctor, error) {
    var d Doctor
    err := db.QueryRow("SELECT email, degree FROM Doctor WHERE email=$1", email).
        Scan(&d.Email, &d.Degree)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &d, nil
}

func CreateDoctor(db *sql.DB, d *Doctor) error {
    _, err := db.Exec("INSERT INTO Doctor (email, degree) VALUES ($1, $2)",
        d.Email, d.Degree)
    return err
}

func UpdateDoctor(db *sql.DB, d *Doctor) error {
    _, err := db.Exec("UPDATE Doctor SET degree=$1 WHERE email=$2",
        d.Degree, d.Email)
    return err
}

func DeleteDoctor(db *sql.DB, email string) error {
    _, err := db.Exec("DELETE FROM Doctor WHERE email=$1", email)
    return err
}
