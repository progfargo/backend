package office

import (
	"database/sql"
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/page/office/office_lib"
)

func MakeMainOffice(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("office", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("productionId", -1)
	ctx.ReadCargo()

	productionId := ctx.Cargo.Int("productionId")
	rec, err := office_lib.GetOfficeRec(productionId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/office"))
			return
		}

		panic(err)
	}

	if rec.IsMainOffice == "yes" {
		ctx.Msg.Warning(ctx.T("Record is already a main office."))
		ctx.Redirect(ctx.U("/office"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := fmt.Sprintf(`update office set isMainOffice = 'no'`)

	res, err := tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStr = fmt.Sprintf(`update office set
								isMainOffice = 'yes'
							where
								officeId = ?`)

	res, err = tx.Exec(sqlStr, productionId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		ctx.Redirect(ctx.U("/office"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/office"))
}
