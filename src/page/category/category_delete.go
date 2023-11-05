package category

import (
	"net/http"

	"backend/src/app"
	"backend/src/lib/context"

	"github.com/go-sql-driver/mysql"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("category", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("categoryId", -1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/category"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	categoryId := ctx.Cargo.Int("categoryId")

	sqlStr := `delete from
					category
				where
					categoryId = ?`

	res, err := tx.Exec(sqlStr, categoryId)
	if err != nil {
		tx.Rollback()

		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1451 {
				ctx.Msg.Warning(ctx.T("Could not delete parent record."))
				ctx.Redirect(ctx.U("/category"))
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/category"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/category"))
}
