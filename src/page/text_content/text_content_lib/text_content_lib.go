package text_content_lib

import (
	"backend/src/app"
)

type TextContentRec struct {
	TextContentId int64
	Name          string
	Exp           string
	Body          string
}

func GetTextContentRec(textId int64) (*TextContentRec, error) {
	sqlStr := `select
					textContentId,
					name,
					exp,
					body
				from
					textContent
				where
					textContentId = ?`

	row := app.Db.QueryRow(sqlStr, textId)

	rec := new(TextContentRec)
	err := row.Scan(&rec.TextContentId, &rec.Name, &rec.Exp, &rec.Body)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetTextContentRecByName(name string) (*TextContentRec, error) {
	sqlStr := `select
					textContentId,
					name,
					exp,
					body
				from
					textContent
				where
					name = ?`

	row := app.Db.QueryRow(sqlStr, name)

	rec := new(TextContentRec)
	err := row.Scan(&rec.TextContentId, &rec.Name, &rec.Exp, &rec.Body)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetTextContentList() []*TextContentRec {
	sqlStr := `select
					textContentId,
					name,
					exp,
					body
				 from
					textContent
				order by name`

	rows, err := app.Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*TextContentRec, 0, 100)
	for rows.Next() {
		rec := new(TextContentRec)
		if err = rows.Scan(&rec.TextContentId, &rec.Name, &rec.Exp, &rec.Body); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
