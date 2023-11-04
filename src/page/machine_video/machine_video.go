package machine_video

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/machine_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_video/machine_video_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddInt("urlId", -1)
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

	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/machine_video/machine_video.js")

	browseMid(ctx, machineRec)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "machineVideoPage"

	ctx.AddHtml("pageName", &str)
	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx, machineRec *machine_lib.MachineRec) {
	updateRight := ctx.IsRight("machine", "update")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("CNC Machines"), ctx.T("Urls"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine", "key", "stat", "pn", "catId", "manId")))

	var urlStr string
	if updateRight {
		buf.Add(content.NewButton(ctx, ctx.U("/machine_video_insert", "machineId", "key", "pn", "stat", "catId", "manId")))
	}

	buf.Add("</div>")
	buf.Add("</div>")

	machineMenu := machine_menu.New("machineId", "key", "stat", "pn", "catId", "manId")
	machineMenu.Set(ctx, "machine_video")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(machineMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s / %s</th>", ctx.T("Title"), ctx.T("Video Url"))
	buf.Add("<th class=\"fixedZero\">%s</th>", ctx.T("Preview"))

	if updateRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	machineId := ctx.Cargo.Int("machineId")
	machineVideoList := machine_video_lib.GetMachineVideoList(machineId)
	if len(machineVideoList) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	var videoHash string
	var partList []string

	for _, row := range machineVideoList {
		ctx.Cargo.SetInt("urlId", row.MachineVideoId)

		buf.Add("<tr>")
		buf.Add("<td>")
		buf.Add("<p class=\"videoTitle\">%s</p>", util.ScrStr(row.Title))
		buf.Add("<p class=\"videoUrl\">")
		buf.Add("<a href=\"%s\" target=\"_blank\">%s</a>", util.ScrStr(row.Url), util.ScrStr(row.Url))
		buf.Add("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>", ctx.T("Copy link url."))
		buf.Add("</p>")
		buf.Add("</td>")

		partList = strings.Split(row.Url, "?v=")
		if len(partList) == 2 {
			videoHash = partList[1]
			buf.Add("<td>%s</td>", videoLink(videoHash))
		} else {
			buf.Add("<td>%s</td>", ctx.T("no video hash"))
		}

		if updateRight {
			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			urlStr = ctx.U("/machine_video_update", "urlId", "machineId", "key", "stat", "pn", "catId", "manId")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Edit record."), ctx.T("Edit"))

			urlStr = ctx.U("/machine_video_delete", "urlId", "machineId", "key", "stat", "pn", "catId", "manId")
			buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Delete record."), ctx.T("Delete"))

			buf.Add("</div>")
			buf.Add("</td>")
		}

		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

func videoLink(videoHash string) string {
	return fmt.Sprintf(`<iframe width='400' height='235' src="https://www.youtube.com/embed/%s?&theme=dark&autohide=2"frameborder="0" allowfullscreen></iframe><div style='font-size: 0.8em'></div>`, videoHash)
}
