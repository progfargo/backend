package help

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/help/help_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("help", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("helpId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()

	helpId := ctx.Cargo.Int("helpId")
	rec, err := help_lib.GetHelpRec(helpId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/help", "key", "pn"))
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
					help
				where
					helpId = ?`

	res, err := tx.Exec(sqlStr, helpId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/help", "key", "pn"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/help", "key", "pn"))
}

func deleteConfirm(ctx *context.Ctx, rec *help_lib.HelpRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Help"), ctx.T("Delete Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/help", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	title := rec.Title
	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("Title:"), util.ScrStr(title))

	summary := rec.Summary
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Summary:"), summary)

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
	urlStr := ctx.U("/help_delete", "helpId", "key", "pn", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "help")

	ctx.Render("default.html")
}
