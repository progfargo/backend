package tran_lib

import (
	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

type TranRec struct {
	TranId int64
	En     string
	Tr     string
}

func GetTranRec(tranId int64) (*TranRec, error) {
	sqlStr := `select
					*
				from
					tran
				where
					tranId = ?`

	row := app.Db.QueryRow(sqlStr, tranId)

	rec := new(TranRec)
	err := row.Scan(&rec.TranId, &rec.En, &rec.Tr)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountTran(key string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*) from tran")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (en like('%%%s%%') or tr like ('%%%s%%'))", key, key)
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetTranPage(ctx *context.Ctx, key string, pageNo int64) []*TranRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select tranId, en, tr from tran")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (en like('%%%s%%') or tr like ('%%%s%%'))", key, key)
	}

	sqlBuf.Add("order by en")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*TranRec, 0, 100)
	for rows.Next() {
		rec := new(TranRec)
		if err = rows.Scan(&rec.TranId, &rec.En, &rec.Tr); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetTranList() []*TranRec {
	rv := make([]*TranRec, 0, 400)

	sqlStr := `select
					en,
					tr
				from
					tran`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		rec := new(TranRec)
		err = rows.Scan(&rec.En, &rec.Tr)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetTranMap() map[string]*TranRec {
	rv := make(map[string]*TranRec, 400)
	sqlStr := `select
					en,
					tr
				from
					tran`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		rec := new(TranRec)
		err = rows.Scan(&rec.En, &rec.Tr)
		if err != nil {
			panic(err)
		}

		rv[rec.En] = rec
	}

	return rv
}
