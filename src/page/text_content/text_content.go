package text_content

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/text_content/text_content_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("text_content", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("textId", -1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "text_content")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	insertRight := ctx.IsRight("text_content", "insert")
	updateRight := ctx.IsRight("text_content", "update")
	setRight := ctx.IsRight("text_content", "set")
	deleteRight := ctx.IsRight("text_content", "delete")
	browseRight := ctx.IsRight("text_content", "browse")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Text Content")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")
		buf.Add(content.NewButton(ctx, ctx.U("/text_content_insert")))
		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s</th>", ctx.T("Name"))
	buf.Add("<th>%s</th>", ctx.T("Explanation"))

	if updateRight || deleteRight || browseRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	textList := text_content_lib.GetTextContentList()
	if len(textList) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	for _, row := range textList {
		ctx.Cargo.SetInt("textId", row.TextContentId)

		buf.Add("<tr>")
		buf.Add("<td>%s</td>", row.Name)
		buf.Add("<td>%s</td>", row.Exp)

		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		urlStr := ctx.U("/text_content_display", "textId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Display record."), ctx.T("Display"))

		if setRight {
			urlStr := ctx.U("/text_content_set", "textId")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Set record value."), ctx.T("Set"))
		}

		if updateRight {
			urlStr := ctx.U("/text_content_update", "textId")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Update record."), ctx.T("Update"))
		}

		if deleteRight {
			urlStr := ctx.U("/text_content_delete", "textId")
			buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Delete record."), ctx.T("Delete"))
		}

		buf.Add("</div>")

		buf.Add("</td>")
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>") //col
	buf.Add("</div>") //row

	ctx.AddHtml("midContent", buf.String())
}
