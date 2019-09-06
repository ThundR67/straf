package straf

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

//NewSchemaBuilder is used get a new schema builder
func NewSchemaBuilder(graphQLType graphql.Output, object interface{}) *SchemaBuilder {
	builder := SchemaBuilder{
		GraphQLType: graphQLType,
		Object:      object,
	}
	builder.Init()
	return &builder
}

//SchemaBuilder is used to build a schema based on a struct
type SchemaBuilder struct {
	GraphQLType graphql.Output
	Object      interface{}
	Schema      graphql.Fields
	args        graphql.FieldConfigArgument
}

//Init initializes
func (schemaBuilder *SchemaBuilder) Init() error {
	args, err := getArgs(schemaBuilder.Object)
	schemaBuilder.args = args
	schemaBuilder.Schema = graphql.Fields{}
	return err
}

//AddFunction adds a function
func (schemaBuilder *SchemaBuilder) AddFunction(
	name,
	description string,
	function func(graphql.ResolveParams) (interface{}, error)) {

	schemaBuilder.Schema[name] = &graphql.Field{
		Type:        schemaBuilder.GraphQLType,
		Description: description,
		Args:        schemaBuilder.args,
		Resolve:     function,
	}
}

func getArgs(object interface{}) (graphql.FieldConfigArgument, error) {
	objectType := reflect.TypeOf(object)
	output := graphql.FieldConfigArgument{}

	for i := 0; i < objectType.NumField(); i++ {
		currentField := objectType.Field(i)
		identifier, ok := currentField.Tag.Lookup("isArg")

		if identifier == "true" && ok {
			fieldType, err := getFieldType(currentField)
			if err != nil {
				return graphql.FieldConfigArgument{}, err
			}
			output[currentField.Name] = &graphql.ArgumentConfig{
				Type:        fieldType,
				Description: getTagValue(currentField, "description"),
			}

		}
	}

	return output, nil
}
