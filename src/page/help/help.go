package help

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/help/help_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("help", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("helpId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, "/help")

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "help")

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")

	totalRows := help_lib.CountHelp(key)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("help", "insert")
	updateRight := ctx.IsRight("help", "update")
	deleteRight := ctx.IsRight("help", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Help")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		if insertRight {
			buf.Add(content.NewButton(ctx, ctx.U("/help_insert", "key", "pn")))
		}

		buf.Add("</div>")
		buf.Add("</div>")
	}

	if totalRows > 0 {
		helpList := help_lib.GetHelpPage(ctx, key, pageNo)

		var title, summary string
		for _, row := range helpList {
			ctx.Cargo.SetInt("helpId", row.HelpId)

			title = row.Title
			title = util.ScrStr(row.Title)

			summary = row.Summary

			if key != "" {
				title = content.Find(title, key)
				summary = content.Find(summary, key)
			}

			buf.Add("<div class=\"col\">")
			buf.Add("<h3>%s</h3>", title)
			buf.Add(summary)
			buf.Add("</div>")

			buf.Add("<div class=\"col\">")

			if updateRight || deleteRight {
				buf.Add("<div class=\"buttonGroupFixed\">")

				if updateRight {
					urlStr := ctx.U("/help_update", "helpId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Edit record."), ctx.T("Edit"))

					urlStr = ctx.U("/help_image", "helpId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Help images."), ctx.T("Images"))
				}

				if deleteRight {
					urlStr := ctx.U("/help_delete", "helpId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Delete record."), ctx.T("Delete"))
				}

				buf.Add("</div>")
			}

			buf.Add("</div>")
		}
	}

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/help", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
