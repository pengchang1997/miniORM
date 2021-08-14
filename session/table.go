package session

import (
	"fmt"
	"miniORM/log"
	"miniORM/schema"
	"reflect"
	"strings"
)

// 为refTable赋值
func (session *Session) Model(value interface{}) *Session {
	// 如果refTable为nil或者传入的对象类型发生了变化，则重新解析对象
	if session.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(session.refTable.Model) {
		session.refTable = schema.Parse(value, session.dialect)
	}

	return session
}

// 获取refTable
func (session *Session) RefTable() *schema.Schema {
	if session.refTable == nil {
		log.Error("Model is not set")
	}

	return session.refTable
}

// 创建表
func (session *Session) CreateTable() error {
	table := session.RefTable()
	var columns []string
	for _, field := range session.refTable.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	desc := strings.Join(columns, ",")
	_, err := session.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, desc)).Exec()
	return err
}

// 删除表
func (session *Session) DropTable() error {
	_, err := session.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", session.RefTable().Name)).Exec()
	return err
}

// 判断表是否存在
func (session *Session) HasTable() bool {
	sql, values := session.dialect.TableExistSQL(session.RefTable().Name)
	row := session.Raw(sql, values...).QueryRow()
	var temp string
	_ = row.Scan(temp)
	return temp == session.RefTable().Name
}
