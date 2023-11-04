package text_content_image

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/text_content_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/text_content/text_content_lib"
	"backend/src/page/text_content_image/text_content_image_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("text_content", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("textId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddStr("stat", "")
	ctx.ReadCargo()

	browseMid(ctx)

	content.Default(ctx)
	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/text_content_image/text_content_image.js")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "text_content")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "textContentImagePage"

	ctx.AddHtml("pageName", &str)
	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	textContentId := ctx.Cargo.Int("textId")
	rec, err := text_content_lib.GetTextContentRec(textContentId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/text_content", "key", "pn", "stat"))
			return
		}

		panic(err)
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Text Content"), ctx.T("Images")))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/text_content", "key", "pn", "stat")))
	buf.Add(content.NewButton(ctx, ctx.U("/text_content_image_insert", "textId", "key", "pn", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	textMenu := text_content_menu.New("textId")
	textMenu.Set(ctx, "text_content_image")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(textMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("Text Content Information:"))
	buf.Add("<tbody>")

	buf.Add("<tr><th class=\"fixedMiddle\">%s</th><td>%s</td></tr>", ctx.T("Title:"), util.ScrStr(rec.Name))

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<caption>%s</caption>", ctx.T("Text Content Images:"))
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedMiddle\">%s</th>", ctx.T("Image"))
	buf.Add("<th>%s</th>", ctx.T("Image Info"))
	buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	textContentImageList := text_content_image_lib.GetTextContentImageList(rec.TextContentId)
	if len(textContentImageList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	var alt string
	for _, row := range textContentImageList {
		ctx.Cargo.SetInt("imgId", row.TextContentImageId)

		buf.Add("<tr>")

		urlStr := ctx.U("/text_content_image_small", "imgId")
		buf.Add("<td class=\"textContentImage\"><img src=\"%s\" alt=\"%s\" title=\"%s\"></td>", urlStr, util.ScrStr(row.ImgName),
			util.ScrStr(row.ImgName))

		buf.Add("<td>")
		buf.Add("<div class=\"imageInfo\">")

		urlStr = ctx.U("/text_content_image_original", "imgId")
		buf.Add("<strong>%s:</strong><a href=\"%s\" target=\"_blank\">%s</a>", util.PadRight(ctx.T("Url"), "&nbsp;", 7), urlStr, urlStr)
		buf.Add("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>", ctx.T("Copy link url."))

		urlStr = ctx.U("/text_content_image_download", "imgId")
		buf.Add("<a class=\"downloadLink\" href=\"%s\" download><i class=\"fas fa-file-download\"></i></a></br>", urlStr)
		buf.Add("<strong>%s:</strong> %d*%d px<br>", util.PadRight(ctx.T("W/H"), "&nbsp;", 7), row.ImgWidth, row.ImgHeight)

		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Name"), "&nbsp;", 7), util.ScrStr(row.ImgName))
		buf.Add("<strong>%s:</strong> %s bytes<br>", util.PadRight(ctx.T("Size"), "&nbsp;", 7), util.FormatInt(row.ImgSize))
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Mime"), "&nbsp", 7), util.ScrStr(row.ImgMime))

		alt = util.NullToString(row.Alt)
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Alt"), "&nbsp;", 7), util.ScrStr(alt))

		buf.Add("</div>")

		buf.Add("</td>")

		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFlex\">")

		urlStr = ctx.U("/text_content_image_update_info", "textId", "imgId", "key", "pn", "stat")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit image info."), ctx.T("Edit Info"))

		urlStr = ctx.U("/text_content_image_update", "textId", "imgId", "key", "pn", "stat")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit image."), ctx.T("Edit Image"))

		urlStr = ctx.U("/text_content_image_delete", "textId", "imgId", "key", "pn", "stat")
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
