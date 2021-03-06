package gomysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", "cgctest:cgc123456@tcp(182.254.185.142:3306)/godb1?charset=utf8")
	checkErr(err)
}

func Write(time string, imei string, packet string) {
	if DB == nil {
		return
	}
	stmt, err := DB.Prepare("insert into gomysql(time,imei,packet)values(?,?,?)")
	checkErr(err)
	defer stmt.Close()
	if result, err := stmt.Exec(time, imei, packet); err == nil {
		if id, err := result.LastInsertId(); err == nil {
			fmt.Println("insert id : ", id)
		}
	}
}

func Read(str string) {
	if DB == nil {
		return
	}
	buf := strings.Join([]string{"SELECT * FROM gomysql where imei ='", str, "'"}, "")

	rows, err := DB.Query(buf)
	checkErr(err)
	for rows.Next() {
		var time, imei, packet string
		var id int
		rows.Columns()
		err = rows.Scan(&id, &time, &imei, &packet)
		checkErr(err)
		fmt.Println(id, time, imei, packet)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
