package context

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"backend/src/app"
	"backend/src/content/css"
	"backend/src/content/js"
	"backend/src/lib/cargo"
	"backend/src/lib/user"
)

type Ctx struct {
	Config    app.ConfigType
	Css       *css.Css
	Js        js.Js
	Msg       *MessageList
	htmlOut   map[string]*string
	jsonOut   map[string]*string
	Rw        http.ResponseWriter
	Req       *http.Request
	SessionId string
	User      *user.UserRec
	RightList map[string]bool
	Session   *sessionType
	Cargo     cargo.CargoList
	Url       *url.URL
	lang      string
}

func NewContext(rw http.ResponseWriter, req *http.Request) *Ctx {
	ctx := new(Ctx)

	ctx.Rw = rw
	ctx.Req = req
	ctx.Config = app.CopyConfig()
	ctx.Cargo = cargo.NewCargo()
	ctx.Cargo.AddStr("ln", ctx.Config.Str("lang"))

	parsedUrl, err := url.ParseRequestURI(ctx.Req.RequestURI)
	if err != nil {
		panic("Can not parse request url.")
	}

	ctx.Url = parsedUrl

	//start session
	ctx.Session = ctx.NewSession()
	ctx.readSession()

	if ctx.SessionId != "" {
		ctx.User = user.New(ctx, ctx.Session.UserId)
	}

	ctx.Msg = NewMessageList()
	ctx.Css = css.New()
	ctx.Js = js.New()

	ctx.htmlOut = make(map[string]*string, 50)
	ctx.jsonOut = make(map[string]*string, 50)
	ctx.lang = ctx.Config.Str("lang")

	ctx.readMessages()
	ctx.clearMessages()

	return ctx
}

func (ctx *Ctx) SetLang() {
	ln := ctx.Cargo.Str("ln")

	if ln == "en" || ln == "tr" {
		ctx.lang = ln
	}
}

func (ctx *Ctx) SetLangTo(lang string) {
	if lang != "en" && lang != "tr" {
		panic(ctx.T("Unknown site language:") + " " + lang)
	}

	ctx.lang = lang
	ctx.Cargo.SetStr("ln", lang)
}

func (ctx *Ctx) Lang() string {
	return ctx.lang
}

func (ctx *Ctx) ReadCargo() {
	values := ctx.Url.Query()

	for k, _ := range ctx.Cargo {
		str := values.Get(k)
		if str != "" {
			ctx.Cargo.SetConvert(k, str)
		}
	}

	ctx.SetLang()
}

func (ctx *Ctx) AddHtml(name string, cnt *string) {
	if _, ok := ctx.htmlOut[name]; ok {
		panic("Formatted content already exists: " + name)
	}

	ctx.htmlOut[name] = cnt
}

func (ctx *Ctx) AddJson(name string, cnt *string) {
	if _, ok := ctx.jsonOut[name]; ok {
		panic("Formatted content already exists: " + name)
	}

	ctx.jsonOut[name] = cnt
}

func (ctx *Ctx) L() string {
	return strings.Title(ctx.lang)
}

func (ctx *Ctx) T(str string) string {

	if ctx.lang == "en" {
		return str
	}

	if val, ok := app.Tran[str]; ok && val != "" {
		return val
	}

	return "[" + str + "]"
}

func (ctx *Ctx) Redirect(urlStr string) {
	ctx.saveMessages()
	ctx.SaveSession()

	http.Redirect(ctx.Rw, ctx.Req, urlStr, 302)
}

func (ctx *Ctx) Render(pageTmp string) {
	ctx.SaveSession()

	ctx.AddHtml("css", ctx.Css.Format())
	ctx.AddHtml("js", ctx.Js.Format())

	messageStr := ctx.Msg.Format(ctx)
	if len(*messageStr) > 0 {
		ctx.AddHtml("message", messageStr)
	}

	lang := ctx.Config.Str("lang")
	ctx.AddHtml("lang", &lang)

	err := app.Tmpl.ExecuteTemplate(ctx.Rw, pageTmp, ctx.htmlOut)
	if err != nil {
		panic(err)
	}
}

func (ctx *Ctx) RenderAjax() {
	ajaxStatus := "success"
	ctx.AddJson("status", &ajaxStatus)

	messageStr := ctx.Msg.Format(ctx)
	if len(*messageStr) > 0 {
		ctx.AddJson("message", messageStr)
	}

	lang := ctx.Config.Str("lang")
	ctx.AddJson("lang", &lang)

	jsonStr, err := json.Marshal(ctx.jsonOut)
	if err != nil {
		panic(err)
	}

	ctx.Rw.Header().Set("Content-Type", "application/json")
	ctx.Rw.Header().Set("Content-Length", strconv.Itoa(len(jsonStr)))
	ctx.Rw.Write(jsonStr)
}

func (ctx *Ctx) RenderAjaxError() {
	ajaxStatus := "error"
	ctx.AddJson("status", &ajaxStatus)

	messageStr := ctx.Msg.Format(ctx)
	if len(*messageStr) > 0 {
		ctx.AddJson("message", messageStr)
	}

	lang := ctx.Config.Str("lang")
	ctx.AddJson("lang", &lang)

	jsonStr, err := json.Marshal(ctx.jsonOut)
	if err != nil {
		panic(err)
	}

	ctx.Rw.Header().Set("Content-Type", "application/json")
	ctx.Rw.Header().Set("Content-Length", strconv.Itoa(len(jsonStr)))
	ctx.Rw.Write(jsonStr)
}

func (ctx *Ctx) RenderEmail(pageTmp string) (string, error) {
	lang := ctx.Config.Str("lang")
	ctx.AddHtml("emailLang", &lang)

	buf := new(bytes.Buffer)
	err := app.Tmpl.ExecuteTemplate(buf, pageTmp, ctx.htmlOut)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (ctx *Ctx) TotalPage(totalRows, pageLen int64) int64 {
	var r int64

	if totalRows%pageLen > 0 {
		r = 1
	} else {
		r = 0
	}

	totalPage := totalRows/pageLen + r

	return totalPage
}

func (ctx *Ctx) TouchPageNo(pageNo, totalRows, pageLen int64) int64 {
	totalPage := ctx.TotalPage(totalRows, pageLen)
	if pageNo < 1 {
		return 1
	} else if pageNo > totalPage {
		return totalPage
	}

	return pageNo
}

func (ctx *Ctx) IsRight(name, function string) bool {
	if ctx.User == nil {
		return false
	}

	return ctx.User.IsRight(name, function)
}

func (ctx *Ctx) IsLoggedIn() bool {
	return ctx.SessionId != ""
}

func (ctx *Ctx) IsSuperuser() bool {
	if ctx.User == nil {
		return false
	}

	return ctx.User.IsSuperUser()

}

func (ctx *Ctx) U(urlStr string, args ...string) string {
	return ctx.Cargo.MakeUrl(urlStr, args...)
}
