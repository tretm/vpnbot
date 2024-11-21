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
	linksTable  = "`links`"
	la          = "l"
	linksFields = la + ".`id`, " +
		la + ".`user_id`," +
		la + ".`link`," +
		la + ".`vpn_link`," +
		la + ".`vpn_link_id`," +
		la + ".`vpn_link_password`," +
		la + ".`state`," +
		la + ".`time_end`," +
		la + ".`time_create`," +
		la + ".`time_update`"

	linksInsertFields = "`user_id`," +
		"`link`," +
		"`vpn_link`," +
		"`vpn_link_id`," +
		"`vpn_link_password`," +
		"`state`," +
		"`time_end`"

	layout = "2006-01-02 15:04:05"
)
const FieldCount = 10

type LinkStorage struct {
	db *sql.DB
}

func NewLinkStorage(conn *sql.DB) *LinkStorage {
	return &LinkStorage{db: conn}
}

// **********************************************************************************************************
func (ls *LinkStorage) Find(filter *db.LinkFilter, orderFilter *db.OrderLinks, offset, limit int) ([]*models.Link, error) {

	query := "SELECT " + linksFields + " FROM " + linksTable + " " + la + " "

	where := prepareFindLinkFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}

	order := prepareOrderLinks(orderFilter)
	query += order

	if limit == 0 {
		limit = db.DefaultLimit
	}

	query += " LIMIT " + strconv.Itoa(int(offset)) + ", " + strconv.Itoa(int(limit))

	stmt, err := ls.db.Prepare(query)
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

	return scanLinks(rows, limit)

}
func prepareOrderLinks(of *db.OrderLinks) string {
	desc := ""
	var order []string

	condition := false
	if of.Date {
		order = append(order, "time_create")
		condition = true
	}
	if of.Desk {
		desc = "desc"

	}
	if condition {
		ob := " order by %s %s "
		return fmt.Sprintf(ob, strings.Join(order, ","), desc)
	} else {
		return ""
	}
}

func prepareFindLinkFilter(filter *db.LinkFilter) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter.Id != "" {
		id, err := strconv.Atoi(filter.Id)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, la+".`id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}
	if filter.UserId != "" {
		id, err := strconv.ParseInt(filter.UserId, 10, 64)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, la+".`user_id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}
	if filter.Link != "" {
		stmt.Conditions = append(stmt.Conditions, la+".`link` like ? ")
		stmt.Params = append(stmt.Params, filter.Link)
	}

	if filter.VpnLink != "" {
		stmt.Conditions = append(stmt.Conditions, la+".`vpn_link` like ? ")
		stmt.Params = append(stmt.Params, filter.VpnLink)
	}
	if filter.VpnLinkId != "" {
		stmt.Conditions = append(stmt.Conditions, la+".`vpn_link_id` = ? ")
		stmt.Params = append(stmt.Params, filter.VpnLinkId)
	}
	if filter.VpnLinkPassword != "" {
		stmt.Conditions = append(stmt.Conditions, la+".`vpn_link_password` like ? ")
		stmt.Params = append(stmt.Params, filter.VpnLinkPassword)
	}
	if filter.State != "" {
		stmt.Conditions = append(stmt.Conditions, la+".`state` = ? ")
		stmt.Params = append(stmt.Params, filter.State)
	}
	if filter.TimeEnd != "" {
		t, err := time.Parse(layout, filter.TimeEnd+" 00:00:00")
		if err != nil {
			log.Println("storages.prepareFindLinkFilter, filter.Date:", err)
		} else {
			stmt.Conditions = append(stmt.Conditions, la+".`time_end` < ? ")
			stmt.Params = append(stmt.Params, t)
		}
	}

	// if filter.DateEnd != "" {
	// 	t, err := time.Parse(layout, filter.DateEnd+" 23:59:59")
	// 	if err != nil {
	// 		log.Println("storages.prepareFindLinkFilter, filter.Date:", err)
	// 	} else {
	// 		stmt.Conditions = append(stmt.Conditions, la+".`date` <= ? ")
	// 		stmt.Params = append(stmt.Params, t)
	// 	}
	// }
	return stmt
}

// **********************************************************************************************************
func (ls *LinkStorage) Insert(tx *sql.Tx, item *models.Link) (*sql.Tx, int64, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s  (%s) VALUES (?, ?, ?, ?, ?, ?, ?)", linksTable, linksInsertFields)

	params := []interface{}{
		item.UserId,
		item.Link,
		item.VpnLink,
		item.VpnLinkId,
		item.VpnPassword,
		item.State,
		item.TimeEnd,
	}
	return queryExecTx(tx, ls.db, insertQuery, params, INSERT)

}

// **********************************************************************************************************
func (ls *LinkStorage) Update(tx *sql.Tx, item *models.Link, filterWhere *db.LinkFilter) (*sql.Tx, int64, error) {
	set := updateSetLinkFilter(item)
	setVal := strings.Join(set.Conditions[:], ", ")
	where := prepareFindLinkFilter(filterWhere)
	whereVal := strings.ReplaceAll(strings.Join(where.Conditions[:], " AND "), la+".", "")
	query := fmt.Sprintf(
		"UPDATE %s  SET %s WHERE %s",
		linksTable,
		setVal,
		whereVal,
	)
	set.Params = append(set.Params, where.Params...)
	return queryExecTx(tx, ls.db, query, set.Params, UPDATE)

}

func updateSetLinkFilter(filter *models.Link) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter.UserId > 0 {
		stmt.Conditions = append(stmt.Conditions, "`user_id` = ? ")
		stmt.Params = append(stmt.Params, filter.UserId)
	}
	if filter.State != "" {
		stmt.Conditions = append(stmt.Conditions, "`state` = ? ")
		stmt.Params = append(stmt.Params, filter.State)
	}
	if filter.Link != "" {
		stmt.Conditions = append(stmt.Conditions, "`link` = ? ")
		stmt.Params = append(stmt.Params, filter.Link)

	}
	if filter.VpnLink != "" {
		stmt.Conditions = append(stmt.Conditions, "`vpn_link` = ? ")
		stmt.Params = append(stmt.Params, filter.VpnLink)

	}
	if filter.VpnLinkId != "" {
		stmt.Conditions = append(stmt.Conditions, "`vpn_link_id` = ? ")
		stmt.Params = append(stmt.Params, filter.VpnLinkId)
	}
	if filter.VpnPassword != "" {
		stmt.Conditions = append(stmt.Conditions, "`vpn_link_password` = ? ")
		stmt.Params = append(stmt.Params, filter.VpnPassword)
	}
	if !filter.TimeEnd.IsZero() {
		stmt.Conditions = append(stmt.Conditions, "`time_end` = ? ")
		stmt.Params = append(stmt.Params, filter.TimeEnd)
	}

	return stmt
}

//**********************************************************************************************************

func (ls *LinkStorage) Delete(tx *sql.Tx, filter *db.LinkFilter) (*sql.Tx, error) {
	query := fmt.Sprintf("DELETE FROM `%s` ", linksTable)
	where := prepareFindLinkFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}
	tx, _, err := queryExecTx(
		tx,
		ls.db,
		query,
		where.Params,
		DELETE,
	)
	return tx, err

}

//**********************************************************************************************************

func scanLinks(rows *sql.Rows, cap int) ([]*models.Link, error) {
	list := make([]*models.Link, 0, cap)
	for rows.Next() {
		g, err := scanLink(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, g)
	}
	if len(list) == 0 {
		return list, errors.New("links not found")
	}
	return list, nil
}

func scanLink(row Scanner) (*models.Link, error) {
	item := &models.Link{}
	var (
		// userId      sql.NullString
		link        sql.NullString
		VpnLink     sql.NullString
		VpnPassword sql.NullString
		state       sql.NullString
	)
	err := row.Scan(
		&item.Id,
		&item.UserId,
		&link,
		&VpnLink,
		&item.VpnLinkId,
		&VpnPassword,
		&state,
		&item.TimeEnd,
		&item.TimeCreate,
		&item.TimeUpdate,
	)

	if link.Valid {
		item.Link = link.String
	}
	if VpnLink.Valid {
		item.VpnLink = VpnLink.String
	}
	if VpnPassword.Valid {
		item.VpnPassword = VpnPassword.String
	}
	if state.Valid {
		item.State = state.String
	}
	return item, err
}
