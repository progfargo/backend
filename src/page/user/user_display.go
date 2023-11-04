package user

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/user"
	"backend/src/lib/util"
	"backend/src/page/user/user_lib"
)

func Display(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("rid", -1)
	ctx.Cargo.AddStr("stat", "default")
	ctx.ReadCargo()

	userId := ctx.Cargo.Int("userId")
	rec, err := user_lib.GetUserRec(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning("Record could not be found.")
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	displayUserRec(ctx, rec)

	ctx.Render("default.html")
}

func displayUserRec(ctx *context.Ctx, rec *user_lib.UserRec) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Users"), ctx.T("Display Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("User Information:"))
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("User Name:"), util.ScrStr(rec.Name))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Login Name:"), util.ScrStr(rec.Login))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Email:"), util.ScrStr(rec.Email))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Status:"), user_lib.StatusToLabel(ctx, rec.Status))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("User Roles:"))
	buf.Add("<tbody>")

	libUserRec := user.New(ctx, rec.UserId)
	for k, _ := range libUserRec.RoleList {
		buf.Add("<tr>")
		buf.Add("<td>%s</td>", k)
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
