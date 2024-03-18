package gen

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	QueryFolder  = "./db/queries/"
	SchemaFolder = "./db/schemas/"
)

type TableDescriber interface{}
type Column struct {
	Name     string
	DataType string
}

// GoToSQL generates SQL files from struct definitions for each table in the input slice.
//
// Parameters:
// - v: a slice of TableDescriber interfaces representing different tables.
// Returns an error if any issues occur during the process.
func GoToSQL(v []TableDescriber) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errList []error

	for i := range v {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			tableName, columns, err := structToColumn(v[i])
			if err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				log.Printf("Error converting struct to columns for table %s: %v", tableName, err)
				return
			}

			createTableSql := columnToSqlCreateTable(tableName, columns)
			if err := writeToFile(fmt.Sprintf(SchemaFolder+"%s.sql", tableName), createTableSql); err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				log.Printf("Error writing create table SQL to file for table %s: %v", tableName, err)
				return
			}

			crudSql := columnToSqlCrud(tableName, columns)
			if err := writeToFile(fmt.Sprintf(QueryFolder+"%s.sql", tableName), crudSql); err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				log.Printf("Error writing CRUD SQL to file for table %s: %v", tableName, err)
				return
			}
		}(i)
	}

	wg.Wait()

	if len(errList) > 0 {
		log.Printf("Encountered errors while processing tables: %v", errList)
		return errors.New("encountered errors while processing tables")
	}

	return nil
}

// structToColumn converts a struct to a table name and a slice of columns.
//
// It takes a parameter `v` of type `interface{}` which represents the struct to be converted.
// The function returns three values: `tableName` of type `string`, `columns` of type `[]Column`, and `err` of type `error`.
func structToColumn(v interface{}) (tableName string, columns []Column, err error) {
	structType := reflect.TypeOf(v)
	structName := toSnakeCase(singularToPlural(structType.Name()))

	numFields := structType.NumField()
	columns = make([]Column, 0, numFields)

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		sqlType, err := goTypeToSQLType(field.Type)
		if err != nil {
			log.Println("Error converting Go type to SQL type:", err)
			return "", nil, err
		}
		columns = append(columns, Column{
			Name:     toSnakeCase(field.Name),
			DataType: sqlType,
		})
	}

	return structName, columns, nil
}

// goTypeToSQLType converts a Go type to the corresponding SQL type.
//
// Parameter:
//
//	goType - the Go type to be converted to SQL type
//
// Return:
//
//	string - the corresponding SQL type
//	error - an error if the Go type is unsupported
func goTypeToSQLType(goType reflect.Type) (string, error) {
	switch goType.Kind() {
	case reflect.String:
		return "VARCHAR", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INT", nil
	case reflect.Float32, reflect.Float64:
		return "FLOAT", nil
	default:
		return "", fmt.Errorf("unsupported Go type: %v", goType)
	}
}

// columnToSqlCreateTable creates a SQL CREATE TABLE statement.
//
// tableName string, columns []Column
// string
func columnToSqlCreateTable(tableName string, columns []Column) string {
	var sb strings.Builder

	sb.WriteString("CREATE TABLE ")
	sb.WriteString(tableName)
	sb.WriteString(" (\n")

	for i, col := range columns {
		sb.WriteString("    ")
		sb.WriteString(col.Name)
		sb.WriteString(" ")
		sb.WriteString(col.DataType)

		if i < len(columns)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}

	sb.WriteString(");")

	return sb.String()
}

// columnToString generates an array of SQL column definitions based on the input columns.
//
// columns: slice of Column struct
// []string: array of SQL column definitions
func columnToString(columns []Column) []string {
	sql := make([]string, len(columns))
	for i, column := range columns {
		sql[i] = fmt.Sprintf("%s %s", toSnakeCase(column.Name), column.DataType)
	}
	return sql
}

// columnToSqlCrud generates SQL CRUD statements for the given table name and columns.
//
// tableName: the name of the table
// columns: a slice of Column struct
// string: the generated SQL statements
func columnToSqlCrud(tableName string, columns []Column) string {
	var fetchColumns, insertColumns, updateColumns strings.Builder

	for _, column := range columns {
		fetchColumns.WriteString(column.Name + ", ")
		insertColumns.WriteString(column.Name + ", ")
		updateColumns.WriteString(column.Name + " = @" + column.Name + "::" + column.DataType + ", ")
	}

	// Trim trailing commas
	fetchColumnsStr := strings.TrimSuffix(fetchColumns.String(), ", ")
	insertColumnsStr := strings.TrimSuffix(insertColumns.String(), ", ")
	updateColumnsStr := strings.TrimSuffix(updateColumns.String(), ", ")

	sql := fmt.Sprintf(`-- name: Fetch%s :many
SELECT %s FROM %s;

-- name: Fetch%sByID :one
SELECT %s FROM %s WHERE id = @id::%s;

-- name: Count%s :one
SELECT COUNT(*) FROM %s;

-- name: Insert%s :exec
INSERT INTO %s (%s) VALUES (%s);

-- name: Update%s :exec
UPDATE %s SET %s WHERE id = @id::%s;

-- name: Delete%s :exec
DELETE FROM %s WHERE id = @id::%s;
`, capitalize(tableName), fetchColumnsStr, tableName, capitalize(tableName), fetchColumnsStr, tableName, getDataType(columns, "id"),
		capitalize(tableName), tableName, capitalize(tableName), tableName, insertColumnsStr, getInsertValues(columns),
		capitalize(tableName), tableName, updateColumnsStr, getDataType(columns, "id"), capitalize(tableName), tableName, getDataType(columns, "id"))

	// Add logging statements for debugging
	log.Printf("Generate SQL for table %s\n", tableName)

	return sql
}

// getDataType returns the data type of a specified column name.
//
// Parameters:
// - columns: a slice of Column struct representing the columns to search through.
// - columnName: a string representing the name of the column to find the data type for.
// Return type:
// string - the data type of the specified column name, or "UNKNOWN" if not found.
func getDataType(columns []Column, columnName string) string {
	for _, column := range columns {
		if column.Name == columnName {
			return column.DataType
		}
	}
	return "UNKNOWN" // Default to UNKNOWN if the column is not found
}

// getInsertValues generates a string with insert values for the given columns.
//
// columns is a slice of Column structs.
// string is returned.
func getInsertValues(columns []Column) string {
	var builder strings.Builder
	for _, column := range columns {
		builder.WriteString("@" + column.Name + "::" + column.DataType + ", ")
	}
	// Trim trailing comma
	return strings.TrimSuffix(builder.String(), ", ")
}

// capitalize takes a string and returns the capitalized version of it.
//
// Parameter(s):
//
//	s string - the input string to be capitalized
//
// Return type(s):
//
//	string - the capitalized version of the input string
func capitalize(s string) string {
	return cases.Title(language.English).String(s)
}

// writeToFile extracts the directory path from the filename, creates the directory if it doesn't exist,
// and then creates or opens the file to write the content to.
//
// Parameters:
// - filename string: the name of the file to write to.
// - content string: the content to write to the file.
// Return type:
// - error: an error if any occurs during the file operations.
func writeToFile(filename, content string) error {
	// Extract the directory path from the filename
	dir := filepath.Dir(filename)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Create or open the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	return err
}

// toSnakeCase converts a given input string to snake case.
//
// It takes a string as input and returns the snake case version of the string.
func toSnakeCase(input string) string {
	var buf bytes.Buffer
	inUpper := false
	for i, c := range input {
		if unicode.IsUpper(c) {
			if i > 0 && !inUpper {
				buf.WriteRune('_')
			}
			inUpper = true
		} else {
			inUpper = false
		}
		buf.WriteRune(unicode.ToLower(c))
	}
	return buf.String()
}

// singularToPlural converts a singular verb to its plural form.
//
// Takes a string parameter 'verb' and returns a string.
func singularToPlural(verb string) string {
	if strings.HasSuffix(verb, "y") {
		// For verbs ending with 'y', replace 'y' with 'ies'
		return verb[:len(verb)-1] + "ies"
	} else {
		// For other cases, just add 's'
		return verb + "s"
	}
}
