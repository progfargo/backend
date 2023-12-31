package user

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/page/role/role_lib"
	"backend/src/page/user/user_lib"
)

func UserRoleRevoke(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "role_revoke") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddInt("roleId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("rid", -1)
	ctx.Cargo.AddStr("stat", "default")
	ctx.ReadCargo()

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	userId := ctx.Cargo.Int("userId")
	roleId := ctx.Cargo.Int("roleId")

	if userId == -1 || roleId == -1 {
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
		return
	}

	userRec, err := user_lib.GetUserRec(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("User record could not be found."))
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	roleRec, err := role_lib.GetRoleRec(roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Role record could not be found."))
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	if userRec.Login == "superuser" {
		ctx.Msg.Warning(ctx.T("'superuser' account can not be updated."))
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	if userRec.Login == "testuser" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("'testuser' account can not be updated."))
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	if roleRec.Name == "admin" {
		n := user_lib.CountAdmin()
		if n < 2 {
			tx.Rollback()
			ctx.Msg.Warning(ctx.T("Last admin role can not be revoked."))
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}
	}

	sqlStr := `delete from
					userRole
				where
					userId = ? and
					roleId = ?`

	res, err := tx.Exec(sqlStr, userId, roleId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
		return
	}

	tx.Commit()
	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/user_role", "userId", "key", "pn", "rid", "stat"))
}
