package category

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/combo"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/category/category_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("category", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("categoryId", -1)
	ctx.ReadCargo()

	categoryId := ctx.Cargo.Int("categoryId")
	rec, err := category_lib.GetCategoryRec(categoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/category"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateForm(ctx, rec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/category"))
		return
	}

	name := ctx.Req.PostFormValue("name")
	status := ctx.Req.PostFormValue("status")

	parentIdStr := ctx.Req.PostFormValue("parentId")
	if parentIdStr == "root" {
		parentIdStr = "1"
	}

	parentId, err := strconv.ParseInt(parentIdStr, 10, 64)
	if err != nil {
		parentId = -1
	}

	enum, err := strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
	if err != nil {
		enum = 0
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	pgl := category_lib.NewCategoryList()
	pgl.SetWithTx(ctx, tx)

	parentList := pgl.Tax.GetAllParents(fmt.Sprintf("%d", rec.CategoryId))

	sqlStr := "select * from category where categoryId in (?) for update"
	_, err = tx.Exec(sqlStr, strings.Join(parentList, ","))
	if err != nil {
		tx.Rollback()
		ctx.Msg.Warning("Could not lock category table.")
		ctx.Redirect(ctx.U("/category"))
		return
	}

	if name == "" || parentId <= 0 || status == "default" {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateForm(ctx, rec)
		return
	}

	children := pgl.Tax.GetAllChildren(fmt.Sprintf("%d", rec.CategoryId))
	for _, v := range children {
		if parentIdStr == v {
			tx.Rollback()
			ctx.Msg.Warning(ctx.T("This update can not be done. It creates a loop in categories."))
			updateForm(ctx, rec)
			return
		}
	}

	sqlStr = `update category set
						parentId = ?,
						name = ?,
						status = ?,
						enum = ?
					where
						categoryId = ?`

	res, err := tx.Exec(sqlStr, parentId, name, status, enum, categoryId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				updateForm(ctx, rec)
				return
			} else if err.Number == 1452 {
				ctx.Msg.Warning(ctx.T("Could not find parent record."))
				updateForm(ctx, rec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been changed."))
	ctx.Redirect(ctx.U("/category"))
}

func updateForm(ctx *context.Ctx, rec *category_lib.CategoryRec) {
	content.Include(ctx)

	var name, status string
	var parentId, enum int64
	var err error
	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		status = ctx.Req.PostFormValue("status")

		parentId, err = strconv.ParseInt(ctx.Req.PostFormValue("parentId"), 10, 64)
		if err != nil {
			parentId = -1
		}

		enum, err = strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
		if err != nil {
			enum = 0
		}
	} else {
		name = rec.Name
		status = rec.Status
		parentId = rec.ParentId
		enum = rec.Enum
	}

	parentCombo := combo.NewTaxCombo(`select
											categoryId,
											parentId,
											name,
											enum
										from
											category`, ctx.T("Main Category"))
	parentCombo.Set(ctx)

	statusCombo := combo.NewEnumCombo()
	statusCombo.Add("default", ctx.T("Select Status"))
	statusCombo.Add("active", ctx.T("Active"))
	statusCombo.Add("passive", ctx.T("Passive"))

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")

	buf.Add(content.PageTitle(ctx.T("Machine Categories"), util.ScrStr(name), ctx.T("Edit")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/category")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/category_update", "categoryId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Parent Group:"))
	buf.Add("<select name=\"parentId\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")

	buf.Add(parentCombo.Format(fmt.Sprintf("%d", parentId)))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Name:"))
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\">", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Status:"))
	buf.Add("<select name=\"status\" class=\"formControl\" tabindex=\"3\">")

	buf.Add(statusCombo.Format(status))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("Enumeration:"))
	buf.Add("<input type=\"text\" name=\"enum\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"11\" tabindex=\"4\">", enum)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"5\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"6\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "category")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
