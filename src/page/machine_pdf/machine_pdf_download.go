package machine_pdf

import (
	"fmt"
	"net/http"

	"backend/src/lib/context"
	"backend/src/page/machine_pdf/machine_pdf_lib"
)

func Download(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine_image", "browse") {
		http.Error(rw, ctx.T("bad request."), http.StatusNotFound)
		return
	}

	ctx.Cargo.AddInt("pdfId", -1)
	ctx.ReadCargo()

	pdfId := ctx.Cargo.Int("pdfId")
	rec, err := machine_pdf_lib.GetMachinePdfRec(pdfId)
	if err != nil {
		http.Error(rw, ctx.T("Pdf not found."), http.StatusNotFound)
		return
	}

	if !ctx.IsRight("machine", "browse") {
		http.Error(rw, ctx.T("You don't have right to access this record."), 401)
		return
	}

	data := []byte(rec.PdfData)

	ctx.Rw.Header().Set("Content-Type", rec.PdfMime)
	ctx.Rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	ctx.Rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", rec.PdfName))

	if _, err := ctx.Rw.Write(data); err != nil {
		http.Error(ctx.Rw, ctx.T("Could not write pdf to response writer."), http.StatusNotFound)
	}
}
