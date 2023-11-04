package machine_video

import (
	"database/sql"
	"net/http"
	"strings"

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
			ctx.Msg.Warning(ctx.T("Machine record could not be found."))
			ctx.Redirect(ctx.U("/machine_video", "machineId", "key", "stat", "pn", "catId", "manId"))
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

	title := ctx.Req.PostFormValue("title")
	url := ctx.Req.PostFormValue("url")

	if title == "" || url == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		insertForm(ctx, machineRec)
		return
	}

	if !strings.HasPrefix(url, "https://www.youtube.com/watch?v=") {
		ctx.Msg.Warning(ctx.T("Invalid video url."))
		insertForm(ctx, machineRec)
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `insert into
					machineVideo(machineVideoId, machineId, title, url)
					values(null, ?, ?, ?)`

	_, err = tx.Exec(sqlStr, machineId, title, url)
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

	ctx.Msg.Success(ctx.T("Record has been saved."))
	ctx.Redirect(ctx.U("/machine_video", "machineId", "key", "stat", "pn", "catId", "manId"))
}

func insertForm(ctx *context.Ctx, machineRec *machine_lib.MachineRec) {
	content.Include(ctx)

	var title, url string
	if ctx.Req.Method == "POST" {
		title = ctx.Req.PostFormValue("title")
		url = ctx.Req.PostFormValue("url")
	} else {
		title = ""
		url = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Cnc Machines"), ctx.T("New Video Url"), machineRec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/machine_video", "machineId", "key", "stat", "pn", "catId", "manId")))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/machine_video_insert", "machineId", "key", "stat", "pn", "catId", "manId")
	buf.Add("<form action=\"%s\" method=\"post\">", urlStr)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Title:"))
	buf.Add("<input type=\"text\" name=\"title\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"200\" tabindex=\"1\" autofocus>", util.ScrStr(title))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Video Url:"))
	buf.Add("<input type=\"text\" name=\"url\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"200\" tabindex=\"2\" autofocus>", util.ScrStr(url))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"3\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"4\">%s</button>", ctx.T("Reset"))
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
