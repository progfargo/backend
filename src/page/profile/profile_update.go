package profile

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/user/user_lib"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("profile", "update") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	rec, err := user_lib.GetUserRec(ctx.User.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			panic(ctx.T("User record could not be found."))
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	name := ctx.Req.PostFormValue("name")

	if name == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update user set
					name = ?
				where
					userId = ?`

	res, err := tx.Exec(sqlStr, name)
	if err != nil {
		tx.Rollback()
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
	ctx.Redirect(ctx.U("/profile"))
}

func updateForm(ctx *context.Ctx, rec *user_lib.UserRec) {

	var name string

	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
	} else {
		name = rec.Name
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("User Profile"), ctx.T("Update Record"), util.ScrStr(rec.Name)))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/profile")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/profile_update")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("Login Name:"))
	buf.Add("<input type=\"text\" name=\"login\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" disabled>", util.ScrStr(rec.Login))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("E-Mail:"))
	buf.Add("<input type=\"text\" name=\"email\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" disabled>", util.ScrStr(rec.Email))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("User Name:"))
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"8\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"9\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "profile")

	ctx.Render("default.html")
}
