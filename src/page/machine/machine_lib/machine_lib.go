package machine_lib

import (
	"database/sql"
	"fmt"
	"strings"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/category/category_lib"
)

type MachineRec struct {
	MachineId        int64
	CategoryId       int64
	ManufacturerId   int64
	ManufacturerName string
	MachineImageId   sql.NullInt64
	Name             string
	Model            string
	Exp              sql.NullString
	Location         string
	Yom              int64
	Status           string
	Price            int64
}

func GetMachineRec(machineId int64) (*MachineRec, error) {
	sqlStr := `select
					machine.machineId,
					machine.categoryId,
					machine.manufacturerId,
					manufacturer.name,
					machineImage.machineImageId,
					machine.name,
					machine.model,
					machine.exp,
					machine.location,
					machine.yom,
					machine.Status,
					machine.price
				from
					manufacturer,
					machine left join machineImage on
						machine.machineId = machineImage.machineId and
						machineImage.isHeader = 'yes'
				where
					machine.manufacturerId = manufacturer.manufacturerId and
					machine.machineId = ?`

	row := app.Db.QueryRow(sqlStr, machineId)

	rec := new(MachineRec)
	err := row.Scan(&rec.MachineId, &rec.CategoryId, &rec.ManufacturerId,
		&rec.ManufacturerName, &rec.MachineImageId, &rec.Name, &rec.Model, &rec.Exp,
		&rec.Location, &rec.Yom, &rec.Status, &rec.Price)

	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountMachine(ctx *context.Ctx, key, stat string, catId, manId int64) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					count(*)`)

	fromBuf := util.NewBuf()
	fromBuf.Add("machine")

	conBuf := util.NewBuf()

	if stat != "default" {
		stat = util.DbStr(stat)
		conBuf.Add("(machine.status='%s')", stat)
	}

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(machine.name like('%%%s%%'))`, key)
	}

	if catId >= 1 {
		catList := category_lib.NewCategoryList()
		catList.Set(ctx)

		childList := catList.Tax.GetAllChildren(fmt.Sprintf("%d", catId))
		if len(childList) > 0 {
			conBuf.Add("(machine.categoryId in (%s))", strings.Join(childList, ", "))
		}
	}

	if manId != -1 {
		conBuf.Add("(machine.manufacturerId = %d)", manId)
	}

	sqlBuf.Add("from")
	sqlBuf.Add(*fromBuf.StringSep(", "))

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

func GetMachinePage(ctx *context.Ctx, key, stat string, catId, manId, pageNo int64) []*MachineRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					distinct(machine.machineId),
					machine.categoryId,
					machine.manufacturerId,
					manufacturer.name,
					machineImage.machineImageId,
					machine.name,
					machine.model,
					machine.location,
					machine.yom,
					machine.status,
					machine.price`)

	fromBuf := util.NewBuf()
	fromBuf.Add("manufacturer")
	fromBuf.Add(`machine left outer join machineImage on machine.machineId = machineImage.machineId and
		machineImage.IsHeader = 'yes'`)

	conBuf := util.NewBuf()
	conBuf.Add("(machine.manufacturerId = manufacturer.manufacturerId)")
	//conBuf.Add("()")

	if stat != "default" {
		stat = util.DbStr(stat)
		conBuf.Add("(machine.status='%s')", stat)
	}

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(machine.name like('%%%s%%'))`, key)
	}

	if catId >= 1 {
		catList := category_lib.NewCategoryList()
		catList.Set(ctx)

		childList := catList.Tax.GetAllChildren(fmt.Sprintf("%d", catId))
		if len(childList) > 0 {
			conBuf.Add("(machine.categoryId in (%s))", strings.Join(childList, ", "))
		}
	}

	if manId != -1 {
		conBuf.Add("(machine.manufacturerId = %d)", manId)
	}

	sqlBuf.Add("from")
	sqlBuf.Add(*fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where")
		sqlBuf.Add(*conBuf.StringSep(" and "))
	}

	sqlBuf.Add("order by machineId desc")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*MachineRec, 0, 100)
	for rows.Next() {
		rec := new(MachineRec)
		err = rows.Scan(&rec.MachineId, &rec.CategoryId,
			&rec.ManufacturerId, &rec.ManufacturerName, &rec.MachineImageId,
			&rec.Name, &rec.Model, &rec.Location, &rec.Yom, &rec.Status, &rec.Price)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func StatusToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "pending":
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Pending"))
	case "active":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Active"))
	case "passive":
		return fmt.Sprintf("<span class=\"label labelError labelXs\">%s</span>", ctx.T("Passive"))
	case "sold":
		return fmt.Sprintf("<span class=\"label labelSold labelXs\">%s</span>", ctx.T("Sold"))
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
