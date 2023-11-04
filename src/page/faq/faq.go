package faq

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/faq/faq_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("faq", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("faqId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)
	content.Search(ctx, "/faq")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "faq")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	key := ctx.Cargo.Str("key")
	pageNo := ctx.Cargo.Int("pn")

	totalRows := faq_lib.CountFaq(key)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("faq", "insert")
	updateRight := ctx.IsRight("faq", "update")
	deleteRight := ctx.IsRight("faq", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("FAQ")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")
		buf.Add(content.NewButton(ctx, ctx.U("/faq_insert", "key", "pn")))
		buf.Add("</div>")
		buf.Add("</div>")
	}

	if totalRows > 0 {

		faqList := faq_lib.GetFaqPage(ctx, key, pageNo)

		var question, summary string
		for _, row := range faqList {
			ctx.Cargo.SetInt("faqId", row.FaqId)

			question = row.Question
			summary = row.Summary

			if key != "" {
				question = content.Find(question, key)
				summary = content.Find(summary, key)
			}

			buf.Add("<div class=\"col\">")
			buf.Add("<h3>%s</h3>", question)
			buf.Add(summary)
			buf.Add("</div>")

			if updateRight || deleteRight {
				buf.Add("<div class=\"col\">")
				buf.Add("<div class=\"buttonGroupFixed\">")

				if updateRight {
					urlStr := ctx.U("/faq_update", "faqId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Edit record."), ctx.T("Edit"))
				}

				if deleteRight {
					urlStr := ctx.U("/faq_delete", "faqId", "key", "pn")
					buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
						urlStr, ctx.T("Delete record."), ctx.T("Delete"))
				}

				buf.Add("</div>")
				buf.Add("</div>")
			}
		}
	}

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/faq", "key"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
