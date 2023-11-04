package machine_pdf

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/machine_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
	"backend/src/page/machine_pdf/machine_pdf_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)
	ctx.Cargo.AddInt("pdfId", -1)
	ctx.ReadCargo()

	machineId := ctx.Cargo.Int("machineId")
	machineRec, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		ctx.Msg.Error(ctx.T("You don't have right to access this record."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	browseMid(ctx, machineRec)

	content.Default(ctx)
	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/machine_pdf/machine_pdf.js")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx, machineRec *machine_lib.MachineRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("PDF"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col lg2\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine", "key", "stat", "pn", "catId", "manId")))
	buf.Add(content.NewButton(ctx, ctx.U("/machine_pdf_insert", "machineId", "key", "pn", "stat", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	machineMenu := machine_menu.New("machineId", "key", "stat", "pn", "catId", "manId")
	machineMenu.Set(ctx, "machine_pdf")

	buf.Add("<div class=\"col lg10\">")
	buf.Add(machineMenu.Format(ctx))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedZero\" title=\"%s\">%s</th>", ctx.T("Enumeration"), ctx.T("E."))
	buf.Add("<th>%s</th>", ctx.T("Title"))
	buf.Add("<th>%s</th>", ctx.T("PDF Info"))

	buf.Add("<th class=\"right\">%s</th>", ctx.T("Command"))

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	machinePdfList := machine_pdf_lib.GetMachinePdfList(machineRec.MachineId)
	if len(machinePdfList) == 0 {
		ctx.Msg.Warning("Empty list.")
	}

	var title, pdfUrl, urlStr string
	for _, row := range machinePdfList {
		ctx.Cargo.SetInt("pdfId", row.MachinePdfId)

		buf.Add("<tr>")

		buf.Add("<td>%d</td>", row.Enum)

		title = util.ScrStr(row.Title)
		if title == "" {
			title = "not entered."
		}

		buf.Add("<td>%s</td>", title)

		pdfUrl = ctx.U("/machine_pdf_download", "pdfId")
		buf.Add("<td>")
		buf.Add("<div class=\"pdfInfo\">")

		padSize := 6

		buf.Add("<strong>%s:</strong><a href=\"%s\" target=\"_blank\">%s</a>", util.PadRight(ctx.T("Url"), "&nbsp;", padSize), pdfUrl, pdfUrl)

		buf.Add("<span class=\"copyLink\" title=\"%s\"><i class=\"fas fa-copy\"></i></span>", ctx.T("Copy link url."))

		urlStr = ctx.U("/machine_pdf_download", "pdfId")
		buf.Add("<a class=\"downloadLink\" href=\"%s\" title=\"%s\"><i class=\"fas fa-file-download\"></i></a><br>", urlStr, ctx.T("Download"))

		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Name"), "&nbsp;", padSize), util.ScrStr(row.PdfName))
		buf.Add("<strong>%s:</strong> %s bytes<br>", util.PadRight(ctx.T("Size"), "&nbsp;", padSize), util.FormatInt(row.PdfSize))
		buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Mime"), "&nbsp;", padSize), util.ScrStr(row.PdfMime))

		buf.Add("</div>")
		buf.Add("</td>")

		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFlex\">")

		urlStr = ctx.U("/machine_pdf_update", "machineId", "pdfId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit pdf record."), ctx.T("Edit Pdf"))

		urlStr = ctx.U("/machine_pdf_update_info", "machineId", "pdfId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Edit pdf info."), ctx.T("Edit Pdf Info"))

		urlStr = ctx.U("/machine_pdf_delete", "machineId", "pdfId", "key", "stat", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Delete record."), ctx.T("Delete"))

		buf.Add("</div>")
		buf.Add("</td>")
		buf.Add("</tr>")
	}

	buf.Add("</tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
