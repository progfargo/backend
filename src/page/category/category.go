package category

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/category/category_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("category", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("categoryId", -1)
	ctx.ReadCargo()

	content.Include(ctx)

	browseMid(ctx)

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "category")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "categoryPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	ctx.Css.Add("/asset/css/page/category.css")
	ctx.Js.Add("/asset/js/page/category/category.js")

	insertRight := ctx.IsRight("category", "insert")
	updateRight := ctx.IsRight("category", "update")
	deleteRight := ctx.IsRight("category", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machine Categories")))
	buf.Add("</div>")

	if insertRight {
		buf.Add("<div class=\"col\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		buf.Add(content.NewButton(ctx, ctx.U("/category_insert")))

		buf.Add("<a class=\"expandAll button buttonDefault buttonSm\" title=\"%s\">%s</a>",
			ctx.T("Expand all records."), ctx.T("Expand All"))

		buf.Add("<a class=\"collapseAll button buttonDefault buttonSm\" title=\"%s\">%s</a>",
			ctx.T("Collapse all records."), ctx.T("Collapse All"))

		buf.Add("</div>")
		buf.Add("</div>")
	}

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedZero\" title=\"%s\">%s</th>", ctx.T("Enumeration"), ctx.T("En."))
	buf.Add("<th>%s</th>", ctx.T("Name"))
	buf.Add("<th class=\"fixedSmall center\">%s</th>", ctx.T("Status"))

	if updateRight || deleteRight {
		buf.Add("<th class=\"right fixedZero\">%s</th>", ctx.T("Command"))
	}

	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	pgl := category_lib.NewCategoryList()
	pgl.Set(ctx)

	children := pgl.Tax.GetChildren("root")

	if len(children) == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	} else {
		for _, name := range children {
			FormatItem(ctx, pgl, name, buf, 0, updateRight, deleteRight, make(map[string]bool, 5))
		}
	}

	buf.Add("<tbody>")
	buf.Add("</table>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}

func FormatItem(ctx *context.Ctx, pgl *category_lib.CategoryList, name string,
	buf *util.Buf, tab int, updateRight, deleteRight bool, classList map[string]bool) {
	item := pgl.Tax.GetItem(name)
	data := item.Data.(*category_lib.CategoryItem)

	isParent := pgl.Tax.IsParent(name)

	icon := ""
	parentClass := ""
	if isParent {
		icon = "<span class=\"collapseIcon\"><i class=\"fas fa-minus\"></i></span>" +
			"<span class=\"expandIcon\"><i class=\"fas fa-plus\"></i></span>"
		parentClass = fmt.Sprintf(" parent-%s", item.Name())
	}

	keys := make([]string, 0, len(classList))
	for k := range classList {
		keys = append(keys, k)
	}

	buf.Add("<tr>")

	buf.Add("<td>%d</td>", item.Enum())

	title := ctx.T("Category")

	buf.Add("<td class=\"name %s%s\" title=\"%s\">%s%s %s</td>",
		strings.Join(keys, " "),
		parentClass,
		title,
		strings.Repeat("&nbsp;", tab),
		icon,
		data.CategoryName,
	)

	buf.Add("<td class=\"center\">%s</td>", category_lib.StatusToLabel(ctx, data.Status))

	nameInt, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		panic(err)
	}

	ctx.Cargo.SetInt("categoryId", nameInt)

	var urlStr string
	if updateRight || deleteRight {
		buf.Add("<td class=\"right\">")
		buf.Add("<div class=\"buttonGroupFixed\">")

		buf.Add("<div class=\"editDelete\">")

		if updateRight {
			urlStr = ctx.U("/category_update", "categoryId")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Edit category properties."), ctx.T("Edit"))
		}

		if deleteRight {
			buf.Add("<a href=\"#\" class=\"button buttonError buttonXs buttonDelete\" title=\"%s\">%s</a>",
				ctx.T("Delete record."), ctx.T("Delete"))
		}

		buf.Add("</div>")

		buf.Add("<div class=\"deleteConfirm\">")
		buf.Add(ctx.T("Do you realy want to delete this record?"))

		urlStr = ctx.U("/category_delete", "categoryId")
		buf.Add("<a href=\"%s\" class=\"button buttonSuccess buttonXs\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Delete this record."), ctx.T("Yes"))

		buf.Add("<a class=\"button buttonDefault buttonXs cancelButton\">%s</a>", ctx.T("Cancel"))
		buf.Add("</div>")

		buf.Add("</div>")
		buf.Add("</td>")
	}

	buf.Add("</tr>")

	if isParent {
		classList[fmt.Sprintf("sub-%s", item.Name())] = true

		children := pgl.Tax.GetChildren(name)
		for _, v := range children {
			FormatItem(ctx, pgl, v, buf, tab+6, updateRight, deleteRight, classList)
		}

		delete(classList, fmt.Sprintf("sub-%d", item.Name()))
	}
}
