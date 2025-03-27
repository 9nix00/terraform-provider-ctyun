package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
)

// StructToTFObjectTypes将结构体转换为types.ObjectType类型
func StructToTFObjectTypes(s interface{}) types.ObjectType {
	result := make(map[string]attr.Type)
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" {
			continue
		}
		var fieldType attr.Type
		switch field.Type {
		case reflect.TypeOf(types.String{}):
			fieldType = types.StringType
		case reflect.TypeOf(types.Bool{}):
			fieldType = types.BoolType
		case reflect.TypeOf(types.Int64{}):
			fieldType = types.Int64Type
		case reflect.TypeOf(types.Float64{}):
			fieldType = types.Float64Type
		case reflect.TypeOf(types.List{}):
			// 这里假设列表元素类型为 String，实际可能需要更复杂处理
			fieldType = types.ListType{ElemType: types.StringType}
		case reflect.TypeOf(types.Map{}):
			// 这里假设映射元素类型为 String，实际可能需要更复杂处理
			fieldType = types.MapType{ElemType: types.StringType}
		default:
			continue
		}
		result[tag] = fieldType
	}
	return types.ObjectType{AttrTypes: result}
}
