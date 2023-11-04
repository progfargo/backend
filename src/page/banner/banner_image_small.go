package banner

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strings"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/page/banner/banner_lib"

	"github.com/disintegration/gift"
)

func ImageSmall(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("banner", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("bannerId", -1)
	ctx.ReadCargo()

	bannerId := ctx.Cargo.Int("bannerId")
	rec, err := banner_lib.GetBannerRec(bannerId)
	if err != nil {
		http.Error(rw, ctx.T("Image not found."), 404)
		return
	}

	srcImage, _, err := image.Decode(strings.NewReader(rec.ImgData))
	if err != nil {
		http.Error(rw, err.Error(), 204)
		return
	}

	attr, _, err := image.DecodeConfig(strings.NewReader(rec.ImgData))
	if err != nil {
		http.Error(rw, err.Error(), 204)
		return
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	if attr.Width > 400 {
		g := gift.New(gift.Resize(400, 0, gift.CubicResampling))
		dstImage := image.NewNRGBA(g.Bounds(srcImage.Bounds()))
		g.Draw(dstImage, srcImage)

		err = jpeg.Encode(writer, dstImage, &jpeg.Options{90})
		if err != nil {
			http.Error(rw, err.Error(), 204)
			return
		}

		rec.ImgData = buf.String()
	}

	ctx.Rw.Header().Set("Content-Type", "image/jpeg")
	ctx.Rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(rec.ImgData)))
	ctx.Rw.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", rec.ImgName))

	ctx.Rw.Write(buf.Bytes())
}
