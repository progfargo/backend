package text_content_image_lib

import (
	"database/sql"

	"backend/src/app"
)

type TextContentImageRec struct {
	TextContentImageId int64
	TextContentId      int64
	Alt                sql.NullString
	ImgName            string
	ImgMime            string
	ImgSize            int64
	ImgHeight          int64
	ImgWidth           int64
	ImgData            string
}

func GetTextContentImageRec(imgId int64) (*TextContentImageRec, error) {
	sqlStr := `select
					textContentImage.textContentImageId,
					textContentImage.textContentId,
					textContentImage.alt,
					textContentImage.imgName,
					textContentImage.imgMime,
					textContentImage.imgSize,
					textContentImage.imgHeight,
					textContentImage.imgWidth,
					textContentImage.imgData
				from
					textContentImage
				where
					textContentImageId = ?`

	row := app.Db.QueryRow(sqlStr, imgId)

	rec := new(TextContentImageRec)
	err := row.Scan(&rec.TextContentImageId, &rec.TextContentId, &rec.Alt, &rec.ImgName,
		&rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetTextContentImageList(textId int64) []*TextContentImageRec {
	sqlStr := `select
					textContentImage.textContentImageId,
					textContentImage.textContentId,
					textContentImage.alt,
					textContentImage.imgName,
					textContentImage.imgMime,
					textContentImage.imgSize,
					textContentImage.imgHeight,
					textContentImage.imgWidth,
					textContentImage.imgData
				from
					textContentImage
				where
					textContentImage.textContentId = ?
				order by
					textContentImage.imgType`

	rows, err := app.Db.Query(sqlStr, textId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*TextContentImageRec, 0, 100)
	for rows.Next() {
		rec := new(TextContentImageRec)
		err = rows.Scan(&rec.TextContentImageId, &rec.TextContentId, &rec.Alt, &rec.ImgName,
			&rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
