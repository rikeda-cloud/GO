package carDataDB

type CarData struct {
	ID            int
	FileName      string
	CarSpeed      float64
	CarSteering   float64
	IdealSpeed    float64
	IdealSteering float64
	MarkFlag      int
	CreatedAt     string
}

var (
	CreateCarDataTableSQL = `
		CREATE TABLE IF NOT EXISTS car_data(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_name TEXT UNIQUE,
			car_speed REAL NOT NULL,
			car_steering REAL NOT NULL,
			ideal_speed REAL DEFAULT 0,
			ideal_steering REAL DEFAULT 0,
			mark_flag INTEGER DEFAULT 0,
			created_at TEXT DEFAULT (datetime('now'))
		)`

	InsertCarDataSQL = `
		INSERT INTO
			car_data(file_name, car_speed, car_steering)
		VALUES
			(?, ?, ?);`

	UpdateCarDataSQL = `
		UPDATE car_data
			SET ideal_speed = ?,
				ideal_steering = ?,
				mark_flag = 1
			WHERE
				file_name = ?;`

	SelectNoMarkedCarDataSQL = `
		SELECT *
		FROM car_data
		WHERE
			mark_flag = 0 AND 
			id > ?
		ORDER BY id ASC
		LIMIT 1;`

	DeleteCarDataSQL = `
		DELETE
		FROM car_data
		WHERE
			file_name = ?;`
)
