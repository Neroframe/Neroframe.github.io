package models

import (
	"database/sql"
)

type PatientDisease struct {
	Email       string
	DiseaseCode string
}

func GetAllPatientDiseases(db *sql.DB) ([]PatientDisease, error) {
	rows, err := db.Query("SELECT email, disease_code FROM PatientDisease")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patientDiseases []PatientDisease
	for rows.Next() {
		var pd PatientDisease
		err := rows.Scan(&pd.Email, &pd.DiseaseCode)
		if err != nil {
			return nil, err
		}
		patientDiseases = append(patientDiseases, pd)
	}
	return patientDiseases, nil
}

func GetPatientDisease(db *sql.DB, email, diseaseCode string) (*PatientDisease, error) {
	var pd PatientDisease
	err := db.QueryRow("SELECT email, disease_code FROM PatientDisease WHERE email=$1 AND disease_code=$2", email, diseaseCode).
		Scan(&pd.Email, &pd.DiseaseCode)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pd, nil
}

func CreatePatientDisease(db *sql.DB, pd *PatientDisease) error {
	_, err := db.Exec("INSERT INTO PatientDisease (email, disease_code) VALUES ($1, $2)",
		pd.Email, pd.DiseaseCode)
	return err
}

func UpdatePatientDisease(db *sql.DB, email, oldDiseaseCode, newDiseaseCode string) error {
	query := `UPDATE PatientDisease SET disease_code = $1 WHERE email = $2 AND disease_code = $3`
	_, err := db.Exec(query, newDiseaseCode, email, oldDiseaseCode)
	return err
}

func DeletePatientDisease(db *sql.DB, email, diseaseCode string) error {
	_, err := db.Exec("DELETE FROM PatientDisease WHERE email=$1 AND disease_code=$2", email, diseaseCode)
	return err
}
