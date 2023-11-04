package jumbotron

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/jumbotron/jumbotron_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("jumbotron", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("jumbotronId", -1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()

	jumbotronId := ctx.Cargo.Int("jumbotronId")
	rec, err := jumbotron_lib.GetJumbotronRec(jumbotronId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/jumbotron"))
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
					jumbotron
				where
					jumbotronId = ?`

	res, err := tx.Exec(sqlStr, jumbotronId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/jumbotron"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/jumbotron"))
}

func deleteConfirm(ctx *context.Ctx, rec *jumbotron_lib.JumbotronRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Jumbotron"), ctx.T("Delete Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/jumbotron")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("Title:"), rec.Title)

	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Body:"), rec.Body)

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
	urlStr := ctx.U("/jumbotron_delete", "jumbotronId", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "jumbotron")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
