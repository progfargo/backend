package role

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"backend/src/app"
	"backend/src/content"
	"backend/src/content/left_menu"
	"backend/src/content/top_menu"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/role/role_lib"
)

func RoleRight(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("role", "role_right") {
		app.BadRequest()
	}

	ctx.Cargo.AddInt("roleId", -1)
	ctx.ReadCargo()

	roleId := ctx.Cargo.Int("roleId")
	roleRec, err := role_lib.GetRoleRec(roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Msg.Warning(ctx.T("Record could not be found."))
			ctx.Redirect(ctx.U("/role"))
			return
		}

		panic(err)
	}

	if roleRec.Name == "admin" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("'admin' role can only be updated by superuser."))
		ctx.Redirect(ctx.U("/role"))
		return
	}

	if ctx.Req.Method == "GET" {
		roleRightForm(ctx, roleRec)
		return
	}

	if app.Ini.AppType == "demo" && !ctx.IsSuperuser() {
		ctx.Msg.Warning(ctx.T("This function is not permitted in demo mode."))
		ctx.Redirect(ctx.U("/role"))
		return
	}

	if err := ctx.Req.ParseForm(); err != nil {
		panic(err)
	}

	rightList := ctx.Req.Form["right"]

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := `delete from
					roleRight
				where
					roleId = ?`

	res, err := tx.Exec(sqlStr, roleId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	var formValue []string
	var ok bool
	var key string

	if len(rightList) > 0 {
		sqlBuf := new(util.Buf)
		sqlBuf.Add("insert into roleRight(roleId, pageName, funcName)")

		values := make([]string, 0, 60)
		for _, val := range rightList {
			formValue = strings.Split(val, ":")

			if len(formValue) < 2 {
				ctx.Msg.Warning(ctx.T("Wrong value received from user right form."))
				tx.Rollback()
				ctx.Redirect(ctx.U("/role"))
				return
			}

			key = app.MakeKey(util.DbStr(formValue[0]), util.DbStr(formValue[1]))
			if _, ok = app.UserRightMap[key]; !ok {
				tx.Rollback()
				ctx.Msg.Warning(ctx.T("Could not find user right:") + " " + key)
				ctx.Redirect(ctx.U("/role_right", "roleId"))
				return
			}

			values = append(values, fmt.Sprintf("(%d, '%s', '%s')",
				roleId, util.DbStr(formValue[0]), util.DbStr(formValue[1])))
		}

		sqlBuf.Add("values" + strings.Join(values, ",\n"))

		res, err = tx.Exec(*sqlBuf.String())
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	tx.Commit()

	if n, _ := res.RowsAffected(); n == 0 {
		tx.Rollback()
		ctx.Msg.Warning(ctx.T("You did not change the record."))
		roleRightForm(ctx, roleRec)
		return
	}

	ctx.Msg.Success(ctx.T("Role rights has been changed."))
	ctx.Redirect(ctx.U("/role"))
}

func roleRightForm(ctx *context.Ctx, rec *role_lib.RoleRec) {
	role_lib.SetRoleRightList(ctx)
	content.Include(ctx)

	ctx.Css.Add("/asset/css/page/role_right.css")
	ctx.Js.Add("/asset/js/page/role/role_right.js")

	roleId := ctx.Cargo.Int("roleId")

	roleRightList, err := role_lib.GetRoleRight(roleId)
	if err != nil {
		panic(err)
	}

	buf := util.NewBuf()

	buf.Add("<div class=\"row\">")
	buf.Add("<div class=\"col\">")
	buf.Add(content.PageTitle(ctx.T("Roles"), ctx.T("Role Rights"), rec.Name))
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")
	buf.Add("<div class=\"buttonGroupFixed\">")
	buf.Add(content.BackButton(ctx, ctx.U("/role")))

	buf.Add("<button type=\"button\" class=\"expandAll button buttonDefault buttonSm\">%s</button>", ctx.T("Expand All"))
	buf.Add("<button type=\"button\" class=\"collapseAll button buttonDefault buttonSm\">%s</button>", ctx.T("Collapse All"))
	buf.Add("</div>")
	buf.Add("</div>")

	buf.Add("<div class=\"col\">")

	urlStr := ctx.U("/role_right", "roleId")
	buf.Add("<form class=\"formRole\" action=\"%s\" method=\"post\">", urlStr)

	tabIndex := 1

	for _, rt := range app.UserRightList.List {
		tbuf := util.NewBuf()
		tbuf.Add("<table>")
		tbuf.Add("<thead>")

		tbuf.Add("<tr class=\"roleHeader\">")
		tbuf.Add("<th class=\"fixedMiddle\">%s</th>", rt.Title)
		tbuf.Add("<th>")
		tbuf.Add("<i class=\"fa fa-caret-down accordionIcon\"></i>")
		tbuf.Add("<button type=\"button\" class=\"selectAll button buttonDefault buttonXs\">%s</button>", ctx.T("All"))
		tbuf.Add("<button type=\"button\" class=\"selectNone button buttonDefault buttonXs\">%s</button>", ctx.T("None"))
		tbuf.Add("<button type=\"button\" class=\"selectReverse button buttonDefault buttonXs\">%s</button>", ctx.T("Reverse"))
		tbuf.Add("</th>")
		tbuf.Add("</tr>")
		tbuf.Add("</thead>")

		tbuf.Add("<tbody>")
		tbuf.Add("<tr><th class=\"fixedMiddle\">%s</th><th>%s</th></tr>",
			ctx.T("Right Name"), ctx.T("Explanation"))

		isEmpty := true

		for _, r := range rt.List {
			isEmpty = false

			tbuf.Add("<tr>")

			key := app.MakeKey(r.PageName, r.FuncName)
			tbuf.Add("<td>")

			var class string
			if roleRightList[app.MakeKey(r.PageName, r.FuncName)] {
				class = " checked"
			} else {
				class = ""
			}

			tbuf.Add("<input type=\"checkbox\" name=\"right\" value=\"%s\" tabindex=\"%d\" %s><span class=\"radioLabel\">%s</span>",
				key, tabIndex, class, r.FuncName)

			tbuf.Add("</td>")

			tbuf.Add("<td>%s</td>", r.Exp)
			tbuf.Add("</tr>")

			tabIndex += 1
		}

		tbuf.Add("</tbody>")
		tbuf.Add("</table>")

		if !isEmpty {
			buf.Add(*tbuf.String())
		}
	}

	buf.Add("<div class=\"formGroup formCommand\">")
	tabIndex++
	buf.Add("<button type=\"submit\" class=\"button buttonPrimary buttonSm\" tabindex=\"%d\">%s</button>",
		tabIndex, ctx.T("Submit"))

	tabIndex++
	buf.Add("<button type=\"reset\" class=\"button buttonDefault buttonSm\" tabindex=\"%d\">%s</button>",
		tabIndex, ctx.T("Reset"))
	buf.Add("</div>")

	buf.Add("</form>")

	buf.Add("</div>")
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	content.Default(ctx)

	lmenu := left_menu.New()
	lmenu.Set(ctx, "role")

	tmenu := top_menu.New()
	tmenu.Set(ctx)

	str := "roleRightPage"
	ctx.AddHtml("pageName", &str)

	ctx.Render("default.html")
}
