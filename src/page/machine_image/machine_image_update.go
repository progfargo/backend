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
	"backend/src/page/machine_image/machine_image_lib"

	"github.com/go-sql-driver/mysql"
)

func UpdateImage(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddInt("imgId", -1)
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
			ctx.Msg.Warning(ctx.T("Machine record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if !ctx.IsRight("machine", "update") {
		ctx.Msg.Error(ctx.T("You don't have right to access this record."))
		ctx.Redirect(ctx.U("/machine", "key", "stat", "pn", "catId", "manId"))
		return
	}

	imgId := ctx.Cargo.Int("imgId")
	imgRec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Machine image record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
			return
		}

		panic(err)
	}

	if ctx.Req.Method == "GET" {
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/machine_image", "machineId"))
		return
	}

	err = ctx.Req.ParseMultipartForm(1024 * 1024 * 4)
	if err != nil {
		ctx.Msg.Error(err.Error())
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	mpForm := ctx.Req.MultipartForm
	image := mpForm.File["image"]
	if len(image) == 0 {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		updateImageForm(ctx, machineRec, imgRec)
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
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	inFile, err := imgFile.Open()
	if err != nil {
		ctx.Msg.Warning(ctx.T("Can not open uploaded file."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	defer inFile.Close()

	imgData := bytes.NewBuffer(nil)
	io.Copy(imgData, inFile)

	imgSize := imgData.Len()

	if imgSize == 0 { //calculate this later
		ctx.Msg.Warning(ctx.T("You have sent a file with zero size."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	imgConfig, err := jpeg.DecodeConfig(imgData)
	if err != nil {
		ctx.Msg.Warning(ctx.T("The image is not a valid jpeg file."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	if imgConfig.Width > imgConfig.Height && (imgConfig.Width < 1280 || imgConfig.Height < 800) {
		ctx.Msg.Warning(ctx.T("Minimum image size for a landscape image:") + "1280px*800px")
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	if imgConfig.Width < imgConfig.Height && (imgConfig.Width < 800 || imgConfig.Height < 1280) {
		ctx.Msg.Warning(ctx.T("Minimum image size for a portrait image:") + "800px*1280px")
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	imgData = bytes.NewBuffer(nil)
	inFile.Seek(0, os.SEEK_SET)
	io.Copy(imgData, inFile)

	tx, err = app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update machineImage set
					isCropped = ?,
					imgName = ?,
					imgMime = ?,
					imgSize = ?,
					imgHeight = ?,
					imgWidth = ?,
					imgData = ?
				where
					machineImageId = ?`

	res, err := tx.Exec(sqlStr, "no", imgName, imgMime, imgSize, imgConfig.Height, imgConfig.Width, imgData.String(), imgId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				updateImageForm(ctx, machineRec, imgRec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		updateImageForm(ctx, machineRec, imgRec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId"))
}

func updateImageForm(ctx *context.Ctx, machineRec *machine_lib.MachineRec, imgRec *machine_image_lib.MachineImageRec) {
	content.Include(ctx)
	ctx.Css.Add("/asset/jquery_modal/jquery.modal.css")
	ctx.Js.Add("/asset/jquery_modal/jquery.modal.js")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines"), ctx.T("Update Image"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_image", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	smallImageUrl := ctx.U("/machine_image_small", "imgId")
	popupUrl := ctx.U("/machine_popup", "imgId")

	buf.Add("<p>")
	buf.Add("<img class=\"smallImage\" src=\"%s\" alt=\"\">", smallImageUrl)
	buf.Add("<a class=\"popupImage\" href=\"%s\"></a>", popupUrl)
	buf.Add("</p>")

	urlStr := ctx.U("/machine_image_update", "machineId", "imgId", "key", "stat", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\" enctype=\"multipart/form-data\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Image:"))
	buf.Add("<input type=\"file\" name=\"image\" class=\"formControl\"" +
		" tabindex=\"1\" autofocus>")
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

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	ctx.Render("default.html")
}
