package tran

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/tran/tran_lib"
)

func Export(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if (!ctx.IsRight("tran", "export") || app.Ini.AppType != "prod") && !ctx.IsSuperuser() {
		app.BadRequest()
	}

	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	tranList := tran_lib.GetTranList()
	if len(tranList) == 0 {
		ctx.Msg.Warning(ctx.T("Translation list is empty."))
		ctx.Redirect(ctx.U("/tran", "key", "pn"))
	}

	jsonStr, err := json.Marshal(tranList)
	if err != nil {
		panic(err)
	}

	ctx.Rw.Header().Set("Content-Type", "application/json")
	ctx.Rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonStr)))

	fileName := fmt.Sprintf("translation_%s.json", util.Int64ToDateStr(util.Now()))
	ctx.Rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	if _, err := ctx.Rw.Write(jsonStr); err != nil {
		http.Error(ctx.Rw, ctx.T("Could not write json data to response writer."), http.StatusNotFound)
	}
}
