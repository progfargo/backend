package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"backend/src/app"
	"backend/src/content"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/lib/watch"
	"backend/src/page/banner"
	"backend/src/page/category"
	"backend/src/page/config"
	"backend/src/page/faq"
	"backend/src/page/help"
	"backend/src/page/help_image"
	"backend/src/page/jumbotron"
	"backend/src/page/login"
	"backend/src/page/logout"
	"backend/src/page/machine"
	"backend/src/page/machine_image"
	"backend/src/page/machine_pdf"
	"backend/src/page/machine_video"
	"backend/src/page/manufacturer"
	"backend/src/page/news"
	"backend/src/page/news_image"
	"backend/src/page/office"
	"backend/src/page/profile"
	"backend/src/page/role"
	"backend/src/page/site_error"
	"backend/src/page/social"
	"backend/src/page/text_content"
	"backend/src/page/text_content_image"
	"backend/src/page/tran"
	"backend/src/page/user"
	"backend/src/page/welcome"

	"github.com/dchest/captcha"
)

type HandleFuncType func(http.ResponseWriter, *http.Request)

func main() {
	cron()

	if app.Ini.Debug && app.Ini.IsLocal {
		watch.WatchLess()
	}

	//static constent
	fsAdmin := http.FileServer(http.Dir(app.Ini.HomeDir + "/asset/"))
	http.Handle("/asset/", http.StripPrefix("/asset/", fsAdmin))

	http.Handle("/captcha/", captcha.Server(200, 90))

	//admin
	http.HandleFunc("/tran", makeHandle(tran.Browse))
	http.HandleFunc("/tran_insert", makeHandle(tran.Insert))
	http.HandleFunc("/tran_update", makeHandle(tran.Update))
	http.HandleFunc("/tran_delete", makeHandle(tran.Delete))
	http.HandleFunc("/tran_sync", makeHandle(tran.Sync))
	http.HandleFunc("/tran_export", makeHandle(tran.Export))
	http.HandleFunc("/tran_import", makeHandle(tran.Import))

	http.HandleFunc("/role", makeHandle(role.Browse))
	http.HandleFunc("/role_insert", makeHandle(role.Insert))
	http.HandleFunc("/role_update", makeHandle(role.Update))
	http.HandleFunc("/role_delete", makeHandle(role.Delete))
	http.HandleFunc("/role_right", makeHandle(role.RoleRight))

	http.HandleFunc("/user", makeHandle(user.Browse))
	http.HandleFunc("/user_insert", makeHandle(user.Insert))
	http.HandleFunc("/user_update", makeHandle(user.Update))
	http.HandleFunc("/user_update_pass", makeHandle(user.UpdatePass))
	http.HandleFunc("/user_delete", makeHandle(user.Delete))
	http.HandleFunc("/user_role", makeHandle(user.UserRole))
	http.HandleFunc("/user_role_revoke", makeHandle(user.UserRoleRevoke))
	http.HandleFunc("/user_role_grant", makeHandle(user.UserRoleGrant))
	http.HandleFunc("/user_display", makeHandle(user.Display))

	http.HandleFunc("/config", makeHandle(config.Browse))
	http.HandleFunc("/config_insert", makeHandle(config.Insert))
	http.HandleFunc("/config_update", makeHandle(config.Update))
	http.HandleFunc("/config_delete", makeHandle(config.Delete))
	http.HandleFunc("/config_set", makeHandle(config.Set))

	http.HandleFunc("/site_error", makeHandle(site_error.Browse))
	http.HandleFunc("/site_error_delete", makeHandle(site_error.Delete))
	http.HandleFunc("/site_error_delete_all", makeHandle(site_error.DeleteAll))

	http.HandleFunc("/login", makeHandle(login.Browse))
	http.HandleFunc("/welcome", makeHandle(welcome.Browse))
	http.HandleFunc("/logout", makeHandle(logout.Browse))

	http.HandleFunc("/profile", makeHandle(profile.Browse))
	http.HandleFunc("/profile_update", makeHandle(profile.Update))
	http.HandleFunc("/profile_update_pass", makeHandle(profile.UpdatePass))

	http.HandleFunc("/banner", makeHandle(banner.Browse))
	http.HandleFunc("/banner_insert", makeHandle(banner.Insert))
	http.HandleFunc("/banner_update_info", makeHandle(banner.UpdateInfo))
	http.HandleFunc("/banner_update_image", makeHandle(banner.UpdateImage))
	http.HandleFunc("/banner_delete", makeHandle(banner.Delete))
	http.HandleFunc("/banner_image_small", makeHandle(banner.ImageSmall))
	http.HandleFunc("/banner_image_middle", makeHandle(banner.ImageMiddle))
	http.HandleFunc("/banner_image_original", makeHandle(banner.ImageOriginal))
	http.HandleFunc("/banner_image_download", makeHandle(banner.Download))
	http.HandleFunc("/banner_crop_image", makeHandle(banner.CropImage))

	http.HandleFunc("/help", makeHandle(help.Browse))
	http.HandleFunc("/help_insert", makeHandle(help.Insert))
	http.HandleFunc("/help_update", makeHandle(help.Update))
	http.HandleFunc("/help_delete", makeHandle(help.Delete))

	http.HandleFunc("/help_image", makeHandle(help_image.Browse))
	http.HandleFunc("/help_image_insert", makeHandle(help_image.Insert))
	http.HandleFunc("/help_image_update_image", makeHandle(help_image.UpdateImage))
	http.HandleFunc("/help_image_delete", makeHandle(help_image.Delete))
	http.HandleFunc("/help_image_download", makeHandle(help_image.Download))
	http.HandleFunc("/help_image_small", makeHandle(help_image.ImageSmall))
	http.HandleFunc("/help_image_middle", makeHandle(help_image.ImageMiddle))
	http.HandleFunc("/help_image_original", makeHandle(help_image.ImageOriginal))

	http.HandleFunc("/jumbotron", makeHandle(jumbotron.Browse))
	http.HandleFunc("/jumbotron_insert", makeHandle(jumbotron.Insert))
	http.HandleFunc("/jumbotron_update", makeHandle(jumbotron.Update))
	http.HandleFunc("/jumbotron_delete", makeHandle(jumbotron.Delete))

	http.HandleFunc("/faq", makeHandle(faq.Browse))
	http.HandleFunc("/faq_insert", makeHandle(faq.Insert))
	http.HandleFunc("/faq_update", makeHandle(faq.Update))
	http.HandleFunc("/faq_delete", makeHandle(faq.Delete))

	http.HandleFunc("/news", makeHandle(news.Browse))
	http.HandleFunc("/news_insert", makeHandle(news.Insert))
	http.HandleFunc("/news_update", makeHandle(news.Update))
	http.HandleFunc("/news_delete", makeHandle(news.Delete))

	http.HandleFunc("/social", makeHandle(social.Browse))
	http.HandleFunc("/social_insert", makeHandle(social.Insert))
	http.HandleFunc("/social_update", makeHandle(social.Update))
	http.HandleFunc("/social_delete", makeHandle(social.Delete))

	http.HandleFunc("/news_image", makeHandle(news_image.Browse))
	http.HandleFunc("/news_image_insert", makeHandle(news_image.Insert))
	http.HandleFunc("/news_image_update", makeHandle(news_image.UpdateImage))
	http.HandleFunc("/news_image_update_info", makeHandle(news_image.UpdateImageInfo))
	http.HandleFunc("/news_image_delete", makeHandle(news_image.Delete))
	http.HandleFunc("/news_image_download", makeHandle(news_image.Download))
	http.HandleFunc("/news_image_small", makeHandle(news_image.ImageSmall))
	http.HandleFunc("/news_image_middle", makeHandle(news_image.ImageMiddle))
	http.HandleFunc("/news_image_original", makeHandle(news_image.ImageOriginal))

	http.HandleFunc("/text_content", makeHandle(text_content.Browse))
	http.HandleFunc("/text_content_insert", makeHandle(text_content.Insert))
	http.HandleFunc("/text_content_update", makeHandle(text_content.Update))
	http.HandleFunc("/text_content_set", makeHandle(text_content.Set))
	http.HandleFunc("/text_content_delete", makeHandle(text_content.Delete))
	http.HandleFunc("/text_content_display", makeHandle(text_content.Display))

	http.HandleFunc("/text_content_image", makeHandle(text_content_image.Browse))
	http.HandleFunc("/text_content_image_insert", makeHandle(text_content_image.Insert))
	http.HandleFunc("/text_content_image_update", makeHandle(text_content_image.UpdateImage))
	http.HandleFunc("/text_content_image_update_info", makeHandle(text_content_image.UpdateImageInfo))
	http.HandleFunc("/text_content_image_delete", makeHandle(text_content_image.Delete))
	http.HandleFunc("/text_content_image_download", makeHandle(text_content_image.Download))
	http.HandleFunc("/text_content_image_small", makeHandle(text_content_image.ImageSmall))
	http.HandleFunc("/text_content_image_middle", makeHandle(text_content_image.ImageMiddle))
	http.HandleFunc("/text_content_image_original", makeHandle(text_content_image.ImageOriginal))

	http.HandleFunc("/manufacturer", makeHandle(manufacturer.Browse))
	http.HandleFunc("/manufacturer_insert", makeHandle(manufacturer.Insert))
	http.HandleFunc("/manufacturer_update", makeHandle(manufacturer.Update))
	http.HandleFunc("/manufacturer_delete", makeHandle(manufacturer.Delete))

	http.HandleFunc("/machine_video", makeHandle(machine_video.Browse))
	http.HandleFunc("/machine_video_insert", makeHandle(machine_video.Insert))
	http.HandleFunc("/machine_video_update", makeHandle(machine_video.Update))
	http.HandleFunc("/machine_video_delete", makeHandle(machine_video.Delete))

	http.HandleFunc("/machine", makeHandle(machine.Browse))
	http.HandleFunc("/machine_display", makeHandle(machine.Display))
	http.HandleFunc("/machine_popup", makeHandle(machine.Popup))
	http.HandleFunc("/machine_insert", makeHandle(machine.Insert))
	http.HandleFunc("/machine_update", makeHandle(machine.Update))
	http.HandleFunc("/machine_delete", makeHandle(machine.Delete))

	http.HandleFunc("/machine_image", makeHandle(machine_image.Browse))
	http.HandleFunc("/machine_image_insert", makeHandle(machine_image.Insert))
	http.HandleFunc("/machine_image_update", makeHandle(machine_image.UpdateImage))
	http.HandleFunc("/machine_image_update_info", makeHandle(machine_image.UpdateImageInfo))
	http.HandleFunc("/machine_image_crop", makeHandle(machine_image.CropImage))
	http.HandleFunc("/machine_image_delete", makeHandle(machine_image.Delete))
	http.HandleFunc("/machine_image_download", makeHandle(machine_image.Download))
	http.HandleFunc("/machine_image_small", makeHandle(machine_image.ImageSmall))
	http.HandleFunc("/machine_image_middle", makeHandle(machine_image.ImageMiddle))
	http.HandleFunc("/machine_image_original", makeHandle(machine_image.ImageOriginal))
	http.HandleFunc("/machine_image_make_header", makeHandle(machine_image.MakeHeaderImage))
	http.HandleFunc("/machine_image_rotate", makeHandle(machine_image.ImageRotate))

	http.HandleFunc("/machine_pdf", makeHandle(machine_pdf.Browse))
	http.HandleFunc("/machine_pdf_insert", makeHandle(machine_pdf.Insert))
	http.HandleFunc("/machine_pdf_update", makeHandle(machine_pdf.UpdatePdf))
	http.HandleFunc("/machine_pdf_update_info", makeHandle(machine_pdf.UpdatePdfInfo))
	http.HandleFunc("/machine_pdf_delete", makeHandle(machine_pdf.Delete))
	http.HandleFunc("/machine_pdf_download", makeHandle(machine_pdf.Download))

	http.HandleFunc("/category", makeHandle(category.Browse))
	http.HandleFunc("/category_insert", makeHandle(category.Insert))
	http.HandleFunc("/category_update", makeHandle(category.Update))
	http.HandleFunc("/category_delete", makeHandle(category.Delete))

	http.HandleFunc("/office", makeHandle(office.Browse))
	http.HandleFunc("/office_insert", makeHandle(office.Insert))
	http.HandleFunc("/office_update", makeHandle(office.Update))
	http.HandleFunc("/office_make_main", makeHandle(office.MakeMainOffice))
	http.HandleFunc("/office_delete", makeHandle(office.Delete))
	http.HandleFunc("/office_display", makeHandle(office.Display))

	http.HandleFunc("/", makeHandle(login.Browse))

	//timeout for old browsers and old mobile devices
	srv := &http.Server{
		Addr:         app.Ini.Host + ":" + app.Ini.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("%s starting server on port: %s\n", app.SiteName, app.Ini.Port)
	log.Printf("%s starting server on port: %s\n", app.SiteName, app.Ini.Port)
	err := srv.ListenAndServe()

	if err != nil {
		println(err.Error())
		log.Println(err.Error())
	}
}

func makeHandle(hf HandleFuncType) HandleFuncType {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch err := err.(type) {
				case string:
					fatalError(rw, req, err)
				case error:
					fatalError(rw, req, err.Error())
				case *app.BadRequestError:
					badRequest(rw, req)
				case *app.NotFoundError:
					notFound(rw, req)
				default:
					fmt.Printf("%s %v", "unknown error", err)
					log.Printf("%s %v\n", "unknown error", err)
				}
			}
		}()

		if app.Ini.Cache {
			rw.Header().Add("Cache-Control", "public, max-age=3600, must-revalidate")
		} else {
			rw.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		}

		hf(rw, req)
	}
}

func BasicAuth(hf HandleFuncType, username, password, realm string) HandleFuncType {

	return func(rw http.ResponseWriter, req *http.Request) {

		user, pass, ok := req.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			rw.Header().Set("WWW-Authenticate", "Basic realm=\""+realm+"\"")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorised."))
			return
		}

		hf(rw, req)
	}
}

func fatalError(rw http.ResponseWriter, req *http.Request, msg string) {
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		list := strings.Split(string(debug.Stack()[:]), "\n")
		errorMsg := fmt.Sprintf("%s %s", app.T("Fatal Error:"), app.T(msg))

		if app.Ini.Debug {
			errorMsg += "\n" + strings.Join(list[6:], "\n")
		}

		http.Error(rw, errorMsg, http.StatusInternalServerError)
		return
	}

	buf := new(util.Buf)

	buf.Add("<!DOCTYPE html>")
	buf.Add("<html lang=\"en\">")

	buf.Add("<head>")
	buf.Add("<meta charset=\"utf-8\">")
	buf.Add("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	buf.Add("<link rel=\"icon\" href=\"/asset/img/favicon.png\">")
	buf.Add("<link rel=\"stylesheet\" href=\"/asset/css/style.css\">")
	buf.Add("<link rel=\"stylesheet\" href=\"/asset/css/page/fatal_error.css\">")
	buf.Add("<title>%s</title>", app.T("Fatal Error"))
	buf.Add("</head>")

	buf.Add("<body id=\"fatalErrorPage\">")

	buf.Add("<main>")
	buf.Add("<div class=\"row\">")

	buf.Add("<div class=\"col colSpan md1 lg2 xl3\"></div>")

	buf.Add("<div class=\"col md10 lg8 xl6 panel panelError\">")

	buf.Add("<div class=\"panelHeading\">")
	buf.Add("<h3>%s</h3>", app.T("Fatal Error"))
	buf.Add("</div>")

	buf.Add("<div class=\"panelBody\">")
	buf.Add("<p>%s</p>", app.T(msg))

	list := strings.Split(string(debug.Stack()[:]), "\n")
	buf.Add("<p>%s</p>", strings.Join(list[6:], "<br>"))
	buf.Add("</div>")

	buf.Add("</div>") //end of col

	buf.Add("<div class=\"col colSpan md1 lg2 xl3\"></div>")

	buf.Add("</div>") //end of row

	buf.Add("</main>")

	buf.Add("</body>")
	buf.Add("</html>")

	fmt.Fprintf(rw, *buf.String())
}

func badRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := context.NewContext(rw, req)

	ctx.Rw.WriteHeader(http.StatusBadRequest)

	buf := util.NewBuf()
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h3>%s</h3>", app.T("Bad request."))
	buf.Add("<p>%s</p>", app.T("You don't have right to access this site."))
	buf.Add("<p><a href=\"/login\">%s</a></p>", app.T("Login page."))
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	str := "endPage"
	ctx.AddHtml("pageName", &str)

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/end.css")
	content.Default(ctx)

	ctx.Render("end.html")
}

func notFound(rw http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	ctx := context.NewContext(rw, req)
	ctx.Rw.WriteHeader(http.StatusNotFound)

	buf := util.NewBuf()
	buf.Add("<div class=\"callout calloutError\">")
	buf.Add("<h3>%s</h3>", app.T("404 error."))
	buf.Add("<p><img src=\"/asset/img/404.jpg\" width=\"100%%\"></p>")
	buf.Add("<p>%s</p>", app.T("The page you are looking for does not exists."))
	buf.Add("</div>")

	ctx.AddHtml("midContent", buf.String())

	str := "endPage"
	ctx.AddHtml("pageName", &str)

	content.Include(ctx)
	ctx.Css.Add("/asset/css/page/end.css")
	content.Default(ctx)

	ctx.Render("end.html")
}

func cron() {
	//clear old sessions
	go func() {
		c := time.Tick(time.Duration(app.Ini.CookieExpires) * time.Second)
		for _ = range c {
			now := time.Now()
			epoch := now.Unix()

			tx, err := app.Db.Begin()
			if err != nil {
				panic(err)
			}

			sqlStr := `delete from 	session
							where (? - session.recordTime) > ? and
							sessionType = 'web'`
			_, err = tx.Exec(sqlStr, epoch, app.Ini.CookieExpires)
			if err != nil {
				tx.Rollback()
				log.Printf("%s %s", "error while deleting expired web sessions:", err.Error())
				return
			}

			tx.Commit()
		}
	}()
}
