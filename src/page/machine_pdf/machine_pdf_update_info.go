package machine_pdf

import (
	"database/sql"
	"net/http"
	"strconv"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_pdf/machine_pdf_lib"
)

func UpdatePdfInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddInt("pdfId", -1)
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
	pdfRec, err := machine_pdf_lib.GetMachinePdfRec(pdfId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine pdf record could not be found."))
			ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updatePdfInfoForm(ctx, machineRec, pdfRec)
		return
	}

	title := ctx.Req.PostFormValue("title")
	enum, err := strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
	if err != nil {
		enum = 0
	}

	if title == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updatePdfInfoForm(ctx, machineRec, pdfRec)
		return
	}

	sqlStr := `update machinePdf set
					title = ?,
					enum = ?
				where
					machinePdfId = ?`

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	res, err := tx.Exec(sqlStr, title, enum, pdfId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updatePdfInfoForm(ctx, machineRec, pdfRec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
}

func updatePdfInfoForm(ctx *context.Ctx, machineRec *machine_lib.MachineRec, pdfRec *machine_pdf_lib.MachinePdfRec) {
	content.Include(ctx)

	var title string
	var enum int64
	var err error
	if ctx.Req.Method == "POST" {
		title = ctx.Req.PostFormValue("title")

		enum, err = strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
		if err != nil {
			enum = 0
		}
	} else {
		title = pdfRec.Title
		enum = pdfRec.Enum
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("Update Pdf Info"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_pdf_update_info", "machineId", "pdfId", "key", "stat", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Title:"))
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"200\" tabindex=\"1\" autofocus>", util.ScrStr(title))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label title=\"Enumeration\">%s</label>", ctx.T("Enum:"))
	buf.Add("<input type=\"text\" name=\"enum\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"10\" tabindex=\"2\">", enum)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"3\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"4\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)
	content.Include(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
