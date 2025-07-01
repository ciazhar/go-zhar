package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"strconv"
	"strings"
	"sync"
)

type Repository interface {
	CountTable(ctx context.Context) (int, error)
	CreateTable(ctx context.Context) error
	ImportCSV(ctx context.Context, ioR io.Reader) error
	GetAll(ctx context.Context) ([]model.HealthData, error)
	GetAllOptimizedCountAll(ctx context.Context, ch chan<- []string) error
	GetAllOptimizedParallel(ctx context.Context, ch chan<- []string) error
	GetAllOptimizedCopy(ctx context.Context, dst io.Writer) error
}

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) CountTable(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM health_data").Scan(&count)
	if err != nil {
		return count, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *repository) GetAllOptimizedCopy(ctx context.Context, dst io.Writer) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	// Make sure to include header using a UNION ALL
	query := `
COPY (
	SELECT
		year, gender, age, location,
		race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
		hypertension, heart_disease, smoking_history,
		bmi, hba1c_level, blood_glucose_level, diabetes
	FROM health_data
) TO STDOUT WITH CSV DELIMITER ',' HEADER
`

	_, err = conn.Conn().PgConn().CopyTo(ctx, dst, query)
	if err != nil {
		return fmt.Errorf("copy to: %w", err)
	}
	return nil
}

func (r *repository) GetAllOptimizedParallel(ctx context.Context, ch chan<- []string) error {

	defer close(ch)

	const pageSize = 10000
	offset := 0

	// Pool to reuse string slices
	rowPool := &sync.Pool{
		New: func() any {
			// Length = column count (16), capacity gives room to reuse
			return make([]string, 0, 16)
		},
	}

	for {
		query := `
			SELECT year, gender, age, location,
			       race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
			       hypertension, heart_disease, smoking_history,
			       bmi, hba1c_level, blood_glucose_level, diabetes
			FROM health_data
			LIMIT $1 OFFSET $2
		`

		rows, err := r.pool.Query(ctx, query, pageSize, offset)
		if err != nil {
			return fmt.Errorf("query: %w", err)
		}

		count := 0
		for rows.Next() {
			var h model.HealthData
			if err := rows.Scan(
				&h.Year,
				&h.Gender,
				&h.Age,
				&h.Location,
				&h.RaceAfricanAmerican,
				&h.RaceAsian,
				&h.RaceCaucasian,
				&h.RaceHispanic,
				&h.RaceOther,
				&h.Hypertension,
				&h.HeartDisease,
				&h.SmokingHistory,
				&h.BMI,
				&h.Hba1cLevel,
				&h.BloodGlucoseLevel,
				&h.Diabetes,
			); err != nil {
				rows.Close()
				return fmt.Errorf("scan: %w", err)
			}

			count++
			select {
			case <-ctx.Done():
				rows.Close()
				return ctx.Err()
			default:
				// Reuse slice from pool
				row := rowPool.Get().([]string)[:0]

				row = append(row,
					strconv.Itoa(h.Year), h.Gender, strconv.Itoa(h.Age), h.Location,
					strconv.FormatBool(h.RaceAfricanAmerican), strconv.FormatBool(h.RaceAsian), strconv.FormatBool(h.RaceCaucasian),
					strconv.FormatBool(h.RaceHispanic), strconv.FormatBool(h.RaceOther),
					strconv.FormatBool(h.Hypertension), strconv.FormatBool(h.HeartDisease), h.SmokingHistory,
					strconv.FormatFloat(h.BMI, 'f', -1, 64), strconv.FormatFloat(h.Hba1cLevel, 'f', -1, 64),
					strconv.FormatFloat(h.BloodGlucoseLevel, 'f', -1, 64), strconv.FormatBool(h.Diabetes),
				)

				ch <- row

				// Optionally return to pool later in consumer
				// rowPool.Put(row)
			}
		}
		rows.Close()

		if count < pageSize {
			break // Last page
		}
		offset += pageSize
	}

	return nil
}

func (r *repository) GetAllOptimizedCountAll(ctx context.Context, ch chan<- []string) error {

	defer close(ch)

	rows, err := r.pool.Query(ctx, "SELECT * FROM health_data")
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		var h model.HealthData
		if err := rows.Scan(
			&h.Year,
			&h.Gender,
			&h.Age,
			&h.Location,
			&h.RaceAfricanAmerican,
			&h.RaceAsian,
			&h.RaceCaucasian,
			&h.RaceHispanic,
			&h.RaceOther,
			&h.Hypertension,
			&h.HeartDisease,
			&h.SmokingHistory,
			&h.BMI,
			&h.Hba1cLevel,
			&h.BloodGlucoseLevel,
			&h.Diabetes,
		); err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- []string{strconv.Itoa(h.Year), h.Gender, strconv.Itoa(h.Age), h.Location,
			strconv.FormatBool(h.RaceAfricanAmerican), strconv.FormatBool(h.RaceAsian), strconv.FormatBool(h.RaceCaucasian), strconv.FormatBool(h.RaceHispanic), strconv.FormatBool(h.RaceOther),
			strconv.FormatBool(h.Hypertension), strconv.FormatBool(h.HeartDisease), h.SmokingHistory, strconv.FormatFloat(h.BMI, 'f', -1, 64), strconv.FormatFloat(h.Hba1cLevel, 'f', -1, 64), strconv.FormatFloat(h.BloodGlucoseLevel, 'f', -1, 64), strconv.FormatBool(h.Diabetes)}:
		}

	}
	return nil
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) CreateTable(ctx context.Context) error {

	// Create the table if it doesn't exist
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS health_data (
			year INT,
			gender TEXT,
			age INT,
			location TEXT,
			race_african_american BOOLEAN,
			race_asian BOOLEAN,
			race_caucasian BOOLEAN,
			race_hispanic BOOLEAN,
			race_other BOOLEAN,
			hypertension BOOLEAN,
			heart_disease BOOLEAN,
			smoking_history TEXT,
			bmi REAL,
			hba1c_level REAL,
			blood_glucose_level REAL,
			diabetes BOOLEAN
		);
	`)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	return nil
}

func (r *repository) ImportCSV(ctx context.Context, ioR io.Reader) error {
	reader := csv.NewReader(ioR)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	// Validate expected header
	expected := []string{
		"year", "gender", "age", "location",
		"race:AfricanAmerican", "race:Asian", "race:Caucasian", "race:Hispanic", "race:Other",
		"hypertension", "heart_disease", "smoking_history",
		"bmi", "hbA1c_level", "blood_glucose_level", "diabetes",
	}
	if len(header) != len(expected) {
		return fmt.Errorf("invalid header: got %d columns", len(header))
	}

	// Prepare insert statement
	insertSQL := `
		INSERT INTO health_data (
			year, gender, age, location,
			race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
			hypertension, heart_disease, smoking_history,
			bmi, hba1c_level, blood_glucose_level, diabetes
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("read row: %w", err)
		}

		// Convert row fields
		year, _ := strconv.Atoi(row[0])
		age, _ := strconv.Atoi(row[2])
		raceFlags := make([]bool, 5)
		for i := 0; i < 5; i++ {
			raceFlags[i] = row[4+i] == "1" || strings.ToLower(row[4+i]) == "true"
		}
		hypertension := row[9] == "1" || strings.ToLower(row[9]) == "true"
		heartDisease := row[10] == "1" || strings.ToLower(row[10]) == "true"
		bmi, _ := strconv.ParseFloat(row[12], 32)
		hba1c, _ := strconv.ParseFloat(row[13], 32)
		glucose, _ := strconv.ParseFloat(row[14], 32)
		diabetes := row[15] == "1" || strings.ToLower(row[15]) == "true"

		_, err = r.pool.Exec(ctx, insertSQL,
			year,
			row[1], // gender
			age,
			row[3], // location
			raceFlags[0], raceFlags[1], raceFlags[2], raceFlags[3], raceFlags[4],
			hypertension,
			heartDisease,
			row[11], // smoking_history
			bmi,
			hba1c,
			glucose,
			diabetes,
		)

		if err != nil {
			return fmt.Errorf("insert row: %w", err)
		}
	}
	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]model.HealthData, error) {

	rows, err := r.pool.Query(ctx, "SELECT * FROM health_data")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	var healthData []model.HealthData
	for rows.Next() {
		var h model.HealthData
		if err := rows.Scan(
			&h.Year,
			&h.Gender,
			&h.Age,
			&h.Location,
			&h.RaceAfricanAmerican,
			&h.RaceAsian,
			&h.RaceCaucasian,
			&h.RaceHispanic,
			&h.RaceOther,
			&h.Hypertension,
			&h.HeartDisease,
			&h.SmokingHistory,
			&h.BMI,
			&h.Hba1cLevel,
			&h.BloodGlucoseLevel,
			&h.Diabetes,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		healthData = append(healthData, h)
	}
	return healthData, nil
}
