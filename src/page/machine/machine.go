package machine

import (
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/combo"
	"backend/src/content/left_menu"
	"backend/src/content/ruler"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/machine/machine_lib"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("machine", "browse") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("machineId", -1)
	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddStr("stat", "default")
	ctx.Cargo.AddInt("pn", 1)
	ctx.Cargo.AddInt("imgId", -1)
	ctx.Cargo.AddInt("catId", 1)
	ctx.Cargo.AddInt("manId", -1)
	ctx.ReadCargo()

	content.Include(ctx)
	ctx.Css.Add("/asset/jquery_modal/jquery.modal.css")
	ctx.Js.Add("/asset/jquery_modal/jquery.modal.js")
	ctx.Js.Add("/asset/js/page/machine/machine.js")

	browseMid(ctx)

	content.Default(ctx)

	content.Search(ctx, "/machine", "stat")

	lmenu := left_menu.New()
	lmenu.Set(ctx, "machine")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "machinePage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}

func browseMid(ctx *context.Ctx) {
	key := ctx.Cargo.Str("key")
	stat := ctx.Cargo.Str("stat")
	pageNo := ctx.Cargo.Int("pn")
	catId := ctx.Cargo.Int("catId")
	manId := ctx.Cargo.Int("manId")

	totalRows := machine_lib.CountMachine(ctx, key, stat, catId, manId)
	if totalRows == 0 {
		ctx.Msg.Warning(ctx.T("Empty list."))
	}

	pageLen := ctx.Config.Int("pageLen")
	pageNo = ctx.TouchPageNo(pageNo, totalRows, pageLen)

	insertRight := ctx.IsRight("machine", "insert")
	updateRight := ctx.IsRight("machine", "update")
	deleteRight := ctx.IsRight("machine", "delete")

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Machines")))
	buf.Add("</div>")

	var urlStr string
	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")

	if insertRight {
		buf.Add(content.NewButton(ctx, ctx.U("/machine_insert", "key", "stat", "pn", "catId", "manId")))
	}

	if stat == "default" {
		urlStr = ctx.U("/machine?stat=passive", "key", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Show passive records."), ctx.T("Passive"))
	} else if stat == "passive" {
		urlStr = ctx.U("/machine", "key", "pn", "catId", "manId")
		buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm\" title=\"%s\">%s</a>",
			urlStr, ctx.T("Show all records."), ctx.T("All"))
	}

	categoryCombo := combo.NewTaxCombo(`select
											categoryId,
											parentId,
											name,
											enum
										from
											category`, ctx.T("Main Category"))
	categoryCombo.Set(ctx)

	//category form
	buf.Add("<form id=\"categoryForm\" class=\"formInline\" action=\"/machine\" method=\"get\">")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<select name=\"catId\" class=\"formControl\">")

	buf.Add(categoryCombo.Format(fmt.Sprintf("%d", catId)))

	buf.Add("</select>")
	buf.Add("</div>")
	buf.Add(content.HiddenCargo(ctx, "key", "stat", "ln", "manId"))
	buf.Add("</form>")

	manufacturerCombo := combo.NewCombo(`select
											manufacturerId,
											name
										from
											manufacturer
										order by name`, ctx.T("Manufacturer"))
	manufacturerCombo.Set()

	//manufacturer form
	buf.Add("<form id=\"manufacturerForm\" class=\"formInline\" action=\"/machine\" method=\"get\">")

	buf.Add("<div class=\"formGroup\">")
	buf.Add("<select name=\"manId\" class=\"formControl\">")

	buf.Add(manufacturerCombo.Format(manId))

	buf.Add("</select>")
	buf.Add("</div>")
	buf.Add(content.HiddenCargo(ctx, "key", "stat", "ln", "catId", "manId"))
	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	buf.Add("<table>")
	buf.Add("<thead>")
	buf.Add("<tr>")
	buf.Add("<th class=\"fixedZero\">%s</th>", ctx.T("Id"))
	buf.Add("<th>%s</th>", ctx.T("Name"))
	buf.Add("<th>%s</th>", ctx.T("Image"))
	buf.Add("<th class=\"fixedZero\">%s</th>", ctx.T("Price"))
	buf.Add("<th class=\"fixedZero\">%s</th>", ctx.T("Status"))
	buf.Add("<th class=\"right fixedZero\">%s</th>", ctx.T("Command"))
	buf.Add("</tr>")
	buf.Add("</thead>")

	buf.Add("<tbody>")

	if totalRows > 0 {
		machineList := machine_lib.GetMachinePage(ctx, key, stat, catId, manId, pageNo)

		var name string
		var machineImageId int64
		var smallImageUrl, popupUrl string
		var padSize = 9

		for _, row := range machineList {
			ctx.Cargo.SetInt("machineId", row.MachineId)

			name = util.ScrStr(row.Name)

			if key != "" {
				name = content.Find(name, key)
			}

			buf.Add("<tr>")

			buf.Add("<td>%d</td>", row.MachineId)

			buf.Add("<td>")
			urlStr = ctx.U("/machine_display", "key", "stat", "pn", "catId", "manId", "machineId")
			buf.Add("<a href=\"%s\">%s</a>", urlStr, name)

			buf.Add("<div class=\"machineInfo\">")

			buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Maker"), "&nbsp;", padSize), util.ScrStr(row.ManufacturerName))
			buf.Add("<strong>%s:</strong> %d<br>", util.PadRight(ctx.T("Year"), "&nbsp;", padSize), row.Yom)
			buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Model"), "&nbsp;", padSize), row.Model)
			buf.Add("<strong>%s:</strong> %s<br>", util.PadRight(ctx.T("Location"), "&nbsp;", padSize), row.Location)
			buf.Add("</div>") //machineInfo

			buf.Add("</td>")

			machineImageId = util.NullToInt64(row.MachineImageId)

			if machineImageId != 0 {
				ctx.Cargo.SetInt("imgId", machineImageId)
				smallImageUrl = ctx.U("/machine_image_small", "imgId")
				popupUrl = ctx.U("/machine_popup", "imgId")

				buf.Add("<td>")
				buf.Add("<img class=\"smallImage\" src=\"%s\" alt=\"\">", smallImageUrl)
				buf.Add("<a class=\"popupImage\" href=\"%s\"></a>", popupUrl)
				buf.Add("</td>")
			} else {
				buf.Add("<td></td>")
			}

			buf.Add("<td class=\"right\">%d &#x24;</td>", row.Price)
			buf.Add("<td>%s</td>", machine_lib.StatusToLabel(ctx, row.Status))

			buf.Add("<td class=\"right\">")
			buf.Add("<div class=\"buttonGroupFixed\">")

			urlStr = ctx.U("/machine_display", "key", "stat", "pn", "catId", "manId", "machineId")
			buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
				urlStr, ctx.T("Display record."), ctx.T("Display"))

			if updateRight {
				urlStr = ctx.U("/machine_update", "machineId", "key", "stat", "pn", "catId", "manId")
				buf.Add("<a href=\"%s\" class=\"button buttonDefault buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Edit record."), ctx.T("Edit"))
			}

			if deleteRight {
				urlStr = ctx.U("/machine_delete", "machineId", "key", "stat", "pn", "catId", "manId")
				buf.Add("<a href=\"%s\" class=\"button buttonError buttonXs\" title=\"%s\">%s</a>",
					urlStr, ctx.T("Delete record."), ctx.T("Delete"))
			}

			buf.Add("</div>")
			buf.Add("</td>")

			buf.Add("</tr>")
		}
	}

	buf.Add("</tbody>")
	buf.Add("</table>")
	buf.Add("</div>")

	totalPage := ctx.TotalPage(totalRows, pageLen)
	if totalPage > 1 {
		buf.Add("<div class=\"col\">")
		ruler := ruler.NewRuler(totalPage, pageNo, ctx.U("/machine", "key", "stat", "catId", "manId"))
		ruler.Set(ctx)
		buf.Add(ruler.Format())
		buf.Add("</div>")
	}

	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())
}
