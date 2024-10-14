package benchs

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"strconv"
	"strings"
)

var pgxdb *pgx.Conn

const (
	pgxInsertBaseSQL   = `INSERT INTO models (name, title, fax, web, age, "right", counter) VALUES `
	pgxInsertValuesSQL = `($1, $2, $3, $4, $5, $6, $7)`
	pgxInsertSQL       = pgxInsertBaseSQL + pgxInsertValuesSQL
	pgxUpdateSQL       = `UPDATE models SET name = $1, title = $2, fax = $3, web = $4, age = $5, "right" = $6, counter = $7 WHERE id = $8`
	pgxSelectSQL       = `SELECT id, name, title, fax, web, age, "right", counter FROM models WHERE id = $1`
	pgxSelectMultiSQL  = `SELECT id, name, title, fax, web, age, "right", counter FROM models WHERE id > 0 LIMIT 100`
)

func init() {
	st := NewSuite("pgx")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, PgxInsert)
		st.AddBenchmark("BulkInsert 100 row", 500*ORM_MULTI, 0, PgxInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, PgxUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, 0, PgxRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, PgxReadSlice)

		//pgx, _ = sql.Open("postgres", benchs.ORM_SOURCE)
		conn, err := pgx.Connect(context.Background(), ORM_SOURCE)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		pgxdb = conn
	}
}

func PgxInsert(b *B) {
	var m *Model

	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		// pq dose not support the LastInsertId method.
		_, err := pgxdb.Exec(context.Background(), pgxInsertSQL, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Right, m.Counter)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func pgxInsert(m *Model) error {
	// pq dose not support the LastInsertId method.
	_, err := pgxdb.Exec(context.Background(), pgxInsertSQL, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Right, m.Counter)
	if err != nil {
		return err
	}
	return nil
}

func PgxInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
	})

	var valuesSQL string
	counter := 1
	for i := 0; i < 100; i++ {
		hoge := ""
		for j := 0; j < 7; j++ {
			if j != 6 {
				hoge += "$" + strconv.Itoa(counter) + ","
			} else {
				hoge += "$" + strconv.Itoa(counter)
			}
			counter++

		}
		if i != 99 {
			valuesSQL += "(" + hoge + "),"
		} else {
			valuesSQL += "(" + hoge + ")"
		}
	}

	for i := 0; i < b.N; i++ {
		nFields := 7
		query := pgxInsertBaseSQL + valuesSQL
		args := make([]interface{}, len(ms)*nFields)
		for j := range ms {
			offset := j * nFields
			args[offset+0] = ms[j].Name
			args[offset+1] = ms[j].Title
			args[offset+2] = ms[j].Fax
			args[offset+3] = ms[j].Web
			args[offset+4] = ms[j].Age
			args[offset+5] = ms[j].Right
			args[offset+6] = ms[j].Counter
		}
		// pq dose not support the LastInsertId method.
		_, err := pgxdb.Exec(context.Background(), query, args...)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgxUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		rawInsert(m)
	})

	for i := 0; i < b.N; i++ {
		_, err := pgxdb.Exec(context.Background(), pgxUpdateSQL, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Right, m.Counter, m.Id)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgxRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		rawInsert(m)
	})

	for i := 0; i < b.N; i++ {
		var mout Model
		err := pgxdb.QueryRow(context.Background(), pgxSelectSQL, 1).Scan(
			&mout.Id,
			&mout.Name,
			&mout.Title,
			&mout.Fax,
			&mout.Web,
			&mout.Age,
			&mout.Right,
			&mout.Counter,
		)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgxReadSlice(b *B) {
	var m = NewModel()
	wrapExecute(b, func() {
		var err error
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			err = pgxInsert(m)
			if err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var j int
		models := make([]Model, b.L)
		rows, err := pgxdb.Query(context.Background(), strings.Replace(pgxSelectMultiSQL, "100", strconv.Itoa(b.L), -1))
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
		for j = 0; rows.Next() && j < len(models); j++ {
			err = rows.Scan(
				&models[j].Id,
				&models[j].Name,
				&models[j].Title,
				&models[j].Fax,
				&models[j].Web,
				&models[j].Age,
				&models[j].Right,
				&models[j].Counter,
			)
			if err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
		models = models[:j]
		if err = rows.Err(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
		rows.Close()
	}
}
