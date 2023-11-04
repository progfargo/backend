package login

import (
	"database/sql"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/banner/banner_lib"
	"backend/src/page/user/user_lib"

	"github.com/dchest/captcha"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if req.URL.Path != "/login" && req.URL.Path != "/" {
		app.NotFound()
	}

	ctx.ReadCargo()

	if ctx.IsLoggedIn() {
		ctx.Msg.Warning(ctx.T("You are already logged in."))
		ctx.Redirect(ctx.U("/welcome"))
		return
	}

	if ctx.Req.Method == "GET" {
		loginForm(ctx)
		return
	}

	login := ctx.Req.PostFormValue("login")
	password := ctx.Req.PostFormValue("password")
	captchaId := ctx.Req.PostFormValue("captchaId")
	captchaAnswer := ctx.Req.PostFormValue("captchaAnswer")

	if login == "" || password == "" {
		ctx.Msg.Warning(ctx.T("You have left one or more fields empty."))
		loginForm(ctx)
		return
	}

	if !captcha.VerifyString(captchaId, captchaAnswer) {
		ctx.Msg.Warning(ctx.T("Wrong security answer. Please try again."))
		loginForm(ctx)
		return
	}

	password = util.PasswordHash(password)
	rec, err := user_lib.GetUserRecByLogin(login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Wrong user name or password. Please try again."))
			loginForm(ctx)
			return
		}

		panic(err)
	}

	if rec.Status == "blocked" {
		ctx.Msg.Warning(ctx.T("Your account is not active. Please contact administrator to have your" +
			" account activated."))
		loginForm(ctx)
		return
	}

	ctx.Session.UserId = rec.UserId
	ctx.Session.UserName = rec.Name

	ctx.CreateSession()

	ctx.Redirect(ctx.U("/welcome"))
}

func loginForm(ctx *context.Ctx) {
	content.Include(ctx)
	ctx.Js.Add("/asset/js/page/login/login.js")
	ctx.Css.Add("/asset/css/page/login.css")

	var login, password string
	if ctx.Req.Method == "POST" {
		login = ctx.Req.PostFormValue("login")
		password = ctx.Req.PostFormValue("password")
	} else {
		login = ""
		password = ""
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row login\">")
	buf.Add("<div class=\"col colSpan md2 lg3 xl4\"></div>")

	buf.Add("<div class=\"col md8 lg6 xl4 panel panelPrimary\">")

	buf.Add("<div class=\"panelHeading\">")
	buf.Add("<h3>%s</h3>", ctx.T("Login"))
	buf.Add("</div>")

	buf.Add("<div class=\"panelBody\">")

	buf.Add("<form action=\"/login\" method=\"post\">")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Login Name:"))
	buf.Add("<input type=\"text\" name=\"login\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"100\" tabindex=\"1\" autofocus>", util.ScrStr(login))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Password:"))
	buf.Add("<input type=\"password\" name=\"password\" class=\"formControl\""+
		" value=\"%s\" maxlength=\"30\" tabindex=\"2\">", util.ScrStr(password))
	buf.Add("</div>")

	captchaId := captcha.NewLen(4)

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<label class=\"required\">%s</label>", ctx.T("Security Question:"))
	buf.Add("<input type=\"text\" name=\"captchaAnswer\"" +
		" class=\"formControl\" value=\"\" maxlength=\"10\" tabindex=\"3\">")
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<p><img src=\"/captcha/%s.png\" alt=\"%s\"></p>", captchaId, ctx.T("Captcha image."))
	buf.Add("<input type=hidden name=captchaId value=\"%s\">", captchaId)
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup formCommand\">")
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"4\">%s</button>", ctx.T("Submit"))
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"5\">%s</button>", ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<a href=\"/forgot_pass\">%s</a>", ctx.T("I forgot my password."))
	buf.Add("</div>")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<a href=\"/register\">%s</a>", ctx.T("Create account."))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>") //panelBody
	buf.Add("</div>") //col
	buf.Add("<div class=\"col colSpan md2 lg3 xl4\"></div>")
	buf.Add("</div>") //row

	backgroundList := banner_lib.GetBannerIdList()
	buf.Add("<div class=\"backgroundIdList\">")

	for _, v := range backgroundList {
		buf.Add("<span>/banner_image_original?bannerId=%d</span>", v.BannerId)
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	str := "loginPage"
	ctx.AddHtml("pageName", &str)

	content.Default(ctx)

	ctx.Render("login.html")
}
