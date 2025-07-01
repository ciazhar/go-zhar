package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/db_ch_csv_zip_http/internal/model"
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
	GetAllOptimizedParallel(ctx context.Context, ch chan<- []string) error
	GetAllOptimizedCopy(ctx context.Context, dst io.Writer) error
	GetAllOptimizedCountAll(ctx context.Context, ch chan<- []string) error
}

type repository struct {
	conn clickhouse.Conn
}

func (r *repository) GetAllOptimizedCopy(ctx context.Context, dst io.Writer) error {
	query := `
        SELECT 
            year, gender, age, location,
            race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
            hypertension, heart_disease, smoking_history, bmi, hba1c_level,
            blood_glucose_level, diabetes
        FROM health_data
    `

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer := csv.NewWriter(dst)
	defer writer.Flush()

	for rows.Next() {
		var (
			year                int32
			gender              string
			age                 int32
			location            string
			raceAfricanAmerican bool
			raceAsian           bool
			raceCaucasian       bool
			raceHispanic        bool
			raceOther           bool
			hypertension        bool
			heartDisease        bool
			smokingHistory      string
			bmi                 float32
			hba1cLevel          float32
			bloodGlucoseLevel   float32
			diabetes            bool
		)

		if err := rows.Scan(
			&year, &gender, &age, &location,
			&raceAfricanAmerican, &raceAsian, &raceCaucasian, &raceHispanic, &raceOther,
			&hypertension, &heartDisease, &smokingHistory,
			&bmi, &hba1cLevel, &bloodGlucoseLevel, &diabetes,
		); err != nil {
			return err
		}

		record := []string{
			fmt.Sprint(year),
			gender,
			fmt.Sprint(age),
			location,
			fmt.Sprint(raceAfricanAmerican),
			fmt.Sprint(raceAsian),
			fmt.Sprint(raceCaucasian),
			fmt.Sprint(raceHispanic),
			fmt.Sprint(raceOther),
			fmt.Sprint(hypertension),
			fmt.Sprint(heartDisease),
			smokingHistory,
			fmt.Sprintf("%.2f", bmi),
			fmt.Sprintf("%.2f", hba1cLevel),
			fmt.Sprintf("%.2f", bloodGlucoseLevel),
			fmt.Sprint(diabetes),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (r *repository) GetAllOptimizedParallel(ctx context.Context, ch chan<- []string) error {

	defer close(ch)

	partitions := []int{2016, 2017, 2018, 2019, 2020, 2021, 2022} // example partitions by year

	var wg sync.WaitGroup
	errCh := make(chan error, len(partitions))

	for _, year := range partitions {
		wg.Add(1)
		go func() {
			err := func(startYear int) error {
				defer wg.Done()
				rows, err := r.conn.Query(ctx, `
				SELECT year, gender, age, location,
					   race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
					   hypertension, heart_disease, smoking_history, bmi, hba1c_level,
					   blood_glucose_level, diabetes
				FROM health_data WHERE year >= ? AND year < ?
			`, startYear, startYear+1)
				if err != nil {

					errCh <- err
					return fmt.Errorf("error executing query: %w", err)
				}
				defer rows.Close()

				for rows.Next() {
					var (
						year                                                                   int32
						gender                                                                 string
						age                                                                    int32
						location                                                               string
						raceAfricanAmerican, raceAsian, raceCaucasian, raceHispanic, raceOther bool
						hypertension, heartDisease                                             bool
						smokingHistory                                                         string
						bmi, hba1cLevel, bloodGlucoseLevel                                     float32
						diabetes                                                               bool
					)

					if err := rows.Scan(
						&year, &gender, &age, &location,
						&raceAfricanAmerican, &raceAsian, &raceCaucasian, &raceHispanic, &raceOther,
						&hypertension, &heartDisease, &smokingHistory,
						&bmi, &hba1cLevel, &bloodGlucoseLevel, &diabetes,
					); err != nil {
						return fmt.Errorf("error scanning row: %w", err)
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					case ch <- []string{fmt.Sprint(year), gender, fmt.Sprint(age), location,
						fmt.Sprint(raceAfricanAmerican), fmt.Sprint(raceAsian), fmt.Sprint(raceCaucasian),
						fmt.Sprint(raceHispanic), fmt.Sprint(raceOther),
						fmt.Sprint(hypertension), fmt.Sprint(heartDisease), smokingHistory,
						fmt.Sprintf("%.2f", bmi), fmt.Sprintf("%.2f", hba1cLevel),
						fmt.Sprintf("%.2f", bloodGlucoseLevel), fmt.Sprint(diabetes)}:
					}
				}
				return nil
			}(year)
			if err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)
	return <-errCh
}

func (r *repository) GetAllOptimizedCountAll(ctx context.Context, ch chan<- []string) error {

	defer close(ch)

	rows, err := r.conn.Query(ctx, `
		SELECT year, gender, age, location,
			   race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
			   hypertension, heart_disease, smoking_history, bmi, hba1c_level,
			   blood_glucose_level, diabetes
		FROM health_data
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			year                                                                   int32
			gender                                                                 string
			age                                                                    int32
			location                                                               string
			raceAfricanAmerican, raceAsian, raceCaucasian, raceHispanic, raceOther bool
			hypertension, heartDisease                                             bool
			smokingHistory                                                         string
			bmi, hba1cLevel, bloodGlucoseLevel                                     float32
			diabetes                                                               bool
		)

		if err := rows.Scan(
			&year, &gender, &age, &location,
			&raceAfricanAmerican, &raceAsian, &raceCaucasian, &raceHispanic, &raceOther,
			&hypertension, &heartDisease, &smokingHistory,
			&bmi, &hba1cLevel, &bloodGlucoseLevel, &diabetes,
		); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- []string{fmt.Sprint(year), gender, fmt.Sprint(age), location,
			fmt.Sprint(raceAfricanAmerican), fmt.Sprint(raceAsian), fmt.Sprint(raceCaucasian),
			fmt.Sprint(raceHispanic), fmt.Sprint(raceOther),
			fmt.Sprint(hypertension), fmt.Sprint(heartDisease), smokingHistory,
			fmt.Sprintf("%.2f", bmi), fmt.Sprintf("%.2f", hba1cLevel),
			fmt.Sprintf("%.2f", bloodGlucoseLevel), fmt.Sprint(diabetes)}:
		}
	}
	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]model.HealthData, error) {
	query := "SELECT year, gender, age, location, race_african_american, race_asian, race_caucasian, race_hispanic, race_other, hypertension, heart_disease, smoking_history, bmi, hba1c_level, blood_glucose_level, diabetes FROM health_data"
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.HealthData
	for rows.Next() {
		var hd model.HealthData
		if err := rows.Scan(
			&hd.Year,
			&hd.Gender,
			&hd.Age,
			&hd.Location,
			&hd.RaceAfricanAmerican,
			&hd.RaceAsian,
			&hd.RaceCaucasian,
			&hd.RaceHispanic,
			&hd.RaceOther,
			&hd.Hypertension,
			&hd.HeartDisease,
			&hd.SmokingHistory,
			&hd.BMI,
			&hd.Hba1cLevel,
			&hd.BloodGlucoseLevel,
			&hd.Diabetes,
		); err != nil {
			return nil, err
		}
		results = append(results, hd)
	}
	return results, rows.Err()
}

func (r *repository) Stream(ctx context.Context) (<-chan model.HealthData, <-chan error) {
	dataCh := make(chan model.HealthData)
	errCh := make(chan error, 1)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		query := "SELECT year, gender, age, location, race_african_american, race_asian, race_caucasian, race_hispanic, race_other, hypertension, heart_disease, smoking_history, bmi, hba1c_level, blood_glucose_level, diabetes FROM health_data"
		rows, err := r.conn.Query(ctx, query)
		if err != nil {
			errCh <- err
			return
		}
		defer rows.Close()

		for rows.Next() {
			var hd model.HealthData
			if err := rows.Scan(
				&hd.Year,
				&hd.Gender,
				&hd.Age,
				&hd.Location,
				&hd.RaceAfricanAmerican,
				&hd.RaceAsian,
				&hd.RaceCaucasian,
				&hd.RaceHispanic,
				&hd.RaceOther,
				&hd.Hypertension,
				&hd.HeartDisease,
				&hd.SmokingHistory,
				&hd.BMI,
				&hd.Hba1cLevel,
				&hd.BloodGlucoseLevel,
				&hd.Diabetes,
			); err != nil {
				errCh <- err
				return
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case dataCh <- hd:
			}
		}
		if err := rows.Err(); err != nil {
			errCh <- err
		}
	}()

	return dataCh, errCh
}

func (r *repository) StreamPage(ctx context.Context, limit, offset int) (<-chan model.HealthData, <-chan error) {
	dataCh := make(chan model.HealthData)
	errCh := make(chan error, 1)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		query := "SELECT ... FROM health_data ORDER BY year, location LIMIT ? OFFSET ?"
		rows, err := r.conn.Query(ctx, query, limit, offset)
		if err != nil {
			errCh <- err
			return
		}
		defer rows.Close()

		for rows.Next() {
			var hd model.HealthData
			if err := rows.Scan(); err != nil {
				errCh <- err
				return
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case dataCh <- hd:
			}
		}

		if err := rows.Err(); err != nil {
			errCh <- err
		}
	}()

	return dataCh, errCh
}

func (r *repository) StreamCSV(ctx context.Context) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		query := `SELECT year, gender, age, location, race_african_american, race_asian, race_caucasian, race_hispanic, race_other,
                  hypertension, heart_disease, smoking_history, bmi, hba1c_level, blood_glucose_level, diabetes FROM health_data`

		rows, err := r.conn.Query(ctx, query)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer rows.Close()

		csvWriter := csv.NewWriter(pw)
		defer csvWriter.Flush()

		// Write CSV header
		csvWriter.Write([]string{
			"year", "gender", "age", "location",
			"race_african_american", "race_asian", "race_caucasian", "race_hispanic", "race_other",
			"hypertension", "heart_disease", "smoking_history", "bmi", "hba1c_level", "blood_glucose_level", "diabetes",
		})

		for rows.Next() {
			var hd model.HealthData
			err := rows.Scan(
				&hd.Year, &hd.Gender, &hd.Age, &hd.Location,
				&hd.RaceAfricanAmerican, &hd.RaceAsian, &hd.RaceCaucasian, &hd.RaceHispanic, &hd.RaceOther,
				&hd.Hypertension, &hd.HeartDisease, &hd.SmokingHistory, &hd.BMI, &hd.Hba1cLevel, &hd.BloodGlucoseLevel, &hd.Diabetes,
			)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			record := []string{
				strconv.Itoa(int(hd.Year)), hd.Gender, strconv.Itoa(int(hd.Age)), hd.Location,
				strconv.FormatBool(hd.RaceAfricanAmerican), strconv.FormatBool(hd.RaceAsian), strconv.FormatBool(hd.RaceCaucasian), strconv.FormatBool(hd.RaceHispanic), strconv.FormatBool(hd.RaceOther),
				strconv.FormatBool(hd.Hypertension), strconv.FormatBool(hd.HeartDisease), hd.SmokingHistory,
				fmt.Sprintf("%f", hd.BMI), fmt.Sprintf("%f", hd.Hba1cLevel), fmt.Sprintf("%f", hd.BloodGlucoseLevel), strconv.FormatBool(hd.Diabetes),
			}
			if err := csvWriter.Write(record); err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		if err := rows.Err(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	return pr, nil
}

func (r *repository) CountTable(ctx context.Context) (int, error) {
	var count int
	err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM health_data").Scan(&count)
	if err != nil {
		return count, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

func (r *repository) CreateTable(ctx context.Context) error {

	// Create the table if it doesn't exist
	err := r.conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS health_data (
			year Int32,
			gender String,
			age Int32,
			location String,
			race_african_american UInt8,
			race_asian UInt8,
			race_caucasian UInt8,
			race_hispanic UInt8,
			race_other UInt8,
			hypertension UInt8,
			heart_disease UInt8,
			smoking_history String,
			bmi Float32,
			hba1c_level Float32,
			blood_glucose_level Float32,
			diabetes UInt8
		) ENGINE = MergeTree()
		ORDER BY (year, location, gender);
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

		err = r.conn.Exec(ctx, insertSQL,
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

func NewRepository(conn clickhouse.Conn) Repository {
	return &repository{conn: conn}
}
