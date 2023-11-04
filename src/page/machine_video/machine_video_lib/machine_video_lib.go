package machine_video_lib

import (
	"backend/src/app"
)

type MachineVideoRec struct {
	MachineVideoId int64
	MachineId      int64
	Title          string
	Url            string
}

func GetMachineVideoRec(urlId int64) (*MachineVideoRec, error) {
	sqlStr := `select
					machineVideoId,
					machineId,
					title,
					url
				from
					machineVideo
				where
					machineVideoId = ?`

	row := app.Db.QueryRow(sqlStr, urlId)

	rec := new(MachineVideoRec)
	err := row.Scan(&rec.MachineVideoId, &rec.MachineId, &rec.Title, &rec.Url)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetMachineVideoList(machineId int64) []*MachineVideoRec {
	sqlStr := `select
					machineVideoId,
					machineId,
					title,
					url
				from
					machineVideo
				where
					machineId = ?
				order by url`

	rows, err := app.Db.Query(sqlStr, machineId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*MachineVideoRec, 0, 100)
	for rows.Next() {
		rec := new(MachineVideoRec)
		if err = rows.Scan(&rec.MachineVideoId, &rec.MachineId, &rec.Title, &rec.Url); err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}
