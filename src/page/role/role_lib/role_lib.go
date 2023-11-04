package role_lib

import (
	"backend/src/app"
	"backend/src/lib/context"
)

type RoleRec struct {
	RoleId int64
	Name   string
	Exp    string
}

type RoleRightRec struct {
	RoleId   int64
	PageName string
	FuncName string
}

func GetRoleRec(roleId int64) (*RoleRec, error) {
	sqlStr := `select
					role.roleId,
					role.name,
					role.exp
				from
					role
				where
					role.roleId = ?`

	row := app.Db.QueryRow(sqlStr, roleId)

	rec := new(RoleRec)
	err := row.Scan(&rec.RoleId, &rec.Name, &rec.Exp)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetRoleList() []*RoleRec {
	sqlStr := `select
					role.roleId,
					role.name,
					role.exp
				from
					role
				order by role.name`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*RoleRec, 0, 20)
	for rows.Next() {
		rec := new(RoleRec)
		if err = rows.Scan(&rec.RoleId, &rec.Name, &rec.Exp); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func CountRoleByUser(userId int64) (int64, error) {
	sqlStr := `select
					count(*)
				from
					role,
					userRole
				where
					role.roleId = userRole.roleId and
					userRole.userId = ?`

	rows := app.Db.QueryRow(sqlStr, userId)

	var rv int64
	if err := rows.Scan(&rv); err != nil {
		return 0, err
	}

	return rv, nil
}

func GetRoleByUser(userId int64) ([]*RoleRec, error) {
	sqlStr := `select
					role.roleId,
					role.name,
					role.exp
				from
					role,
					userRole
				where
					role.roleId = userRole.roleId and
					userRole.userId = ?
				order by
					role.name`

	rows, err := app.Db.Query(sqlStr, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rv := make([]*RoleRec, 0, 10)
	for rows.Next() {
		rec := new(RoleRec)
		if err = rows.Scan(&rec.RoleId, &rec.Name, &rec.Exp); err != nil {
			return nil, err
		}

		rv = append(rv, rec)
	}

	return rv, nil
}

func GetRoleRight(roleId int64) (map[string]bool, error) {
	sqlStr := `select
					roleId,
					pageName,
					funcName
				from
					roleRight
				where
					roleId = ?
				order by
					pageName, funcName`

	rows, err := app.Db.Query(sqlStr, roleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rv := make(map[string]bool, 50)
	for rows.Next() {
		rec := new(RoleRightRec)
		err = rows.Scan(&rec.RoleId, &rec.PageName, &rec.FuncName)
		if err != nil {
			return nil, err
		}

		rv[app.MakeKey(rec.PageName, rec.FuncName)] = true
	}

	return rv, nil
}

func SetRoleRightList(ctx *context.Ctx) {
	app.UserRightList = app.NewRightList()
	var rl = app.UserRightList
	var tab *app.RightTab

	//login
	tab = app.NewRightTab("Site")
	tab.Add(app.NewRight("site", "login", ctx.T("Can login admin page.")))
	tab.Add(app.NewRight("site", "select_user", ctx.T("Can change effective user.")))
	rl.Add("site", tab)

	//translation
	tab = app.NewRightTab("Translation Table")
	tab.Add(app.NewRight("tran", "browse", ctx.T("Can browse translation records.")))
	tab.Add(app.NewRight("tran", "insert", ctx.T("Can insert translation records.")))
	tab.Add(app.NewRight("tran", "update", ctx.T("Can update translation records.")))
	tab.Add(app.NewRight("tran", "delete", ctx.T("Can delete translation records.")))
	tab.Add(app.NewRight("tran", "synchronize", ctx.T("Can synchronize translation units with source code.")))
	tab.Add(app.NewRight("tran", "export", ctx.T("Can export translation table into a json file.")))
	tab.Add(app.NewRight("tran", "import", ctx.T("Can import translation table from json file.")))
	rl.Add("tran", tab)

	//user
	tab = app.NewRightTab("Users")
	tab.Add(app.NewRight("user", "browse", ctx.T("Can browse user records.")))
	tab.Add(app.NewRight("user", "insert", ctx.T("Can insert user record.")))
	tab.Add(app.NewRight("user", "update", ctx.T("Can update user records.")))
	tab.Add(app.NewRight("user", "confirm", ctx.T("Can confirm new user accounts.")))
	tab.Add(app.NewRight("user", "update_pass", ctx.T("Can update user passowrd.")))
	tab.Add(app.NewRight("user", "delete", ctx.T("Can delete user records.")))
	tab.Add(app.NewRight("user", "role_browse", ctx.T("Can browse user roles.")))
	tab.Add(app.NewRight("user", "role_revoke", ctx.T("Can revoke user roles.")))
	tab.Add(app.NewRight("user", "role_grant", ctx.T("Can grant user roles.")))
	rl.Add("user", tab)

	//role
	tab = app.NewRightTab("Roles")
	tab.Add(app.NewRight("role", "browse", ctx.T("Can browse role records.")))
	tab.Add(app.NewRight("role", "insert", ctx.T("Can insert role record.")))
	tab.Add(app.NewRight("role", "update", ctx.T("Can update role records.")))
	tab.Add(app.NewRight("role", "delete", ctx.T("Can delete role records.")))
	tab.Add(app.NewRight("role", "role_right", ctx.T("Can update role rights.")))
	rl.Add("role", tab)

	//config
	tab = app.NewRightTab("Configuration")
	tab.Add(app.NewRight("config", "browse", ctx.T("Can browse configuration records.")))
	tab.Add(app.NewRight("config", "set", ctx.T("Can set configuration record value.")))
	rl.Add("config", tab)

	//profile
	tab = app.NewRightTab("Profile")
	tab.Add(app.NewRight("profile", "browse", ctx.T("Can browse own profile record.")))
	tab.Add(app.NewRight("profile", "update", ctx.T("Can update own profile record.")))
	tab.Add(app.NewRight("profile", "update_password", ctx.T("Can change own password.")))
	rl.Add("profile", tab)

	//site error
	tab = app.NewRightTab("Site Error")
	tab.Add(app.NewRight("site_error", "browse", ctx.T("Can browse site error record.")))
	tab.Add(app.NewRight("site_error", "delete", ctx.T("Can delete site error record.")))
	rl.Add("site_error", tab)

	//banner
	tab = app.NewRightTab("Banner")
	tab.Add(app.NewRight("banner", "browse", ctx.T("Can browse banner records.")))
	tab.Add(app.NewRight("banner", "insert", ctx.T("Can insert banner records.")))
	tab.Add(app.NewRight("banner", "update", ctx.T("Can update banner records.")))
	tab.Add(app.NewRight("banner", "delete", ctx.T("Can delete banner records.")))
	rl.Add("banner", tab)

	//help
	tab = app.NewRightTab("Help")
	tab.Add(app.NewRight("help", "browse", ctx.T("Can browse help records.")))
	tab.Add(app.NewRight("help", "insert", ctx.T("Can insert help records.")))
	tab.Add(app.NewRight("help", "update", ctx.T("Can update help records.")))
	tab.Add(app.NewRight("help", "delete", ctx.T("Can delete help records.")))
	rl.Add("help", tab)

	//news
	tab = app.NewRightTab("News")
	tab.Add(app.NewRight("news", "browse", ctx.T("Can browse news records.")))
	tab.Add(app.NewRight("news", "insert", ctx.T("Can insert news records.")))
	tab.Add(app.NewRight("news", "update", ctx.T("Can update news records.")))
	tab.Add(app.NewRight("news", "delete", ctx.T("Can delete news records.")))
	rl.Add("news", tab)

	//faq
	tab = app.NewRightTab("FAQ")
	tab.Add(app.NewRight("faq", "browse", ctx.T("Can browse faq records.")))
	tab.Add(app.NewRight("faq", "insert", ctx.T("Can insert faq records.")))
	tab.Add(app.NewRight("faq", "update", ctx.T("Can update faq records.")))
	tab.Add(app.NewRight("faq", "delete", ctx.T("Can delete faq records.")))
	rl.Add("faq", tab)

	//text content
	tab = app.NewRightTab("Text Content")
	tab.Add(app.NewRight("text_content", "browse", ctx.T("Can browse text content records.")))
	tab.Add(app.NewRight("text_content", "insert", ctx.T("Can insert text content records.")))
	tab.Add(app.NewRight("text_content", "update", ctx.T("Can update text content records.")))
	tab.Add(app.NewRight("text_content", "set", ctx.T("Can set text content records.")))
	tab.Add(app.NewRight("text_content", "delete", ctx.T("Can delete text content records.")))
	rl.Add("text_content", tab)

	//manufacturer
	tab = app.NewRightTab("Manufacturer")
	tab.Add(app.NewRight("manufacturer", "browse", ctx.T("Can browse manufacturer records.")))
	tab.Add(app.NewRight("manufacturer", "insert", ctx.T("Can insert manufacturer records.")))
	tab.Add(app.NewRight("manufacturer", "update", ctx.T("can update manufacturer records.")))
	tab.Add(app.NewRight("manufacturer", "delete", ctx.T("Can delete manufacturer records.")))
	rl.Add("manufacturer", tab)

	//machine
	tab = app.NewRightTab("Machine")
	tab.Add(app.NewRight("machine", "browse", ctx.T("Can browse own machine records.")))
	tab.Add(app.NewRight("machine", "insert", ctx.T("Can insert machine records.")))
	tab.Add(app.NewRight("machine", "update", ctx.T("can update own machine records.")))
	tab.Add(app.NewRight("machine", "delete", ctx.T("Can delete own machine records.")))
	rl.Add("machine", tab)

	//copyright
	tab = app.NewRightTab("Privacy Policy")
	tab.Add(app.NewRight("privacy_policy", "browse", ctx.T("Can browse privacy policy text.")))
	rl.Add("privacy_policy", tab)

	app.UserRightMap = app.UserRightList.GetRightMap() //to set AppRightMap
}
