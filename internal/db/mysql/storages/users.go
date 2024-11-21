package storages

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"
)

const (
	usersTable  = "`users`"
	ua          = "u"
	usersFields = ua + ".`id`, " +
		ua + ".`user_id`," +
		ua + ".`user_name`," +
		ua + ".`message_type`," +
		ua + ".`command`," +
		ua + ".`status`," +
		ua + ".`history_user_name`," +
		ua + ".`balance`," +
		ua + ".`balance_all_time`," +
		ua + ".`auto_pay`," +
		ua + ".`test_used`," +
		ua + ".`last_msg_id`," +
		ua + ".`referal_id`," +
		ua + ".`time_create`," +
		ua + ".`time_update`"

	usersInsertFields = "`user_id`," +
		"`user_name`," +
		"`message_type`," +
		"`command`," +
		"`status`," +
		"`history_user_name`," +
		"`balance`," +
		"`balance_all_time`," +
		"`auto_pay`," +
		"`test_used`," +
		"`last_msg_id`," +
		"`referal_id`"
)
const FieldUsersCount = 11

type UsersStorage struct {
	db *sql.DB
}

func NewUsersStorage(conn *sql.DB) *UsersStorage {
	return &UsersStorage{db: conn}
}

// **********************************************************************************************************
func (ls *UsersStorage) Find(filter *db.UserFilter, orderFilter *db.OrderByUsers, offset, limit int) ([]*models.User, error) {

	query := "SELECT " + usersFields + " FROM " + usersTable + " " + ua + " "

	where := prepareFindUserFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}
	order := prepareOrderUsers(orderFilter)
	query += order
	if limit == 0 {
		limit = db.DefaultLimit
	}
	query += " LIMIT " + strconv.Itoa(int(offset)) + ", " + strconv.Itoa(int(limit))

	stmt, err := ls.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	// fmt.Println(query)
	rows, err := stmt.Query(where.Params...)
	defer func() {
		_ = rows.Close()
		_ = stmt.Close()
	}()
	if err != nil {
		return nil, err
	}

	return scanUsers(rows, limit)

}
func prepareFindUserFilter(filter *db.UserFilter) *QueryStmt {
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
			stmt.Conditions = append(stmt.Conditions, ua+".`id` = ? ")
			stmt.Params = append(stmt.Params, id)
		}
	}

	if filter.UserId != "" {
		// id, err := strconv.ParseInt(filter.UserId, 10, 64)
		// if err == nil {
		stmt.Conditions = append(stmt.Conditions, ua+".`user_id` = ? ")
		stmt.Params = append(stmt.Params, filter.UserId)
		// }
	}
	if filter.UserName != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`user_name` like ? ")
		stmt.Params = append(stmt.Params, filter.UserName)
	}
	if filter.Password != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`password` like ? ")
		stmt.Params = append(stmt.Params, filter.Password)
	}
	if filter.MessageType != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`message_type` = ? ")
		stmt.Params = append(stmt.Params, filter.MessageType)
	}
	if filter.Command != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`command` like ? ")
		stmt.Params = append(stmt.Params, filter.Command)
	}
	if filter.Lang != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`lang` like ? ")
		stmt.Params = append(stmt.Params, filter.Lang)
	}

	// "`is_deanon`," +
	if filter.HistoryUserName != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`history_user_name` like ? ")
		stmt.Params = append(stmt.Params, filter.HistoryUserName)
	}
	if filter.WhateDescription != "" {
		wd, err := strconv.Atoi(filter.WhateDescription)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, ua+".`whate_description` = ? ")
			stmt.Params = append(stmt.Params, wd)
		}
	}
	if filter.LinkDescription != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`link_description` like ? ")
		stmt.Params = append(stmt.Params, filter.LinkDescription)
	}
	if filter.City != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`city` like ? ")
		stmt.Params = append(stmt.Params, filter.City)
	}
	if filter.Phone != "" {
		stmt.Conditions = append(stmt.Conditions, ua+".`phone` like ? ")
		stmt.Params = append(stmt.Params, filter.Phone)
	}
	if filter.IsDeanon != "" {
		isd, err := strconv.Atoi(filter.WhateDescription)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, ua+".`is_deanon` like ? ")
			stmt.Params = append(stmt.Params, isd)
		}
	}
	if filter.Status != "" {
		status, err := strconv.Atoi(filter.Id)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, ua+".`status` = ? ")
			stmt.Params = append(stmt.Params, status)
		}
	}
	if filter.Role != "" {
		role, err := strconv.Atoi(filter.Id)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, ua+".`role` = ? ")
			stmt.Params = append(stmt.Params, role)
		}
	}
	if filter.Count != "" {
		count, err := strconv.Atoi(filter.Count)
		if err == nil {
			stmt.Conditions = append(stmt.Conditions, ua+".`count` > ? ")
			stmt.Params = append(stmt.Params, count)
		}
	} else {
		stmt.Conditions = append(stmt.Conditions, ua+".`count` > ? ")
		stmt.Params = append(stmt.Params, 0)
	}
	return stmt
}

func prepareOrderUsers(of *db.OrderByUsers) string {
	desc := ""
	var order []string

	condition := false
	if of.Count {
		order = append(order, "count")
		condition = true
	}
	if of.Id {
		order = append(order, "id")
		condition = true
	}
	if of.TimeCreate {
		order = append(order, "time_create")
		condition = true
	}
	if of.TimeDeanon {
		order = append(order, "time_deanon")
		condition = true
	}
	if of.TimeUpdate {
		order = append(order, "time_update")
		condition = true
	}
	if of.Desc {
		desc = "desc"
	}
	if condition {
		ob := " order by %s %s "
		return fmt.Sprintf(ob, strings.Join(order, ","), desc)
	} else {
		return ""
	}
}

func (ls *UsersStorage) FindOne(userID int64) (*models.User, error) {

	query := "SELECT " + usersFields + " FROM " + usersTable + " " + ua + " WHERE " + ua + ".`user_id`=?"

	stmt, err := ls.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(userID)
	defer func() {
		_ = stmt.Close()
	}()
	if err != nil {
		return nil, err
	}
	return scanUser(row)
}
func (us *UsersStorage) UpdateUser(tx *sql.Tx, user *models.User, userID int64) (*sql.Tx, error) {
	set := updateSetUserFilter(user)
	setVal := strings.Join(set.Conditions[:], ", ")

	query := fmt.Sprintf(
		"UPDATE %s  SET %s WHERE `user_id`=? ;",
		usersTable,
		setVal,
	)
	set.Params = append(set.Params, userID)

	tx, _, err := queryExecTx(tx, us.db, query, set.Params, UPDATE)

	return tx, err

}

// **********************************************************************************************************
func (us *UsersStorage) Insert(tx *sql.Tx, item *models.User) (*sql.Tx, int64, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s  (%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", usersTable, usersInsertFields)
	qs := insertUserFilter(item)

	return queryExecTx(tx, us.db, insertQuery, qs.Params, INSERT)

}

func insertUserFilter(filter *models.User) *QueryStmt {
	stmt := &QueryStmt{
		Params: make([]interface{}, 0, FieldCount),
	}

	if filter.UserId > 0 {
		stmt.Params = append(stmt.Params, filter.UserId)
	} else {
		return stmt
	}

	if filter.UserName != models.DEFOULTVALUE {
		stmt.Params = append(stmt.Params, filter.UserName)
	} else {
		stmt.Params = append(stmt.Params, "")
	}

	if filter.MessageType >= 0 {
		stmt.Params = append(stmt.Params, filter.MessageType)
	} else {
		stmt.Params = append(stmt.Params, 0)
	}
	if filter.Command != models.DEFOULTVALUE {
		stmt.Params = append(stmt.Params, filter.Command)
	} else {
		stmt.Params = append(stmt.Params, "")
	}
	if filter.Status >= 0 {
		stmt.Params = append(stmt.Params, filter.Status)
	} else {
		stmt.Params = append(stmt.Params, 0)
	}
	if filter.HistoryUserName != models.DEFOULTVALUE {
		stmt.Params = append(stmt.Params, filter.HistoryUserName)
	} else {
		stmt.Params = append(stmt.Params, "")
	}
	if filter.Balance > 0 {
		stmt.Params = append(stmt.Params, filter.Balance)
	} else {
		stmt.Params = append(stmt.Params, 0)
	}
	if filter.BalanceAllTime > 0 {
		stmt.Params = append(stmt.Params, filter.BalanceAllTime)
	} else {
		stmt.Params = append(stmt.Params, 0)
	}

	stmt.Params = append(stmt.Params, filter.AutoPay)
	stmt.Params = append(stmt.Params, filter.TestUsed)

	if filter.LastMsgId > 0 {
		stmt.Params = append(stmt.Params, filter.LastMsgId)
	} else {
		stmt.Params = append(stmt.Params, 0)
	}
	if filter.ReferalId > 0 {
		stmt.Params = append(stmt.Params, filter.ReferalId)
	} else {
		stmt.Params = append(stmt.Params, 0)
	}

	return stmt
}

// **********************************************************************************************************
func (us *UsersStorage) Update(tx *sql.Tx, item *models.User, filterWhere *db.UserFilter) (*sql.Tx, int64, error) {
	set := updateSetUserFilter(item)
	setVal := strings.Join(set.Conditions[:], ", ")
	where := prepareFindUserFilter(filterWhere)
	whereVal := strings.ReplaceAll(strings.Join(where.Conditions[:], " AND "), ua+".", "")
	query := fmt.Sprintf(
		"UPDATE %s  SET %s WHERE %s",
		usersTable,
		setVal,
		whereVal,
	)
	set.Params = append(set.Params, where.Params...)

	return queryExecTx(tx, us.db, query, set.Params, UPDATE)

}
func updateSetUserFilter(filter *models.User) *QueryStmt {
	stmt := &QueryStmt{
		Conditions: make([]string, 0, FieldCount),
		Params:     make([]interface{}, 0, FieldCount),
	}
	if filter.UserId > 0 {
		stmt.Conditions = append(stmt.Conditions, "`user_id` = ? ")
		stmt.Params = append(stmt.Params, filter.UserId)
	}
	if filter.UserName != models.DEFOULTVALUE {
		stmt.Conditions = append(stmt.Conditions, "`user_name` = ? ")
		stmt.Params = append(stmt.Params, filter.UserName)
	}
	if filter.MessageType >= 0 {
		stmt.Conditions = append(stmt.Conditions, "`message_type` = ? ")
		stmt.Params = append(stmt.Params, filter.MessageType)
	}
	if filter.Command != models.DEFOULTVALUE {
		stmt.Conditions = append(stmt.Conditions, "`command` = ? ")
		stmt.Params = append(stmt.Params, filter.Command)
	}
	if filter.Status >= 0 {
		stmt.Conditions = append(stmt.Conditions, "`status` = ? ")
		stmt.Params = append(stmt.Params, filter.Status)
	}
	if filter.HistoryUserName != models.DEFOULTVALUE {
		stmt.Conditions = append(stmt.Conditions, "`history_user_name` = ? ")
		stmt.Params = append(stmt.Params, filter.HistoryUserName)
	}
	if filter.Balance >= 0 {
		stmt.Conditions = append(stmt.Conditions, "`balance` = ? ")
		stmt.Params = append(stmt.Params, filter.Balance)
	}
	if filter.BalanceAllTime >= 0 {
		stmt.Conditions = append(stmt.Conditions, "`balance_all_time` = ? ")
		stmt.Params = append(stmt.Params, filter.BalanceAllTime)
	}
	if filter.ReferalId > 0 {
		stmt.Conditions = append(stmt.Conditions, "`referal_id` = ? ")
		stmt.Params = append(stmt.Params, filter.ReferalId)
	}
	if filter.LastMsgId > 0 {
		stmt.Conditions = append(stmt.Conditions, "`last_msg_id` = ? ")
		stmt.Params = append(stmt.Params, filter.LastMsgId)
	}
	stmt.Conditions = append(stmt.Conditions, "`auto_pay` = ? ")
	stmt.Params = append(stmt.Params, filter.AutoPay)

	stmt.Conditions = append(stmt.Conditions, "`test_used` = ? ")
	stmt.Params = append(stmt.Params, filter.TestUsed)

	return stmt
}

//**********************************************************************************************************

func (us *UsersStorage) Delete(tx *sql.Tx, filter *db.UserFilter) (*sql.Tx, error) {
	query := fmt.Sprintf("DELETE FROM %s ", usersTable)
	where := prepareFindUserFilter(filter)
	if len(where.Conditions) > 0 {
		query += " WHERE " + strings.Join(where.Conditions[:], " AND ")
	}
	tx, _, err := queryExecTx(
		tx,
		us.db,
		query,
		where.Params,
		DELETE,
	)
	return tx, err

}

//**********************************************************************************************************

func scanUsers(rows *sql.Rows, cap int) ([]*models.User, error) {
	users := make([]*models.User, 0, cap)
	for rows.Next() {
		g, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, g)
	}
	return users, nil
}

func scanUser(row Scanner) (*models.User, error) {
	item := &models.User{}
	var (
		timeCreate, timeUpdate sql.NullTime
		// password               sql.NullString
		// imgPath     sql.NullString
		// editedPath  sql.NullString
		// description sql.NullString
		// note        sql.NullString
	)

	err := row.Scan(
		&item.Id,
		&item.UserId,
		&item.UserName,
		&item.MessageType,
		&item.Command,
		&item.Status,
		&item.HistoryUserName,
		&item.Balance,
		&item.BalanceAllTime,
		&item.AutoPay,
		&item.TestUsed,
		&item.LastMsgId,
		&item.ReferalId,
		&timeCreate,
		&timeUpdate,
	)

	if timeCreate.Valid {
		item.TimeCreate = timeCreate.Time

	}
	if timeUpdate.Valid {
		item.TimeUpdate = timeUpdate.Time

	}

	return item, err
}
