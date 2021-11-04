package db

import (
	"context"
	"errors"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

var Pool *sqlitex.Pool

func createTables(conn *sqlite.Conn) error {
	createQuestionsTable := `CREATE TABLE IF NOT EXISTS questions(
		id INTEGER UNIQUE PRIMARY KEY,
		body TEXT NOT NULL,
		weight REAL NOT NULL,
		type TEXT NOT NULL
	)`

	createOptionsTable := `CREATE TABLE IF NOT EXISTS options(
		id INTEGER UNIQUE PRIMARY KEY,
		body TEXT NOT NULL,
		weight REAL NOT NULL,
		questionId INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE
	)`

	createAnswersTable := `CREATE TABLE IF NOT EXISTS answers(
		questionId INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		questionWeight REAL NOT NULL,
		questionType TEXT NOT NULL,
		selectedOptionId INTEGER REFERENCES options(id) ON DELETE CASCADE,
		selectedOptionWeight INTEGER,
		enteredText TEXT
	)`

	err := sqlitex.Exec(conn, createQuestionsTable, nil)
	if err != nil {
		return err
	}

	err = sqlitex.Exec(conn, createOptionsTable, nil)
	if err != nil {
		return err
	}

	err = sqlitex.Exec(conn, createAnswersTable, nil)
	if err != nil {
		return err
	}

	return nil
}

func insertQuestionsAndOptions(conn *sqlite.Conn) error {
	insertQuestions := `INSERT INTO questions(id, body, weight, type) VALUES(?, ?, ?, ?);`
	insertOptions := `INSERT INTO options(id, body, weight, questionId) VALUES(?, ?, ?, ?);`

	err := sqlitex.Exec(conn, insertQuestions, nil, 100, "Where does the sun set?", 0.5, "ChoiceQuestion")
	if err != nil {
		return err
	}
	err = sqlitex.Exec(conn, insertQuestions, nil, 101, "What is your favorite food?", 1, "TextQuestion")
	if err != nil {
		return err
	}

	err = sqlitex.Exec(conn, insertOptions, nil, 200, "East", 0, 100)
	if err != nil {
		return err
	}
	err = sqlitex.Exec(conn, insertOptions, nil, 201, "West", 1, 100)
	if err != nil {
		return err
	}

	return nil
}

func dropTables(conn *sqlite.Conn) {
	dropQuestions := `DROP TABLE IF EXISTS questions;`
	dropOptions := `DROP TABLE IF EXISTS options;`
	dropAnswers := `DROP TABLE IF EXISTS answers;`

	sqlitex.Exec(conn, dropAnswers, nil)
	sqlitex.Exec(conn, dropOptions, nil)
	sqlitex.Exec(conn, dropQuestions, nil)
}

func InitializeDB() error {
	var err error
	Pool, err = sqlitex.Open("./database.db", 0, 10)
	if err != nil {
		return err
	}

	conn := Pool.Get(context.Background())
	if conn == nil {
		return errors.New("Failed to established initial db connection")
	}
	defer Pool.Put(conn)

	dropTables(conn)

	err = createTables(conn)
	if err != nil {
		return err
	}

	err = insertQuestionsAndOptions(conn)
	if err != nil {
		return err
	}

	return nil
}
