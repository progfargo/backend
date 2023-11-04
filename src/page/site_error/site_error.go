package site_error

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/site_error/site_error_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("site_error", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("errorId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	content.Include(ctx)
	content.Default(ctx)

	browseMid(ctx)

	content.Search(ctx, "/site_error")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "site_error")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "siteErrorPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")

	totalRows := site_error_lib.CountSiteError(key)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	deleteRight := ctx.IsRight("site_error", "delete")
	deleteAllRight := ctx.IsRight("site_error", "delete_all")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Site Errors")))
	buf.Add("</div>")

	if deleteAllRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		urlStr := ctx.U("/site_error_delete_all", "key", "pn")
		buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm deleteAllButton\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Delete all records."), ctx.T("Delete All"))

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedSmall\">%s</th>", ctx.T("Date/Time"))
	buf.Add("<th>%s</th>", ctx.T("Title"))
	buf.Add("<th>%s</th>", ctx.T("Messasge"))

	if deleteRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	if totalRows > 0 {
		siteErrorList := site_error_lib.GetSiteErrorPage(ctx, key, pageNo)

		var title, message string
		for _, row := range siteErrorList {
			ctx.Cargo.SetInt("errorId", row.SiteErrorId)

			title = util.ScrStr(row.Title)
			message = util.ScrStr(row.Message)

			if key != "" {
				title = content.Find(title, key)
				message = content.Find(message, key)
			}

			buf.Add("<tr>")
			buf.Add("<td class=\"date\">%s</td>", util.Int64ToTimeStr(row.DateTime))
			buf.Add("<td>%s</td>", title)
			buf.Add("<td>%s</td>", message)

			if deleteRight {
				buf.Add("<td class=\"right\">")
				buf.Add("<div class=\"buttonGroupFixed\">")

				urlStr := ctx.U("/site_error_delete", "errorId", "key", "pn")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))

				buf.Add("</td>")
			}

			buf.Add("</div>")
			buf.Add("</tr>")
		}
	}

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/site_error", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
