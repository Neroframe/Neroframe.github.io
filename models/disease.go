package models

import (
    "database/sql"
)

type Disease struct {
    DiseaseCode string
    Pathogen    string
    Description string
    ID          int
}

func GetAllDiseases(db *sql.DB) ([]Disease, error) {
    rows, err := db.Query("SELECT disease_code, pathogen, description, id FROM Disease")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var diseases []Disease
    for rows.Next() {
        var d Disease
        err := rows.Scan(&d.DiseaseCode, &d.Pathogen, &d.Description, &d.ID)
        if err != nil {
            return nil, err
        }
        diseases = append(diseases, d)
    }
    return diseases, nil
}

func GetDisease(db *sql.DB, diseaseCode string) (*Disease, error) {
    var d Disease
    err := db.QueryRow("SELECT disease_code, pathogen, description, id FROM Disease WHERE disease_code=$1", diseaseCode).
        Scan(&d.DiseaseCode, &d.Pathogen, &d.Description, &d.ID)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &d, nil
}

func CreateDisease(db *sql.DB, d *Disease) error {
    _, err := db.Exec("INSERT INTO Disease (disease_code, pathogen, description, id) VALUES ($1, $2, $3, $4)",
        d.DiseaseCode, d.Pathogen, d.Description, d.ID)
    return err
}

func UpdateDisease(db *sql.DB, d *Disease) error {
    _, err := db.Exec("UPDATE Disease SET pathogen=$1, description=$2, id=$3 WHERE disease_code=$4",
        d.Pathogen, d.Description, d.ID, d.DiseaseCode)
    return err
}

func DeleteDisease(db *sql.DB, diseaseCode string) error {
    _, err := db.Exec("DELETE FROM Disease WHERE disease_code=$1", diseaseCode)
    return err
}
