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

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("profile", "browse") {
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

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "profile")

	displayProfile(ctx, rec)

	ctx.Render("default.html")
}

func displayProfile(ctx *context.Ctx, rec *user_lib.UserRec) {

	ctx.Js.Add("/asset/js/page/profile/profile.js")
	buf := util.NewBuf()

	updateRight := ctx.IsRight("profile", "update")
	updatePasswordRight := ctx.IsRight("profile", "update_password")

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("User Profile")))
	buf.Add("</div>")

	if updateRight || updatePasswordRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		if updateRight {
			urlStr := ctx.U("/profile_update")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Edit record."), ctx.T("Edit"))
		}

		if updatePasswordRight {
			urlStr := ctx.U("/profile_update_pass")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Change User password."), ctx.T("Change Password"))
		}

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("User Name:"), util.ScrStr(rec.Name))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Login Name:"), util.ScrStr(rec.Login))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Email:"), util.ScrStr(rec.Email))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
