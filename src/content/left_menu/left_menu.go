package left_menu

import (
	"fmt"
	"strings"

	"backend/src/lib/context"
	"backend/src/lib/tax"
	"backend/src/lib/util"
)

type leftMenuItem struct {
	url           string
	label         string
	icon          string
	isActive      bool
	isSubmenuOpen bool
}

type leftMenu struct {
	tax *tax.Tax
}

func New() *leftMenu {
	rv := new(leftMenu)
	rv.tax = tax.New()
	rv.Add("root", "end", 0, true, "", "", "")

	return rv
}

func (lm *leftMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string) {
	lm.tax.Add(name, parent, enum, isVisible, &leftMenuItem{url, label, icon, false, false})
}

func (lm *leftMenu) Set(ctx *context.Ctx, name ...string) {
	lm.Add("admin", "root", 10, ctx.IsRight("user", "browse"), ctx.U("/user"), ctx.T("Admin"), "fas fa-gear")

	lm.Add("user_", "admin", 10, ctx.IsRight("user", "browse"), ctx.U("/user"), ctx.T("User Management"), "fas fa-user")
	lm.Add("user", "user_", 10, ctx.IsRight("user", "browse"), ctx.U("/user"), ctx.T("Users"), "fas fa-user-cog")
	lm.Add("role", "user_", 20, ctx.IsRight("role", "browse"), ctx.U("/role"), ctx.T("Roles"), "fas fa-users-gear")

	lm.Add("config", "admin", 30, ctx.IsRight("config", "browse"), ctx.U("/config"), ctx.T("App Configuration"), "fas fa-gauge")
	lm.Add("tran", "admin", 40, ctx.IsRight("tran", "browse"), ctx.U("/tran"), ctx.T("Translation Table"), "fas fa-globe")

	lm.Add("category", "root", 20, ctx.IsRight("category", "browse"), ctx.U("/category"), ctx.T("Categories"), "fas fa-sitemap")

	lm.Add("machine", "root", 30, ctx.IsRight("machine", "browse"), ctx.U("/machine"),
		ctx.T("Machines"), "fas fa-list")
	lm.Add("manufacturer", "machine", 10, ctx.IsRight("manufacturer", "browse_own"), ctx.U("/manufacturer"),
		ctx.T("Manufacturers"), "fas fa-copyright")

	lm.Add("content", "root", 50, ctx.IsRight("news", "browse"), ctx.U("/news"), ctx.T("Content"), "fas fa-code")
	lm.Add("jumbotron", "content", 10, ctx.IsRight("jumbotron", "browse"), ctx.U("/jumbotron"), ctx.T("Jumbotron"), "fas fa-square")
	lm.Add("social", "content", 20, ctx.IsRight("social", "browse"), ctx.U("/social"), ctx.T("Social Links"), "fas fa-share-alt")
	lm.Add("office", "content", 20, ctx.IsRight("office", "browse"), ctx.U("/office"), ctx.T("Offices"), "fas fa-map-marker-alt")
	lm.Add("news", "content", 30, ctx.IsRight("news", "browse"), ctx.U("/news"), ctx.T("News"), "fas fa-newspaper")
	lm.Add("faq", "content", 40, ctx.IsRight("faq", "browse"), ctx.U("/faq"), ctx.T("FAQ"), "fas fa-circle-question")
	lm.Add("banner", "content", 50, ctx.IsRight("banner", "browse"), ctx.U("/banner"), ctx.T("Banner"), "fas fa-images")
	lm.Add("text_content", "content", 60, ctx.IsRight("text_content", "browse"), ctx.U("/text_content"), ctx.T("Text Content"), "fas fa-file-lines")

	lm.Add("site_error", "root", 80, ctx.IsRight("site_error", "browse"), ctx.U("/site_error"), ctx.T("Site Errors"), "fas fa-bug")

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		lm.setActive(name[0])
	}

	lm.tax.SortChildren()
	lm.reduce("root")
	lm.format(ctx)
}

func (lm *leftMenu) setActive(name string) {
	allParents := lm.tax.GetAllParents(name)
	allParents = append(allParents, name)

	for _, val := range allParents {
		item := lm.tax.GetItem(val)
		data := item.Data.(*leftMenuItem)
		data.isSubmenuOpen = true

		if val == name {
			data.isActive = true
		}
	}
}

func (lm *leftMenu) reduce(name string) {
	item := lm.tax.GetItem(name)

	if lm.tax.IsParent(name) {
		children := lm.tax.GetChildren(name)
		for _, val := range children {
			lm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !lm.tax.IsParent(name) {
			lm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = lm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := lm.tax.GetItem(firstChild)

		nameData := item.Data.(*leftMenuItem)
		firstChildData := firstChildItem.Data.(*leftMenuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		lm.tax.Delete(name)
	}
}

func (lm *leftMenu) format(ctx *context.Ctx) {
	rv := util.NewBuf()

	children := lm.tax.GetChildren("root")

	if len(children) == 0 {
		return
	}

	rv.Add("<nav>")
	rv.Add("<ul>")

	for _, v := range children {
		lm.FormatItem(v, rv)
	}

	rv.Add("</ul>")
	rv.Add("</nav>")

	ctx.AddHtml("leftMenu", rv.String())
}

func (lm *leftMenu) FormatItem(name string, rv *util.Buf) {
	item := lm.tax.GetItem(name)
	data := item.Data.(*leftMenuItem)

	classList := make([]string, 0, 5)
	if data.isActive {
		classList = append(classList, "active")
	}

	icon := ""
	if data.icon != "" {
		icon = fmt.Sprintf("<i class=\"%s fa-fw menuIcon left\"></i>", data.icon)
	}

	if lm.tax.IsParent(name) {
		classList = append(classList, "subMenu")

		var subMenuShow, subMenuIcon string
		if data.isSubmenuOpen {
			subMenuShow = " class=\"subMenuShow\""
			subMenuIcon = "<i class=\"subMenuButton fas fa-minus fa-fw\">"
		} else {
			subMenuShow = ""
			subMenuIcon = "<i class=\"subMenuButton fas fa-plus fa-fw\">"
		}

		classStr := ""
		if len(classList) > 0 {
			classStr = fmt.Sprintf(" class=\"%s\"", strings.Join(classList, " "))

		}

		rv.Add("<li%s>", classStr)
		rv.Add("<a href=\"%s\">%s%s%s"+
			"</i></a>", data.url, icon, data.label, subMenuIcon)
		rv.Add("<ul%s>", subMenuShow)

		children := lm.tax.GetChildren(name)
		for _, v := range children {
			lm.FormatItem(v, rv)
		}

		rv.Add("</ul>")
		rv.Add("</li>")
	} else {
		classStr := ""
		if len(classList) > 0 {
			classStr = fmt.Sprintf(" class=\"%s\"", strings.Join(classList, " "))
		}

		rv.Add("<li%s>", classStr)
		rv.Add("<a href=\"%s\">%s%s</a>", data.url, icon, data.label)
		rv.Add("</li>")
	}
}
