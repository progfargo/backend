package text_content_image

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
	"backend/src/page/text_content/text_content_lib"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	updateRight := ctx.IsRight("text_content", "update")
	if !updateRight {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("textId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddStr("stat", "")
	ctx.ReadCargo()

	textContentId := ctx.Cargo.Int("textId")
	_, err := text_content_lib.GetTextContentRec(textContentId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/text_content_image", "textId", "key", "pn", "stat"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
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
		insertForm(ctx)
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
		insertForm(ctx)
		return
	}

	inFile, err := imgFile.Open()
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
		insertForm(ctx)
		return
	}

	defer inFile.Close()

	imgData := bytes.NewBuffer(nil)
	io.Copy(imgData, inFile)
	imgSize := imgData.Len()

	if imgSize == 0 {
		ctx.Msg.Warning(ctx.T("You have sent a file with zero size."))
		insertForm(ctx)
		return
	}

	imgConfig, err := jpeg.DecodeConfig(imgData)
	if err != nil {
		ctx.Msg.Warning(ctx.T("The image is not a valid jpeg file."))
		insertForm(ctx)
		return
	}

	imgData = bytes.NewBuffer(nil)
	inFile.Seek(0, os.SEEK_SET)
	io.Copy(imgData, inFile)

	tx, err = app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					textContentImage(textContentImageId, textContentId,
					imgType, imgName, imgMime, imgSize, imgHeight, imgWidth, imgData)
					values(null, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(sqlStr, textContentId, "alternative", imgName, imgMime, imgSize, imgConfig.Height,
		imgConfig.Width, imgData.String())
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

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/text_content_image", "textId", "key", "pn", "stat"))
}

func insertForm(ctx *context.Ctx) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Text Content"), ctx.T("Images"), ctx.T("New Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/text_content_image", "textId", "key", "pn", "stat")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/text_content_image_insert", "textId", "key", "pn", "stat")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Image:"))
	buf.Add("<input type=\"file\" name=\"image\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")
	buf.Add("<span class=\"helpBlock\">%s</span>", ctx.T("Image type must be jpeg."))
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
	lmenu.Set(ctx, "text_content")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
