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
	ipStoriesTable  = "`ip_stories`"
	ia              = "i"
	ipStoriesFields = ia + ".`id`, " +
		ia + ".`ip`," +
		ia + ".`user_agent`," +
		ia + ".`country`," +
		ia + ".`city`," +
		ia + ".`provider`," +
		ia + ".`company`," +
		ia + ".`link`," +
		ia + ".`date`"
	// 	Id       int       `json:"id"`
	// Ip       string    `json:"ip"`
	// Country  string    `json:"country"`
	// City     string    `json:"city"`
	// Provider string    `json:"provider"`
	// Company  string    `json:"company"`
	// Link     string    `json:"link"`
	// Date     time.Time `json:"date"`
	ipStoriesInsertFields = "`ip`," +
		"`user_agent`," +
		"`country`," +
		"`city`," +
		"`provider`," +
		"`company`," +
		"`link`"
)
const FieldIpStoriesCount = 8

type IpStoriesStorage struct {
	db *sql.DB
}

func NewIpStoriesStorage(conn *sql.DB) *IpStoriesStorage {
	return &IpStoriesStorage{db: conn}
}

// **********************************************************************************************************
func (ls *IpStoriesStorage) Find(filter *db.IpStoryFilter, offset, limit int) ([]*models.IpStory, error) {

	query := "SELECT " + ipStoriesFields + " FROM " + ipStoriesTable + " " + ia + " "

	where := prepareFindIpStoriesFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}

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

	return scanIpStories(rows, limit)

}

func prepareIpUsrFilter(filter *db.IpWithUserFileter) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter == nil {
		return stmt
	}

	if filter.UserId != "" {
		stmt.Conditions = append(stmt.Conditions, la+".`user_id` = ? ")
		stmt.Params = append(stmt.Params, filter.UserId)
	}
	if filter.Link != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`link` = ? ")
		stmt.Params = append(stmt.Params, filter.Link)
	}

	return stmt
}

func prepareFindIpStoriesFilter(filter *db.IpStoryFilter) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter == nil {
		return stmt
	}
	if filter.Id != "" {
		id, err := strconv.Atoi(filter.Id)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, ia+".`id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}
	if filter.Ip != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`ip` = ? ")
		stmt.Params = append(stmt.Params, filter.Ip)
	}
	if filter.UserAgent != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`user_agent` = ? ")
		stmt.Params = append(stmt.Params, filter.UserAgent)
	}
	if filter.Link != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`link` like ? ")
		stmt.Params = append(stmt.Params, filter.Link)
	}
	if filter.Country != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`country` like ? ")
		stmt.Params = append(stmt.Params, filter.Country)
	}
	if filter.Company != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`company` like ? ")
		stmt.Params = append(stmt.Params, filter.Company)
	}
	if filter.Provider != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`provider` like ? ")
		stmt.Params = append(stmt.Params, filter.Provider)
	}
	if filter.City != "" {
		stmt.Conditions = append(stmt.Conditions, ia+".`city` like ? ")
		stmt.Params = append(stmt.Params, filter.City)
	}
	if filter.DateStart != "" {
		t, err := time.Parse(layout, filter.DateStart)
		if err != nil {
			log.Println("filter.Date:", err)
		} else {
			stmt.Conditions = append(stmt.Conditions, ia+".`date` >= ? ")
			stmt.Params = append(stmt.Params, t)
		}
	}
	if filter.DateEnd != "" {
		t, err := time.Parse(layout, filter.DateEnd)
		if err != nil {
			log.Println("filter.Date:", err)
		} else {
			stmt.Conditions = append(stmt.Conditions, ia+".`date` <= ? ")
			stmt.Params = append(stmt.Params, t)
		}
	}
	return stmt
}

func prepareOrderIp(of *db.OrderIpStoryes, ta string) string {
	desc := ""
	var order []string

	condition := false
	if of.Date {
		order = append(order, ta+".`date`")
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

// **********************************************************************************************************
func (ls *IpStoriesStorage) Insert(tx *sql.Tx, item *models.IpStory) (*sql.Tx, int64, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s  (%s) VALUES (?, ?, ?, ?, ?, ?, ?)", ipStoriesTable, ipStoriesInsertFields)

	params := []interface{}{
		item.Ip,
		item.UserAgent,
		item.Country,
		item.City,
		item.Provider,
		item.Company,
		item.Link,
	}
	return queryExecTx(tx, ls.db, insertQuery, params, INSERT)

}

// **********************************************************************************************************
func (ls *IpStoriesStorage) Update(tx *sql.Tx, item *models.IpStory, filterWhere *db.IpStoryFilter) (*sql.Tx, int64, error) {
	set := updateSetIpStoriesFilter(item)
	setVal := strings.Join(set.Conditions[:], ", ")
	where := prepareFindIpStoriesFilter(filterWhere)
	whereVal := strings.ReplaceAll(strings.Join(where.Conditions[:], " AND "), ia+".", "")
	query := fmt.Sprintf(
		"UPDATE %s  SET %s WHERE %s",
		ipStoriesTable,
		setVal,
		whereVal,
	)
	set.Params = append(set.Params, where.Params...)

	return queryExecTx(tx, ls.db, query, set.Params, UPDATE)

}

func updateSetIpStoriesFilter(filter *models.IpStory) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter.Ip != "" {
		stmt.Conditions = append(stmt.Conditions, "`ip` = ? ")
		stmt.Params = append(stmt.Params, filter.Ip)
	}
	if filter.UserAgent != "" {
		stmt.Conditions = append(stmt.Conditions, "`user_agent` = ? ")
		stmt.Params = append(stmt.Params, filter.Ip)
	}
	if filter.Link != "" {
		stmt.Conditions = append(stmt.Conditions, "`link` = ? ")
		stmt.Params = append(stmt.Params, filter.Link)
	}
	if filter.Company != "" {
		stmt.Conditions = append(stmt.Conditions, "`company` = ? ")
		stmt.Params = append(stmt.Params, filter.Company)

	}
	if filter.Country != "" {
		stmt.Conditions = append(stmt.Conditions, "`country` = ? ")
		stmt.Params = append(stmt.Params, filter.Country)
	}
	if filter.City != "" {
		stmt.Conditions = append(stmt.Conditions, "`city` = ? ")
		stmt.Params = append(stmt.Params, filter.City)
	}
	if filter.Provider != "" {
		stmt.Conditions = append(stmt.Conditions, "`provider` = ? ")
		stmt.Params = append(stmt.Params, filter.Provider)

	}

	return stmt
}

//**********************************************************************************************************

func (ls *IpStoriesStorage) Delete(tx *sql.Tx, filter *db.IpStoryFilter) (*sql.Tx, error) {
	query := fmt.Sprintf("DELETE FROM `%s` ", linksTable)
	where := prepareFindIpStoriesFilter(filter)
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

func scanIpStories(rows *sql.Rows, cap int) ([]*models.IpStory, error) {
	ip := make([]*models.IpStory, 0, cap)
	for rows.Next() {
		g, err := scanIpStory(rows)
		if err != nil {
			return nil, err
		}
		ip = append(ip, g)
	}
	return ip, nil
}
func scanIpStory(row Scanner) (*models.IpStory, error) {
	item := &models.IpStory{}
	var (
		userAgent sql.NullString
		// imgPath     sql.NullString
		// editedPath  sql.NullString
		// description sql.NullString
		// note        sql.NullString
	)
	err := row.Scan(
		&item.Id,
		&item.Ip,
		&userAgent,
		&item.Country,
		&item.City,
		&item.Provider,
		&item.Company,
		&item.Link,
		&item.Date,
	)

	if userAgent.Valid {
		item.UserAgent = userAgent.String
	}
	return item, err
}
