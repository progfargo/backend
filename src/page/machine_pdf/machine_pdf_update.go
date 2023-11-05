package machine_pdf

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"path/filepath"
	"regexp"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_pdf/machine_pdf_lib"

	"github.com/go-sql-driver/mysql"
)

func UpdatePdf(rw http.ResponseWriter, req *http.Request) {
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
	imgRec, err := machine_pdf_lib.GetMachinePdfRec(pdfId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine image record could not be found."))
			ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
		return
	}

	err = ctx.Req.ParseMultipartForm(1024 * 1024 * 4)
	if err != nil {
		panic(err)
	}

	pdfForm := ctx.Req.MultipartForm
	pdf := pdfForm.File["pdf"]
	if len(pdf) == 0 {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	pdfFile := pdf[0]
	pdfName := pdfFile.Filename
	pdfMime := pdfFile.Header["Content-Type"][0]

	ext := filepath.Ext(pdfName)
	isPdf, err := regexp.MatchString("(?i)^\\.pdf$", ext)
	if !isPdf {
		ctx.Msg.Warning(ctx.T("You can only send pdf files."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	inFile, err := pdfFile.Open()
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	defer inFile.Close()

	pdfData := bytes.NewBuffer(nil)
	io.Copy(pdfData, inFile)

	pdfSize := pdfData.Len()

	tx, err = app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update machinePdf set
					pdfName = ?,
					pdfMime = ?,
					pdfSize = ?,
					pdfData = ?
				where
					machinePdfId = ?`

	res, err := tx.Exec(sqlStr, pdfName, pdfMime, pdfSize, pdfData.String(), pdfId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				updateImageForm(ctx, machineRec, imgRec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
}

func updateImageForm(ctx *context.Ctx, machineRec *machine_lib.MachineRec, imgRec *machine_pdf_lib.MachinePdfRec) {
	content.Include(ctx)
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("Update PDF"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_pdf_update", "machineId", "pdfId", "key", "stat", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("PDF:"))
	buf.Add("<input type=\"file\" name=\"pdf\" class=\"formControl\"" +
		" tabindex=\"1\" multiple=\"multiple\" autofocus>")
	buf.Add("<span class=\"helpBlock\">%s</span>", ctx.T("File type must be pdf."))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"2\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"3\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

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
