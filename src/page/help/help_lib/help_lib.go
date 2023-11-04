package help_lib

import (
	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

type HelpRec struct {
	HelpId  int64
	Title   string
	Summary string
	Body    string
}

func GetHelpRec(helpId int64) (*HelpRec, error) {
	sqlStr := `select
					helpId,
					title,
					summary,
					body
				from
					help
				where
					helpId = ?`

	row := app.Db.QueryRow(sqlStr, helpId)

	rec := new(HelpRec)
	err := row.Scan(&rec.HelpId, &rec.Title, &rec.Summary, &rec.Body)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountHelp(key string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*) from help")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add(`where (title like('%%%s%%') or
			summary like ('%%%s%%'))`, key, key)
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetHelpPage(ctx *context.Ctx, key string, pageNo int64) []*HelpRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					helpId,
					title,
					summary,
					body
				 from
					help`)

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add(`where (title like('%%%s%%') or
				summary like ('%%%s%%'))`, key, key)
	}

	sqlBuf.Add("order by title")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*HelpRec, 0, 100)
	for rows.Next() {
		rec := new(HelpRec)
		if err = rows.Scan(&rec.HelpId, &rec.Title, &rec.Summary, &rec.Body); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
