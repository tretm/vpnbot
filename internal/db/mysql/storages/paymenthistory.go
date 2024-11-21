package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"
)

const (
	paymentHistoryTable  = "`payment_history`"
	pht                  = "ph"
	paymentHistoryFields = pht + ".`id`, " +
		pht + ".`user_id`," +
		pht + ".`amount`," +
		pht + ".`transaction_type`," +
		pht + ".`comment`," +
		pht + ".`time_create`"

	paymentHistoryInsertFields = "`user_id`," +
		"`amount`," +
		"`transaction_type`," +
		"`comment`"
)

type PaymentHistoryStorage struct {
	db *sql.DB
}

func NewPaymentHistoryStorage(conn *sql.DB) *PaymentHistoryStorage {
	return &PaymentHistoryStorage{db: conn}
}

func (ph *PaymentHistoryStorage) FindOne(userId int64) (*models.PaymentHistory, error) {

	query := "SELECT " + paymentHistoryFields + " FROM " + paymentHistoryTable + " " + pht + " "
	id := strconv.FormatInt(userId, 10)
	where := prepareFindPaymentHistoryFilter(&db.PaymentHistoryFilter{UserId: id})
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}

	stmt, err := ph.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(where.Params...)
	defer func() {
		_ = rows.Close()
		_ = stmt.Close()
	}()
	if err != nil {
		return nil, err
	}

	return scanPH(rows)
}

func prepareFindPaymentHistoryFilter(filter *db.PaymentHistoryFilter) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter.Id != "" {
		id, err := strconv.Atoi(filter.Id)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, pht+".`id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}
	if filter.UserId != "" {
		id, err := strconv.ParseInt(filter.UserId, 10, 64)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, pht+".`user_id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}
	if filter.Amount != "" {
		stmt.Conditions = append(stmt.Conditions, pht+".`amount` like ? ")
		stmt.Params = append(stmt.Params, filter.Amount)
	}

	if filter.TransactionType != "" {
		stmt.Conditions = append(stmt.Conditions, pht+".`transaction_type` like ? ")
		stmt.Params = append(stmt.Params, filter.TransactionType)
	}
	if filter.Comment != "" {
		stmt.Conditions = append(stmt.Conditions, pht+".`comment` = ? ")
		stmt.Params = append(stmt.Params, filter.Comment)
	}

	if filter.TimeCreate != "" {
		t, err := time.Parse(layout, filter.TimeCreate+" 00:00:00")
		if err != nil {
			log.Println("storages.prepareFindPaymentHistoryFilter, filter.Date:", err)
		} else {
			stmt.Conditions = append(stmt.Conditions, pht+".`time_create` <= ? ")
			stmt.Params = append(stmt.Params, t)
		}
	}

	return stmt
}
func scanPHs(rows *sql.Rows, cap int) ([]*models.PaymentHistory, error) {
	phs := make([]*models.PaymentHistory, 0, cap)
	for rows.Next() {
		g, err := scanPH(rows)
		if err != nil {
			return nil, err
		}
		phs = append(phs, g)
	}
	if len(phs) == 0 {
		return phs, errors.New("links not found")
	}
	return phs, nil
}

func scanPH(row Scanner) (*models.PaymentHistory, error) {
	item := &models.PaymentHistory{}
	var (
		comment         sql.NullString
		transactionType sql.NullString
		date            sql.NullTime
	)
	err := row.Scan(
		&item.Id,
		&item.UserId,
		&item.Amount,
		&transactionType,
		&comment,
		&date,
	)

	if err != nil {
		return nil, err
	}

	if comment.Valid {
		item.Comment = comment.String
	}

	if transactionType.Valid {
		item.TransactionType = transactionType.String
	}

	if date.Valid {
		item.TimeCreate = date.Time
	}

	return item, err
}
func (ph *PaymentHistoryStorage) Find(filter *db.PaymentHistoryFilter, orderFilter *db.OrderByPaymentHistory, offset, limit int) ([]*models.PaymentHistory, error) {

	query := "SELECT " + paymentHistoryFields + " FROM " + paymentHistoryTable + " " + pht + " "

	where := prepareFindPaymentHistoryFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}

	order := preparePaymentHistory(orderFilter)
	query += order

	if limit == 0 {
		limit = db.DefaultLimit
	}

	query += " LIMIT " + strconv.Itoa(int(offset)) + ", " + strconv.Itoa(int(limit))

	stmt, err := ph.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(where.Params...)
	defer func() {
		_ = rows.Close()
		_ = stmt.Close()
	}()
	if err != nil {
		return nil, err
	}

	return scanPHs(rows, limit)
}

func preparePaymentHistory(of *db.OrderByPaymentHistory) string {
	desc := ""
	var order []string

	condition := false
	if of.Amount {
		order = append(order, "amount")
		condition = true
	}
	if of.TimeCreate {
		desc = "time_create"
		condition = true
	}
	if of.UserId {
		order = append(order, "user_id")
		condition = true
	}
	if condition {
		ob := " order by %s %s "
		return fmt.Sprintf(ob, strings.Join(order, ","), desc)
	} else {
		return ""
	}
}
func (ph *PaymentHistoryStorage) Insert(tx *sql.Tx, item *models.PaymentHistory) (*sql.Tx, int64, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s  (%s) VALUES (?, ?, ?, ?)", paymentHistoryTable, paymentHistoryInsertFields)

	params := []interface{}{
		item.UserId,
		item.Amount,
		item.TransactionType,
		item.Comment,
	}
	return queryExecTx(tx, ph.db, insertQuery, params, INSERT)
}
func (ph *PaymentHistoryStorage) Delete(tx *sql.Tx, pHist *models.PaymentHistory) (*sql.Tx, error) {
	return nil, nil
}
