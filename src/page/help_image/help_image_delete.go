package help_image

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/help_image/help_image_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("help", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("helpId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	imgId := ctx.Cargo.Int("imgId")
	rec, err := help_image_lib.GetHelpImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.RenderAjaxError()
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
					helpImage
				where
					helpImageId = ?`

	res, err := tx.Exec(sqlStr, imgId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/help_image", "helpId", "key", "pn"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/help_image", "helpId", "key", "pn"))

}

func deleteConfirm(ctx *context.Ctx, rec *help_image_lib.HelpImageRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Help"), ctx.T("Images"), ctx.T("New Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/help_image", "helpId", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixed-middle\">%s</th>", ctx.T("Image:"))

	imgUrl := ctx.U("/help_image_small", "imgId")
	buf.Add("<td><img src=\"%s\" alt=\"%s\"></td>", imgUrl, rec.ImgName)
	buf.Add("</tr>")

	buf.Add("<tr><th>%s</th><td>%s bytes</td></tr>", ctx.T("Size:"), util.FormatInt(rec.ImgSize))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Name:"), util.ScrStr(rec.ImgName))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Mime:"), util.ScrStr(rec.ImgMime))
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
	urlStr := ctx.U("/help_image_delete", "helpId", "imgId", "confirm", "key", "pn")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)
	content.Include(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "help")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
