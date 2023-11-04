package machine_image_lib

import (
	"database/sql"
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
)

type MachineImageRec struct {
	MachineImageId int64
	MachineId      int64
	Enum           int64
	IsHeader       string
	Alt            sql.NullString
	IsCropped      string
	ImgName        string
	ImgMime        string
	ImgSize        int64
	ImgHeight      int64
	ImgWidth       int64
	ImgData        string
}

func GetMachineImageRec(imgId int64) (*MachineImageRec, error) {
	sqlStr := `select
					machineImage.machineImageId,
					machineImage.machineId,
					machineImage.enum,
					machineImage.isHeader,
					machineImage.alt,
					machineImage.isCropped,
					machineImage.imgName,
					machineImage.imgMime,
					machineImage.imgSize,
					machineImage.imgHeight,
					machineImage.imgWidth,
					machineImage.imgData
				from
					machineImage,
					machine
				where
					machine.machineId = machineImage.machineId and
					machineImage.machineImageId = ?`

	row := app.Db.QueryRow(sqlStr, imgId)

	rec := new(MachineImageRec)
	err := row.Scan(&rec.MachineImageId, &rec.MachineId, &rec.Enum,
		&rec.IsHeader, &rec.Alt, &rec.IsCropped, &rec.ImgName, &rec.ImgMime, &rec.ImgSize,
		&rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetMachineImageList(machineId int64) []*MachineImageRec {
	sqlStr := `select
					machineImage.machineImageId,
					machineImage.machineId,
					machineImage.enum,
					machineImage.isHeader,
					machineImage.alt,
					machineImage.isCropped,
					machineImage.imgName,
					machineImage.imgMime,
					machineImage.imgSize,
					machineImage.imgHeight,
					machineImage.imgWidth
				from
					machineImage,
					machine
				where
					machineImage.machineId = machine.machineId and
					machineImage.machineId = ?
				order by
					machineImage.isHeader, machineImage.enum`

	rows, err := app.Db.Query(sqlStr, machineId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*MachineImageRec, 0, 100)
	for rows.Next() {
		rec := new(MachineImageRec)
		err = rows.Scan(&rec.MachineImageId, &rec.MachineId, &rec.Enum,
			&rec.IsHeader, &rec.Alt, &rec.IsCropped, &rec.ImgName, &rec.ImgMime,
			&rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth)
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

func IsCroppedToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "yes":
		return "<i class=\"fas fa-check\"></i>"
	case "no":
		return "<i class=\"fas fa-times\"></i>"
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
