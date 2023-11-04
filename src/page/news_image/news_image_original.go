package news_image

import (
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/page/news_image/news_image_lib"
)

func ImageOriginal(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("news", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	imgId := ctx.Cargo.Int("imgId")
	rec, err := news_image_lib.GetNewsImageRec(imgId)
	if err != nil {
		http.Error(rw, ctx.T("Image not found."), 404)
		return
	}

	ctx.Rw.Header().Set("Content-Type", rec.ImgMime)
	ctx.Rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(rec.ImgData)))
	ctx.Rw.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", rec.ImgName))

	ctx.Rw.Write([]byte(rec.ImgData))
}
