package goinception

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"testing"
)

func TestGoInception_CheckAuditSql(t *testing.T) {

	g := NewGoInception(Options{
		User:     "root",
		Password: "root",
		Host:     "10.13.84.220:4000",
	})

	sql := ` use  test;
   create table t1(id int primary key);
   alter table t1 add index idx_id (id);
   create table t2(jid int primary key);`

	r, s, err := g.parseSQLs(sql)
	if err != nil {
		panic(err)
	}
	rr, _ := json.Marshal(r)
	ss, _ := json.Marshal(s)
	fmt.Println(string(rr))
	fmt.Println(string(ss))
}

func TestNewGoInception(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(10.13.84.220:4000)/")
	defer db.Close()

	sql := `/*--user=root;--password=root;--host=10.13.84.220;--port=4000;--check=true;*/
    inception_magic_start;
	inception set osc_on = 1;
--    create table t1(id int primary key);
--    alter table t1 add index idx_id (id);
--    create table t2(jid int primary key);
   inception_magic_commit;`

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	fmt.Println(cols)

	type row struct {
		OrderId      int     `db:"order_id"`
		Stage        *string `db:"stage"`
		ErrorLevel   *int    `db:"error_level"`
		StageStatus  *string `db:"stage_status"`
		ErrorMsg     *string `db:"error_message"`
		Sql          *string `db:"sql"`
		AffectedRows int     `db:"affected_rows"`
		Sequence     *string `db:"sequence"`
		BackupDbName *string `db:"backup_dbname"`
		ExecuteTime  *string `db:"execute_time"`
		SqlSha       *string `db:"sqlsha1"`
		BackupTime   *string `db:"backup_time"`
	}
	got := []row{}
	for rows.Next() {
		var r row
		err = rows.Scan(&r.OrderId, &r.Stage, &r.ErrorLevel, &r.StageStatus, &r.ErrorMsg, &r.Sql, &r.AffectedRows, &r.Sequence, &r.BackupDbName, &r.ExecuteTime, &r.SqlSha, &r.BackupTime)
		if err != nil {
			t.Fatalf("Scan: %v", err)
		}
		got = append(got, r)
	}
	err = rows.Err()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	res, _ := json.Marshal(got)
	fmt.Println(string(res))
}

type executeContext struct {
	dbName string
}

func TestNewGoInception2(t *testing.T) {
	s := &executeContext{dbName: "wer"}
	execStr, _ := json.Marshal(s)
	fmt.Println(fmt.Sprintf("[inception] execCtx : %s", string(execStr)))
}

func TestNewGoInception3(t *testing.T) {
	sql := "GET /my_index/_search {};"
	headerIndex := strings.Index(sql, "{")
	if headerIndex < 5 {
		return
	}

	header := sql[:headerIndex]
	body := sql[headerIndex:]
	if !json.Valid([]byte(body)) {
		return
	}
	methodIndex := strings.Index(header, " ")
	method := header[:methodIndex]
	path := header[methodIndex:]

	fmt.Println(fmt.Sprintf("%s:%s:%s", method, path, body))

}
