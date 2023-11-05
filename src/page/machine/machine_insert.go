package machine

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	if !ctx.IsRight("machine", "insert") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "default")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)
	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	categoryIdStr := ctx.Req.PostFormValue("categoryId")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		categoryId = -1
	}

	manufacturerId, err := strconv.ParseInt(ctx.Req.PostFormValue("manufacturerId"), 10, 64)
	if err != nil {
		manufacturerId = -1
	}

	name := ctx.Req.PostFormValue("name")
	model := ctx.Req.PostFormValue("model")
	exp := ctx.Req.PostFormValue("exp")
	location := ctx.Req.PostFormValue("location")

	yom, err := strconv.ParseInt(ctx.Req.PostFormValue("yom"), 10, 64)
	if err != nil {
		yom = -1
	}

	status := ctx.Req.PostFormValue("status")
	priceStr := ctx.Req.PostFormValue("price")

	if categoryId <= 1 || manufacturerId == -1 || name == "" || model == "" || yom == -1 || status == "default" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx)
		return
	}

	if yom < app.MIN_YOM || yom > int64(time.Now().Year()) {
		ctx.Msg.Warning(ctx.T("Invalid year of manufacture."))
		insertForm(ctx)
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil || price < 0 {
		ctx.Msg.Warning("Invalid price.")
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert machine(machineId, categoryId, manufacturerId, name, model, exp, location, yom, status, price)
					values(null, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.Exec(sqlStr, categoryId, manufacturerId, name, model, exp, location, yom, status, price)

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
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("Could not read last insert id."))
		insertForm(ctx)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))

	ctx.Cargo.SetInt("machineId", lid)
	ctx.Redirect(ctx.U("/machine_update", "key", "stat", "pn", "catId", "machineId", "manId"))
}

func insertForm(ctx *context.Ctx) {
	content.Include(ctx)
	ctx.Js.Add("/asset/tinymce/tinymce.min.js")
	ctx.Js.Add("/asset/tinymce/tinymce_func.js")
	ctx.Js.Add("/asset/js/page/machine/machine_insert.js")

	stat := ctx.Cargo.Str("stat")

	var categoryId, manufacturerId, yom int64
	var name, model, exp, location, status string
	var price int64
	var err error

	if ctx.Req.Method == "POST" {
		categoryId, err = strconv.ParseInt(ctx.Req.PostFormValue("categoryId"), 10, 64)
		if err != nil {
			categoryId = -1
		}

		manufacturerId, err = strconv.ParseInt(ctx.Req.PostFormValue("manufacturerId"), 10, 64)
		if err != nil {
			manufacturerId = -1
		}

		name = ctx.Req.PostFormValue("name")
		model = ctx.Req.PostFormValue("model")
		exp = ctx.Req.PostFormValue("exp")
		location = ctx.Req.PostFormValue("location")

		yom, err = strconv.ParseInt(ctx.Req.PostFormValue("yom"), 10, 64)
		if err != nil {
			yom = -1
		}

		status = ctx.Req.PostFormValue("status")

		price, err = strconv.ParseInt(ctx.Req.PostFormValue("price"), 10, 64)
		if err != nil {
			price = 0
		}
	} else {
		categoryId = ctx.Cargo.Int("catId")
		manufacturerId = ctx.Cargo.Int("manId")
		name = ""
		model = ""
		exp = ""
		location = ""
		yom = 0
		status = stat
		price = 0
	}

	categoryCombo := combo.NewTaxCombo(`select
											categoryId,
											parentId,
											name,
											enum
										from
											category`, ctx.T("Main Category"))
	categoryCombo.Set(ctx)

	if categoryCombo.IsEmpty() {
		ctx.Msg.Warning(ctx.T("Category list is empty. You should enter at least one category first."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	manufacturerCombo := combo.NewCombo(`select
											manufacturerId,
											name
										from
											manufacturer
										order by name`, ctx.T("Manufacturer"))
	manufacturerCombo.Set()

	if manufacturerCombo.IsEmpty() {
		ctx.Msg.Warning(ctx.T("Manufacturer list is empty. You should enter at least one manufacturer first."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	statusCombo := combo.NewEnumCombo()
	statusCombo.Add("default", ctx.T("Select Status"))
	statusCombo.Add("pending", ctx.T("Pending"))
	statusCombo.Add("active", ctx.T("Active"))
	statusCombo.Add("passive", ctx.T("Passive"))
	statusCombo.Add("sold", ctx.T("Sold"))

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("New Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_insert", "key", "pn", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Category:"))
	buf.Add("<select name=\"categoryId\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")

	buf.Add(categoryCombo.Format(fmt.Sprintf("%d", categoryId)))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Name:"))
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"2\">", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Manufacturer:"))
	buf.Add("<select name=\"manufacturerId\" class=\"formControl\"" +
		" tabindex=\"3\">")

	buf.Add(manufacturerCombo.Format(manufacturerId))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Model:"))
	buf.Add("<input type=\"text\" name=\"model\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"4\">", util.ScrStr(model))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("Explanation:"))
	buf.Add("<textarea name=\"exp\" id=\"exp\" class=\"formControl\""+
		" tabindex=\"5\">%s</textarea>", exp)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("Location:"))
	buf.Add("<input type=\"text\" name=\"location\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"6\">", util.ScrStr(location))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Year of Manufacture:"))
	buf.Add("<input type=\"text\" name=\"yom\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"4\" tabindex=\"7\">", yom)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Status:"))
	buf.Add("<select name=\"status\" class=\"formControl\" tabindex=\"8\">")

	buf.Add(statusCombo.Format(status))

	buf.Add("</select>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label>%s</label>", ctx.T("Price:"))
	buf.Add("<input type=\"text\" name=\"price\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"10\" tabindex=\"9\">", price)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"10\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"11\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "machineInsertPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}
