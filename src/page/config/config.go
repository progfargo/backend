package config

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/config/config_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("config", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("configId", -1)
	ctx.Cargo.AddInt("gid", -1)
	ctx.ReadCargo()

	content.Include(ctx)
	content.Default(ctx)

	browseMid(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "config")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	configList := config_lib.GetConfigList()
	if len(configList) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	insertRight := ctx.IsRight("config", "insert")
	updateRight := ctx.IsRight("config", "update")
	deleteRight := ctx.IsRight("config", "delete")
	setRight := ctx.IsRight("config", "set")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Configuration")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		buf.Add(content.NewButton(ctx, ctx.U("/config_insert")))

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s</th>", ctx.T("Title"))
	buf.Add("<th>%s</th>", ctx.T("Value"))

	if updateRight || deleteRight {
		buf.Add("<th>%s</th>", ctx.T("Enum"))
	}

	buf.Add("<th>%s</th>", ctx.T("Group Name"))
	buf.Add("<th>%s</th>", ctx.T("Explanation"))

	if setRight || updateRight || deleteRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	var name, title, value, groupName, exp string
	for _, row := range configList {
		ctx.Cargo.SetInt("configId", row.ConfigId)

		name = util.ScrStr(row.Name)
		title = util.ScrStr(row.Title)
		value = util.ScrStr(row.Value)
		exp = util.ScrStr(row.Exp)

		buf.Add("<tr>")
		buf.Add("<td title=\"%s\">%s</td>", name, title)
		buf.Add("<td>%s</td>", value)

		if updateRight || deleteRight {
			buf.Add("<td>%d</td>", row.Enum)
		}

		buf.Add("<td>%s</td>", groupName)
		buf.Add("<td>%s</td>", exp)

		if setRight || updateRight || deleteRight {
			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if setRight {
				urlStr := ctx.U("/config_set", "configId", "gid")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Set value."), ctx.T("Set Value"))
			}

			if updateRight {
				urlStr := ctx.U("/config_update", "configId", "gid")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit record."), ctx.T("Edit"))
			}

			if deleteRight {
				urlStr := ctx.U("/config_delete", "configId", "gid")
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
