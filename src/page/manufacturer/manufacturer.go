package manufacturer

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/manufacturer/manufacturer_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("manufacturer", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("manufacturerId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, "/manufacturer")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "manufacturer")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")

	totalRows := manufacturer_lib.CountManufacturer(key)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("manufacturer", "insert")
	updateRight := ctx.IsRight("manufacturer", "update")
	deleteRight := ctx.IsRight("manufacturer", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Manufacturers")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		if insertRight {
			buf.Add(content.NewButton(ctx, ctx.U("/manufacturer_insert", "key", "pn")))
		}

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s</th>", ctx.T("Name"))

	if updateRight || deleteRight {
		buf.Add("<th class=\"right fixedSmall\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	if totalRows > 0 {
		manufacturerList := manufacturer_lib.GetManufacturerPage(ctx, key, pageNo)

		var name string
		for _, row := range manufacturerList {
			ctx.Cargo.SetInt("manufacturerId", row.ManufacturerId)

			name = util.ScrStr(row.Name)

			if key != "" {
				name = content.Find(name, key)
			}

			buf.Add("<tr>")
			buf.Add("<td>%s</td>", name)

			if updateRight || deleteRight {
				buf.Add("<td class=\"right\">")
				buf.Add("<div class=\"buttonGroupFixed\">")

				if updateRight {
					urlStr := ctx.U("/manufacturer_update", "manufacturerId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Edit record."), ctx.T("Edit"))
				}

				if deleteRight {
					urlStr := ctx.U("/manufacturer_delete", "manufacturerId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Delete record."), ctx.T("Delete"))
				}

				buf.Add("</div>")
				buf.Add("</td>")
			}

			buf.Add("</tr>")
		}
	}

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/manufacturer", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
