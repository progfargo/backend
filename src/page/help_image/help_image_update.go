package help_image

import (
	"bytes"
	"database/sql"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/help/help_lib"
	"backend/src/page/help_image/help_image_lib"

	"github.com/go-sql-driver/mysql"
)

func UpdateImage(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("help", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("helpId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.ReadCargo()

	helpId := ctx.Cargo.Int("helpId")
	helpRec, err := help_lib.GetHelpRec(helpId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Help record could not be found."))
			ctx.Redirect(ctx.U("/help_image", "helpId", "key", "pn"))
			return
		}

		panic(err)
	}

	imgId := ctx.Cargo.Int("imgId")
	imgRec, err := help_image_lib.GetHelpImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Help image record could not be found."))
			ctx.Redirect(ctx.U("/help_image", "helpId", "key", "pn"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/help_image", "helpId", "key", "pn"))
		return
	}

	err = ctx.Req.ParseMultipartForm(1024 * 1024 * 2)
	if err != nil {
		panic(err)
	}

	mpForm := ctx.Req.MultipartForm
	image := mpForm.File["image"]
	if len(image) == 0 {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	imgFile := image[0]
	imgName := imgFile.Filename
	imgMime := imgFile.Header["Content-Type"][0]

	ext := filepath.Ext(imgName)
	isJpeg, err := regexp.MatchString("(?i)^\\.jpeg$", ext)
	isJpg, err := regexp.MatchString("(?i)^\\.jpg$", ext)
	if !isJpeg && !isJpg {
		ctx.Msg.Warning(ctx.T("You can only send jpeg images."))
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	inFile, err := imgFile.Open()
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	defer inFile.Close()

	imgData := bytes.NewBuffer(nil)
	io.Copy(imgData, inFile)
	imgSize := imgData.Len()

	if imgSize == 0 {
		ctx.Msg.Warning(ctx.T("You have sent a file with zero size."))
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	imgConfig, err := jpeg.DecodeConfig(imgData)
	if err != nil {
		ctx.Msg.Warning(ctx.T("The image is not a valid jpeg file."))
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	imgData = bytes.NewBuffer(nil)
	inFile.Seek(0, os.SEEK_SET)
	io.Copy(imgData, inFile)

	tx, err = app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update helpImage set
					imgName = ?,
					imgMime = ?,
					imgSize = ?,
					imgHeight = ?,
					imgWidth = ?,
					imgData = ?
				where
					helpImageId = ?`

	res, err := tx.Exec(sqlStr, imgName, imgMime, imgSize, imgConfig.Height, imgConfig.Width, imgData.String(), imgId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				updateImageForm(ctx, helpRec, imgRec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateImageForm(ctx, helpRec, imgRec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/help_image", "helpId", "key", "pn"))
}

func updateImageForm(ctx *context.Ctx, helpRec *help_lib.HelpRec, imgRec *help_image_lib.HelpImageRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Help"), ctx.T("Images"), ctx.T("Edit")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/help_image", "helpId", "key", "pn")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/help_image_small", "imgId")
	buf.Add("<p><img src=\"%s\" alt=\"\" title=\"%s\"></p>", urlStr,
		util.ScrStr(imgRec.ImgName))

	urlStr = ctx.U("/help_image_update_image", "helpId", "key", "pn", "imgId")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Image:"))
	buf.Add("<input type=\"file\" name=\"image\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")
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
	lmenu.Set(ctx)

	tmenu := top_menu.New()
	tmenu.Set(ctx, "help")

	ctx.Render("default.html")
}
