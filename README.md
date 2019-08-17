# Straf
Convert Golang Struct To GraphQL Object On The Fly


## Example
```go
type UserExtra struct {
    Age int
    Gender string
}

type User struct {
    UserID int
    Username string
    Extra UserExtra
}


func main() {
    //GetGraphQLObject will convert golang struct to a graphQL object
    userType, err := straf.GetGraphQLObject(User{})

    //You can then use userType in your graphQL schema
}
```
