package news

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/news/news_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("news", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("newsId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("stat", "")
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, "/news")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "news")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")
	stat := ctx.Cargo.Str("stat")

	totalRows := news_lib.CountNews(key, stat)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("news", "insert")
	updateRight := ctx.IsRight("news", "update")
	deleteRight := ctx.IsRight("news", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("News")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	if insertRight {
		buf.Add(content.NewButton(ctx, ctx.U("/news_insert", "key", "pn", "stat")))
	}

	if stat == "" {
		urlStr := ctx.U("/news?stat=draft", "key", "pn")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Draft records."), ctx.T("Show Draft records."))
	} else if stat == "draft" {
		urlStr := ctx.U("/news", "key", "pn")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
			urlStr, ctx.T("All records."), ctx.T("Show All records."))
	}

	buf.Add("</div>")
	buf.Add("</div>")

	if totalRows > 0 {
		newsList := news_lib.GetNewsPage(ctx, key, stat, pageNo)

		var title, statusIcon, summary string
		for _, row := range newsList {
			ctx.Cargo.SetInt("newsId", row.NewsId)

			title = row.Title
			summary = row.Summary

			if key != "" {
				title = content.Find(title, key)
				summary = content.Find(summary, key)
			}

			if row.Status == "published" {
				statusIcon = "<i class=\"fas fa-eye fa-fw\"></i>"
			} else {
				statusIcon = "<i class=\"fas fa-eye-slash fa-fw\"></i>"
			}

			buf.Add("<div class=\"col\">")
			buf.Add("<h3>(%s) %s %s</h3>",
				util.Int64ToDateStr(row.RecordDate), title, statusIcon)
			buf.Add(summary)
			buf.Add("</div>")

			buf.Add("<div class=\"col\">")

			if updateRight || deleteRight {
				buf.Add("<div class=\"buttonGroupFixed\">")

				if updateRight {
					urlStr := ctx.U("/news_update", "newsId", "key", "stat", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Edit record."), ctx.T("Edit"))

					urlStr = ctx.U("/news_image", "newsId", "key", "stat", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("News images."), ctx.T("Images"))
				}

				if deleteRight {
					urlStr := ctx.U("/news_delete", "newsId", "key", "stat", "pn")
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
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/news", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
