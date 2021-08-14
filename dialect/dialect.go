package dialect

import "reflect"

type Dialect interface {
	// 将Go语言的类型转换为某个数据库的数据类型
	DataTypeOf(t reflect.Value) string

	// 返回某个表是否存在的SQL语句，参数是表名
	TableExistSQL(tableName string) (string, []interface{})
}

var dialectsMap = map[string]Dialect{}

// 注册dialect实例
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 获取dialect实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
