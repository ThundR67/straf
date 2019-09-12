package straf

import (
	"errors"
	"reflect"

	"github.com/graphql-go/graphql"
)

//GetGraphQLObject Converts struct into graphql object
func GetGraphQLObject(object interface{}) (*graphql.Object, error) {
	objectType := reflect.TypeOf(object)
	fields := convertStruct(objectType)

	output := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   objectType.Name(),
			Fields: fields,
		},
	)

	return output, nil
}

//convertStructToObject converts simple struct to graphql object
func convertStructToObject(
	objectType reflect.Type) *graphql.Object {

	fields := convertStruct(objectType)

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   objectType.Name(),
			Fields: fields,
		},
	)
}

//convertStruct converts struct to graphql fields
func convertStruct(objectType reflect.Type) graphql.Fields {
	fields := graphql.Fields{}

	for i := 0; i < objectType.NumField(); i++ {
		currentField := objectType.Field(i)
		fieldType := getFieldType(currentField)

		fields[currentField.Name] = &graphql.Field{
			Name:              currentField.Name,
			Type:              fieldType,
			DeprecationReason: getTagValue(currentField, "deprecationReason"),
			Description:       getTagValue(currentField, "description"),
		}
	}

	return fields
}

//getFieldType Converts object to a graphQL field type
func getFieldType(object reflect.StructField) graphql.Output {

	isID, ok := object.Tag.Lookup("unique")
	if isID == "true" && ok {
		return graphql.ID
	}

	objectType := object.Type
	if objectType.Kind() == reflect.Struct {
		return convertStructToObject(objectType)

	} else if objectType.Kind() == reflect.Slice &&
		objectType.Elem().Kind() == reflect.Struct {

		elemType := convertStructToObject(objectType.Elem())
		return graphql.NewList(elemType)

	} else if objectType.Kind() == reflect.Slice {
		elemType, _ := convertSimpleType(objectType.Elem())
		return graphql.NewList(elemType)
	}

	output, _ := convertSimpleType(objectType)
	return output
}

//convertSimpleType converts simple type  to graphql field
func convertSimpleType(objectType reflect.Type) (*graphql.Scalar, error) {

	typeMap := map[reflect.Kind]*graphql.Scalar{
		reflect.String:  graphql.String,
		reflect.Bool:    graphql.Boolean,
		reflect.Int:     graphql.Int,
		reflect.Int8:    graphql.Int,
		reflect.Int16:   graphql.Int,
		reflect.Int32:   graphql.Int,
		reflect.Int64:   graphql.Int,
		reflect.Float32: graphql.Float,
		reflect.Float64: graphql.Float,
	}

	graphqlType, ok := typeMap[objectType.Kind()]

	if !ok {
		return &graphql.Scalar{}, errors.New("Invalid Type")
	}

	return graphqlType, nil
}

//getTagValue returns tag value of a struct
func getTagValue(objectType reflect.StructField, tagName string) string {
	tag := objectType.Tag
	value, ok := tag.Lookup(tagName)
	if !ok {
		return ""
	}
	return value
}
