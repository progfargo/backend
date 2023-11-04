package banner

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/banner/banner_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("banner", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("bannerId", -1)
	ctx.ReadCargo()

	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/banner/banner.js")

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "banner")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "bannerPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	insertRight := ctx.IsRight("banner", "insert")
	updateRight := ctx.IsRight("banner", "update")
	deleteRight := ctx.IsRight("banner", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Banner")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")
		buf.Add(content.NewButton(ctx, ctx.U("/banner_insert")))
		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedZero\" title=\"Enumeration\">%s</th>", ctx.T("En."))
	buf.Add("<th class=\"fixedZero\">%s</th>", ctx.T("Image"))
	buf.Add("<th>%s</th>", ctx.T("Image Info"))
	buf.Add("<th>%s</th>", ctx.T("Status"))

	if updateRight || deleteRight {
		buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	bannerList := banner_lib.GetBannerList()
	if len(bannerList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	var urlStr string
	for _, row := range bannerList {
		ctx.Cargo.SetInt("bannerId", row.BannerId)

		buf.Add("<tr>")
		buf.Add("<td>%d</td>", row.Enum)
		urlStr = ctx.U("/banner_image_small", "bannerId")
		buf.Add("<td><img src=\"%s\" alt=\"%s\" title=\"%s\"></td>", urlStr, util.ScrStr(row.ImgName),
			util.ScrStr(row.ImgName))

		buf.Add("<td>")
		buf.Add("<div class=\"imageInfo\">")

		urlStr = ctx.U("/banner_image_original", "bannerId")
		buf.Add("<strong>%s:</strong><a href=\"%s\" target=\"_blank\">%s</a>", util.PadRight(ctx.T("Url"), "&nbsp;", 7), urlStr, urlStr)
		buf.Add("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>", ctx.T("Copy link url."))

		urlStr = ctx.U("/banner_image_download", "bannerId")
		buf.Add("<a class=\"downloadLink\" href=\"%s\" download><i class=\"fas fa-file-download\"></i></a></br>", urlStr)
		buf.Add("<strong>%s:</strong> %d*%d px<br>", util.PadRight(ctx.T("W/H"), "&nbsp;", 7), row.ImgWidth, row.ImgHeight)

		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Name"), "&nbsp;", 7), util.ScrStr(row.ImgName))
		buf.Add("<strong>%s:</strong> %s bytes<br>", util.PadRight(ctx.T("Size"), "&nbsp;", 7), util.FormatInt(row.ImgSize))
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Mime"), "&nbsp", 7), util.ScrStr(row.ImgMime))

		buf.Add("</div>")
		buf.Add("</td>")

		buf.Add("<td>%s</td>", banner_lib.StatusToLabel(ctx, row.Status))

		if updateRight || deleteRight {
			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if updateRight {
				urlStr = ctx.U("/banner_update_info", "bannerId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit image information."), ctx.T("Edit Info"))

				urlStr = ctx.U("/banner_update_image", "bannerId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit image information."), ctx.T("Edit Image"))

				urlStr = ctx.U("/banner_crop_image", "bannerId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Crop image."), ctx.T("Crop"))
			}

			if deleteRight {
				urlStr = ctx.U("/banner_delete", "bannerId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))
			}

			buf.Add("</div>")
		}

		buf.Add("</td>")
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
