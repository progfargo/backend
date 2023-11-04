package news_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

type NewsRec struct {
	NewsId     int64
	RecordDate int64
	Status     string
	Title      string
	Summary    string
	Body       string
}

func CountNews(key, status string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*) from news")

	conBuf := util.NewBuf()
	if key != "" {
		key = util.DbStr(key)
		conBuf.Add("(title like('%%%s%%') or summary like ('%%%s%%'))", key, key)
	}

	if status != "" {
		conBuf.Add("(status = '%s')", status)
	}

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where")
		sqlBuf.Add(*conBuf.StringSep("and"))
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetNewsRec(newsId int64) (*NewsRec, error) {
	sqlStr := `select
					newsId,
					recordDate,
					status,
					title,
					summary,
					body
				from
					news
				where
					newsId = ?`

	row := app.Db.QueryRow(sqlStr, newsId)

	rec := new(NewsRec)
	err := row.Scan(&rec.NewsId, &rec.RecordDate, &rec.Status, &rec.Title, &rec.Summary, &rec.Body)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetNewsPage(ctx *context.Ctx, key, stat string, pageNo int64) []*NewsRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					newsId,
					recordDate,
					status,
					title,
					summary,
					body
				 from
					news`)

	conBuf := util.NewBuf()

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add("(title like('%%%s%%') or summary like ('%%%s%%'))", key, key)
	}

	if stat != "" {
		stat = util.DbStr(stat)
		conBuf.Add("(status = '%s')", stat)
	}

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where " + *conBuf.StringSep("and"))
	}

	sqlBuf.Add("order by recordDate desc")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*NewsRec, 0, 100)
	for rows.Next() {
		rec := new(NewsRec)
		if err = rows.Scan(&rec.NewsId, &rec.RecordDate, &rec.Status, &rec.Title, &rec.Summary, &rec.Body); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func StatusToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "published":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Published"))
	case "draft":
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">%s</span>", ctx.T("Draft"))
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
