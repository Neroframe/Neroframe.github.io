package models

import (
    "database/sql"
)

type Record struct {
    Email         string
    CName         string
    DiseaseCode   string
    TotalDeaths   int
    TotalPatients int
}

func GetAllRecords(db *sql.DB) ([]Record, error) {
    rows, err := db.Query("SELECT email, cname, disease_code, total_deaths, total_patients FROM Record")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var records []Record
    for rows.Next() {
        var r Record
        err := rows.Scan(&r.Email, &r.CName, &r.DiseaseCode, &r.TotalDeaths, &r.TotalPatients)
        if err != nil {
            return nil, err
        }
        records = append(records, r)
    }
    return records, nil
}

func GetRecord(db *sql.DB, email, cname, diseaseCode string) (*Record, error) {
    var r Record
    err := db.QueryRow("SELECT email, cname, disease_code, total_deaths, total_patients FROM Record WHERE email=$1 AND cname=$2 AND disease_code=$3",
        email, cname, diseaseCode).
        Scan(&r.Email, &r.CName, &r.DiseaseCode, &r.TotalDeaths, &r.TotalPatients)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &r, nil
}

func CreateRecord(db *sql.DB, r *Record) error {
    _, err := db.Exec("INSERT INTO Record (email, cname, disease_code, total_deaths, total_patients) VALUES ($1, $2, $3, $4, $5)",
        r.Email, r.CName, r.DiseaseCode, r.TotalDeaths, r.TotalPatients)
    return err
}

func UpdateRecord(db *sql.DB, r *Record) error {
    _, err := db.Exec("UPDATE Record SET total_deaths=$1, total_patients=$2 WHERE email=$3 AND cname=$4 AND disease_code=$5",
        r.TotalDeaths, r.TotalPatients, r.Email, r.CName, r.DiseaseCode)
    return err
}

func DeleteRecord(db *sql.DB, email, cname, diseaseCode string) error {
    _, err := db.Exec("DELETE FROM Record WHERE email=$1 AND cname=$2 AND disease_code=$3", email, cname, diseaseCode)
    return err
}
