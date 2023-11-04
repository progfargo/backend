package machine_pdf

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_pdf/machine_pdf_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddInt("pdfId", -1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)

	ctx.ReadCargo()

	machineId := ctx.Cargo.Int("machineId")
	machineRec, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine record could not be found."))
			ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		ctx.Msg.Error(ctx.T("You don't have right to access this record."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	pdfId := ctx.Cargo.Int("pdfId")
	machinePdfRec, err := machine_pdf_lib.GetMachinePdfRec(pdfId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteConfirm(ctx, machineRec, machinePdfRec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					machinePdf
				where
					machinePdfId = ?`

	res, err := tx.Exec(sqlStr, pdfId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))

}

func deleteConfirm(ctx *context.Ctx, machineRec *machine_lib.MachineRec, machinePdfRec *machine_pdf_lib.MachinePdfRec) {
	content.Include(ctx)
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("Delete PDF"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("File name:"), util.ScrStr(machinePdfRec.PdfName))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Mime:"), util.ScrStr(machinePdfRec.PdfMime))
	buf.Add("<tr><th>%s</th><td>%s bytes</td></tr>", ctx.T("Size:"), util.FormatInt(machinePdfRec.PdfSize))
	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h4>%s</h4>", ctx.T("Please confirm:"))
	buf.Add("<p>%s</p>", ctx.T("Do you realy want to delete this record?"))
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.Cargo.SetStr("confirm", "yes")
	urlStr := ctx.U("/machine_pdf_delete", "machineId", "pdfId", "confirm", "key", "pn", "stat", "manId")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
