package session

import (
	"database/sql"
	"miniORM/log"
	"strings"
)

type Session struct {
	// 使用sql.Open()方法连接数据库成功之后返回的指针
	db *sql.DB

	// SQL语句
	sql strings.Builder

	// SQL语句中占位符的对应值
	sqlVars []interface{}
}

// 创建一个Session
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// 清空Session
func (session *Session) Clear() {
	session.sql.Reset()
	session.sqlVars = nil
}

// 获取Session中的sql.DB指针
func (session *Session) DB() *sql.DB {
	return session.db
}

func (session *Session) Raw(sql string, values ...interface{}) *Session {
	session.sql.WriteString(sql)
	session.sql.WriteString(" ")
	session.sqlVars = append(session.sqlVars, values...)
	return session
}

// 封装原生方法
// 封装原生方法的目的在于统一日志打印，包括执行的SQL语句和错误日志
// 另外，这些方法在SQL语句执行完成后会清空(s *Session).sql和(s *Session).sqlVars两个变量
// 这样做可以复用Session，即开启一次会话后可以多次执行SQL语句

// 封装Exec方法
func (session *Session) Exec() (result sql.Result, err error) {
	defer session.Clear()
	log.Info(session.sql.String(), session.sqlVars)
	if result, err = session.DB().Exec(session.sql.String(), session.sqlVars...); err != nil {
		log.Error(err)
	}

	return
}

// 封装QueryRow方法
func (session *Session) QueryRow() *sql.Row {
	defer session.Clear()
	log.Info(session.sql.String(), session.sqlVars)
	return session.DB().QueryRow(session.sql.String(), session.sqlVars...)
}

// 封装Query方法
func (session *Session) QueryRows() (rows *sql.Rows, err error) {
	defer session.Clear()
	log.Info(session.sql.String(), session.sqlVars)
	if rows, err = session.DB().Query(session.sql.String(), session.sqlVars...); err != nil {
		log.Error(err)
	}

	return
}
