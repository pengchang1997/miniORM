package schema

import (
	"go/ast"
	"miniORM/dialect"
	"reflect"
)

// 字段定义
type Field struct {
	// 字段名
	Name string

	// 字段类型
	Type string

	// 约束条件
	Tag string
}

// 表定义
type Schema struct {
	// 被映射的对象
	Model interface{}

	// 表名
	Name string

	// 表中的字段
	Fields []*Field

	// 字段名
	FieldNames []string

	// 存储字段名和字段之间的映射关系
	fieldMap map[string]*Field
}

// 根据字段名获取字段
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 将任意的对象解析为Schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 获取对象的类型
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	// 创建Schema实例
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 只考虑非匿名字段和导出字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}

	return schema
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
