package models

import (
    "database/sql"
)

type PublicServant struct {
    Email      string
    Department string
}

func GetAllPublicServants(db *sql.DB) ([]PublicServant, error) {
    rows, err := db.Query("SELECT email, department FROM PublicServant")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var publicServants []PublicServant
    for rows.Next() {
        var ps PublicServant
        err := rows.Scan(&ps.Email, &ps.Department)
        if err != nil {
            return nil, err
        }
        publicServants = append(publicServants, ps)
    }
    return publicServants, nil
}

func GetPublicServant(db *sql.DB, email string) (*PublicServant, error) {
    var ps PublicServant
    err := db.QueryRow("SELECT email, department FROM PublicServant WHERE email=$1", email).
        Scan(&ps.Email, &ps.Department)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &ps, nil
}

func CreatePublicServant(db *sql.DB, ps *PublicServant) error {
    _, err := db.Exec("INSERT INTO PublicServant (email, department) VALUES ($1, $2)",
        ps.Email, ps.Department)
    return err
}

func UpdatePublicServant(db *sql.DB, ps *PublicServant) error {
    _, err := db.Exec("UPDATE PublicServant SET department=$1 WHERE email=$2",
        ps.Department, ps.Email)
    return err
}

func DeletePublicServant(db *sql.DB, email string) error {
    _, err := db.Exec("DELETE FROM PublicServant WHERE email=$1", email)
    return err
}
