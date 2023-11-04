package tran

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/tran/tran_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("tran", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("tranId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, "/tran")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "tran")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")

	totalRows := tran_lib.CountTran(key)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("tran", "insert")
	updateRight := ctx.IsRight("tran", "update")
	importRight := ctx.IsRight("tran", "import")
	exportRight := ctx.IsRight("tran", "export")
	deleteRight := ctx.IsRight("tran", "delete")
	synchronizeRight := ctx.IsRight("tran", "synchronize") && app.Ini.IsLocal

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Translation Table")))
	buf.Add("</div>")

	if insertRight || synchronizeRight || exportRight || importRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		var urlStr string
		if insertRight {
			buf.Add(content.NewButton(ctx, ctx.U("/tran_insert", "key", "pn")))
		}

		if synchronizeRight {
			urlStr = ctx.U("/tran_sync", "key", "pn")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Synchronize translation units with source code."), ctx.T("Synchronize"))
		}

		if (importRight && app.Ini.AppType == "prod") || ctx.IsSuperuser() {
			urlStr = ctx.U("/tran_export", "key", "pn")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Export all translation units as json file."), ctx.T("Export"))
		}

		if (exportRight && app.Ini.AppType == "prod") || ctx.IsSuperuser() {
			urlStr = ctx.U("/tran_import", "key", "pn")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Import all translation units from a json file."), ctx.T("Import"))
		}

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s</th>", ctx.T("English"))
	buf.Add("<th>%s</th>", ctx.T("Turkish"))

	if updateRight || deleteRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	if totalRows > 0 {
		tranList := tran_lib.GetTranPage(ctx, key, pageNo)

		var en, tr string
		for _, row := range tranList {
			ctx.Cargo.SetInt("tranId", row.TranId)

			en = util.ScrStr(row.En)
			tr = util.ScrStr(row.Tr)

			if key != "" {
				en = content.Find(en, key)
				tr = content.Find(tr, key)
			}

			buf.Add("<tr>")
			buf.Add("<td>%s</td>", en)
			buf.Add("<td>%s</td>", tr)

			if updateRight || deleteRight {
				buf.Add("<td class=\"right\">")
				buf.Add("<div class=\"buttonGroupFixed\">")

				if updateRight {
					urlStr := ctx.U("/tran_update", "tranId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Edit record."), ctx.T("Edit"))
				}

				if deleteRight {
					urlStr := ctx.U("/tran_delete", "tranId", "key", "pn")
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
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/tran", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
