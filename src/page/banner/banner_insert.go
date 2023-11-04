package banner

import (
	"bytes"
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

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("banner", "insert") {
		app.BadRequest()
	}

	ctx.ReadCargo()

	if ctx.Req.Method == "GET" {
		insertForm(ctx)
		return
	}

	err := ctx.Req.ParseMultipartForm(1024 * 1024 * 4)
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

	if imgConfig.Height < 900 || imgConfig.Width < 1600 {
		ctx.Msg.Warning(ctx.T("Minimum image size (w/h):") + " " + "1600/900")
		insertForm(ctx)
		return
	}

	imgData = bytes.NewBuffer(nil)
	inFile.Seek(0, os.SEEK_SET)
	io.Copy(imgData, inFile)

	sqlStr := `insert into
					banner(bannerId, enum, imgName, imgMime, imgSize, imgHeight, imgWidth, imgData)
					values(null, ?, ?, ?, ?, ?, ?, ?)`

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(sqlStr, 0, imgName, imgMime, imgSize, imgConfig.Height, imgConfig.Width, imgData.String())
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
	ctx.Redirect(ctx.U("/banner"))
}

func insertForm(ctx *context.Ctx) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Banner"), ctx.T("Delete Record")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/banner")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/banner_insert")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Image:"))
	buf.Add("<input type=\"file\" name=\"image\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")
	buf.Add("<span class=\"helpBlock\">%s 1600/900 %s</span>", ctx.T("Minimum image size:"), ctx.T("(w/h)"))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"2\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"3\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Include(ctx)
	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "banner")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
