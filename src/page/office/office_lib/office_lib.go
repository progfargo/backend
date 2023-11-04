package office_lib

import (
	"fmt"

	"backend/src/app"
	"backend/src/lib/context"
)

type OfficeRec struct {
	OfficeId     int64
	Name         string
	City         string
	Address      string
	Telephone    string
	Email        string
	MapLink      string
	IsMainOffice string
	ImgName      string
	ImgMime      string
	ImgSize      int64
	ImgData      string
}

func GetOfficeRec(officeId int64) (*OfficeRec, error) {
	sqlStr := `select
					officeId,
					name,
					city,
					address,
					telephone,
					email,
					mapLink,
					isMainOffice
				from
					office
				where
					officeId = ?`

	row := app.Db.QueryRow(sqlStr, officeId)

	rec := new(OfficeRec)
	err := row.Scan(&rec.OfficeId, &rec.Name, &rec.City, &rec.Address,
		&rec.Telephone, &rec.Email, &rec.MapLink, &rec.IsMainOffice)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetOfficeList() []*OfficeRec {
	sqlStr := `select
					officeId,
					name,
					city,
					address,
					telephone,
					email,
					isMainOffice
				from
					office
				order by name`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*OfficeRec, 0, 100)
	for rows.Next() {
		rec := new(OfficeRec)
		err := rows.Scan(&rec.OfficeId, &rec.Name, &rec.City, &rec.Address,
			&rec.Telephone, &rec.Email, &rec.IsMainOffice)

		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetMainOfficeRec() (*OfficeRec, error) {
	sqlStr := `select
					officeId,
					name,
					city,
					address,
					telephone,
					email,
					mapLink,
					imgName,
					imgMime,
					imgSize,
					imgData
				from
					office
				where
					isMainOffice = 'yes'`

	row := app.Db.QueryRow(sqlStr)

	rec := new(OfficeRec)
	err := row.Scan(&rec.OfficeId, &rec.Name, &rec.City, &rec.Address,
		&rec.Telephone, &rec.Email, &rec.MapLink, &rec.ImgName,
		&rec.ImgMime, &rec.ImgSize, &rec.ImgData)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func IsMainOfficeToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "yes":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Yes"))
	case "no":
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">%s</span>", ctx.T("No"))
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
