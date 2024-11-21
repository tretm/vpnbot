package storages

import (
	"database/sql"
)

const LIMIT = 100

const (
	INSERT string = "insert"
	UPDATE string = "update"
	DELETE string = "delete"
)

type QueryStmt struct {
	Conditions []string
	Params     []interface{}
}

type ExecStmt struct {
	Set []string
	QueryStmt
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func queryExecTx(tx *sql.Tx, db *sql.DB, query string, Params []interface{}, method string) (*sql.Tx, int64, error) {

	if tx == nil {
		txx, err := db.Begin()
		if err != nil {
			return tx, 0, err
		}
		tx = txx
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return tx, 0, err
	}
	defer func() { _ = stmt.Close() }()
	result, err := stmt.Exec(Params...)
	if err != nil {
		return tx, 0, err
	}
	//fmt.Println(Params...)
	if method == INSERT {
		r, err := result.RowsAffected()
		return tx, r, err
	} else if method == UPDATE || method == DELETE {
		r, err := result.RowsAffected()
		return tx, r, err
	}
	r, err := result.LastInsertId()
	return tx, r, err
}

func queryExec(db *sql.DB, query string, Params []interface{}, method string) (int64, error) {

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer func() { _ = stmt.Close() }()
	result, err := stmt.Exec(Params...)
	if err != nil {
		return 0, err
	}
	//fmt.Println(Params...)
	if method == INSERT {
		return result.RowsAffected()
	} else if method == UPDATE || method == DELETE {
		return result.RowsAffected()
	}
	return result.LastInsertId()
}

func countOffset(pageNumber int) int {
	return (pageNumber - 1) * LIMIT
}
