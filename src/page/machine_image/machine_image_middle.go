package machine_image

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
	"backend/src/page/machine_image/machine_image_lib"

	"github.com/disintegration/gift"
)

func ImageMiddle(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	imgId := ctx.Cargo.Int("imgId")
	rec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		http.Error(rw, ctx.T("Image not found."), 404)
		return
	}

	if !ctx.IsRight("machine", "browse") {
		http.Error(rw, ctx.T("You don't have right to access this record."), 401)
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

	if attr.Width > 500 {
		g := gift.New(gift.Resize(500, 0, gift.CubicResampling))
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
