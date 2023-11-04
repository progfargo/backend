package content

import (
	"fmt"
	"regexp"
	"strings"

	"backend/src/lib/context"
	"backend/src/lib/util"
)

func Include(ctx *context.Ctx) {
	ctx.Css.Add("/asset/fontawesome/css/all.css")
	ctx.Css.Add("/asset/css/style.css")

	ctx.Js.Add("/asset/js/jquery-2.1.3.min.js")
	ctx.Js.Add("/asset/js/flex_watch.js")
	ctx.Js.Add("/asset/js/common.js")
}

func Default(ctx *context.Ctx) {
	HtmlTitle(ctx, "backend")
	SiteTitle(ctx)

	if ctx.IsLoggedIn() {
		UserInfo(ctx)
	}

	Logo(ctx)
	FooterNote(ctx)
}

func HtmlTitle(ctx *context.Ctx, str string) {
	ctx.AddHtml("htmlTitle", &str)
}

func SiteTitle(ctx *context.Ctx) {
	str := fmt.Sprintf("<div class=\"siteTitle\"><h1><span>backend</span> %s</h1></div>", ctx.T("sample"))
	ctx.AddHtml("siteTitle", &str)
}

func FooterNote(ctx *context.Ctx) {
	str := "<p><small>&copy; 2022 backend ver:1.0.0</small></p>"

	ctx.AddHtml("footerNote", &str)
}

func Search(ctx *context.Ctx, url string, args ...string) {
	buf := util.NewBuf()

	key := ctx.Cargo.Str("key")
	buf.Add("<form class=\"formInline formSearch\" action=\"%s\" method=\"get\">", url)
	buf.Add("<input type=\"text\" name=\"key\" placeholder=\"%s\" value=\"%s\">", "search word...", key)
	buf.Add(HiddenCargo(ctx, args...))
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\">%s</button>", "Search")
	buf.Add("</form>")

	ctx.AddHtml("search", buf.String())
}

func UserInfo(ctx *context.Ctx) {
	buf := util.NewBuf()
	buf.Add("<div class=\"userInfo\">")

	buf.Add("<p><strong>%s</strong> <span title=\"%s\">%s</span></p>", ctx.T("User:"), ctx.User.Name, ctx.User.Login)

	buf.Add("</div>")

	ctx.AddHtml("userInfo", buf.String())
}

func Logo(ctx *context.Ctx) {
	buf := util.NewBuf()

	buf.Add("<div class=\"logo\">")
	buf.Add("<img src=\"/asset/img/base_logo.png\" title=\"to the moon\" alt=\"backend logo\">")
	buf.Add("</div>")

	ctx.AddHtml("logo", buf.String())
}

func HiddenCargo(ctx *context.Ctx, args ...string) string {
	buf := util.NewBuf()
	var val string

	for _, name := range args {
		if !ctx.Cargo.IsDefault(name) {
			val = ctx.Cargo.Str(name)
			buf.Add("<input type=\"hidden\" name=\"%s\" value=\"%s\">", name, val)
		}
	}

	return *buf.String()
}

func Find(str, key string) string {
	key = strings.ReplaceAll(key, "%", "")
	key = strings.ReplaceAll(key, "?", "")

	re := regexp.MustCompile("(?i)" + key)
	return re.ReplaceAllStringFunc(str, replaceFunc)
}

func replaceFunc(str string) string {
	return fmt.Sprintf("<span class=\"find\">%s</span>", str)
}

func PageTitle(args ...string) string {
	buf := util.NewBuf()
	buf.Add("<div class=\"pageTitle\">")
	buf.Add("<h3>%s</h3>", strings.Join(args, " <i class=\"fa fa-fw fa-angle-right\"></i> "))
	buf.Add("</div>")

	return *buf.String()
}

func BackButton(ctx *context.Ctx, urlStr string) string {
	return fmt.Sprintf("<a href=\"%s\" class=\"button buttonAccent buttonSm\">%s</a>",
		urlStr, ctx.T("Back"))
}

func NewButton(ctx *context.Ctx, urlStr string) string {
	return fmt.Sprintf("<a href=\"%s\" class=\"button buttonPrimary buttonSm\" title=\"%s\">%s</a>",
		urlStr, ctx.T("New record."), ctx.T("New"))
}

func End(ctx *context.Ctx, calloutType, calloutTitle, htmlTitle string, msg ...string) {
	buf := util.NewBuf()

	buf.Add("<div class=\"callout callout%s\">", strings.Title(calloutType))
	buf.Add("<h4>%s</h4>", calloutTitle)

	for _, v := range msg {
		buf.Add("<p>%s</p>", v)
	}

	buf.Add("</div>")

	Include(ctx)
	ctx.Css.Add("/asset/css/page/end.css")
	ctx.AddHtml("midContent", buf.String())

	SiteTitle(ctx)
	HtmlTitle(ctx, htmlTitle)
	FooterNote(ctx)

	str := "endPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("end.html")
}
