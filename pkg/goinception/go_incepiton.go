package goinception

import (
	"database/sql"
	"fmt"
)

type Options struct {
	User     string // goinception 用户
	Password string // goinception 密码
	Host     string // goinception 服务地址
}

type GoInception struct {
	option Options
}

func NewGoInception(options Options) *GoInception {
	return &GoInception{
		option: options,
	}
}

//sql := `/*--user=root;--password=;--host=10.13.255.184;--port=3306;--check=true;*/
//  inception_magic_start;
//  use  test;
//  create table t1(id int primary key);
//  alter table t1 add index idx_id (id);
//  create table t2(jid int primary key);
//  inception_magic_commit;`

const InceptionStart = "inception_magic_start;"
const InceptionCommit = "inception_magic_commit;"

func (g *GoInception) CheckAuditSql(db DbInfo, sql string) ([]InceptionRow, error) {

	prefix := fmt.Sprintf("/*--user=%s;--password=%s;--host=%s;--port=%s;--check=true;*/", db.User, db.Password, db.Host, db.Port)
	sqlStr := fmt.Sprintf("%s%s%s%s", prefix, InceptionStart, sql, InceptionCommit)
	return g.doSql(sqlStr)
}

func (g *GoInception) ExecSql(db DbInfo, sql string, backUp bool) ([]InceptionRow, error) {

	prefix := fmt.Sprintf("/*--user=%s;--password=%s;--host=%s;--port=%s;--execute=true;--backup=%t;*/", db.User, db.Password, db.Host, db.Port, backUp)
	sqlStr := fmt.Sprintf("%s%s%s%s", prefix, InceptionStart, sql, InceptionCommit)
	return g.doSql(sqlStr)
}

func (g *GoInception) doSql(sqlStr string) ([]InceptionRow, error) {
	dateSourceName := fmt.Sprintf("%s:%s@tcp(%s)/", g.option.User, g.option.Password, g.option.Host)
	db, err := sql.Open("mysql", dateSourceName)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	got := []InceptionRow{}
	for rows.Next() {
		var r InceptionRow
		err = rows.Scan(&r.OrderId, &r.Stage, &r.ErrorLevel, &r.StageStatus, &r.ErrorMsg, &r.Sql, &r.AffectedRows, &r.Sequence, &r.BackupDbName, &r.ExecuteTime, &r.SqlSha, &r.BackupTime)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return got, nil
}

type DbInfo struct {
	User     string
	Password string
	Host     string
	Port     string
}

type InceptionRow struct {
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
