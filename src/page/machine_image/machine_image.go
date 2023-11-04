package machine_image

import (
	"database/sql"
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/machine_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_image/machine_image_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
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
	machineRec, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		ctx.Msg.Error(ctx.T("You don't have right to access this record."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	browseMid(ctx, machineRec)

	content.Include(ctx)
	ctx.Css.Add("/asset/jquery_modal/jquery.modal.css")
	ctx.Js.Add("/asset/jquery_modal/jquery.modal.js")
	ctx.Css.Add("/asset/css/page/machine_image.css")
	ctx.Js.Add("/asset/js/page/machine_image/machine_image.js")

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "machineImagePage"

	ctx.AddHtml("pageName", &str)
	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx, machineRec *machine_lib.MachineRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("Images"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine", "key", "stat", "pn", "catId", "manId")))
	buf.Add(content.NewButton(ctx, ctx.U("/machine_image_insert", "machineId", "key", "pn", "stat", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	machineMenu := machine_menu.New("machineId", "key", "stat", "pn", "catId", "manId")
	machineMenu.Set(ctx, "machine_image")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(machineMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedZero\" title=\"%s\">%s</th>", ctx.T("Enumeration"), ctx.T("E."))
	buf.Add("<th>%s</th>", ctx.T("Image"))
	buf.Add("<th>")
	buf.Add(ctx.T("Image Info"))
	buf.AddLater("<span class=\"allCropped\">(%s %s)</span>", ctx.T("All cropped:"), "<i class=\"fas fa-check\"></i>")
	buf.Add("</th>")

	buf.Add("<th title=\"%s\">%s</th>", ctx.T("Header image."), ctx.T("Is Header"))
	buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	machineImageList := machine_image_lib.GetMachineImageList(machineRec.MachineId)
	if len(machineImageList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	var alt, smallImageUrl, popupUrl, rotateLink90, rotateLink270, rotateUrl90, rotateUrl270 string
	padSize := 8
	allCropped := true

	for _, row := range machineImageList {
		ctx.Cargo.SetInt("imgId", row.MachineImageId)

		buf.Add("<tr>")

		buf.Add("<td>%d</td>", row.Enum)

		smallImageUrl = ctx.U("/machine_image_small", "imgId")
		popupUrl = ctx.U("/machine_popup", "imgId")

		buf.Add("<td>")
		buf.Add("<img id=\"%d\" class=\"smallImage\" src=\"%s\" alt=\"\" title=\"\">", row.MachineImageId, smallImageUrl)
		buf.Add("<a class=\"popupImage\" href=\"%s\"></a>", popupUrl)
		buf.Add("</td>")

		buf.Add("<td>")
		buf.Add("<div class=\"imageInfo\">")

		urlStr := ctx.U("/machine_image_original", "imgId")

		buf.Add("<strong>%s:</strong><a href=\"%s\" target=\"_blank\"> %s</a>", util.PadRight(ctx.T("Url"), "&nbsp;", padSize), urlStr, urlStr)

		buf.Add("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>", ctx.T("Copy link url."))

		urlStr = ctx.U("/machine_image_download", "imgId")
		buf.Add("<a class=\"downloadLink\" href=\"%s\" title=\"%s\"><i class=\"fas fa-file-download\"></i></a><br>", urlStr, ctx.T("Download"))

		buf.Add("<strong>%s:</strong> %d*%d px<br>", util.PadRight(ctx.T("W/H"), "&nbsp;", padSize), row.ImgWidth, row.ImgHeight)

		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Name"), "&nbsp;", padSize), util.ScrStr(row.ImgName))
		buf.Add("<strong>%s:</strong> %s bytes<br>", util.PadRight(ctx.T("Size"), "&nbsp;", padSize), util.FormatInt(row.ImgSize))
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Mime"), "&nbsp;", padSize), util.ScrStr(row.ImgMime))

		alt = util.NullToString(row.Alt)
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Alt"), "&nbsp;", padSize), util.ScrStr(alt))

		rotateUrl90 = ctx.U("/machine_image_rotate?angle=90", "imgId")
		rotateUrl270 = ctx.U("/machine_image_rotate?angle=270", "imgId")
		rotateLink90 = fmt.Sprintf("<a href=\"%s\" title=\"%s\" class=\"rotateLink\"><i class=\"fas fa-redo\"></i></a>", rotateUrl90, ctx.T("Rotate 90 deg."))
		rotateLink270 = fmt.Sprintf("<a href=\"%s\" title=\"%s\" class=\"rotateLink\"><i class=\"fas fa-undo\"></i></a>", rotateUrl270, ctx.T("Rotate -90 deg."))
		buf.Add("<strong>%s:</strong> %s %s<br>", util.PadRight(ctx.T("Rotate"), "&nbsp;", padSize), rotateLink270, rotateLink90)

		if row.IsCropped == "no" {
			allCropped = false
		}

		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Cropped"), "&nbsp;", padSize), machine_image_lib.IsCroppedToLabel(ctx, row.IsCropped))
		buf.Add("</div>")
		buf.Add("</td>")

		buf.Add("<td class=\"center\">")
		buf.Add(machine_image_lib.IsHeaderToLabel(ctx, row.IsHeader))
		buf.Add("</td>")

		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFlex\">")

		urlStr = ctx.U("/machine_image_make_header", "machineId", "imgId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Select this record as header image."), ctx.T("Make Header"))

		urlStr = ctx.U("/machine_image_update", "machineId", "imgId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit record."), ctx.T("Edit Image"))

		urlStr = ctx.U("/machine_image_update_info", "machineId", "imgId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit image information."), ctx.T("Edit Info"))

		urlStr = ctx.U("/machine_image_crop", "machineId", "imgId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Crop image."), ctx.T("Crop"))

		urlStr = ctx.U("/machine_image_delete", "machineId", "imgId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Delete record."), ctx.T("Delete"))

		buf.Add("</div>")
		buf.Add("</td>")
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	if allCropped && len(machineImageList) > 0 {
		buf.Forge()
	}

	ctx.AddHtml("midContent", buf.String())
}
