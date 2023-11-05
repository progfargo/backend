package machine_image

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_image/machine_image_lib"
)

func Delete(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddStr("confirm", "no")
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)

	ctx.ReadCargo()

	machineId := ctx.Cargo.Int("machineId")
	machineRec, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		ctx.Msg.Error(ctx.T("You don't have right to access this record."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	imgId := ctx.Cargo.Int("imgId")
	machineImageRec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if ctx.Cargo.Str("confirm") != "yes" {
		deleteConfirm(ctx, machineRec, machineImageRec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/machine_image", "machineId"))
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					machineImage
				where
					machineImageId = ?`

	res, err := tx.Exec(sqlStr, imgId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Record could not be found."))
		ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been deleted."))
	ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))

}

func deleteConfirm(ctx *context.Ctx, machineRec *machine_lib.MachineRec, machineImageRec *machine_image_lib.MachineImageRec) {
	content.Include(ctx)
	ctx.Css.Add("/asset/jquery_modal/jquery.modal.css")
	ctx.Js.Add("/asset/jquery_modal/jquery.modal.js")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("machines"), ctx.T("Delete Image"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<table>")
	buf.Add("<tbody>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixed-middle\">%s</th>", ctx.T("Image:"))

	smallImageUrl := ctx.U("/machine_image_small", "imgId")
	popupUrl := ctx.U("/machine_popup", "imgId")

	buf.Add("<td>")
	buf.Add("<img class=\"smallImage\" src=\"%s\" alt=\"\">", smallImageUrl)
	buf.Add("<a class=\"popupImage\" href=\"%s\"></a>", popupUrl)
	buf.Add("</td>")

	buf.Add("</tr>")

	buf.Add("<tr><th>%s</th><td>%s bytes</td></tr>", ctx.T("Size:"), util.FormatInt(machineImageRec.ImgSize))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Name:"), util.ScrStr(machineImageRec.ImgName))
	buf.Add("<tr><th>%s</th><td>%s</td></tr>", ctx.T("Mime:"), util.ScrStr(machineImageRec.ImgMime))
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
	urlStr := ctx.U("/machine_image_delete", "machineId", "imgId", "confirm", "key", "pn", "stat", "manId")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"confirmCommand\">")
	buf.Add("<a href=\"%s\" class=\"button buttonError buttonSm\">%s</a>", urlStr, ctx.T("Yes"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
