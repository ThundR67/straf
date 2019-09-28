package straf

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

type user2 struct {
	userID   string `isArg:"true"`
	username string `isArg:"true"`
	age      int
}

var database = []user2{}
var middlewareCame = false

func middleware(
	function func(graphql.ResolveParams) (interface{}, error),
	params graphql.ResolveParams) (interface{}, error) {

	middlewareCame = true
	fmt.Println("True")
	return function(params)
}

func handler(params graphql.ResolveParams) (interface{}, error) {
	fmt.Println(params.Args)
	return nil, nil
}

func TestSchemaBuilder(t *testing.T) {
	assert := assert.New(t)
	graphQLObject, _ := GetGraphQLObject(user2{})
	builder := NewSchemaBuilder(graphQLObject, user2{}, middleware)
	assert.NotNil(builder.middleware)
	assert.IsType(&SchemaBuilder{}, builder)
	builder.AddFunction("create", "des", handler)

	assert.Contains(builder.Schema, "create")
	assert.Contains(builder.Schema["create"].Args, "userID")
	assert.Contains(builder.Schema["create"].Args, "username")
	assert.NotContains(builder.Schema["create"].Args, "age")
	assert.Equal(builder.Schema["create"].Description, "des")
	assert.Equal(builder.Schema["create"].Type, graphQLObject)
	builder.Schema["create"].Resolve(graphql.ResolveParams{})
	assert.True(middlewareCame)

	builder = NewSchemaBuilder(graphQLObject, user2{})
	assert.Nil(builder.middleware)
	builder.AddFunction("create", "des", handler)
}
