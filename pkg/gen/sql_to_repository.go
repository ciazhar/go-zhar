package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const sqlcYAMLContent = `version: 2
sql:
  - engine: "postgresql"
    queries: "../db/queries/"
    schema: "../db/schemas/"
    gen:
      go:
        package: "gen"
        out: "../internal/gen/repository"
        sql_package: "pgx/v5"
        emit_empty_slices: true
        emit_json_tags: true
        json_tags_case_style: "snake"
        overrides:
          - db_type: "uuid"
            go_type:
              import: "github.com/google/uuid"
              type: "UUID"
`

const sqlcYAMLFilePath = "./configs/sqlc.yaml"

func generateSQLCFile() error {

	// Extract directory from file path
	sqlcDirectory := filepath.Dir(sqlcYAMLFilePath)

	// Check if the directory exists, if not, create it
	if _, err := os.Stat(sqlcDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(sqlcDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", sqlcDirectory, err)
		}
		fmt.Printf("Directory %s created successfully.\n", sqlcDirectory)
	}

	// Check if the file already exists
	if _, err := os.Stat(sqlcYAMLFilePath); os.IsNotExist(err) {
		// File does not exist, create it with the specified content
		err := ioutil.WriteFile(sqlcYAMLFilePath, []byte(sqlcYAMLContent), 0644)
		if err != nil {
			return fmt.Errorf("error creating %s: %v", sqlcYAMLFilePath, err)
		}
		fmt.Printf("%s created successfully.\n", sqlcYAMLFilePath)
	} else if err != nil {
		// Some error occurred while checking file existence
		return fmt.Errorf("error checking %s existence: %v", sqlcYAMLFilePath, err)
	} else {
		// File already exists
		fmt.Printf("%s already exists.\n", sqlcYAMLFilePath)
	}

	return nil
}

func executeSQLCGenerate() error {
	// Set the working directory to "configs"
	if err := os.Chdir("configs"); err != nil {
		return err
	}

	// Set up the command
	cmd := exec.Command("sqlc", "generate")

	// Redirect standard output and standard error to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running 'sqlc generate': %v", err)
	}

	fmt.Println("sqlc generate completed successfully.")
	return nil
}

func SQLToRepository() error {
	err := generateSQLCFile()
	if err != nil {
		return err
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(SCHEMA_FOLDER, os.ModePerm); err != nil {
		return err
	}

	err = executeSQLCGenerate()
	if err != nil {
		return err
	}

	return nil
}
