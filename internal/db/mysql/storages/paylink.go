package storages

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"
)

const (
	payPayLinkTable = "`pay_link`"
	pl              = "pl"
	payLinkFields   = pl + ".`id`, " +
		pl + ".`pay_id`," +
		pl + ".`amount`," +
		pl + ".`user_id`," +
		pl + ".`status`," +
		pl + ".`date`"

	payLinkInsertFields = "`pay_id`," +
		"`amount`," +
		"`user_id`," +
		"`status`"
)

type PayLinkStorage struct {
	db *sql.DB
}

func NewPayLinkStorage(conn *sql.DB) *PayLinkStorage {
	return &PayLinkStorage{db: conn}
}

func (ph *PayLinkStorage) Find(filter *db.PayLinkFilter) (*models.PayLink, error) {

	query := "SELECT " + payLinkFields + " FROM " + payPayLinkTable + " " + pl + " "

	where := prepareFindPayLinkFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}

	stmt, err := ph.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows := stmt.QueryRow(where.Params...)
	defer func() {
		// _ = rows.Close()
		_ = stmt.Close()
	}()
	if err != nil {
		return nil, err
	}

	return scanPL(rows)
}

func prepareFindPayLinkFilter(filter *db.PayLinkFilter) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}

	if filter.UserId != "" {
		id, err := strconv.ParseInt(filter.UserId, 10, 64)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, pl+".`user_id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}
	if filter.Amount != "" {
		stmt.Conditions = append(stmt.Conditions, pl+".`amount` = ? ")
		stmt.Params = append(stmt.Params, filter.Amount)
	}

	if filter.PayId != "" {
		stmt.Conditions = append(stmt.Conditions, pl+".`pay_id` = ? ")
		stmt.Params = append(stmt.Params, filter.PayId)
	}
	if filter.Status != "" {
		stmt.Conditions = append(stmt.Conditions, pl+".`status` = ? ")
		stmt.Params = append(stmt.Params, filter.Status)
	}

	if filter.Date != "" {
		t, err := time.Parse(layout, filter.Date+" 00:00:00")
		if err != nil {
			log.Println("storages.prepareFindPaymentHistoryFilter, filter.Date:", err)
		} else {
			stmt.Conditions = append(stmt.Conditions, pl+".`date` <= ? ")
			stmt.Params = append(stmt.Params, t)
		}
	}

	return stmt
}

func scanPL(row Scanner) (*models.PayLink, error) {
	item := &models.PayLink{}
	var (
		date sql.NullTime
	)

	err := row.Scan(
		&item.Id,
		&item.PayId,
		&item.Amount,
		&item.UserId,
		&item.Status,
		&date,
	)

	if err != nil {
		return nil, err
	}

	if date.Valid {
		item.Date = date.Time
	}

	return item, err
}

// func preparePayLink(of *db.OrderByPayLink) string {
// 	desc := ""
// 	var order []string

//		condition := false
//		if of.Amount {
//			order = append(order, "amount")
//			condition = true
//		}
//		if of.TimeCreate {
//			desc = "time_create"
//			condition = true
//		}
//		if of.UserId {
//			order = append(order, "user_id")
//			condition = true
//		}
//		if condition {
//			ob := " order by %s %s "
//			return fmt.Sprintf(ob, strings.Join(order, ","), desc)
//		} else {
//			return ""
//		}
//	}
func (ph *PayLinkStorage) Insert(tx *sql.Tx, item *models.PayLink) (*sql.Tx, int64, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s  (%s) VALUES (?, ?, ?, ?)", payPayLinkTable, payLinkInsertFields)

	params := []interface{}{
		item.PayId,
		item.Amount,
		item.UserId,
		item.Status,
	}
	return queryExecTx(tx, ph.db, insertQuery, params, INSERT)
}
func (ph *PayLinkStorage) Delete(tx *sql.Tx, filter *db.PayLinkFilter) (*sql.Tx, error) {
	query := fmt.Sprintf("DELETE FROM %s as %s", payPayLinkTable, pl)
	where := prepareFindPayLinkFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}
	tx, _, err := queryExecTx(
		tx,
		ph.db,
		query,
		where.Params,
		DELETE,
	)
	return tx, err

}
