package machine_image

import (
	"database/sql"
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_image/machine_image_lib"
)

func MakeHeaderImage(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)
	ctx.ReadCargo()

	machineId := ctx.Cargo.Int("machineId")
	_, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		http.Error(rw, ctx.T("You don't have right to access this record."), 401)
		return
	}

	imgId := ctx.Cargo.Int("imgId")
	imgRec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine image record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if imgRec.IsHeader == "yes" {
		ctx.Msg.Warning(ctx.T("Record is already a header image."))
		ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := fmt.Sprintf(`update machineImage set
								isHeader = 'no'
							where
								machineId = ?`)

	res, err := tx.Exec(sqlStr, machineId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sqlStr = fmt.Sprintf(`update machineImage set
								isHeader = 'yes'
							where
								machineImageId = ?`)

	res, err = tx.Exec(sqlStr, imgId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
}
