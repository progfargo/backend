package office

import (
	"net/http"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/office/office_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("office", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("productionId", -1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "office")

	tmenu := top_menu.New()
	tmenu.Set(ctx, "root")

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	insertRight := ctx.IsRight("office", "insert")
	updateRight := ctx.IsRight("office", "update")
	deleteRight := ctx.IsRight("office", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Offices")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")
		buf.Add(content.NewButton(ctx, ctx.U("/office_insert")))
		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s</th>", ctx.T("Name"))
	buf.Add("<th>%s</th>", ctx.T("City"))
	buf.Add("<th>%s</th>", ctx.T("Address"))
	buf.Add("<th>%s</th>", ctx.T("Telephone"))
	buf.Add("<th>%s</th>", ctx.T("Email"))
	buf.Add("<th title=\"%s\">%s</th>", ctx.T("Is main office?"), ctx.T("Is Main"))

	if updateRight || deleteRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	officeList := office_lib.GetOfficeList()
	if len(officeList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	var urlStr string
	var name, city, address, telephone, email string
	for _, row := range officeList {
		ctx.Cargo.SetInt("productionId", row.OfficeId)

		name = util.ScrStr(row.Name)
		city = util.ScrStr(row.City)
		address = util.ScrStr(row.Address)
		telephone = util.ScrStr(row.Telephone)
		email = util.ScrStr(row.Email)

		buf.Add("<tr>")

		urlStr = ctx.U("/office_display", "productionId")
		buf.Add("<td><a href=\"%s\">%s</a></td>", urlStr, name)

		buf.Add("<td>%s</td>", city)
		buf.Add("<td>%s</td>", strings.Replace(address, "\n", "<br>", -1))
		buf.Add("<td>%s</td>", telephone)
		buf.Add("<td>%s</td>", email)
		buf.Add("<td>%s</td>", office_lib.IsMainOfficeToLabel(ctx, row.IsMainOffice))

		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		urlStr = ctx.U("/office_display", "productionId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Display record."), ctx.T("Display"))

		if updateRight || deleteRight {
			if updateRight {
				urlStr = ctx.U("/office_make_main", "productionId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Make this the main office."), ctx.T("Make Main"))

				urlStr = ctx.U("/office_update", "productionId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit record."), ctx.T("Edit"))
			}

			if deleteRight {
				urlStr = ctx.U("/office_delete", "productionId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))
			}
		}

		buf.Add("</div>")
		buf.Add("</td>")
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
