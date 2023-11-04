package office

import (
	"database/sql"
	"net/http"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/office/office_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("office", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("productionId", -1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.ReadCargo()

	productionId := ctx.Cargo.Int("productionId")
	rec, err := office_lib.GetOfficeRec(productionId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/office"))
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
					office
				where
					officeId = ?`

	res, err := tx.Exec(sqlStr, productionId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/office"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/office"))
}

func deleteConfirm(ctx *context.Ctx, rec *office_lib.OfficeRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Offices"), ctx.T("Delete Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/office")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	name := util.ScrStr(rec.Name)
	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("Name:"), name)

	city := util.ScrStr(rec.City)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("City:"), city)

	address := util.ScrStr(rec.Address)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Addres:"), strings.Replace(address, "\n", "<br>", -1))

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
	urlStr := ctx.U("/office_delete", "productionId", "confirm")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "office")

	tmenu := top_menu.New()
	tmenu.Set(ctx, "root")

	ctx.Render("default.html")
}
