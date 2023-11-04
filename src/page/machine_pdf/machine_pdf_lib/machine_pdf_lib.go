package machine_pdf_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
)

type MachinePdfRec struct {
	MachinePdfId int64
	MachineId    int64
	Title        string
	Enum         int64
	PdfName      string
	PdfMime      string
	PdfSize      int64
	PdfData      string
}

func GetMachinePdfRec(pdfId int64) (*MachinePdfRec, error) {
	sqlStr := `select
					machinePdf.machinePdfId,
					machinePdf.machineId,
					machinePdf.title,
					machinePdf.enum,
					machinePdf.pdfName,
					machinePdf.pdfMime,
					machinePdf.pdfSize,
					machinePdf.pdfData
				from
					machine,
					machinePdf
				where
					machine.machineId = machinePdf.machineId and
					machinePdf.machinePdfId = ?`

	row := app.Db.QueryRow(sqlStr, pdfId)

	rec := new(MachinePdfRec)
	err := row.Scan(&rec.MachinePdfId, &rec.MachineId, &rec.Title, &rec.Enum,
		&rec.PdfName, &rec.PdfMime, &rec.PdfSize, &rec.PdfData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetMachinePdfList(machineId int64) []*MachinePdfRec {
	sqlStr := `select
					machinePdf.machinePdfId,
					machinePdf.machineId,
					machinePdf.title,
					machinePdf.enum,
					machinePdf.pdfName,
					machinePdf.pdfMime,
					machinePdf.pdfSize
				from
					machine,
					machinePdf
				where
					machine.machineId = machinePdf.machineId and
					machinePdf.machineId = ?
				order by machinePdf.enum`

	rows, err := app.Db.Query(sqlStr, machineId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*MachinePdfRec, 0, 100)
	for rows.Next() {
		rec := new(MachinePdfRec)
		err = rows.Scan(&rec.MachinePdfId, &rec.MachineId,
			&rec.Title, &rec.Enum, &rec.PdfName, &rec.PdfMime, &rec.PdfSize)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func IsHeaderToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "yes":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Yes"))
	case "no":
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">%s</span>", ctx.T("No"))
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
