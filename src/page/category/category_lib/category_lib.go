package category_lib

import (
	"database/sql"
	"fmt"
	"sort"

	"backend/src/app"
	"backend/src/content/combo"
	"backend/src/lib/context"
	"backend/src/lib/tax"
	"backend/src/lib/util"
)

type CategoryRec struct {
	CategoryId int64
	ParentId   int64
	Name       string
	Status     string
	Enum       int64
	IsAdded    bool
}

//to sort the category dir for the machine category page.
type CategoryRecList []*CategoryRec

func (crl CategoryRecList) Len() int {
	return len(crl)
}

func (crl CategoryRecList) Less(i, j int) bool {
	return crl[i].Name < crl[j].Name
}

func (crl CategoryRecList) Swap(i, j int) {
	crl[i], crl[j] = crl[j], crl[i]
}

func GetCategoryRec(categoryId int64) (*CategoryRec, error) {
	sqlStr := `select
					category.CategoryId,
					category.parentId,
					category.name,
					category.status,
					category.enum
				from
					category
				where
					category.categoryId = ?`

	row := app.Db.QueryRow(sqlStr, categoryId)

	rec := new(CategoryRec)
	err := row.Scan(&rec.CategoryId, &rec.ParentId, &rec.Name, &rec.Status, &rec.Enum)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

type CategoryItem struct {
	CategoryName string
	Status       string
	IsActive     bool
}

type CategoryList struct {
	Tax *tax.Tax
}

func NewCategoryList() *CategoryList {
	rv := new(CategoryList)
	rv.Tax = tax.New()
	return rv
}

func (pgl *CategoryList) add(name, parent string, enum int64, isVisible bool, categoryName, status string) {
	pgl.Tax.Add(name, parent, enum, isVisible,
		&CategoryItem{categoryName, status, false})
}

func (pgl *CategoryList) SetActive(name string) {
	allParents := pgl.Tax.GetAllParents(name)

	for _, val := range allParents {
		if name == val {
			item := pgl.Tax.GetItem(val)
			item.Data.(*CategoryItem).IsActive = true
		}
	}
}

func (pgl *CategoryList) SetWithTx(ctx *context.Ctx, tx *sql.Tx) {
	sqlBuf := util.NewBuf()

	sqlBuf.Add(`select
					categoryId,
					parentId,
					name,
					status,
					enum
				from
					category`)

	rows, err := tx.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	list := make(map[int64]*CategoryRec, 10)

	for rows.Next() {
		rec := new(CategoryRec)
		err = rows.Scan(&rec.CategoryId, &rec.ParentId, &rec.Name, &rec.Status, &rec.Enum)
		if err != nil {
			panic(err)
		}

		list[rec.CategoryId] = rec
	}

	for k, _ := range list {
		pgl.addItemToTax(ctx, list, k)
	}

	pgl.Tax.SortChildren()
}

func (pgl *CategoryList) Set(ctx *context.Ctx) {
	sqlBuf := util.NewBuf()

	sqlBuf.Add(`select
					categoryId,
					parentId,
					name,
					status,
					enum
				from
					category`)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	list := make(map[int64]*CategoryRec, 10)

	for rows.Next() {
		rec := new(CategoryRec)
		err = rows.Scan(&rec.CategoryId, &rec.ParentId, &rec.Name, &rec.Status, &rec.Enum)
		if err != nil {
			panic(err)
		}

		list[rec.CategoryId] = rec
	}

	for k, _ := range list {
		pgl.addItemToTax(ctx, list, k)
	}

	pgl.Tax.SortChildren()
}

func (pgl *CategoryList) addItemToTax(ctx *context.Ctx, list map[int64]*CategoryRec, item int64) {

	if item == 0 || list[item].IsAdded {
		return
	}

	if !pgl.Tax.IsExists(tax.IntToName(list[item].ParentId)) {
		pgl.addItemToTax(ctx, list, list[item].ParentId)
	}

	var name = list[item].Name
	var status = list[item].Status

	pgl.add(tax.IntToName(item), tax.IntToName(list[item].ParentId), list[item].Enum, true, name, status)
	list[item].IsAdded = true
}

func GetCategoryByMachine(machineId int64, categoryCombo *combo.TaxCombo) ([]*CategoryRec, error) {

	sqlStr := `select
					category.categoryId,
					category.name
				from
					category,
					machineCategory
				where
					category.categoryId = machineCategory.categoryId and
					machineCategory.machineId = ?
				order by
					category.name`

	rows, err := app.Db.Query(sqlStr, machineId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rv := make(CategoryRecList, 0, 10)
	for rows.Next() {
		rec := new(CategoryRec)
		if err = rows.Scan(&rec.CategoryId, &rec.Name); err != nil {
			return nil, err
		}

		rv = append(rv, rec)
	}

	for i := 0; i < len(rv); i++ {
		rv[i].Name = categoryCombo.FormatDir(fmt.Sprintf("%d", rv[i].CategoryId))
	}

	sort.Sort(rv)

	return rv, nil
}

func StatusToLabel(ctx *context.Ctx, str string) string {
	switch str {
	case "active":
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">%s</span>", ctx.T("Active"))
	case "passive":
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">%s</span>", ctx.T("Passive"))
	default:
		return fmt.Sprintf("<span class=\"label labelError labelXs\">%s</span>", ctx.T("Unknown"))
	}
}
