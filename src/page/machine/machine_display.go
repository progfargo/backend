package machine

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/machine_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
)

func Display(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	machineId := ctx.Cargo.Int("machineId")
	rec, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	displayMachine(ctx, rec)
}

func displayMachine(ctx *context.Ctx, rec *machine_lib.MachineRec) {
	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/machine_display.css")
	ctx.Css.Add("/asset/jquery_modal/jquery.modal.css")
	ctx.Js.Add("/asset/jquery_modal/jquery.modal.js")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("Display Record"), util.ScrStr(rec.Name)))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	machineMenu := machine_menu.New("machineId", "key", "stat", "pn", "catId", "manId")
	machineMenu.Set(ctx, "machine_display")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(machineMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")

	machineImageId := util.NullToInt64(rec.MachineImageId)

	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">%s</th>", ctx.T("Image:"))

	if machineImageId != 0 {
		ctx.Cargo.SetInt("imgId", machineImageId)
		smallImageUrl := ctx.U("/machine_image_small", "imgId")
		popupUrl := ctx.U("/machine_popup", "imgId")

		buf.Add("<td>")
		buf.Add("<img class=\"smallImage\" src=\"%s\" alt=\"\">", smallImageUrl)
		buf.Add("<a class=\"popupImage\" href=\"%s\"></a>", popupUrl)
		buf.Add("</td>")
	} else {
		buf.Add("<td></td>")
	}

	buf.Add("</tr>")

	name := util.ScrStr(rec.Name)
	model := util.ScrStr(rec.Model)
	manufacturer := util.ScrStr(rec.ManufacturerName)
	exp := util.NullToString(rec.Exp)
	location := util.ScrStr(rec.Location)

	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Name:"), name)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Model:"), model)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Manufacturer:"), manufacturer)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Explanation:"), exp)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Location:"), location)
	buf.Add("<tr><th title=\"%s\">%s</th><td>%d</td></tr>", ctx.T("Yom:"), ctx.T("Year of Manufacture"), rec.Yom)
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Status:"), machine_lib.StatusToLabel(ctx, rec.Status))
	buf.Add("<tr><th>%s</th><td>%d &#x24;</td></tr>", ctx.T("Price:"), rec.Price)

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")
	buf.Add("</div>") //row

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "machineDisplayPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}
