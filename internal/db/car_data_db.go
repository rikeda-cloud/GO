package carDataDB

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateCarDataTableIf() error {
	db, err := sql.Open(DBMS, DBName)
	if err != nil {
		return err
	}
	_, err = db.Exec(CreateCarDataTableSQL)
	return err
}

func InsertCarData(file_name string, car_speed float64, car_steering float64) error {
	db, err := sql.Open(DBMS, DBName)
	if err != nil {
		return err
	}
	_, err = db.Exec(InsertCarDataSQL, file_name, car_speed, car_steering)
	return err
}

func UpdateCarData(file_name string, ideal_speed, ideal_steering float64) error {
	db, err := sql.Open(DBMS, DBName)
	if err != nil {
		return err
	}
	_, err = db.Exec(UpdateCarDataSQL, ideal_speed, ideal_steering, file_name)
	return err
}

func SelectNoMarkedCarData(prevId int64) (*CarData, error) {
	db, err := sql.Open(DBMS, DBName)
	if err != nil {
		return nil, err
	}
	row := db.QueryRow(SelectNoMarkedCarDataSQL, prevId)

	var carData CarData
	err = row.Scan(
		&carData.ID,
		&carData.FileName,
		&carData.CarSpeed,
		&carData.CarSteering,
		&carData.IdealSpeed,
		&carData.IdealSteering,
		&carData.MarkFlag,
		&carData.CreatedAt,
	)

	// (row.size() == 0) IS Error.
	if err != nil {
		return nil, err
	}

	return &carData, nil
}

func DeleteCarData(file_name string) error {
	db, err := sql.Open(DBMS, DBName)
	if err != nil {
		return err
	}
	_, err = db.Exec(DeleteCarDataSQL, file_name)
	return err
}
