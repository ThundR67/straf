package straf

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

type middlewareType func(func(graphql.ResolveParams) (interface{}, error), graphql.ResolveParams) (interface{}, error)

// NewSchemaBuilder is used get a new schema builder
func NewSchemaBuilder(
	graphQLType graphql.Output,
	object interface{},
	middlewareArg ...middlewareType) *SchemaBuilder {

	var middleware middlewareType

	if middlewareArg != nil {
		middleware = middlewareArg[0]
	}

	builder := SchemaBuilder{
		GraphQLType: graphQLType,
		Object:      object,
		middleware:  middleware,
	}

	builder.Init()
	return &builder
}

// SchemaBuilder is used to build a schema based on a struct
type SchemaBuilder struct {
	GraphQLType graphql.Output
	Object      interface{}
	Schema      graphql.Fields
	args        graphql.FieldConfigArgument
	middleware  middlewareType
}

// Init initializes
func (schemaBuilder *SchemaBuilder) Init() {
	schemaBuilder.args = graphql.FieldConfigArgument{}
	schemaBuilder.AddArgumentsFromStruct(schemaBuilder.Object)
	schemaBuilder.Schema = graphql.Fields{}
}

//AddArgumentsFromStruct is used to add arguments from a struct
func (schemaBuilder *SchemaBuilder) AddArgumentsFromStruct(object interface{}) {
	for key, value := range getArgs(object) {
		schemaBuilder.args[key] = value
	}
}

// AddFunction adds a function
func (schemaBuilder *SchemaBuilder) AddFunction(
	name,
	description string,
	function func(graphql.ResolveParams) (interface{}, error)) {

	var functionToAdd func(graphql.ResolveParams) (interface{}, error)

	if schemaBuilder.middleware != nil {
		functionToAdd = func(params graphql.ResolveParams) (interface{}, error) {
			return schemaBuilder.middleware(function, params)
		}
	} else {
		functionToAdd = function
	}

	schemaBuilder.Schema[name] = &graphql.Field{
		Type:        schemaBuilder.GraphQLType,
		Description: description,
		Args:        schemaBuilder.args,
		Resolve:     functionToAdd,
	}
}

func getArgs(object interface{}) graphql.FieldConfigArgument {
	objectType := reflect.TypeOf(object)
	output := graphql.FieldConfigArgument{}

	for i := 0; i < objectType.NumField(); i++ {
		currentField := objectType.Field(i)
		identifier, ok := currentField.Tag.Lookup("isArg")

		if identifier == "true" && ok {
			fieldType := getFieldType(currentField)
			output[currentField.Name] = &graphql.ArgumentConfig{
				Type:        fieldType,
				Description: getTagValue(currentField, "description"),
			}

		}
	}

	return output
}
