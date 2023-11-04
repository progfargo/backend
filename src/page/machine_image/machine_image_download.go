package machine_image

import (
	"fmt"
	"net/http"

	"backend/src/lib/context"
	"backend/src/page/machine_image/machine_image_lib"
)

func Download(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine_image", "browse") {
		http.Error(rw, ctx.T("bad request."), http.StatusNotFound)
		return
	}

	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	imgId := ctx.Cargo.Int("imgId")
	rec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		http.Error(rw, ctx.T("image not found."), http.StatusNotFound)
		return
	}

	if !ctx.IsRight("machine", "browse") {
		http.Error(rw, ctx.T("You don't have right to access this record."), 401)
		return
	}

	data := []byte(rec.ImgData)

	ctx.Rw.Header().Set("Content-Type", rec.ImgMime)
	ctx.Rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	ctx.Rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", rec.ImgName))

	if _, err := ctx.Rw.Write(data); err != nil {
		http.Error(ctx.Rw, ctx.T("Could not write image to response writer."), http.StatusNotFound)
	}
}
