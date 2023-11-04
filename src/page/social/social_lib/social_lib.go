package social_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
)

type SocialRec struct {
	SocialId int64
	Link     string
	Title    string
	Target   string
	Icon     string
}

func GetSocialRec(socialId int64) (*SocialRec, error) {
	sqlStr := `select
					*
				from
					social
				where
					socialId = ?`

	row := app.Db.QueryRow(sqlStr, socialId)

	rv := new(SocialRec)
	err := row.Scan(&rv.SocialId, &rv.Link, &rv.Title, &rv.Target, &rv.Icon)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func GetSocialList(ctx *context.Ctx) []*SocialRec {

	sortOrder := ctx.Cargo.Str("sortOrder")

	orderStr := ""
	if sortOrder != "" {
		orderStr = "order by title"
	}

	sqlStr := fmt.Sprintf(`select
					*
				from
					social
				%s`, orderStr)

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*SocialRec, 0, 100)
	for rows.Next() {
		rec := new(SocialRec)
		if err = rows.Scan(&rec.SocialId, &rec.Link, &rec.Title, &rec.Target, &rec.Icon); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
