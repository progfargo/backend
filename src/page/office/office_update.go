package office

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/office/office_lib"

	"github.com/go-sql-driver/mysql"
)

func Update(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("office", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("productionId", -1)
	ctx.ReadCargo()

	productionId := ctx.Cargo.Int("productionId")
	rec, err := office_lib.GetOfficeRec(productionId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/office"))
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
		ctx.Redirect(ctx.U("/office"))
		return
	}

	name := ctx.Req.PostFormValue("name")
	city := ctx.Req.PostFormValue("city")
	address := ctx.Req.PostFormValue("address")
	telephone := ctx.Req.PostFormValue("telephone")
	email := ctx.Req.PostFormValue("email")
	mapLink := ctx.Req.PostFormValue("mapLink")

	if name == "" || city == "" || address == "" || telephone == "" || email == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateForm(ctx, rec)
		return
	}

	err = util.IsValidEmail(ctx, email)
	if err != nil {
		ctx.Msg.Warning(err.Error())
		updateForm(ctx, rec)
		return
	}

	fontFile := app.Ini.HomeDir + "/asset/font/watermark_font.ttf"
	emailBuf, err := util.EmailToPng(email, 300, 40, fontFile, 19, 0, 0, 0, 0)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not convert email address to a png image."))
		ctx.Msg.Error(err.Error())
		insertForm(ctx)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update office set
					name = ?,
					city = ?,
					address = ?,
					telephone = ?,
					email = ?,
					imgName = ?,
					imgMime = ?,
					imgSize = ?,
					imgData = ?,
					mapLink = ?
				where
					officeId = ?`

	res, err := tx.Exec(sqlStr, name, city, address, telephone, email,
		email, "image/png", emailBuf.Len(), emailBuf.Bytes(), mapLink, productionId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
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
	ctx.Redirect(ctx.U("/office"))
}

func updateForm(ctx *context.Ctx, rec *office_lib.OfficeRec) {
	content.Include(ctx)

	var name, city, address, telephone, email, mapLink string
	if ctx.Req.Method == "POST" {
		name = ctx.Req.PostFormValue("name")
		city = ctx.Req.PostFormValue("city")
		address = ctx.Req.PostFormValue("address")
		telephone = ctx.Req.PostFormValue("telephone")
		email = ctx.Req.PostFormValue("email")
		mapLink = ctx.Req.PostFormValue("mapLink")
	} else {
		name = rec.Name
		city = rec.City
		address = rec.Address
		telephone = rec.Telephone
		email = rec.Email
		mapLink = rec.MapLink
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Offices"), ctx.T("Update Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/office")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/office_update", "productionId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Name:"))
	buf.Add("<input type=\"text\" name=\"name\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"50\" tabindex=\"1\" autofocus>", util.ScrStr(name))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("City:"))
	buf.Add("<input type=\"text\" name=\"city\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"50\" tabindex=\"2\">", util.ScrStr(city))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Address:"))
	buf.Add("<textarea name=\"address\"  class=\"formControl\" rows=\"7\""+
		" tabindex=\"3\">%s</textarea>", util.ScrStr(address))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Telephone:"))
	buf.Add("<input type=\"text\" name=\"telephone\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"20\" tabindex=\"4\">", util.ScrStr(telephone))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Email:"))
	buf.Add("<input type=\"text\" name=\"email\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"50\" tabindex=\"5\">", util.ScrStr(email))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Map Link:"))
	buf.Add("<input type=\"text\" name=\"mapLink\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"200\" tabindex=\"6\">", util.ScrStr(mapLink))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"7\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"8\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "office")

	tmenu := top_menu.New()
	tmenu.Set(ctx, "root")

	ctx.Render("default.html")
}
