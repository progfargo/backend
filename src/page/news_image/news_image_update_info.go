package news_image

import (
	"database/sql"
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/news/news_lib"
	"backend/src/page/news_image/news_image_lib"
)

func UpdateImageInfo(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("news", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("newsId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddStr("stat", "")
	ctx.ReadCargo()

	newsId := ctx.Cargo.Int("newsId")
	_, err := news_lib.GetNewsRec(newsId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("News record could not be found."))
			ctx.Redirect(ctx.U("/news_image", "newsId", "key", "pn", "stat"))
			return
		}

		panic(err)
	}

	imgId := ctx.Cargo.Int("imgId")
	imgRec, err := news_image_lib.GetNewsImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("News image record could not be found."))
			ctx.Redirect(ctx.U("/news_image", "newsId", "key", "pn", "stat"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateImageInfoForm(ctx, imgRec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/news_image", "newsId", "key", "pn", "stat"))
		return
	}

	alt := ctx.Req.PostFormValue("alt")

	if alt == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateImageInfoForm(ctx, imgRec)
		return
	}

	sqlStr := fmt.Sprintf(`update newsImage set
							alt%s = ?
						where
							newsImageId = ?`, ctx.L())

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	res, err := tx.Exec(sqlStr, alt, imgId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateImageInfoForm(ctx, imgRec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/news_image", "newsId", "key", "pn", "stat"))
}

func updateImageInfoForm(ctx *context.Ctx, imgRec *news_image_lib.NewsImageRec) {
	content.Include(ctx)

	buf := util.NewBuf()

	var alt string
	if ctx.Req.Method == "POST" {
		alt = ctx.Req.PostFormValue("alt")
	} else {
		alt = util.NullToString(imgRec.Alt)
	}

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("News"), ctx.T("Images"), ctx.T("Update Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/news_image", "newsId", "key", "pn", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/news_image_update_info", "newsId", "key", "pn", "stat", "imgId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Alternative Text:"))
	buf.Add("<input type=\"text\" name=\"alt\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"250\" tabindex=\"1\" autofocus>", util.ScrStr(alt))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"2\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"3\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)
	content.Include(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "news")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
