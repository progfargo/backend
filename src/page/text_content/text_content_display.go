package text_content

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/text_content_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/text_content/text_content_lib"
)

func Display(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("text_content", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("textId", -1)

	ctx.ReadCargo()

	textId := ctx.Cargo.Int("textId")
	rec, err := text_content_lib.GetTextContentRec(textId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/text_content"))
			return
		}

		panic(err)
	}

	displayText(ctx, rec)
	ctx.Redirect(ctx.U("/text_content"))
}

func displayText(ctx *context.Ctx, rec *text_content_lib.TextContentRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Text Content"), ctx.T("Display Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/text_content")))
	buf.Add("</div>")
	buf.Add("</div>")

	textMenu := text_content_menu.New("textId")
	textMenu.Set(ctx, "text_content_display")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(textMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td> %s</td></tr>", ctx.T("Name:"), util.ScrStr(rec.Name))
	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td> %s</td></tr>", ctx.T("Explanation:"), util.ScrStr(rec.Exp))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Body:"), rec.Body)

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	buf.Add("</div>") //row

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "text_content")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
