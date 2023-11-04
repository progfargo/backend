package text_content

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/text_content_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/text_content/text_content_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("text_content", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("textId", -1)
	ctx.ReadCargo()

	textId := ctx.Cargo.Int("textId")
	rec, err := text_content_lib.GetTextContentRec(textId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/text_content"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	name := ctx.Req.PostFormValue("name")
	exp := ctx.Req.PostFormValue("exp")
	body := ctx.Req.PostFormValue("body")

	if name == "" || exp == "" || body == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update textContent set
						name = ?,
						exp = ?,
						body = ?
					where
						textContentId = ?`

	res, err := tx.Exec(sqlStr, name, exp, body, textId)
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

	ctx.Msg.Success(ctx.T("Record has been changed."))
	ctx.Redirect(ctx.U("/text_content_update", "textId"))
}

func updateForm(ctx *context.Ctx, rec *text_content_lib.TextContentRec) {
	content.Include(ctx)

	ctx.Js.Add("/asset/tinymce/tinymce.min.js")
	ctx.Js.Add("/asset/tinymce/tinymce_func.js")
	ctx.Js.Add("/asset/js/page/text_content/text_content_update.js")

	var name, exp, body string
	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		exp = ctx.Req.PostFormValue("exp")
		body = ctx.Req.PostFormValue("body")
	} else {
		name = rec.Name
		exp = rec.Exp
		body = rec.Body
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Text Content"), ctx.T("Update Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/text_content")))
	buf.Add("</div>")
	buf.Add("</div>")

	textMenu := text_content_menu.New("textId")
	textMenu.Set(ctx, "text_content_update")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(textMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/text_content_update", "textId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Name:"))
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\">", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Explanation:"))
	buf.Add("<input type=\"text\" name=\"exp\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"2\">", util.ScrStr(exp))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Body:"))
	buf.Add("<textarea name=\"body\" id=\"body\" class=\"formControl\""+
		" tabindex=\"3\">%s</textarea>", body)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"4\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"5\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "text_content")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
