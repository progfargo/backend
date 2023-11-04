package role

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/role/role_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("role", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("roleId", -1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "role")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	ctx.Js.Add("/asset/js/page/role/role_right.js")

	insertRight := ctx.IsRight("role", "insert")
	updateRight := ctx.IsRight("role", "update")
	deleteRight := ctx.IsRight("role", "delete")
	roleRightRight := ctx.IsRight("role", "role_right")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Roles")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")

	if insertRight {
		buf.Add(content.NewButton(ctx, ctx.U("/role_insert")))
	}

	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">%s</th>", ctx.T("Role Name"))
	buf.Add("<th>%s</th>", ctx.T("Explanation"))

	if updateRight || deleteRight || roleRightRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	roleList := role_lib.GetRoleList()
	if len(roleList) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	var name, exp string
	for _, row := range roleList {
		ctx.Cargo.SetInt("roleId", row.RoleId)

		name = util.ScrStr(row.Name)
		exp = util.ScrStr(row.Exp)

		buf.Add("<tr>")
		buf.Add("<td>%s</td>", name)
		buf.Add("<td>%s</td>", exp)

		if updateRight || deleteRight || roleRightRight {
			buf.Add("<td class=\"right fixedSmall\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if updateRight {
				urlStr := ctx.U("/role_update", "roleId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit record."), ctx.T("Edit"))
			}

			if updateRight {
				urlStr := ctx.U("/role_right", "roleId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit role rights."), ctx.T("Role Rights"))
			}

			if deleteRight {
				urlStr := ctx.U("/role_delete", "roleId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))
			}

			buf.Add("</div>")
			buf.Add("</td>")
		}

		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
