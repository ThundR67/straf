package straf

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

type extra struct {
	Age     int
	FavNums []float32
}

type colors struct {
	Name     string
	HexValue string
}

type user struct {
	UserID    int
	ExtraData extra
	FavColors []colors
}

//Tests generation of graphQL Object from struct
func TestGraphQLObjectGen(t *testing.T) {
	assert := assert.New(t)

	graphQLObject, err := GetGraphQLObject(user{})
	assert.NoError(err, "GetGraphQLObject Returned Error")
	extraType, err := GetGraphQLObject(extra{})
	assert.NoError(err, "GetGraphQLObject Returned Error")
	colorType, err := GetGraphQLObject(colors{})
	assert.NoError(err, "GetGraphQLObject Returned Error")

	testField(
		"UserID",
		*graphQLObject.Fields()["UserID"],
		graphql.Int,
		*assert,
	)

	testField(
		"ExtraData",
		*graphQLObject.Fields()["ExtraData"],
		extraType,
		*assert,
	)

	testField(
		"FavColors",
		*graphQLObject.Fields()["FavColors"],
		graphql.NewList(colorType),
		*assert,
	)

}

func testField(name string,
	definition graphql.FieldDefinition,
	graphqlType graphql.Output,
	assert assert.Assertions) {

	assert.Equal(definition.Name, name)
	assert.Equal(definition.Type, graphqlType)

}
