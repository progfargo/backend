package machine

import (
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

func Popup(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	displayPopup(ctx)
}

func displayPopup(ctx *context.Ctx) {
	content.Include(ctx)

	buf := util.NewBuf()

	urlStr := ctx.U("/machine_image_middle", "imgId")
	buf.Add("<div>")
	buf.Add("<img src=\"%s\" alt=\"\">", urlStr)
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	ctx.Render("machine_popup.html")
}
