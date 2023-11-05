package machine_image

import (
	"bufio"
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine_image/machine_image_lib"

	"github.com/disintegration/gift"
	"github.com/go-sql-driver/mysql"
)

const LONGSIZE = 1280
const SHORTSIZE = 800

const MAXWIDTH = 1600
const MAXHEIGHT = 1000

func CropImage(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "update") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddInt("machineId", -1)
	ctx.ReadCargo()

	imgId := ctx.Cargo.Int("imgId")
	rec, err := machine_image_lib.GetMachineImageRec(imgId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Image record could not be found."))
			ctx.Redirect(ctx.U("/machine_image", "machineId"))
			return
		}

		panic(err)
	}

	var landscape bool
	if rec.ImgWidth > rec.ImgHeight {
		landscape = true
	}

	if ctx.Req.Method == "GET" {
		cropImageForm(ctx, rec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/machine_image", "machineId"))
		return
	}

	xStr := ctx.Req.PostFormValue("x[]")
	x, err := strconv.ParseInt(xStr, 10, 64)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Invalid left value for crop image."))
		cropImageForm(ctx, rec)
		return
	}

	yStr := ctx.Req.PostFormValue("y[]")
	y, err := strconv.ParseInt(yStr, 10, 64)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Invalid top value for crop image."))
		cropImageForm(ctx, rec)
		return
	}

	widthStr := ctx.Req.PostFormValue("width[]")
	width, err := strconv.ParseInt(widthStr, 10, 64)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Invalid width value for crop image."))
		cropImageForm(ctx, rec)
		return
	}

	heightStr := ctx.Req.PostFormValue("height[]")
	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		ctx.Msg.Warning(ctx.T("Invalid height value for crop image."))
		cropImageForm(ctx, rec)
		return
	}

	if landscape {
		if width < LONGSIZE || height < SHORTSIZE {
			ctx.Msg.Warning(ctx.T("Invalid width and height value for image crop."))
			cropImageForm(ctx, rec)
			return
		}
	} else {
		if width < SHORTSIZE || height < LONGSIZE {
			ctx.Msg.Warning(ctx.T("Invalid width and height value for image crop."))
			cropImageForm(ctx, rec)
			return
		}
	}

	srcImage, _, err := image.Decode(strings.NewReader(rec.ImgData))
	if err != nil {
		ctx.Msg.Error(err.Error())
		cropImageForm(ctx, rec)
		return
	}

	//crop
	g := gift.New(gift.Crop(image.Rect(int(x), int(y), int(width+x), int(height+y))))
	cropImage := image.NewNRGBA(g.Bounds(srcImage.Bounds()))
	g.Draw(cropImage, srcImage)

	var isResize bool
	if (landscape && (width > MAXWIDTH || height > MAXHEIGHT)) ||
		(!landscape && (width > MAXHEIGHT || height > MAXWIDTH)) {
		isResize = true
	}

	//resize if needed
	var resultImage *image.NRGBA

	if isResize {
		if landscape {
			g = gift.New(gift.ResizeToFit(MAXWIDTH, MAXHEIGHT, gift.CubicResampling))
		} else {
			g = gift.New(gift.ResizeToFit(MAXHEIGHT, MAXWIDTH, gift.CubicResampling))
		}

		resizeImage := image.NewNRGBA(g.Bounds(cropImage.Bounds()))
		g.Draw(resizeImage, cropImage)
		resultImage = resizeImage
	} else {
		resultImage = cropImage
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err = jpeg.Encode(writer, resultImage, &jpeg.Options{90})
	if err != nil {
		http.Error(rw, err.Error(), 204)
		return
	}

	imgConfig, err := jpeg.DecodeConfig(strings.NewReader(buf.String()))
	if err != nil {
		ctx.Msg.Warning(ctx.T("The image is not a valid jpeg file."))
		http.Error(rw, err.Error(), 204)
		return
	}

	//save data
	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `update machineImage set
					isCropped = ?,
					imgSize = ?,
					imgHeight = ?,
					imgWidth = ?,
					imgData = ?
				where
					machineImageId = ?`

	res, err := tx.Exec(sqlStr, "yes", buf.Len(), imgConfig.Height,
		imgConfig.Width, buf.String(), imgId)
	if err != nil {
		tx.Rollback()
		if err, ok := err.(*mysql.MySQLError); ok {
			if err.Number == 1062 {
				ctx.Msg.Warning(ctx.T("Duplicate record."))
				cropImageForm(ctx, rec)
				return
			}
		}

		panic(err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		cropImageForm(ctx, rec)
		return
	}

	tx.Commit()

	ctx.Msg.Success(ctx.T("Record has been changed."))
	ctx.Redirect(ctx.U("/machine_image", "machineId"))
}

func cropImageForm(ctx *context.Ctx, imgRec *machine_image_lib.MachineImageRec) {
	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machine Image"), ctx.T("Update Image")))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_image", "machineId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_image_crop", "imgId", "machineId")
	buf.Add("<form id=\"cropForm\" action=\"%s\" method=\"post\">", urlStr)

	urlStr = ctx.U("/machine_image_middle", "imgId")
	buf.Add("<div class=\"image-wrapper\">")
	buf.Add("<img id=\"image\" src=\"%s\" alt=\"\" title=\"%s\" data-width=%d data-height=%d>",
		urlStr, util.ScrStr(imgRec.ImgName), imgRec.ImgWidth, imgRec.ImgHeight)
	buf.Add("<p class=\"imgInfo\">")
	buf.Add("<span class=\"infoTitle\">%s</span>", ctx.T("Crop Dimention:"))

	if imgRec.ImgWidth > imgRec.ImgHeight {
		buf.Add("<span class=\"dimention\"></span> %s: %dpx*%dpx", ctx.T("must be greater than"), LONGSIZE, SHORTSIZE)
	} else {
		buf.Add("<span class=\"dimention\"></span> %s: %dpx*%dpx", ctx.T("must be greater than"), SHORTSIZE, LONGSIZE)
	}

	buf.Add("</p>")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"2\">%s</button>", ctx.T("Crop"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)
	content.Include(ctx)

	ctx.Css.Add("/asset/rcrop/rcrop.css")
	ctx.Js.Add("/asset/rcrop/rcrop.js")
	ctx.Js.Add("/asset/js/page/machine_image/machine_image_crop.js")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "machineImageCropPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}
