package logout

import (
	"fmt"
	"net/http"

	"backend/src/app"
	"backend/src/content"
	"backend/src/lib/context"
)

func Browse(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsLoggedIn() {
		app.BadRequest()
	}

	ctx.DeleteSession()
	content.End(ctx, "success", ctx.T("Goodbye"), app.SiteName,
		fmt.Sprintf("<p>%s</p>", ctx.T("Your session has been closed.")),
		fmt.Sprintf("<a href=\"/login\">%s</a>", ctx.T("Login page.")))
}
