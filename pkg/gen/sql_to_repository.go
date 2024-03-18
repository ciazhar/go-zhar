package gen

import (
	"fmt"
	"log"
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

// generateSQLCFile generates an SQLC file.
//
// No parameters.
// Returns an error.
func generateSQLCFile() error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(sqlcYAMLFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %v", filepath.Dir(sqlcYAMLFilePath), err)
	}

	// Create or overwrite the file with the specified content
	if err := os.WriteFile(sqlcYAMLFilePath, []byte(sqlcYAMLContent), 0644); err != nil {
		return fmt.Errorf("error creating %s: %v", sqlcYAMLFilePath, err)
	}

	return nil
}

// executeSQLCGenerate generates SQL code using sqlc tool.
//
// No parameters.
// Returns an error.
func executeSQLCGenerate() error {
	// Set the working directory to "configs" and run the command
	cmd := exec.Command("sqlc", "generate")
	cmd.Dir = "configs"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running 'sqlc generate': %v", err)
	}

	log.Println("sqlc generate completed successfully.")
	return nil
}

// SQLToRepository generates SQL C file and executes SQL C generate.
//
// No parameters.
// Returns an error.
func SQLToRepository() error {
	if err := generateSQLCFile(); err != nil {
		return err
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(SchemaFolder, os.ModePerm); err != nil {
		return err
	}

	if err := executeSQLCGenerate(); err != nil {
		return err
	}

	return nil
}
