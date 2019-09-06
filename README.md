[![Go Report Card](https://goreportcard.com/badge/github.com/SonicRoshan/straf)](https://goreportcard.com/report/github.com/SonicRoshan/straf)

# Straf
Convert Golang Struct To GraphQL Object On The Fly


## Example

### Converting struct to GraphQL Object
```go
type UserExtra struct {
    Age int `description:"Age of the user"` //You can use description struct tag to add description
    Gender string `deprecationReason:"Some Reason"` // You can use deprecationReason tag to add a deprecation reason
}

type User struct {
    UserID int
    Username string `unique:"true"` // You can use unique tag to define if a field would be unique
    Extra UserExtra
}


func main() {
    //GetGraphQLObject will convert golang struct to a graphQL object
    userType, err := straf.GetGraphQLObject(User{})

    //You can then use userType in your graphQL schema
}
```


### Converting struct to GraphQL Object
```go
type User struct {
    UserID int `isArg:"true"` //You can use isArg tag to define a field as a graphql argument
    Username string `isArg:"true"`
}

var database []User = []User{}

func main() {
    //GetGraphQLObject will convert golang struct to a graphQL object
    userType, err := straf.GetGraphQLObject(User{})

    builder := straf.NewSchemaBuilder(userType, User{})
    builder.AddFunction("CreateUser", 
                        "Adds a user to database",
                        func(params graphql.ResolveParams) (interface{}, error)) {
                            id := params.Args["UserID"]
                            username := params.Args["Username"]
                            database = append(database, User{UserID: id, Username: Username})
                        })
    schema := builder.Schema
    //You can then use this schema
}
```
