package user_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

type UserRec struct {
	UserId    int64
	Name      string
	Login     string
	Email     string
	Password  string
	ResetKey  string
	ResetTime int64
	Status    string
}

func GetUserRec(userId int64) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.password,
					user.status
				from
					user
				where
					user.userId = ?`

	row := app.Db.QueryRow(sqlStr, userId)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email,
		&rec.Password, &rec.Status)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountUser(key string, roleId int64, status string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*)")

	fromBuf := util.NewBuf()
	fromBuf.Add("user")

	conBuf := util.NewBuf()

	if roleId != -1 {
		fromBuf.Add("userRole")
		conBuf.Add("(user.userId = userRole.userId)")
		conBuf.Add("(userRole.roleId = %d)", roleId)
	}

	if status != "default" {
		conBuf.Add("(user.status = '%s')", util.DbStr(status))
	}

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(user.name like('%%%s%%') or
					user.email like('%%%s%%') or
					user.login like('%%%s%%'))`, key, key, key)
	}

	sqlBuf.Add("from " + *fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where " + *conBuf.StringSep(" and "))
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func CountAdmin() int64 {
	sqlStr := `select
					count(*)
				from
					user,
					userRole,
					role
				where
					user.userId = userRole.userId and
					userRole.roleId = role.roleId and
					role.name = 'Admin'`

	row := app.Db.QueryRow(sqlStr)

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func IsUserRoleExists(userId int64, roleName string) bool {
	userRoleList := GetUserRoleList(userId)
	for _, v := range userRoleList {
		if v == roleName {
			return true
		}
	}

	return false
}

func GetUserPage(ctx *context.Ctx, key string, pageNo, roleId int64, status string) []*UserRec {

	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status`)

	fromBuf := util.NewBuf()
	fromBuf.Add("user")

	conBuf := util.NewBuf()

	if roleId != -1 {
		fromBuf.Add("userRole")
		conBuf.Add("(user.userId = userRole.userId)")
		conBuf.Add("(userRole.roleId = %d)", roleId)
	}

	if status != "default" {
		conBuf.Add("(user.status = '%s')", util.DbStr(status))
	}

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(user.name like('%%%s%%') or
					user.email like('%%%s%%') or
					user.login like('%%%s%%'))`, key, key, key)
	}

	sqlBuf.Add("from " + *fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add(" where " + *conBuf.StringSep(" and "))
	}

	sqlBuf.Add("order by user.name")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*UserRec, 0, 100)
	for rows.Next() {
		rec := new(UserRec)
		err = rows.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email, &rec.Status)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetUserRecByLogin(name, password string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.name,
					user.email,
					user.status
				from
					user
				where
					user.login = ? and
					user.password = ?`

	row := app.Db.QueryRow(sqlStr, name, password)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Name, &rec.Email,
		&rec.Status)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRecByEmail(email string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime
				from
					user
				where
					user.email = ?`

	row := app.Db.QueryRow(sqlStr, email)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email,
		&rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRecByResetKey(key string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime
				from
					user
				where
					user.resetKey = ?`

	row := app.Db.QueryRow(sqlStr, key)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login,
		&rec.Email, &rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserUserRecBySessionId(ctx context.Ctx, sessionId string) (*UserRec, error) {
	sqlStr := `select
					user.userId,
					user.name,
					user.login,
					user.email,
					user.status,
					user.resetKey,
					user.resetTime
				from
					user,
					session
				where
					session.sessionId = ? and
					user.userId = session.userId`

	row := app.Db.QueryRow(sqlStr, sessionId)

	rec := new(UserRec)
	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Email,
		&rec.Status, &rec.Status, &rec.ResetKey, &rec.ResetTime)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetUserRoleList(userId int64) []string {
	sqlStr := `select
					role.name
				from
					role,
					userRole
				where
					role.roleId = userRole.roleId and
					userRole.userId = ? 
				order by
					name`

	rows, err := app.Db.Query(sqlStr, userId)
	if err != nil {
		panic("Could not read user roles." + " " + err.Error())
	}

	defer rows.Close()

	var rv []string
	var roleName string
	for rows.Next() {
		if err = rows.Scan(&roleName); err != nil {
			panic(err.Error())
		}

		rv = append(rv, roleName)
	}

	return rv
}

func StatusToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "active":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Active"))
	case "blocked":
		return fmt.Sprintf("<span class=\"label labelError labelXs\">%s</span>", ctx.T("Blocked"))
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
