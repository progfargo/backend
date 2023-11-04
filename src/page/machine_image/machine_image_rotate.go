package machine_image

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"strings"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/page/machine_image/machine_image_lib"

	"github.com/disintegration/gift"
)

func ImageRotate(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddInt("angle", 0)
	ctx.ReadCargo()

	imgId := ctx.Cargo.Int("imgId")
	rec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		ctx.Msg.Error(ctx.T("Image not found."))
		ctx.RenderAjaxError()
		return
	}

	if !ctx.IsRight("machine", "update") {
		http.Error(rw, ctx.T("You don't have right to access this record."), 401)
		return
	}

	angle := ctx.Cargo.Int("angle")
	if angle != 90 && angle != 270 {
		ctx.Msg.Error(ctx.T("Unsupported rotate angle."))
		ctx.RenderAjaxError()
		return
	}

	srcImage, _, err := image.Decode(strings.NewReader(rec.ImgData))
	if err != nil {
		http.Error(rw, err.Error(), 204)
		return
	}

	g := gift.New(gift.Rotate(float32(angle), color.Transparent, gift.CubicInterpolation))
	dstImage := image.NewNRGBA(g.Bounds(srcImage.Bounds()))
	g.Draw(dstImage, srcImage)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err = jpeg.Encode(writer, dstImage, &jpeg.Options{90})
	if err != nil {
		http.Error(rw, err.Error(), 204)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		ctx.Msg.Error(ctx.T(err.Error()))
		ctx.RenderAjaxError()
		return
	}

	sqlStr := `update machineImage set
					imgData = ?
				where
					machineImageId = ?`

	res, err := tx.Exec(sqlStr, buf.String(), imgId)
	if err != nil {
		tx.Rollback()
		ctx.Msg.Error(err.Error())
		ctx.RenderAjaxError()
		return
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Error(err.Error())
		ctx.RenderAjaxError()
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Image has been rotated."))
	ctx.RenderAjax()
}
