package machine_menu

import (
	"fmt"

	"backend/src/lib/context"
	"backend/src/lib/tax"
	"backend/src/lib/util"
)

type menuItem struct {
	url      string
	label    string
	icon     string
	isActive bool
}

type machineMenu struct {
	tax       *tax.Tax
	cargoList []string
}

func New(cargoVar ...string) *machineMenu {
	rv := new(machineMenu)
	rv.tax = tax.New()
	rv.Add("root", "end", 0, true, "", "", "")
	rv.cargoList = cargoVar

	return rv
}

func (mm *machineMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string) {
	mm.tax.Add(name, parent, enum, isVisible, &menuItem{url: url, label: label, icon: icon})
}

func (mm *machineMenu) Set(ctx *context.Ctx, name ...string) {
	mm.Add("machine_display", "root", 10, true, ctx.U("/machine_display", mm.cargoList...), ctx.T("Display"), "fasye")
	mm.Add("machine_update", "root", 20, true, ctx.U("/machine_update", mm.cargoList...), ctx.T("Edit"), "fasdit")
	mm.Add("machine_image", "root", 30, true, ctx.U("/machine_image", mm.cargoList...), ctx.T("Images"), "fas fa-images")
	mm.Add("machine_video", "root", 40, true, ctx.U("/machine_video", mm.cargoList...), ctx.T("Videos"), "fas fa-video")
	mm.Add("machine_pdf", "root", 50, true, ctx.U("/machine_pdf", mm.cargoList...), ctx.T("PDF"), "fas fa-file-pdf")

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		mm.setActive(name[0])
	}

	mm.tax.SortChildren()
	mm.reduce("root")
}

func (mm *machineMenu) setActive(name string) {
	item := mm.tax.GetItem(name)
	data := item.Data.(*menuItem)
	data.isActive = true
}

func (mm *machineMenu) reduce(name string) {
	item := mm.tax.GetItem(name)

	if mm.tax.IsParent(name) {
		children := mm.tax.GetChildren(name)
		for _, val := range children {
			mm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !mm.tax.IsParent(name) {
			mm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = mm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := mm.tax.GetItem(firstChild)

		nameData := item.Data.(*menuItem)
		firstChildData := firstChildItem.Data.(*menuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		mm.tax.Delete(name)
	}
}

func (mm *machineMenu) Format(ctx *context.Ctx) string {
	rv := util.NewBuf()

	children := mm.tax.GetChildren("root")

	if len(children) == 0 {
		return ""
	}

	rv.Add("<div class=\"buttonBar machineMenu\">")

	for _, v := range children {
		item := mm.tax.GetItem(v)
		data := item.Data.(*menuItem)

		class := ""
		if data.isActive {
			class = " selected"
		}

		icon := ""
		if data.icon != "" {
			icon = fmt.Sprintf("<i class=\"%s fa-fw left\"></i>", data.icon)
		}

		rv.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm%s\">%s%s</a>", data.url, class, icon, data.label)
	}

	rv.Add("</div>")

	return *rv.String()
}
