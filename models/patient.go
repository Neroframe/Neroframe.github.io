package models

import (
	"database/sql"
)

type Patient struct {
	Email string
}

func GetAllPatients(db *sql.DB) ([]Patient, error) {
	rows, err := db.Query("SELECT email FROM Patients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var p Patient
		err := rows.Scan(&p.Email)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func GetPatient(db *sql.DB, email string) (*Patient, error) {
	var p Patient
	err := db.QueryRow("SELECT email FROM Patients WHERE email=$1", email).
		Scan(&p.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func CreatePatient(db *sql.DB, p *Patient) error {
	_, err := db.Exec("INSERT INTO Patients (email) VALUES ($1)",
		p.Email)
	return err
}

func UpdatePatient(db *sql.DB, patient *Patient) error {
	query := `UPDATE Patients SET email = $1 WHERE email = $2`
	_, err := db.Exec(query, patient.Email, patient.Email)
	return err
}

func DeletePatient(db *sql.DB, email string) error {
	_, err := db.Exec("DELETE FROM Patients WHERE email=$1", email)
	return err
}
