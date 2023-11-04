package jumbotron

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/jumbotron/jumbotron_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("jumbotron", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("jumbotronId", -1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "jumbotron")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	insertRight := ctx.IsRight("jumbotron", "insert")
	updateRight := ctx.IsRight("jumbotron", "update")
	deleteRight := ctx.IsRight("jumbotron", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Jumbotron")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")
		buf.Add(content.NewButton(ctx, ctx.U("/jumbotron_insert")))
		buf.Add("</div>")
		buf.Add("</div>")
	}

	jumbotronList := jumbotron_lib.GetJumbotronList()
	if len(jumbotronList) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	for _, row := range jumbotronList {
		ctx.Cargo.SetInt("jumbotronId", row.JumbotronId)

		buf.Add("<div class=\"col\">")
		buf.Add("<h3>%s</h3>", row.Title)
		buf.Add(row.Body)
		buf.Add("</div>")

		if updateRight || deleteRight {
			buf.Add("<div class=\"col\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if updateRight {
				urlStr := ctx.U("/jumbotron_update", "jumbotronId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit record."), ctx.T("Edit"))
			}

			if deleteRight {
				urlStr := ctx.U("/jumbotron_delete", "jumbotronId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))
			}

			buf.Add("</div>")
			buf.Add("</div>")
		}
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
