package models

import (
    "database/sql"
)

type DiseaseType struct {
    ID          int
    Description string
}

func GetAllDiseaseTypes(db *sql.DB) ([]DiseaseType, error) {
    rows, err := db.Query("SELECT id, description FROM DiseaseType")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var diseaseTypes []DiseaseType
    for rows.Next() {
        var dt DiseaseType
        err := rows.Scan(&dt.ID, &dt.Description)
        if err != nil {
            return nil, err
        }
        diseaseTypes = append(diseaseTypes, dt)
    }
    return diseaseTypes, nil
}

func GetDiseaseType(db *sql.DB, id int) (*DiseaseType, error) {
    var dt DiseaseType
    err := db.QueryRow("SELECT id, description FROM DiseaseType WHERE id=$1", id).
        Scan(&dt.ID, &dt.Description)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &dt, nil
}

func CreateDiseaseType(db *sql.DB, dt *DiseaseType) error {
    _, err := db.Exec("INSERT INTO DiseaseType (description) VALUES ($1)", dt.Description)
    return err
}

func UpdateDiseaseType(db *sql.DB, dt *DiseaseType) error {
    _, err := db.Exec("UPDATE DiseaseType SET description=$1 WHERE id=$2", dt.Description, dt.ID)
    return err
}

func DeleteDiseaseType(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM DiseaseType WHERE id=$1", id)
    return err
}
