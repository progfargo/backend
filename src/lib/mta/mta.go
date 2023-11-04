package mta

import (
	"bytes"
	gocontext "context"
	"fmt"
	"io/ioutil"

	"backend/src/app"
	"backend/src/content/address"
	"backend/src/lib/context"
	"backend/src/lib/util"

	"github.com/aymerick/douceur/inliner"
	"github.com/nikoksr/notify/service/sendgrid"
)

type contentType int

const (
	TEXT contentType = iota
	BUTTON
)

type text struct {
	data     string
	dataType contentType
}

type message struct {
	subject   string
	title     string
	content   []*text
	to        []string
	address   string
	copyright string
	formatted string
	inline    string
}

func NewMessage() *message {
	rv := new(message)
	rv.content = make([]*text, 0, 5)
	rv.to = make([]string, 0, 5)

	return rv
}

func (msg *message) SetSubject(subject string) {
	msg.subject = subject
}

func (msg *message) SetTitle(title string) {
	msg.title = title
}

func (msg *message) AddTo(str string) {
	msg.to = append(msg.to, str)
}

func (msg *message) AddText(data string) {
	t := new(text)
	t.data = data
	t.dataType = TEXT

	msg.content = append(msg.content, t)
}

func (msg *message) AddButton(urlStr, label string) {
	t := new(text)
	t.data = fmt.Sprintf("<a href=\"%s\" class=\"btn btnPrimary btnSm\">%s</a>", urlStr, label)
	t.dataType = BUTTON

	msg.content = append(msg.content, t)
}

func (msg *message) AddAddress() {
	msg.address, _ = address.Get()
}

func (msg *message) AddCopyright() {
	msg.copyright = "<p><a href=\"http://legendplex.com\" target=\"_blank\">© LEGENDPLEX</a>. All rights reserved</p>"
}

func (msg *message) GetFormatted() string {
	return msg.formatted
}

func (msg *message) GetInline() string {
	return msg.inline
}

func (msg *message) FormatEmail() {

	if msg.title == "" {
		panic("Title is missing in email.")
	}

	buf := util.NewBuf()

	buf.Add("<table>")
	buf.Add("<tbody>")

	buf.Add("<tr class=\"header\">")
	buf.Add("<td>")
	buf.Add("<p class=\"logo\">")
	buf.Add("<a href=\"http://base.de\" target=\"_blank\"><strong>base.de</strong></a>")
	buf.Add("</p>")
	buf.Add("<p class=\"author\">")
	buf.Add("<a href=\"http://ionova.com\" target=\"_blank\">by ionova</a>")
	buf.Add("</p>")
	buf.Add("</td>")
	buf.Add("<tr>")

	buf.Add("<tr class=\"message first\">")
	buf.Add("<td>")
	buf.Add("<h2>%s</h2>", msg.title)
	buf.Add("</td>")
	buf.Add("<tr>")

	var classStr string
	for i, v := range msg.content {
		if i == len(msg.content)-1 {
			classStr = " last"
		} else if i == 0 {
			classStr = " first"
		} else {
			classStr = ""
		}

		if v.dataType == TEXT {
			classStr += " text"
		} else if v.dataType == BUTTON {
			classStr += " button"
		}

		buf.Add("<tr class=\"message%s\">", classStr)
		buf.Add("<td>")
		buf.Add("<p>%s</p>", v.data)
		buf.Add("</td>")
		buf.Add("<tr>")
	}

	buf.Add("<tr class=\"seperator\">")
	buf.Add("<td></td>")
	buf.Add("<tr>")

	buf.Add("<tr class=\"address\">")
	buf.Add("<td>")
	buf.Add(msg.address)
	buf.Add("</td>")
	buf.Add("<tr>")

	buf.Add("<tr class=\"seperator\">")
	buf.Add("<td></td>")
	buf.Add("<tr>")

	buf.Add("<tr class=\"copyright\">")
	buf.Add("<td>")
	buf.Add(msg.copyright)
	buf.Add("</td>")
	buf.Add("<tr>")

	buf.Add("</tbody>")
	buf.Add("</table>")

	msg.formatted = *buf.String()
}

func (msg *message) InlineEmail(ctx *context.Ctx) error {

	str := "ionova.com"
	ctx.AddHtml("emailHtmlTitle", &str)

	if msg.title == "" {
		panic("Title is missing in email.")
	}

	if msg.subject == "" {
		panic("Subject is missing in email.")
	}

	if len(msg.to) == 0 {
		panic("To address is missing in imail.")
	}

	if msg.formatted == "" {
		panic("Formatted email is empty.")
	}

	ctx.AddHtml("emailMidContent", &msg.formatted)

	cssBuf := new(bytes.Buffer)

	cssBuf.Write([]byte("<style>\n"))

	emailCss, err := ioutil.ReadFile(app.Ini.HomeDir + "/asset/css/page/email.css")
	if err != nil {
		return err
	}

	cssBuf.Write(emailCss)

	cssBuf.Write([]byte("\n</style>"))

	cssText := cssBuf.String()

	ctx.AddHtml("cssText", &cssText)

	mailStr, err := ctx.RenderEmail("email.html")
	if err != nil {
		return err
	}

	msg.inline, err = inliner.Inline(mailStr)
	if err != nil {
		return err
	}

	return nil
}

func (msg *message) Send() error {

	mailService := sendgrid.New("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",
		"support@example.com", "base.com notifier")

	mailService.AddReceivers(msg.to...)

	c := gocontext.Background()
	err := mailService.Send(c, msg.subject, msg.GetInline())
	if err != nil {
		return err
	}

	return nil
}
