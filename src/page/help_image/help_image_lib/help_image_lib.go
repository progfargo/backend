package help_image_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
)

type HelpImageRec struct {
	HelpImageId int64
	HelpId      int64
	ImgName     string
	ImgMime     string
	ImgSize     int64
	ImgHeight   int64
	ImgWidth    int64
	ImgData     string
}

func GetHelpImageRec(imgId int64) (*HelpImageRec, error) {
	sqlStr := `select
					helpImage.helpImageId,
					helpImage.helpId,
					helpImage.imgName,
					helpImage.imgMime,
					helpImage.imgSize,
					helpImage.imgHeight,
					helpImage.imgWidth,
					helpImage.imgData
				from
					helpImage
				where
					helpImageId = ?`

	row := app.Db.QueryRow(sqlStr, imgId)

	rec := new(HelpImageRec)
	err := row.Scan(&rec.HelpImageId, &rec.HelpId, &rec.ImgName,
		&rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetHelpImageList(helpId int64) []*HelpImageRec {
	sqlStr := `select
					helpImage.helpImageId,
					helpImage.helpId,
					helpImage.imgName,
					helpImage.imgMime,
					helpImage.imgSize,
					helpImage.imgHeight,
					helpImage.imgWidth,
					helpImage.imgData
				from
					helpImage
				where
					helpImage.helpId = ?
				order by
					helpImage.imgType`

	rows, err := app.Db.Query(sqlStr, helpId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*HelpImageRec, 0, 100)
	for rows.Next() {
		rec := new(HelpImageRec)
		err = rows.Scan(&rec.HelpImageId, &rec.HelpId, &rec.ImgName,
			&rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func TypeToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "base":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Base"))
	case "alternative":
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">%s</span>", ctx.T("Alternative"))
	default:
		return fmt.Sprintf("<span class=\"label labelError labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
