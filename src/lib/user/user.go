package user

import (
	"database/sql"

	"backend/src/app"
	"backend/src/lib/icontext"
)

type UserRec struct {
	UserId    int64
	Name      string
	Login     string
	Email     string
	Password  string
	Status    string
	RoleList  map[string]bool
	RightList map[string]bool
}

func New(ctx icontext.Ctx, id int64) *UserRec {
	rv := new(UserRec)
	rv.RoleList = make(map[string]bool, 10)
	rv.RightList = make(map[string]bool, 50)

	rv.setUserInfo(id)
	rv.setRoleList()
	rv.setRightList()

	return rv
}

func (rec UserRec) RoleLen() int {
	return len(rec.RoleList)
}

func (rec UserRec) HasRole(name string) bool {
	_, ok := rec.RoleList[name]

	return ok
}

func (rec UserRec) IsSuperUser() bool {
	return rec.Login == "superuser"
}

func (rec *UserRec) setUserInfo(id int64) {
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
					user.userId = ?`

	row := app.Db.QueryRow(sqlStr, id)

	err := row.Scan(&rec.UserId, &rec.Name, &rec.Login, &rec.Name,
		&rec.Email, &rec.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			panic("Could not read user info.")
		}

		panic(err)
	}
}

func (rec *UserRec) setRoleList() {
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

	rows, err := app.Db.Query(sqlStr, rec.UserId)
	if err != nil {
		panic("Could not read user roles." + " " + err.Error())
	}

	defer rows.Close()

	var roleName string
	for rows.Next() {
		if err = rows.Scan(&roleName); err != nil {
			panic(err.Error())
		}

		rec.RoleList[roleName] = true
	}
}

func (rec *UserRec) setRightList() {
	sqlStr := `select
					roleRight.pageName,
					roleRight.funcName
				from
					userRole,
					roleRight,
					user
				where
				 	roleRight.roleId = userRole.roleId and
					userRole.userId = user.userId and
					user.userId = ?`

	rows, err := app.Db.Query(sqlStr, rec.UserId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var pageName, funcName string
	for rows.Next() {
		if err = rows.Scan(&pageName, &funcName); err != nil {
			panic(err)
		}

		rec.RightList[app.MakeKey(pageName, funcName)] = true
	}
}

func (rec UserRec) IsRight(pageName, funcName string) bool {
	if rec.Status == "blocked" {
		return false
	}

	if rec.RightList[app.MakeKey(pageName, funcName)] || rec.Login == "superuser" {
		return true
	}

	return false
}
