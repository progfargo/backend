package site_error_lib

import (
	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

type SiteErrorRec struct {
	SiteErrorId int64
	DateTime    int64
	Title       string
	Message     string
}

func GetSiteErrorRec(errorId int64) (*SiteErrorRec, error) {
	sqlStr := `select
					*
				from
					siteError
				where
					siteErrorId = ?`

	row := app.Db.QueryRow(sqlStr, errorId)

	rec := new(SiteErrorRec)
	err := row.Scan(&rec.SiteErrorId, &rec.DateTime, &rec.Title, &rec.Message)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountSiteError(key string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*) from siteError")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (en like('%%%s%%') or gr like ('%%%s%%'))", key, key)
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetSiteErrorPage(ctx *context.Ctx, key string, pageNo int64) []*SiteErrorRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select siteErrorId, dateTime, title, message from siteError")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (title like('%%%s%%') or message like ('%%%s%%'))", key, key)
	}

	sqlBuf.Add("order by dateTime")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*SiteErrorRec, 0, 100)
	for rows.Next() {
		rec := new(SiteErrorRec)
		if err = rows.Scan(&rec.SiteErrorId, &rec.DateTime, &rec.Title, &rec.Message); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
