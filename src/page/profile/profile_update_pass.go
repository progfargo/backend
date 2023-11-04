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

func UpdatePass(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("profile", "update_password") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	rec, err := user_lib.GetUserRec(ctx.User.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			panic("User record could not be found.")
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updatePassForm(ctx, rec)
		return
	}

	curPassword := ctx.Req.PostFormValue("curPassword")
	newPassword := ctx.Req.PostFormValue("newPassword")
	reNewPassword := ctx.Req.PostFormValue("reNewPassword")

	if curPassword == "" || newPassword == "" || reNewPassword == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updatePassForm(ctx, rec)
		return
	}

	if util.PasswordHash(curPassword) != rec.Password {
		ctx.Msg.Warning(ctx.T("You have entered wrong current password. Please try again."))
		updatePassForm(ctx, rec)
		return
	}

	if err := util.IsValidPassword(ctx, newPassword); err != nil {
		ctx.Msg.Warning(err.Error())
		updatePassForm(ctx, rec)
		return
	}

	if newPassword == curPassword {
		ctx.Msg.Warning(ctx.T("You have entered your old password as new password."))
		updatePassForm(ctx, rec)
		return
	}

	if newPassword != reNewPassword {
		ctx.Msg.Warning(ctx.T("New password and retyped new password mismatch."))
		updatePassForm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	newPassword = util.PasswordHash(newPassword)
	sqlStr := `update user set
					password = ?
				where
					userId = ?`

	_, err = tx.Exec(sqlStr, newPassword, ctx.User.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been changed."))
	ctx.Redirect(ctx.U("/profile"))
}

func updatePassForm(ctx *context.Ctx, rec *user_lib.UserRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("User Profile"), ctx.T("Update Password")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/profile")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/profile_update_pass")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Current Password:"))
	buf.Add("<input type=\"password\" name=\"curPassword\"" +
		" class=\"formControl\" value=\"\" maxlength=\"30\" tabindex=\"1\" autofocus>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("New Password:"))
	buf.Add("<input type=\"password\" name=\"newPassword\"" +
		" class=\"formControl\" value=\"\" maxlength=\"30\" tabindex=\"2\">")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Retype New Password:"))
	buf.Add("<input type=\"password\" name=\"reNewPassword\"" +
		" class=\"formControl\" value=\"\" maxlength=\"30\" tabindex=\"3\">")
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
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "profile")

	ctx.Render("default.html")
}
