package faq

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/faq/faq_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("faq", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("faqId", -1)
		ctx.Cargo.AddStr("key", "")
		ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()
	

	faqId := ctx.Cargo.Int("faqId")
	rec, err := faq_lib.GetFaqRec(faqId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/faq", "key", "pn"))
			return
		}

		panic(err)
	}

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteConfirm(ctx, rec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					faq
				where
					faqId = ?`

	res, err := tx.Exec(sqlStr, faqId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/faq", "key", "pn"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/faq", "key", "pn"))
}

func deleteConfirm(ctx *context.Ctx, rec *faq_lib.FaqRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("FAQ"), ctx.T("Delete Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/faq", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("Question:"), rec.Question)

	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Summary:"), rec.Summary)

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h4>%s</h4>", ctx.T("Please confirm:"))
	buf.Add("<p>%s</p>", ctx.T("Do you realy want to delete this record?"))
	buf.Add("</div>")
	buf.Add("</div>")

	ctx.Cargo.SetStr("confirm", "yes")
	urlStr := ctx.U("/faq_delete", "faqId", "key", "pn", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "faq")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}