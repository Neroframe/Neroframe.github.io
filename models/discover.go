package models

import (
    "database/sql"
    "time"
)

type Discover struct {
    CName        string
    DiseaseCode  string
    FirstEncDate time.Time
}

func GetAllDiscovers(db *sql.DB) ([]Discover, error) {
    rows, err := db.Query("SELECT cname, disease_code, first_enc_date FROM Discover")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var discovers []Discover
    for rows.Next() {
        var d Discover
        err := rows.Scan(&d.CName, &d.DiseaseCode, &d.FirstEncDate)
        if err != nil {
            return nil, err
        }
        discovers = append(discovers, d)
    }
    return discovers, nil
}

func GetDiscover(db *sql.DB, cname, diseaseCode string) (*Discover, error) {
    var d Discover
    err := db.QueryRow("SELECT cname, disease_code, first_enc_date FROM Discover WHERE cname=$1 AND disease_code=$2", cname, diseaseCode).
        Scan(&d.CName, &d.DiseaseCode, &d.FirstEncDate)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &d, nil
}

func CreateDiscover(db *sql.DB, d *Discover) error {
    _, err := db.Exec("INSERT INTO Discover (cname, disease_code, first_enc_date) VALUES ($1, $2, $3)",
        d.CName, d.DiseaseCode, d.FirstEncDate)
    return err
}

func UpdateDiscover(db *sql.DB, d *Discover) error {
    _, err := db.Exec("UPDATE Discover SET first_enc_date=$1 WHERE cname=$2 AND disease_code=$3",
        d.FirstEncDate, d.CName, d.DiseaseCode)
    return err
}

func DeleteDiscover(db *sql.DB, cname, diseaseCode string) error {
    _, err := db.Exec("DELETE FROM Discover WHERE cname=$1 AND disease_code=$2", cname, diseaseCode)
    return err
}
