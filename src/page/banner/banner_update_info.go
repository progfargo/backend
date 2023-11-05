package banner

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/banner/banner_lib"
)

func UpdateInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("banner", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("bannerId", -1)
	ctx.ReadCargo()

	bannerId := ctx.Cargo.Int("bannerId")
	rec, err := banner_lib.GetBannerRec(bannerId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Image record could not be found."))
			ctx.Redirect(ctx.U("/banner"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateInfoForm(ctx, rec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/banner"))
		return
	}

	status := ctx.Req.PostFormValue("status")

	enum, err := strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
	if err != nil {
		enum = 0
	}

	if status == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateInfoForm(ctx, rec)
		return
	}

	sqlStr := fmt.Sprintf(`update banner set
					status = ?,
					enum = ?
				where
					banner.bannerId = ?`)

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	res, err := tx.Exec(sqlStr, status, enum, bannerId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateInfoForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been changed."))
	ctx.Redirect(ctx.U("/banner"))
}

func updateInfoForm(ctx *context.Ctx, rec *banner_lib.BannerRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	var status string
	var enum int64
	var err error
	if ctx.Req.Method == "POST" {
		status = ctx.Req.PostFormValue("status")

		enum, err = strconv.ParseInt(ctx.Req.PostFormValue("enum"), 10, 64)
		if err != nil {
			enum = 0
		}
	} else {
		status = rec.Status
		enum = rec.Enum
	}

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Banner"), ctx.T("Update Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/banner")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/banner_update_info", "bannerId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	var activeChecked, passiveChecked string
	if status == "" || status == "active" {
		activeChecked = " checked"
	} else {
		passiveChecked = " checked"
	}

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Status:"))

	buf.Add("<span>")
	buf.Add("<input type=\"radio\" name=\"status\" value=\"active\" tabindex=\"1\" autofocus %s>"+
		"<span class=\"radioLabel\" title=\"%s\">active</span>", activeChecked, ctx.T("Make active."))
	buf.Add("</span>")

	buf.Add("<span>")
	buf.Add("<input type=\"radio\" name=\"status\" value=\"passive\" tabindex=\"2\"%s>"+
		"<span class=\"radioLabel\" title=\"%s\">passive</span>", passiveChecked, ctx.T("Make passive."))
	buf.Add("</span>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label title=\"Enumeration\">%s</label>", ctx.T("Enum:"))
	buf.Add("<input type=\"text\" name=\"enum\" class=\"formControl\""+
		" value=\"%d\" maxlength=\"10\" tabindex=\"3\">", enum)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"4\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"5\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "banner")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
