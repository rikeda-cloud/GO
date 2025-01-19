package carDataDB

import (
	"GO/internal/config"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateCarDataTableIf() error {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(CreateCarDataTableSQL)
	return err
}

func InsertCarData(file_name string, car_speed float64, car_steering float64) error {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(InsertCarDataSQL, file_name, car_speed, car_steering)
	return err
}

func InsertPredictedCarData(file_name string, car_speed float64, car_steering float64) error {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(InsertPredictedCarDataSQL, file_name, car_speed, car_steering)
	return err
}

func UpdateCarData(file_name string, ideal_speed, ideal_steering float64, tags string) error {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(UpdateCarDataSQL, ideal_speed, ideal_steering, tags, file_name)
	return err
}

func SelectNoMarkedCarData(prevId int64) (*CarData, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
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
		&carData.Tags,
		&carData.CreatedAt,
	)

	// (row.size() == 0) IS Error.
	if err != nil {
		return nil, err
	}

	return &carData, nil
}

func SelectPredictedNoMarkedCarData(prevId int64) (*CarData, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return nil, err
	}
	row := db.QueryRow(SelectPredictedNoMarkedCarDataSQL, prevId)

	var carData CarData
	err = row.Scan(
		&carData.ID,
		&carData.FileName,
		&carData.CarSpeed,
		&carData.CarSteering,
		&carData.IdealSpeed,
		&carData.IdealSteering,
		&carData.MarkFlag,
		&carData.Tags,
		&carData.CreatedAt,
	)

	// (row.size() == 0) IS Error.
	if err != nil {
		return nil, err
	}

	return &carData, nil
}

func SelectNextCarData(id int64) (*CarData, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return nil, err
	}
	row := db.QueryRow(SelectNextCarDataSQL, id)

	var carData CarData
	err = row.Scan(
		&carData.ID,
		&carData.FileName,
		&carData.CarSpeed,
		&carData.CarSteering,
		&carData.IdealSpeed,
		&carData.IdealSteering,
		&carData.MarkFlag,
		&carData.Tags,
		&carData.CreatedAt,
	)

	// (row.size() == 0) IS Error.
	if err != nil {
		return nil, err
	}

	return &carData, nil
}

func SelectPrevCarData(id int64) (*CarData, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return nil, err
	}
	row := db.QueryRow(SelectPrevCarDataSQL, id)

	var carData CarData
	err = row.Scan(
		&carData.ID,
		&carData.FileName,
		&carData.CarSpeed,
		&carData.CarSteering,
		&carData.IdealSpeed,
		&carData.IdealSteering,
		&carData.MarkFlag,
		&carData.Tags,
		&carData.CreatedAt,
	)

	// (row.size() == 0) IS Error.
	if err != nil {
		return nil, err
	}

	return &carData, nil
}

func SelectIdFromFileName(fileName string) (int64, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return -1, err
	}
	row := db.QueryRow(SelectIdFromFileNameSQL, fileName)

	var id int64
	err = row.Scan(&id)

	// (row.size() == 0) IS Error.
	if err != nil {
		return -1, err
	}

	return id, nil
}

func SelectRemainImageCount() (int, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return -1, err
	}

	var count int
	err = db.QueryRow(SelectRemainImageCountSQL).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func SelectPredictedRemainImageCount() (int, error) {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return -1, err
	}

	var count int
	err = db.QueryRow(SelectPredictedRemainImageCountSQL).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func DeleteCarData(file_name string) error {
	cfg := config.GetConfig()
	db, err := sql.Open(cfg.Database.DBMS, cfg.Database.FilePath)
	if err != nil {
		return err
	}
	_, err = db.Exec(DeleteCarDataSQL, file_name)
	return err
}
