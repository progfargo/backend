package manufacturer_lib

import (
	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
)

type ManufacturerRec struct {
	ManufacturerId int64
	Name           string
}

func GetManufacturerRec(manufacturerId int64) (*ManufacturerRec, error) {
	sqlStr := `select
					manufacturerId,
					name
				from
					manufacturer
				where
					manufacturerId = ?`

	row := app.Db.QueryRow(sqlStr, manufacturerId)

	rec := new(ManufacturerRec)
	err := row.Scan(&rec.ManufacturerId, &rec.Name)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountManufacturer(key string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add("select count(*) from manufacturer")

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (name like('%%%s%%')", key)
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetManufacturerPage(ctx *context.Ctx, key string, pageNo int64) []*ManufacturerRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					manufacturerId,
					name
				from
					manufacturer`)

	if key != "" {
		key = util.DbStr(key)
		sqlBuf.Add("where (name like('%%%s%%'))", key)
	}

	sqlBuf.Add("order by name")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*ManufacturerRec, 0, 100)
	for rows.Next() {
		rec := new(ManufacturerRec)
		if err = rows.Scan(&rec.ManufacturerId, &rec.Name); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
