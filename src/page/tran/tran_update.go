package tran

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/tran/tran_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("tran", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("tranId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	tranId := ctx.Cargo.Int("tranId")
	rec, err := tran_lib.GetTranRec(tranId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/tran", "key", "pn"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	en := ctx.Req.PostFormValue("en")
	tr := ctx.Req.PostFormValue("tr")

	if en == "" || tr == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update tran set
					en = ?,
					tr = ?
				where
					tranId = ?`

	res, err := tx.Exec(sqlStr, en, tr, tranId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				updateForm(ctx, rec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateForm(ctx, rec)
		return
	}

	tx.Commit()

	app.ReadTran()
	ctx.Msg.Success(ctx.T("Record has been changed."))
	ctx.Redirect(ctx.U("/tran", "key", "pn"))
}

func updateForm(ctx *context.Ctx, rec *tran_lib.TranRec) {
	content.Include(ctx)

	var en, tr string
	if ctx.Req.Method == "POST" {
		en = ctx.Req.PostFormValue("en")
		tr = ctx.Req.PostFormValue("tr")
	} else {
		en = rec.En
		tr = rec.Tr
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Translation Table"), ctx.T("Update Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/tran", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/tran_update", "tranId", "key", "pn")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("English:"))
	buf.Add("<input type=\"text\" name=\"en\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"1\" autofocus>", util.ScrStr(en))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Turkish:"))
	buf.Add("<input type=\"text\" name=\"tr\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"2\">", util.ScrStr(tr))
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

	lmenu := left_menu.New()
	lmenu.Set(ctx, "tran")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
