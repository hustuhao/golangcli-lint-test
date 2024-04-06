package db

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MainDB *sqlx.DB // 主库

func ConnectDB() {
	MainDB = sqlx.MustConnect("mysql", "dsn")
	MainDB.SetMaxOpenConns(100)
}

func InitDB() {
	if err := excuteSqlFile("../../script/db.sql"); err != nil {
		panic(err)
	}
}

func excuteSqlFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	script, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	sqlls := strings.Split(string(script), ";")

	for _, sqll := range sqlls {
		sqll = strings.Trim(sqll, " \t\r\n")
		if sqll == "" {
			continue
		}
		log.Println(sqll)
		r, err := MainDB.Exec(sqll)
		if err != nil {
			return err
		}
		affect, _ := r.RowsAffected()
		log.Println("affect", affect)
	}

	return nil
}

// 根据结构体的db-tag生成sqlx的insert语句
// insert into xxx (`col1`, `col2`) values (:col1, :col2)
func BuildSqlxInsert(table string, s interface{}) string {
	t := reflect.TypeOf(s)
	var v reflect.Value
	if t.Kind() == reflect.Ptr {
		v = reflect.Indirect(reflect.ValueOf(s))
		t = v.Type()
	} else {
		v = reflect.ValueOf(s)
	}
	tags := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag := f.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}

		tags = append(tags, tag)
	}

	cols := make([]string, 0)
	for _, v := range tags {
		col := "`" + v + "`"
		cols = append(cols, col)
	}

	sqll := "insert into `" + table + "` ("

	sqll += strings.Join(cols, ", ")

	sqll += ") values ("

	for i := range tags {
		tags[i] = ":" + tags[i]
	}

	sqll += strings.Join(tags, ", ")

	sqll += ")"

	return sqll
}
