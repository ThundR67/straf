package straf

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
)

//GetGraphQLObject Converts struct into graphql object
func GetGraphQLObject(object interface{}) (*graphql.Object, error) {
	objectType := reflect.TypeOf(object)
	fields, err := convertStruct(objectType)

	output := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   objectType.Name(),
			Fields: fields,
		},
	)

	if err != nil {
		err = fmt.Errorf("Error While Converting Struct To GraphQL Object: %v", err)
		return &graphql.Object{}, err
	}

	return output, nil
}

//convertStructToObject converts simple struct to graphql object
func convertStructToObject(
	objectType reflect.Type) (*graphql.Object, error) {

	fields, err := convertStruct(objectType)
	if err != nil {
		err = fmt.Errorf(
			"Error while converting type %v to graphql fields: %v",
			objectType,
			err,
		)
		return &graphql.Object{}, err
	}

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   objectType.Name(),
			Fields: fields,
		},
	), nil
}

//convertStruct converts struct to graphql fields
func convertStruct(objectType reflect.Type) (graphql.Fields, error) {
	fields := graphql.Fields{}

	for i := 0; i < objectType.NumField(); i++ {
		currentField := objectType.Field(i)

		if currentField.Type.Kind() == reflect.Struct {
			graphqlObject, err := convertStructToObject(currentField.Type)
			if err != nil {
				err = fmt.Errorf(
					"Error while converting type %v to graphql object: %v",
					currentField.Type,
					err,
				)
				return graphql.Fields{}, err
			}

			fields[currentField.Name] = &graphql.Field{
				Name: currentField.Name,
				Type: graphqlObject,
			}
		} else if currentField.Type.Kind() == reflect.Slice &&
			currentField.Type.Elem().Kind() == reflect.Struct {

			graphqlObject, err := convertStructToObject(currentField.Type.Elem())
			if err != nil {
				err = fmt.Errorf(
					"Error while converting slice type %v to graphql object: %v",
					currentField.Type,
					err,
				)
				return graphql.Fields{}, err
			}

			fields[currentField.Name] = &graphql.Field{
				Name: currentField.Name,
				Type: graphql.NewList(graphqlObject),
			}

		} else {
			field, err := convertSimpleType(
				currentField.Type,
				currentField.Type.Kind() == reflect.Slice,
			)
			if err != nil {
				return nil, err
			}
			fields[currentField.Name] = field
		}
	}

	return fields, nil
}

//convertSimpleType converts simple type or slice of simple type to graphql field
func convertSimpleType(objectType reflect.Type, isList bool) (*graphql.Field, error) {

	if isList {
		objectType = objectType.Elem()
	}

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

	if !ok && !isList {
		return &graphql.Field{}, errors.New("Invalid Type")

	} else if !isList {
		return &graphql.Field{
			Type: graphqlType,
		}, nil
	}

	return &graphql.Field{
		Type: graphql.NewList(graphqlType),
	}, nil
}
