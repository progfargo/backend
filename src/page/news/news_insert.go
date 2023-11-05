package news

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

	if !ctx.IsRight("news", "insert") {
		app.BadRequest()
	}

	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("stat", "")
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/news", "key", "stat", "pn"))
		return
	}

	recordDateStr := ctx.Req.PostFormValue("recordDate")
	recordDate, err := util.DateStrToInt64(recordDateStr)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Could not convert date string to int64."))
		insertForm(ctx)
		return
	}

	title := ctx.Req.PostFormValue("title")
	status := ctx.Req.PostFormValue("status")
	summary := ctx.Req.PostFormValue("summary")
	body := ctx.Req.PostFormValue("body")

	if recordDateStr == "" || status == "" || title == "" || summary == "" || body == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					news(newsId, recordDate, status, title, summary, body)
					values(null, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(sqlStr, recordDate, status, title, summary, body)
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
	ctx.Redirect(ctx.U("/news", "key", "stat", "pn"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)

	ctx.Js.Add("/asset/tinymce/tinymce.min.js")
	ctx.Js.Add("/asset/tinymce/tinymce_func.js")

	ctx.Css.Add("/asset/datetimepicker/jquery.datetimepicker.css")
	ctx.Js.Add("/asset/datetimepicker/jquery.datetimepicker.full.js")

	ctx.Js.Add("/asset/js/page/news/news_insert.js")

	var recordDate, title, summary, body, status string
	if ctx.Req.Method == "POST" {
		recordDate = ctx.Req.PostFormValue("recordDate")
		title = ctx.Req.PostFormValue("title")
		summary = ctx.Req.PostFormValue("summary")
		body = ctx.Req.PostFormValue("body")
		status = ctx.Req.PostFormValue("status")
	} else {
		recordDate = util.Int64ToDateStr(util.Now())
		title = ""
		summary = ""
		body = ""
		status = "draft"
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("News"), ctx.T("New Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/news", "key", "stat", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/news_insert", "key", "pn", "stat")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Date:"))
	buf.Add("<input type=\"text\" name=\"recordDate\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"10\" tabindex=\"1\">", util.ScrStr(recordDate))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Title:"))
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"2\">", util.ScrStr(title))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Summary:"))
	buf.Add("<textarea name=\"summary\" id=\"summary\" class=\"formControl\""+
		" tabindex=\"3\">%s</textarea>", summary)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Body:"))
	buf.Add("<textarea name=\"body\" id=\"body\" class=\"formControl\""+
		" tabindex=\"4\">%s</textarea>", body)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Status:"))

	checkedStr := ""
	if status == "draft" {
		checkedStr = " checked"
	}

	buf.Add("<span>")
	buf.Add("<input type=\"radio\" name=\"status\" value=\"draft\"%s><span class=\"radioLabel\">%s</span>", checkedStr, ctx.T("Draft"))
	buf.Add("</span>")

	checkedStr = ""
	if status == "published" {
		checkedStr = " checked"
	}

	buf.Add("<span>")
	buf.Add("<input type=\"radio\" name=\"status\" value=\"published\"%s><span class=\"radioLabel\">%s</span>", checkedStr, ctx.T("Published"))
	buf.Add("</span>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"5\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"6\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "news")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
