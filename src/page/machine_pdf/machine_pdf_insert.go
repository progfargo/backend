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

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("machine", "insert")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
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
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		insertForm(ctx, machineRec)
		return
	}

	title := ctx.Req.PostFormValue("title")

	if title == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx, machineRec)
		return
	}

	err = ctx.Req.ParseMultipartForm(1024 * 1024 * 4)
	if err != nil {
		panic(err)
	}

	mpForm := ctx.Req.MultipartForm
	pdf := mpForm.File["pdf"]
	if len(pdf) == 0 {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx, machineRec)
		return
	}

	for i := 0; i < len(pdf); i++ {
		pdfFile := pdf[i]
		pdfName := pdfFile.Filename
		pdfMime := pdfFile.Header["Content-Type"][0]

		ext := filepath.Ext(pdfName)
		isPdf, err := regexp.MatchString("(?i)^\\.pdf$", ext)
		if !isPdf {
			ctx.Msg.Warning(ctx.T("You can only send pdf files."))
			insertForm(ctx, machineRec)
			return
		}

		inFile, err := pdfFile.Open()
		if err != nil {
			ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
			insertForm(ctx, machineRec)
			return
		}

		defer inFile.Close()

		pdfData := bytes.NewBuffer(nil)
		io.Copy(pdfData, inFile)

		pdfSize := pdfData.Len()

		if pdfSize == 0 {
			ctx.Msg.Warning(ctx.T("You have sent a file with zero size."))
			insertForm(ctx, machineRec)
			return
		}

		tx, err := app.Db.Begin()
		if err != nil {
			panic(err)
		}

		sqlStr := `insert into
						machinePdf(machinePdfId, machineId, title, pdfName, pdfMime, pdfSize, pdfData)
						values(null, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(sqlStr, machineId, title, pdfName, pdfMime, pdfSize, pdfData.String())
		if err != nil {
			tx.Rollback()
			if err, ok := err.(*mysql.MySQLError); ok {
				if err.Number == 1062 {
					ctx.Msg.Warning(ctx.T("Duplicate record."))
					insertForm(ctx, machineRec)
					return
				} else if err.Number == 1452 {
					ctx.Msg.Warning(ctx.T("Could not find parent record."))
					insertForm(ctx, machineRec)
					return
				}
			}

			panic(err)
		}

		tx.Commit()
	}

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId"))
}

func insertForm(ctx *context.Ctx, machineRec *machine_lib.MachineRec) {
	var title string
	if ctx.Req.Method == "POST" {
		title = ctx.Req.PostFormValue("title")
	} else {
		title = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("New PDF"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_pdf", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_pdf_insert", "machineId", "key", "stat", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Title:"))
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"200\" tabindex=\"1\" autofocus>", util.ScrStr(title))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("PDF:"))
	buf.Add("<input type=\"file\" name=\"pdf\" class=\"formControl\"" +
		" tabindex=\"2\" multiple=\"multiple\" autofocus>")
	buf.Add("<span class=\"helpBlock\">%s</span>", ctx.T("File type must be pdf."))
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
