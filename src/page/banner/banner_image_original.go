package banner

import (
	"fmt"
	"net/http"

	"backend/src/lib/context"
	"backend/src/page/banner/banner_lib"
)

func ImageOriginal(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	ctx.Cargo.AddInt("bannerId", -1)
	ctx.ReadCargo()

	bannerId := ctx.Cargo.Int("bannerId")
	rec, err := banner_lib.GetBannerRec(bannerId)
	if err != nil {
		http.Error(rw, ctx.T("Image not found."), 404)
		return
	}

	ctx.Rw.Header().Set("Content-Type", rec.ImgMime)
	ctx.Rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(rec.ImgData)))
	ctx.Rw.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", rec.ImgName))

	ctx.Rw.Write([]byte(rec.ImgData))
}
