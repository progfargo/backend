package machine_image

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
	"backend/src/page/machine/machine_lib"

	"github.com/go-sql-driver/mysql"
)

func Insert(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)
	ctx.ReadCargo()

	machineId := ctx.Cargo.Int("machineId")
	machineRec, err := machine_lib.GetMachineRec(machineId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		http.Error(rw, ctx.T("You don't have right to access this record."), 401)
		return
	}

	if ctx.Req.Method == "GET" {
		insertForm(ctx, machineRec)
		return
	}

	err = ctx.Req.ParseMultipartForm(1024 * 1024 * 4)
	if err != nil {
		ctx.Msg.Error(err.Error())
		insertForm(ctx, machineRec)
		return
	}

	mpForm := ctx.Req.MultipartForm
	image := mpForm.File["image"]
	if len(image) == 0 {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx, machineRec)
		return
	}

	var insertCount = 0

	for i := 0; i < len(image); i++ {
		imgFile := image[i]
		imgName := imgFile.Filename
		imgMime := imgFile.Header["Content-Type"][0]

		ext := filepath.Ext(imgName)
		isJpeg, err := regexp.MatchString("(?i)^\\.jpeg$", ext)
		isJpg, err := regexp.MatchString("(?i)^\\.jpg$", ext)
		if !isJpeg && !isJpg {
			ctx.Msg.Warning(ctx.T("You can only send jpeg images."))
			insertForm(ctx, machineRec)
			return
		}

		inFile, err := imgFile.Open()
		if err != nil {
			ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
			insertForm(ctx, machineRec)
			return
		}

		defer inFile.Close()

		imgData := bytes.NewBuffer(nil)
		io.Copy(imgData, inFile)

		imgSize := imgData.Len()

		if imgSize == 0 {
			ctx.Msg.Warning(ctx.T("You have sent a file with zero size."))
			insertForm(ctx, machineRec)
			return
		}

		imgConfig, err := jpeg.DecodeConfig(imgData)
		if err != nil {
			ctx.Msg.Warning(ctx.T("The image is not a valid jpeg file."))
			insertForm(ctx, machineRec)
			return
		}

		if imgConfig.Width > imgConfig.Height && (imgConfig.Width < 1280 || imgConfig.Height < 800) {
			ctx.Msg.Warning(ctx.T("Minimum image size for a landscape image:") + "1280px*800px")
			continue
		}

		if imgConfig.Width < imgConfig.Height && (imgConfig.Width < 800 || imgConfig.Height < 1280) {
			ctx.Msg.Warning(ctx.T("Minimum image size for a portrait image:") + "800px*1280px")
			continue
		}

		imgData = bytes.NewBuffer(nil)
		inFile.Seek(0, os.SEEK_SET)
		io.Copy(imgData, inFile)

		tx, err := app.Db.Begin()
		if err != nil {
			panic(err)
		}

		sqlStr := `insert into
						machineImage(machineImageId, machineId,
						imgName, imgMime, imgSize, imgHeight, imgWidth, imgData)
						values(null, ?, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(sqlStr, machineId, imgName, imgMime, imgSize, imgConfig.Height,
			imgConfig.Width, imgData.String())
		if err != nil {
			tx.Rollback()
			if err, ok := err.(*mysql.MySQLError); ok {
				if err.Number == 1062 {
					ctx.Msg.Warning(ctx.T("Duplicate record."))
					insertForm(ctx, machineRec)
					return
				} else if err.Number == 1452 {
					ctx.Msg.Warning(ctx.T("Could not find parent record."))
					insertForm(ctx, machineRec)
					return
				}
			}

			panic(err)
		}

		tx.Commit()

		insertCount++
	}

	if insertCount > 0 {
		ctx.Msg.Success(ctx.T("Records has been saved."))
	}

	ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
}

func insertForm(ctx *context.Ctx, machineRec *machine_lib.MachineRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("New Image"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_image_insert", "machineId", "key", "stat", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Image:"))
	buf.Add("<input type=\"file\" name=\"image\" class=\"formControl\"" +
		" tabindex=\"1\" multiple=\"multiple\" autofocus>")
	buf.Add("<span class=\"helpBlock\">%s</span>", ctx.T("Image type must be jpeg."))
	buf.Add("<span class=\"helpBlock\">%s 1280px/800px or 800px/1280px %s</span>", ctx.T("Minimum image size:"), ctx.T("(w/h)"))
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
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
