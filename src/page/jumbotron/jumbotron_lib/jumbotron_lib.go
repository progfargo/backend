package jumbotron_lib

import (
	"backend/src/app"
	"backend/src/lib/util"
)

type JumbotronRec struct {
	JumbotronId int64
	Title       string
	Body        string
}

func GetJumbotronRec(jumbotronId int64) (*JumbotronRec, error) {
	sqlStr := `select
					jumbotronId,
					title,
					body
				from
					jumbotron
				where
					jumbotronId = ?`

	row := app.Db.QueryRow(sqlStr, jumbotronId)

	rec := new(JumbotronRec)
	err := row.Scan(&rec.JumbotronId, &rec.Title, &rec.Body)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetJumbotronList() []*JumbotronRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					jumbotronId,
					title,
					body
				from
					jumbotron`)

	sqlBuf.Add("order by title")

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*JumbotronRec, 0, 100)
	for rows.Next() {
		rec := new(JumbotronRec)
		if err = rows.Scan(&rec.JumbotronId, &rec.Title, &rec.Body); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
