package tran

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"backend/src/app"
	"backend/src/lib/context"
	"backend/src/lib/util"
	"backend/src/page/tran/tran_lib"
)

func Sync(rw http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(rw, req)

	if !ctx.IsRight("tran", "synchronize") {
		app.BadRequest()
	}

	ctx.Cargo.AddStr("key", "")
	ctx.Cargo.AddInt("pn", 1)
	ctx.ReadCargo()

	var tranList map[string]bool

	dirList := []string{app.Ini.HomeDir + "/src",
		app.Ini.HomeDir + "/../backend/src"}
	exclude := []string{app.Ini.HomeDir + "/src/github.com",
		app.Ini.HomeDir + "/src/golang.org",
		app.Ini.HomeDir + "/../backend/src/github.com",
		app.Ini.HomeDir + "/../backend/src/golang.org"}

	fileList := make([]string, 0, 50)
	for _, val := range dirList {
		getFileList(val, exclude, &fileList)
	}

	tranList = getTranUnitsFromSource(fileList)

	recordList := tran_lib.GetTranMap()

	deleted := deleteUnused(tranList, recordList)
	absent := addAbsent(tranList, recordList)

	if absent == 0 && deleted == 0 {
		ctx.Msg.Warning(ctx.T("Translation table is old."))
		ctx.Msg.Warning(ctx.T("If there are some missing texts in translatin table, please contact site admin."))
	} else {
		saveTran(recordList)
		ctx.Msg.Success(ctx.T("Translation table has been synchronized."))
	}

	ctx.Redirect(ctx.U("/tran", "key", "pn"))
}

func getFileList(dir string, exclude []string, rv *[]string) {
	for _, val := range exclude {
		if strings.HasPrefix(dir, val) {
			return
		}
	}

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, val := range fileInfo {
		if val.IsDir() {
			getFileList(dir+"/"+val.Name(), exclude, rv)
		}

		if path.Ext(val.Name()) == ".go" {
			*rv = append(*rv, dir+"/"+val.Name())
		}
	}
}

func getTranUnitsFromSource(fileList []string) map[string]bool {
	ctxTRe := regexp.MustCompile(`(?m)(?s)ctx\.T\("(.*?)"\)`)
	appTRe := regexp.MustCompile(`(?m)(?s)app\.T\("(.*?)"\)`)
	TRe := regexp.MustCompile(`(?m)(?s)T\("(.*?)"\)`)
	reNewline := regexp.MustCompile(`(?s)"\s*?\+\s*?"`)
	rv := make(map[string]bool, 400)

	for _, val := range fileList {
		lines, err := ioutil.ReadFile(val)
		if err != nil {
			panic(err)
		}

		units := ctxTRe.FindAll(lines, -1)
		prefix := fmt.Sprintf("%c%s", 'c', "tx.T(\"")
		var str string
		for _, val := range units {
			str = string(val)

			str = reNewline.ReplaceAllString(str, "")
			str = strings.TrimPrefix(str, prefix)
			str = strings.TrimSuffix(str, "\")")
			rv[str] = true
		}

		units = appTRe.FindAll(lines, -1)
		prefix = fmt.Sprintf("%c%s", 'a', "pp.T(\"")
		for _, val := range units {
			str = string(val)

			str = reNewline.ReplaceAllString(str, "")
			str = strings.TrimPrefix(str, prefix)
			str = strings.TrimSuffix(str, "\")")
			rv[str] = true
		}

		units = TRe.FindAll(lines, -1)
		prefix = fmt.Sprintf("%c%s", 'T', "(\"")
		for _, val := range units {
			str = string(val)

			str = reNewline.ReplaceAllString(str, "")
			str = strings.TrimPrefix(str, prefix)
			str = strings.TrimSuffix(str, "\")")
			rv[str] = true
		}
	}

	return rv
}

func deleteUnused(tl map[string]bool, rl map[string]*tran_lib.TranRec) int64 {
	var rv int64

	for key, _ := range rl {
		if !tl[key] {
			rv++
			delete(rl, key)
		}
	}

	return rv
}

func addAbsent(tl map[string]bool, rl map[string]*tran_lib.TranRec) int64 {
	var rv int64

	for key, _ := range tl {
		if _, ok := rl[key]; !ok {
			rv++
			rl[key] = &tran_lib.TranRec{En: key, Tr: ""}
		}
	}

	return rv
}

func saveTran(rl map[string]*tran_lib.TranRec) {
	if len(rl) == 0 {
		return
	}

	tx, err := app.Db.Begin()
	if err != nil {
		panic(err)
	}

	sqlStr := "delete from tran"
	_, err = tx.Exec(sqlStr)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	buf := util.NewBuf()
	buf.Add("insert into tran(tranId, en, tr) values")

	keys := make([]string, 0, 400)
	for _, val := range rl {
		keys = append(keys,
			fmt.Sprintf("(null, '%s', '%s')", util.DbStr(val.En), util.DbStr(val.Tr)))
	}

	buf.Add(strings.Join(keys, ",\n"))

	_, err = tx.Exec(*buf.String())
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	return
}
