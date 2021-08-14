package miniORM

import (
	"database/sql"
	"miniORM/dialect"
	"miniORM/log"
	"miniORM/session"
)

// Engine负责数据库的连接与关闭
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// 连接数据库
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	// db.Ping()方法检查与数据库的连接是否仍然有效
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	tempDialect, ok := dialect.GetDialect(driver)
	if !ok {
		log.ErrorF("dialect %s Not Found", driver)
		return
	}

	e = &Engine{
		db:      db,
		dialect: tempDialect,
	}

	log.Info("Connect database success")
	return
}

// 关闭与数据库的连接
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}

	log.Info("Close database success")
}

// 创建一个Session
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
