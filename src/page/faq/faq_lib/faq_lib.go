package faq_lib

import (
	"backend/src/app"
		"backend/src/lib/context"
	"backend/src/lib/util"
)

type FaqRec struct {
	FaqId    int64
	Question string
	Summary  string
	Answer   string
}

func GetFaqRec(faqId int64) (*FaqRec, error) {
	sqlStr := `select
					faqId,
					question,
					summary,
					answer
				from
					faq
				where
					faqId = ?`

	row := app.Db.QueryRow(sqlStr, faqId)

	rec := new(FaqRec)
	err := row.Scan(&rec.FaqId, &rec.Question, &rec.Summary, &rec.Answer)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

	
func CountFaq(key string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*) from faq")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (question like('%%%s%%') or summary like ('%%%s%%'))", key, key)
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}


func GetFaqPage(ctx *context.Ctx, key string, pageNo int64) []*FaqRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					faqId,
					question,
					summary,
					answer
				from
					faq`)

		if key != "" {
			key = util.DbStr(key)
			sqlBuf.Add("where (question like('%%%s%%') or summary like ('%%%s%%'))", key, key)
		}

		sqlBuf.Add("order by question")

		pageLen := ctx.Config.Int("pageLen")
		start := (pageNo - 1) * pageLen
		sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*FaqRec, 0, 100)
	for rows.Next() {
		rec := new(FaqRec)
		if err = rows.Scan(&rec.FaqId, &rec.Question, &rec.Summary, &rec.Answer); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}