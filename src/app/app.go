package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/microcosm-cc/bluemonday"
	"github.com/natefinch/lumberjack"
)

const MIN_YOM int64 = 1950

const SUPSER_USER_ID int64 = 30
const BUYER_SEGMENT_ID int64 = 1
const BEGINING_PACKAGE_ID int64 = 2

type BadRequestError error
type NotFoundError error
type MaintenanceError error

func BadRequest() {
	panic(new(BadRequestError))
}

func NotFound() {
	panic(new(NotFoundError))
}

func Maintenance() {
	panic(new(MaintenanceError))
}

var Debug = true
var SiteName = "backend"

var Tmpl *template.Template
var Db *sql.DB
var Tran map[string]string

var StrictHtml *bluemonday.Policy

var TranRe *regexp.Regexp //regular expression for translation

func init() {
	readIni()
	connect()
	ReadTran()
	ReadConfig()

	log.SetOutput(&lumberjack.Logger{
		Filename:   Ini.HomeDir + "/logs/error.log",
		MaxSize:    10, // megabytes
		MaxBackups: 6,
		MaxAge:     28, //days
	})

	funcmap := template.FuncMap{
		"raw": func(str string) template.HTML {
			return template.HTML(str)
		},
	}

	Tmpl = template.Must(template.New("").Funcs(funcmap).ParseGlob(Ini.HomeDir + "/view/*.html"))

	//html policies
	StrictHtml = bluemonday.StrictPolicy()
}

func T(str string) string {
	if config.Str("lang") == "en" {
		return str
	}

	if val, ok := Tran[str]; ok && val != "" {
		return val
	}

	return "[" + str + "]"
}

func connect() {
	conStr := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?autocommit=false",
		Ini.DbUser, Ini.DbPassword, Ini.DbProtocol,
		Ini.DbHost, Ini.DbPort, Ini.DbName)
	db, err := sql.Open("mysql", conStr)

	if err != nil {
		panic("Could not connec to database.")
	}

	db.SetMaxIdleConns(0)
	Db = db
}

func ReadTran() {
	Tran = make(map[string]string, 300)

	type tranRec struct {
		En string
		Tr string
	}

	sqlStr := "select en, tr from tran"

	rows, err := Db.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		r := new(tranRec)
		if err = rows.Scan(&r.En, &r.Tr); err != nil {
			panic(err)
		}

		Tran[r.En] = r.Tr
	}
}
