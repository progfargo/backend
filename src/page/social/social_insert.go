package social

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

	if !ctx.IsRight("social", "insert") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	link := ctx.Req.PostFormValue("link")
	title := ctx.Req.PostFormValue("title")
	target := ctx.Req.PostFormValue("target")
	icon := ctx.Req.PostFormValue("icon")

	if link == "" || title == "" || target == "" || icon == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					social(socialId, link, title, target, icon)
					values(null, ?, ?, ?, ?)`

	_, err = tx.Exec(sqlStr, link, title, target, icon)
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
	ctx.Redirect(ctx.U("/social"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)

	var link, title, target, icon string
	if ctx.Req.Method == "POST" {
		link = ctx.Req.PostFormValue("link")
		title = ctx.Req.PostFormValue("title")
		target = ctx.Req.PostFormValue("target")
		icon = ctx.Req.PostFormValue("icon")
	} else {
		link = ""
		title = ""
		target = ""
		icon = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Social Links"), ctx.T("New Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/social")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/social_insert")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Link:"))
	buf.Add("<input type=\"text\" name=\"link\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"1\" autofocus>", util.ScrStr(link))
	buf.Add("<span class=\"helpBlock\">sample: http://www.example.com</span>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Title:"))
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"2\">", util.ScrStr(title))
	buf.Add("</div>")

	var selfChecked, blankChecked string
	if target == "" || target == "self" {
		selfChecked = " checked"
	} else {
		blankChecked = " checked"
	}

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Target:"))

	buf.Add("<span>")
	buf.Add("<input type=\"radio\" name=\"target\""+
		" value=\"_blank\" tabindex=\"3\"%s><span class=\"radioLabel\" title=\"%s\">blank</span>",
		blankChecked, ctx.T("Show in a new window."))
	buf.Add("</span>")

	buf.Add("<span>")
	buf.Add("<input type=\"radio\" name=\"target\""+
		" value=\"self\" tabindex=\"4\"%s><span class=\"radioLabel\" title=\"%s\">self</span>",
		selfChecked, ctx.T("Show in the same window."))
	buf.Add("</span>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Icon:"))
	buf.Add("<input type=\"text\" name=\"icon\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"50\" tabindex=\"5\">", icon)
	buf.Add("<span class=\"helpBlock\">sample: fab fa-facebook-square</span>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"6\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"7\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "social")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
