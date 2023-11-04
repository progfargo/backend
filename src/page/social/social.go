package social

import (
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/social/social_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("social", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("socialId", -1)
	ctx.Cargo.AddStr("sortOrder", "")
	ctx.ReadCargo()

	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/social/social.js")

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "social")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {

	insertRight := ctx.IsRight("social", "insert")
	updateRight := ctx.IsRight("social", "update")
	deleteRight := ctx.IsRight("social", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Social Links")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		if insertRight {
			buf.Add(content.NewButton(ctx, ctx.U("/social_insert")))
		}

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th>%s</th>", ctx.T("Icon"))

	urlStr := ctx.U("/social?sortOrder=link", "socialId")
	buf.Add("<th><a href=\"%s\">%s</a></th>", urlStr, ctx.T("Link"))
	buf.Add("<th>%s</th>", ctx.T("Title"))
	buf.Add("<th>%s</th>", ctx.T("Target"))

	if updateRight || deleteRight {
		buf.Add("<th class=\"fixedSmall right\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	socialList := social_lib.GetSocialList(ctx)

	if len(socialList) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	var link, title, target, icon string
	for _, row := range socialList {
		ctx.Cargo.SetInt("socialId", row.SocialId)

		icon = util.ScrStr(row.Icon)
		link = util.ScrStr(row.Link)
		title = util.ScrStr(row.Title)
		target = util.ScrStr(row.Target)

		copyLink := fmt.Sprintf("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>",
			ctx.T("Copy link url."))

		buf.Add("<tr>")
		buf.Add("<td title=\"%s\"><i class=\"%s fa-fw\"></i></td>", title, icon)
		buf.Add("<td title=\"%s\"><a href=\"%s\" target=\"_blank\">%s</a>%s</td>", title, link, title, copyLink)
		buf.Add("<td>%s</td>", title)

		var titleStr string
		if target == "blank" {
			titleStr = ctx.T("show in a new window")
		} else {
			titleStr = ctx.T("show in the same window")
		}

		buf.Add("<td title=\"%s\">%s</td>", titleStr, target)

		if updateRight || deleteRight {
			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			if updateRight {
				urlStr := ctx.U("/social_update", "socialId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit record."), ctx.T("Edit"))
			}

			if deleteRight {
				urlStr := ctx.U("/social_delete", "socialId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))
			}

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
