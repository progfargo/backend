package user

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/combo"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/role/role_lib"
	"backend/src/page/user/user_lib"
)

func UserRole(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("user", "role_browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("userId", -1)
	ctx.Cargo.AddInt("roleId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("rid", -1)
	ctx.Cargo.AddStr("stat", "default")
	ctx.ReadCargo()

	userId := ctx.Cargo.Int("userId")
	rec, err := user_lib.GetUserRec(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		panic(err)
	}

	if rec.Login == "superuser" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("'superuser' account can not be updated."))
		ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
		return
	}

	displayUser(ctx, rec)
}

func displayUser(ctx *context.Ctx, userRec *user_lib.UserRec) {

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")

	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Users"), ctx.T("User Roles")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/user", "key", "pn", "rid", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	//user info
	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("User Information:"))
	buf.Add("<tbody>")
	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("User Name:"), util.ScrStr(userRec.Name))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Login:"), util.ScrStr(userRec.Login))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Email:"), util.ScrStr(userRec.Email))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	//grant role
	if ctx.IsRight("user", "role_grant") {
		sqlStr := `select
						roleId,
						name
					from
						role`

		roleCombo := combo.NewCombo(sqlStr, ctx.T("Select Role"))
		roleCombo.Set()
		if roleCombo.IsEmpty() {
			ctx.Msg.Warning(ctx.T("Role list is empty. You should enter at least one role first."))
			ctx.Redirect(ctx.U("/user", "key", "pn", "rid", "stat"))
			return
		}

		buf.Add("<table>")
		buf.Add("<caption>%s</caption>", ctx.T("Grant Role:"))
		buf.Add("<tbody>")
		buf.Add("<tr>")
		buf.Add("<th class=\"fixedMiddle\">%s</th>", ctx.T("Role Name:"))
		buf.Add("<td>")

		urlStr := ctx.U("/user_role_grant", "userId", "key", "pn", "rid", "stat")
		buf.Add("<form class=\"formFlex\" action=\"%s\" method=\"post\">", urlStr)

		buf.Add("<div class=\"formGroup\">")
		buf.Add("<select name=\"roleId\" class=\"formControl\">")

		buf.Add(roleCombo.Format(-1))

		buf.Add("</select>")
		buf.Add("</div>")

		buf.Add("<div class=\"formGroup\">")
		buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\">%s</button>", ctx.T("Grant"))
		buf.Add("</div>")

		buf.Add("</form>")

		buf.Add("</td>")
		buf.Add("</tr>")

		buf.Add("</tbody>")
		buf.Add("</table>")
	}

	buf.Add("</div>")

	revokeRight := ctx.IsRight("user", "role_revoke")

	buf.Add("<div class=\"col\">")

	//user role list
	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("User Role List:"))
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">%s</th>", ctx.T("Role Name"))
	buf.Add("<th>%s</th>", ctx.T("Explanation"))

	if revokeRight {
		buf.AddLater("<th class=\"right fixedZero\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	userRoleList, err := role_lib.GetRoleByUser(userRec.UserId)
	if err != nil {
		panic(err)
	}

	if len(userRoleList) != 0 {
		buf.Forge()

		for _, row := range userRoleList {
			buf.Add("<tr>")
			buf.Add("<td>%s</td><td>%s</td>", row.Name, row.Exp)

			if revokeRight {
				buf.Add("<td class=\"right\">")
				buf.Add("<div class=\"buttonGroupFixed\">")

				buf.Add("<div class=\"revoke\">")
				buf.Add("<a href=\"#\" class=\"button buttonError buttonXs revokeButton\">%s</a>", ctx.T("Revoke"))
				buf.Add("</div>")

				buf.Add("<div class=\"revokeConfirm\">")
				buf.Add(ctx.T("Do you realy want to revoke this role?"))

				ctx.Cargo.SetInt("roleId", row.RoleId)
				urlStr := ctx.U("/user_role_revoke", "userId", "roleId", "key", "pn", "stat")
				buf.Add("<a href=\"%s\" class=\"button buttonSuccess buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Revoke this role."), ctx.T("Yes"))

				buf.Add("<a class=\"button buttonDefault buttonXs cancelButton\">%s</a>", ctx.T("Cancel"))

				buf.Add("</div>") //revokeConfirm
				buf.Add("</div>") //buttonGroupFixed
				buf.Add("</td>")
			}

			buf.Add("</tr>")
		}
	} else {
		buf.Add("<tr><td colspan=\"2\"><span class=\"label label-danger label-sd\">%s</span></td></tr>",
			ctx.T("Empty"))
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/user_role.css")
	ctx.Js.Add("/asset/js/page/user/user_role.js")

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "user")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "userRolePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}
