package banner_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
)

type BannerRec struct {
	BannerId  int64
	Enum      int64
	Status    string
	ImgName   string
	ImgMime   string
	ImgSize   int64
	ImgHeight int64
	ImgWidth  int64
	ImgData   string
}

func GetBannerRec(bannerId int64) (*BannerRec, error) {
	sqlStr := fmt.Sprintf(`select
					banner.bannerId,
					banner.enum,
					banner.status,
					banner.imgName,
					banner.imgMime,
					banner.imgSize,
					banner.imgHeight,
					banner.imgWidth,
					banner.imgData
				from
					banner
				where
					bannerId = ?`)

	row := app.Db.QueryRow(sqlStr, bannerId)

	rec := new(BannerRec)
	err := row.Scan(&rec.BannerId, &rec.Enum, &rec.Status,
		&rec.ImgName, &rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth, &rec.ImgData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetBannerList() []*BannerRec {
	sqlStr := fmt.Sprintf(`select
								banner.bannerId,
								banner.enum,
								banner.status,
								banner.imgName,
								banner.imgMime,
								banner.imgSize,
								banner.imgHeight,
								banner.imgWidth
							from
								banner
							order by enum`)

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*BannerRec, 0, 10)
	for rows.Next() {
		rec := new(BannerRec)
		err := rows.Scan(&rec.BannerId, &rec.Enum, &rec.Status,
			&rec.ImgName, &rec.ImgMime, &rec.ImgSize, &rec.ImgHeight, &rec.ImgWidth)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetBannerIdList() []*BannerRec {
	sqlStr := `select
					banner.bannerId
				from
					banner
				where
					status = 'active'
				order by enum`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*BannerRec, 0, 10)
	for rows.Next() {
		rec := new(BannerRec)
		err := rows.Scan(&rec.BannerId)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func StatusToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "active":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Active"))
	case "passive":
		return fmt.Sprintf("<span class=\"label labelError labelXs\">%s</span>", ctx.T("Passive"))
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
