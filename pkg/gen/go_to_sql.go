package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const (
	QUERY_FOLDER  = "./db/queries/"
	SCHEMA_FOLDER = "./db/schemas/"
)

type TableDescriber interface{}
type Column struct {
	Name     string
	DataType string
}

func GoToSQL(v []TableDescriber) error {
	for i := range v {
		tableName, columns, err := structToColumn(v[i])
		if err != nil {
			return err
		}

		createTableSql := columnToSqlCreateTable(tableName, columns)
		err = writeToFile(fmt.Sprintf(SCHEMA_FOLDER+"%s.sql", tableName), createTableSql)
		if err != nil {
			return err
		}

		crudSql := columnToSqlCrud(tableName, columns)
		err = writeToFile(fmt.Sprintf(QUERY_FOLDER+"%s.sql", tableName), crudSql)
		if err != nil {
			return err
		}

	}
	return nil
}

func structToColumn(v interface{}) (tableName string, columns []Column, error error) {
	// Get the type of the struct
	structType := reflect.TypeOf(v)

	// Get struct name
	structName := structType.Name()

	// Iterate through the struct fields
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Get SQL data type based on Go type
		sqlType, err := goTypeToSQLType(field.Type)
		if err != nil {
			return "", nil, err
		}

		// Add column definition to the list
		columns = append(columns, Column{
			Name:     toSnakeCase(field.Name),
			DataType: sqlType,
		})
	}

	return toSnakeCase(singularToPlural(structName)), columns, nil
}

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

func columnToSqlCreateTable(tableName string, columns []Column) string {
	return fmt.Sprintf("CREATE TABLE %s (\n    %s\n);", tableName, strings.Join(columnToString(columns), ",\n    "))
}

func columnToString(column []Column) (sql []string) {
	for i := range column {
		sql = append(sql, fmt.Sprintf("%s %s", toSnakeCase(column[i].Name), column[i].DataType))
	}
	return sql
}

func columnToSqlCrud(tableName string, columns []Column) string {
	var fetchColumns, insertColumns, updateColumns string

	for _, column := range columns {
		fetchColumns += column.Name + ", "
		insertColumns += column.Name + ", "
		updateColumns += column.Name + " = @" + column.Name + "::" + column.DataType + ", "
	}

	// Trim trailing commas
	fetchColumns = strings.TrimSuffix(fetchColumns, ", ")
	insertColumns = strings.TrimSuffix(insertColumns, ", ")
	updateColumns = strings.TrimSuffix(updateColumns, ", ")

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
`, capitalize(tableName), fetchColumns, tableName, capitalize(tableName), fetchColumns, tableName, getDataType(columns, "id"),
		capitalize(tableName), tableName, capitalize(tableName), tableName, insertColumns, getInsertValues(columns),
		capitalize(tableName), tableName, updateColumns, getDataType(columns, "id"), capitalize(tableName), tableName, getDataType(columns, "id"))

	return sql
}

func getDataType(columns []Column, columnName string) string {
	for _, column := range columns {
		if column.Name == columnName {
			return column.DataType
		}
	}
	return "UNKNOWN" // Default to UNKNOWN if the column is not found
}

func getInsertValues(columns []Column) string {
	var values string
	for _, column := range columns {
		values += "@" + column.Name + "::" + column.DataType + ", "
	}
	// Trim trailing comma
	return strings.TrimSuffix(values, ", ")
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

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

func toSnakeCase(input string) string {
	// Use a regular expression to find all occurrences of uppercase letters
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snakeCase := re.ReplaceAllString(input, "${1}_${2}")

	// Convert the result to lowercase
	return strings.ToLower(snakeCase)
}

func singularToPlural(verb string) string {
	if strings.HasSuffix(verb, "y") {
		// For verbs ending with 'y', replace 'y' with 'ies'
		return verb[:len(verb)-1] + "ies"
	} else {
		// For other cases, just add 's'
		return verb + "s"
	}
}
