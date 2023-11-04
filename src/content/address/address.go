package address

import (
	"strings"

	"backend/src/lib/util"
	"backend/src/page/office/office_lib"
)

func Get() (string, error) {
	buf := util.NewBuf()

	buf.Add("<h2>Contact Us</h2>")

	rec, err := office_lib.GetMainOfficeRec()
	if err != nil {
		return "", err
	}

	buf.Add("<p>")
	buf.Add("%s %s", rec.Address, rec.City)

	if rec.MapLink != "" {
		buf.Add("<a href=\"%s\" target=\"_blank\" class=\"mapLink\"><i class=\"fas fa-map-marker-alt\"></i></a>", rec.MapLink)
	}

	buf.Add("</p>")

	phone := util.ScrStr(rec.Telephone)
	buf.Add("<p><a href=\"tel:%s\">%s</a></p>", strings.Replace(phone, " ", "", -1), phone)

	buf.Add("<p>%s</p>", rec.Email)

	return *buf.String(), nil
}
