package faq

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("faq", "insert") {
		app.BadRequest()
	}

		ctx.Cargo.AddStr("key", "")
		ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	question := ctx.Req.PostFormValue("question")
	summary := ctx.Req.PostFormValue("summary")
	answer := ctx.Req.PostFormValue("answer")

	if question == "" || summary == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					faq(faqId, question, summary, answer)
					values(null, ?, ?, ?)`

	_, err = tx.Exec(sqlStr, question, summary, answer)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				insertForm(ctx)
				return
			}
		}

		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/faq", "key", "pn"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)

	ctx.Js.Add("/asset/tinymce/tinymce.min.js")
	ctx.Js.Add("/asset/tinymce/tinymce_func.js")

	ctx.Js.Add("/asset/js/page/faq/faq_insert.js")

	var question, summary, answer string
	if ctx.Req.Method == "POST" {
		question = ctx.Req.PostFormValue("question")
		summary = ctx.Req.PostFormValue("summary")
		answer = ctx.Req.PostFormValue("answer")
	} else {
		question = ""
		summary = ""
		answer = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Faq"), ctx.T("New Record")))
	buf.Add("</div>")

	
	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/faq", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/faq_insert", "key", "pn")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Question:"))
	buf.Add("<input type=\"text\" name=\"question\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"1\" autofocus>", util.ScrStr(question))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Summary:"))
	buf.Add("<textarea name=\"summary\" id=\"summary\" class=\"formControl\""+
		" tabindex=\"2\">%s</textarea>", summary)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("Answer:"))
	buf.Add("<textarea name=\"answer\" id=\"answer\" class=\"formControl\""+
		" tabindex=\"3\">%s</textarea>", answer)
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
	lmenu.Set(ctx, "faq")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}