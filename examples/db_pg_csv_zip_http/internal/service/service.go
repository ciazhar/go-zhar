package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/model"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/repository"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Service interface {
	ImportCSV(ctx context.Context) error
	ExportAndSendUnoptimized(ctx context.Context) error
	ExportAndSendOptimizedCountAll(ctx context.Context) error
	ExportAndSendOptimizedParallel(ctx context.Context) error
	ExportAndSendOptimizedCopy(ctx context.Context) error
}

type service struct {
	repo repository.Repository
}

func (s service) ImportCSV(ctx context.Context) error {

	err := s.repo.CreateTable(ctx)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	file, err := os.Open("/Users/ciazhar/GolandProjects/go-start-small/datasets/diabetes_clinical_100k/diabetes_dataset.csv")
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	err = s.repo.ImportCSV(ctx, file)
	if err != nil {
		return fmt.Errorf("import csv: %w", err)
	}

	return nil
}

func (s service) ExportAndSendUnoptimized(ctx context.Context) error {
	fmt.Println("Exporting data...")
	all, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get all: %w", err)
	}
	fmt.Printf("retrieving %d records\n", len(all))

	fmt.Println("Exporting to csv...")
	fileName := "health_data.csv"
	if err := s.exportToCsv(fileName, all); err != nil {
		return fmt.Errorf("export to csv: %w", err)
	}

	fmt.Println("Zipping...")
	err = s.zip(err, fileName)
	if err != nil {
		return fmt.Errorf("error in zip: %w", err)
	}

	fmt.Println("Deleting file...")
	if err = os.Remove(fileName); err != nil {
		return fmt.Errorf("remove file: %w", err)
	}

	fmt.Println("Sending email...")
	err = s.httpCall(err)
	if err != nil {
		return fmt.Errorf("http call: %w", err)
	}

	fmt.Println("Cleaning up...")
	if err = os.Remove("health_data.zip"); err != nil {
		return fmt.Errorf("remove file: %w", err)
	}

	return nil
}

func (s service) httpCall(err error) error {
	jsonBody := model.JsonBody{
		FileName:    "health_data.zip",
		ContentType: "application/zip",
	}
	f, err := os.Open("health_data.zip")
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 5000000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read file: %w", err)
	}

	jsonBody.Base64 = encodeBase64(buf[:n])

	jsonBytes, err := json.Marshal(jsonBody)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	resp, err := http.Post("http://localhost:3000/email", "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("error in sending zip file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in sending zip file: %s", resp.Status)
	}
	return nil
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (s service) zip(err error, fileName string) error {
	zipFile, err := os.Create("health_data.zip")
	if err != nil {
		return fmt.Errorf("create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	zipFileHeader, err := zipWriter.Create(fileName)
	if err != nil {
		return fmt.Errorf("create file in zip: %w", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(zipFileHeader, file)
	if err != nil {
		return fmt.Errorf("copy file to zip: %w", err)
	}
	return nil
}

func (s service) exportToCsv(fileName string, all []model.HealthData) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing CSV headers
	headers := []string{"Year", "Gender", "Age", "Location", "RaceAfricanAmerican", "RaceAsian", "RaceCaucasian", "RaceHispanic", "RaceOther", "Hypertension", "HeartDisease", "SmokingHistory", "BMI", "HbA1cLevel", "BloodGlucoseLevel", "Diabetes"}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("write headers: %w", err)
	}

	for _, data := range all {
		row := []string{
			strconv.Itoa(data.Year),
			data.Gender,
			strconv.Itoa(data.Age),
			data.Location,
			strconv.FormatBool(data.RaceAfricanAmerican),
			strconv.FormatBool(data.RaceAsian),
			strconv.FormatBool(data.RaceCaucasian),
			strconv.FormatBool(data.RaceHispanic),
			strconv.FormatBool(data.RaceOther),
			strconv.FormatBool(data.Hypertension),
			strconv.FormatBool(data.HeartDisease),
			data.SmokingHistory,
			strconv.FormatFloat(data.BMI, 'f', -1, 64),
			strconv.FormatFloat(data.Hba1cLevel, 'f', -1, 64),
			strconv.FormatFloat(data.BloodGlucoseLevel, 'f', -1, 64),
			strconv.FormatBool(data.Diabetes),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	return nil
}

// Write CSV to given writer
func (s service) writeCSV(w io.Writer, pw *io.PipeWriter, data <-chan []string) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// header
	if err := csvWriter.Write([]string{
		"year", "gender", "age", "location",
		"race:AfricanAmerican", "race:Asian", "race:Caucasian", "race:Hispanic", "race:Other",
		"hypertension", "heart_disease", "smoking_history",
		"bmi", "hbA1c_level", "blood_glucose_level", "diabetes",
	}); err != nil {
		pw.CloseWithError(err)
		return err
	}

	for row := range data {
		if err := csvWriter.Write(row); err != nil {
			pw.CloseWithError(err)
			return err
		}
	}
	return csvWriter.Error()
}

func (s service) ExportAndSendOptimizedCountAll(ctx context.Context) error {

	// pipe base64 writer
	pr, pw := io.Pipe()

	// base64 buffer to store encoded result temporarily
	var result bytes.Buffer

	go func() {
		defer pw.Close()

		b64Encoder := base64.NewEncoder(base64.StdEncoding, pw)
		defer b64Encoder.Close()

		zipWriter := zip.NewWriter(b64Encoder)
		defer zipWriter.Close()

		zipFileWriter, err := zipWriter.Create("data.csv")
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// CSV writing
		dataChan := make(chan []string, 1000)
		go func() {

			if err = s.repo.GetAllOptimizedCountAll(ctx, dataChan); err != nil {
				pw.CloseWithError(err)
			}

		}()

		if err := s.writeCSV(zipFileWriter, pw, dataChan); err != nil {
			return
		}
	}()

	// Copy the base64 result to a buffer or send it directly to your HTTP request
	if _, err := io.Copy(&result, pr); err != nil {
		return err
	}

	// Send HTTP request
	return s.sendEmailWithAttachment(result.String())
}

func (s service) ExportAndSendOptimizedParallel(ctx context.Context) error {

	// pipe base64 writer
	pr, pw := io.Pipe()

	// base64 buffer to store encoded result temporarily
	var result bytes.Buffer

	go func() {
		defer pw.Close()

		b64Encoder := base64.NewEncoder(base64.StdEncoding, pw)
		defer b64Encoder.Close()

		zipWriter := zip.NewWriter(b64Encoder)
		defer zipWriter.Close()

		zipFileWriter, err := zipWriter.Create("data.csv")
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// CSV writing
		dataChan := make(chan []string, 1000)
		go func() {

			if err = s.repo.GetAllOptimizedParallel(ctx, dataChan); err != nil {
				pw.CloseWithError(err)
			}
		}()

		if err := s.writeCSV(zipFileWriter, pw, dataChan); err != nil {
			return
		}
	}()

	// Copy the base64 result to a buffer or send it directly to your HTTP request
	if _, err := io.Copy(&result, pr); err != nil {
		return err
	}

	// Send HTTP request
	return s.sendEmailWithAttachment(result.String())
}

func (s service) ExportAndSendOptimizedCopy(ctx context.Context) error {
	pr, pw := io.Pipe()
	var result bytes.Buffer

	go func() {
		defer pw.Close()

		b64Encoder := base64.NewEncoder(base64.StdEncoding, pw)
		defer b64Encoder.Close()

		zipWriter := zip.NewWriter(b64Encoder)
		defer zipWriter.Close()

		zipFileWriter, err := zipWriter.Create("health_data.csv")
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// Stream CSV from PostgreSQL using COPY
		if err := s.repo.GetAllOptimizedCopy(ctx, zipFileWriter); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	if _, err := io.Copy(&result, pr); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return s.sendEmailWithAttachment(result.String())
}

func (s service) sendEmailWithAttachment(base64Data string) error {
	jsonBody := model.JsonBody{
		FileName:    "health_data.zip",
		ContentType: "application/zip",
		Base64:      base64Data,
	}

	jsonBytes, err := json.Marshal(jsonBody)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	resp, err := http.Post("http://localhost:3000/email", "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("error in sending zip file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in sending zip file: %s", resp.Status)
	}
	return nil
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}
