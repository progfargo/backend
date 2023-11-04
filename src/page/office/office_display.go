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

func Display(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("office", "delete") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("productionId", -1)
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

	displayRec(ctx, rec)
}

func displayRec(ctx *context.Ctx, rec *office_lib.OfficeRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Offices"), ctx.T("Display Record")))
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

	telephone := util.ScrStr(rec.Telephone)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Telephone:"), telephone)

	email := util.ScrStr(rec.Email)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Email:"), email)

	mapLink := util.ScrStr(rec.MapLink)
	buf.Add("<tr><th>%s</th><td id=\"longitude\">%s</td></tr>", ctx.T("Map Link:"), mapLink)

	buf.Add("<tr><th>%s</th><td id=\"latitude\">%s</td></tr>", ctx.T("Is Main Office:"),
		office_lib.IsMainOfficeToLabel(ctx, rec.IsMainOffice))

	buf.Add("<tr><th>%s</th><td><div id=\"mapid\"></div></td></tr>", ctx.T("Map:"))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "office")

	tmenu := top_menu.New()
	tmenu.Set(ctx, "root")

	str := "officeDisplayPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}
