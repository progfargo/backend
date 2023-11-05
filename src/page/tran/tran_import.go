package tran

import (
	"bytes"
	"encoding/json"
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
	"backend/src/page/tran/tran_lib"
)

func Import(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if (!ctx.IsRight("tran", "export") || app.Ini.AppType != "prod") && !ctx.IsSuperuser() {
		app.BadRequest()
	}

	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		importForm(ctx)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/tran", "key", "pn"))
		return
	}

	err := ctx.Req.ParseMultipartForm(1024 * 1024 * 2)
	if err != nil {
		panic(err)
	}

	mpForm := ctx.Req.MultipartForm
	file := mpForm.File["tranFile"]
	if len(file) == 0 {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		importForm(ctx)
		return
	}

	tranFile := file[0]
	tranFileName := tranFile.Filename
	tranFileMime := tranFile.Header["Content-Type"][0]

	if tranFileMime != "application/json" {
		ctx.Msg.Warning(ctx.T("You can only send json images."))
		importForm(ctx)
		return
	}

	ext := filepath.Ext(tranFileName)
	isJson, err := regexp.MatchString("(?i)^\\.json$", ext)
	if !isJson {
		ctx.Msg.Warning(ctx.T("You can only send json images."))
		importForm(ctx)
		return
	}

	inFile, err := tranFile.Open()
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
		importForm(ctx)
		return
	}

	defer inFile.Close()

	jsonData := bytes.NewBuffer(nil)
	io.Copy(jsonData, inFile)
	jsonDataSize := jsonData.Len()

	if jsonDataSize == 0 {
		ctx.Msg.Warning(ctx.T("You have sent a file with zero size."))
		importForm(ctx)
		return
	}

	var tranList []tran_lib.TranRec

	err = json.Unmarshal(jsonData.Bytes(), &tranList)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not decode json file."))
		importForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := "delete from tran"
	_, err = tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	buf := util.NewBuf()
	buf.Add("insert into tran(tranId, en, tr) values")

	first := true
	var comma string
	for _, v := range tranList {
		if !first {
			comma = ","
		}

		buf.Add("%s(null, '%s', '%s')", comma, util.DbStr(v.En), util.DbStr(v.Tr))
		first = false
	}

	_, err = tx.Exec(*buf.String())
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	app.ReadTran()
	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/tran", "key", "pn"))
}

func importForm(ctx *context.Ctx) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Translation Table"), ctx.T("Import")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/tran", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<div class=\"callout calloutInfo\">")
	buf.Add("<h4>%s</h4>", ctx.T("Please be informed:"))
	buf.Add("<p>%s</p>", ctx.T("This operation overwrites translation table with the content of the file you send."))
	buf.Add("<p>%s</p>", ctx.T("Current translation table will be lost."))
	buf.Add("</div>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/tran_import", "key", "pn")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Translation File:"))
	buf.Add("<input type=\"file\" name=\"tranFile\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")
	buf.Add("<span class=\"helpBlock\">%s</span>", ctx.T("JSON file created with export function."))
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
	lmenu.Set(ctx, "tran")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
