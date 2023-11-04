package category

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/combo"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("category", "insert") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("categoryId", -1)
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	parentIdStr := ctx.Req.PostFormValue("parentId")
	if parentIdStr == "root" {
		parentIdStr = "1"
	}

	parentId, err := strconv.ParseInt(parentIdStr, 10, 64)
	if err != nil {
		parentId = -1
	}

	name := ctx.Req.PostFormValue("name")
	status := ctx.Req.PostFormValue("status")

	enum, err := strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
	if err != nil {
		enum = 0
	}

	if parentId <= 0 || name == "" || status == "default" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					category(categoryId, parentId, name, status, enum)
					values(null, ?, ?, ?, ?)`

	res, err := tx.Exec(sqlStr, parentId, name, status, enum)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				insertForm(ctx)
				return
			} else if err.Number == 1452 {
				ctx.Msg.Warning(ctx.T("Could not find parent record."))
				insertForm(ctx)
				return
			}
		}

		panic(err)
	}

	lid, err := res.LastInsertId()
	if err != nil {
		ctx.Msg.Error(ctx.T("Could not get last insert id."))
		ctx.Msg.Error(err.Error())

		insertForm(ctx)
		return
	}

	tx.Commit()

	ctx.Cargo.SetInt("categoryId", lid)
	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/category"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)

	var parentId, enum int64
	var name, status string
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
		name = ""
		status = "default"
		parentId = -1
		enum = 0
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
	buf.Add(content.PageTitle(ctx.T("Machine Categories"), ctx.T("New Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/category")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/category_insert")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup parentId\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Parent Category:"))
	buf.Add("<select name=\"parentId\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")

	buf.Add(parentCombo.Format(fmt.Sprintf("%s", parentId)))

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
