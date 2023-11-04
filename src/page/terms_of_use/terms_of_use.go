package terms_of_use

import (
	"database/sql"
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

	if !ctx.IsRight("corporate", "update") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	rec, err := text_content_lib.GetTextContentRecByName("terms_of_use")
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/welcome"))
			return
		}

		panic(err)
	}

	displayTermsOfUse(ctx, rec)
	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func displayTermsOfUse(ctx *context.Ctx, rec *text_content_lib.TextContentRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")

	buf.Add(content.PageTitle(ctx.T("Terms of Use")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add(rec.Body)

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}