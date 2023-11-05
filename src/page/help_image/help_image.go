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
	"backend/src/page/help/help_lib"
	"backend/src/page/help_image/help_image_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("help", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("helpId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	browseMid(ctx)

	content.Default(ctx)
	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/help_image/help_image.js")

	lmenu := left_menu.New()
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "help")

	str := "helpImagePage"

	ctx.AddHtml("pageName", &str)
	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

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

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Help"), ctx.T("Images")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/help", "key", "pn")))
	buf.Add(content.NewButton(ctx, ctx.U("/help_image_insert", "helpId", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("Help Information:"))
	buf.Add("<tbody>")

	title := rec.Title
	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("Title:"), util.ScrStr(title))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("Help Images:"))
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">%s</th>", ctx.T("Image"))
	buf.Add("<th>%s</th>", ctx.T("Image Info"))
	buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	helpImageList := help_image_lib.GetHelpImageList(rec.HelpId)
	if len(helpImageList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	for _, row := range helpImageList {
		ctx.Cargo.SetInt("imgId", row.HelpImageId)

		buf.Add("<tr>")

		urlStr := ctx.U("/help_image_small", "imgId")
		buf.Add("<td class=\"helpImage\"><img src=\"%s\" alt=\"%s\" title=\"%s\"></td>", urlStr, util.ScrStr(row.ImgName),
			util.ScrStr(row.ImgName))

		buf.Add("<td>")
		buf.Add("<div class=\"imageInfo\">")

		urlStr = ctx.U("/help_image_original", "key", "pn", "helpId", "imgId")
		buf.Add("<strong>%s:</strong><a href=\"%s\" target=\"_blank\">%s</a>", util.PadRight(ctx.T("Url"), "&nbsp;", 7), urlStr, urlStr)
		buf.Add("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>", ctx.T("Copy link url."))

		urlStr = ctx.U("/help_image_download", "key", "pn", "helpId", "imgId")
		buf.Add("<a class=\"downloadLink\" href=\"%s\" download><i class=\"fas fa-file-download\"></i></a></br>", urlStr)
		buf.Add("<strong>%s:</strong> %d*%d px<br>", util.PadRight(ctx.T("W/H"), "&nbsp;", 7), row.ImgWidth, row.ImgHeight)

		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Name"), "&nbsp;", 7), util.ScrStr(row.ImgName))
		buf.Add("<strong>%s:</strong> %s bytes<br>", util.PadRight(ctx.T("Size"), "&nbsp;", 7), util.FormatInt(row.ImgSize))
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Mime"), "&nbsp", 7), util.ScrStr(row.ImgMime))

		buf.Add("</div>")

		buf.Add("</td>")

		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFlex\">")

		urlStr = ctx.U("/help_image_update_image", "helpId", "imgId", "key", "pn")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit image."), ctx.T("Edit"))

		urlStr = ctx.U("/help_image_delete", "helpId", "imgId", "key", "pn")
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

	ctx.AddHtml("midContent", buf.String())
}
