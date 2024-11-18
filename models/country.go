package models

import (
	"database/sql"
)

type Country struct {
	CName      string
	Population int64
}

func GetAllCountries(db *sql.DB) ([]Country, error) {
	rows, err := db.Query("SELECT cname, population FROM Country")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []Country
	for rows.Next() {
		var country Country
		err := rows.Scan(&country.CName, &country.Population)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func GetCountry(db *sql.DB, cname string) (*Country, error) {
	var country Country
	err := db.QueryRow("SELECT cname, population FROM Country WHERE cname=$1", cname).
		Scan(&country.CName, &country.Population)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func CreateCountry(db *sql.DB, country *Country) error {
	_, err := db.Exec("INSERT INTO Country (cname, population) VALUES ($1, $2)",
		country.CName, country.Population)
	return err
}

func UpdateCountry(db *sql.DB, country *Country) error {
	_, err := db.Exec("UPDATE Country SET population=$1 WHERE cname=$2",
		country.Population, country.CName)
	return err
}

func DeleteCountry(db *sql.DB, cname string) error {
	_, err := db.Exec("DELETE FROM Country WHERE cname=$1", cname)
	return err
}
