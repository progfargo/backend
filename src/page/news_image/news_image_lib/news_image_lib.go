package news_image_lib

import (
	"database/sql"

	"backend/src/app"
)

type NewsImageRec struct {
	NewsImageId int64
	NewsId      int64
	Alt         sql.NullString
	ImgName     string
	ImgMime     string
	ImgSize     int64
	ImgHeight   int64
	ImgWidth    int64
	ImgData     string
}

func GetNewsImageRec(imgId int64) (*NewsImageRec, error) {
	sqlStr := `select
					newsImage.newsImageId,
					newsImage.newsId,
					newsImage.alt,
					newsImage.imgName,
					newsImage.imgMime,
					newsImage.imgSize,
					newsImage.imgHeight,
					newsImage.imgWidth,
					newsImage.imgData
				from
					newsImage
				where
					newsImageId = ?`

	row := app.Db.QueryRow(sqlStr, imgId)

	rec := new(NewsImageRec)
	err := row.Scan(&rec.NewsImageId, &rec.NewsId, &rec.Alt, &rec.ImgName,
		&rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetNewsImageList(newsId int64) []*NewsImageRec {
	sqlStr := `select
					newsImage.newsImageId,
					newsImage.newsId,
					newsImage.alt,
					newsImage.imgName,
					newsImage.imgMime,
					newsImage.imgSize,
					newsImage.imgHeight,
					newsImage.imgWidth,
					newsImage.imgData
				from
					newsImage
				where
					newsImage.newsId = ?
				order by
					newsImage.imgType`

	rows, err := app.Db.Query(sqlStr, newsId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*NewsImageRec, 0, 100)
	for rows.Next() {
		rec := new(NewsImageRec)
		err = rows.Scan(&rec.NewsImageId, &rec.NewsId, &rec.Alt, &rec.ImgName,
			&rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
