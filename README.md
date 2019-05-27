# redact [![GoDoc Badge][badge]][godoc]

`redact` is a Golang library for redacting sensitive values from a struct. This
is useful when you're logging requests and responses on a server, but some
payloads contain secrets or other values that shouldn't be logged.

[badge]: https://godoc.org/github.com/ucarion/redact?status.svg
[godoc]: https://godoc.org/github.com/ucarion/redact

## Example

`redact` provides a `Redact` function that takes a path and the value you want
to mutate. It works on struct elements:

```go
import "github.com/ucarion/redact"

type User struct {
  Name string
  Password string
}

user := User{Name: "John", Password: "letmein"}
fmt.Println(user) // {John letmein}

redact.Redact([]string{"Password"}, &user)
fmt.Println(user) // {John }
```

This works even if the data you want to redact is nested:

```go
type CreateUserRequest struct {
  RequestID string
  User      User
}

req := CreateUserRequest{RequestID: "abc", User: User{Name: "John", Password: "letmein"}}
fmt.Println(req) // {abc {John letmein}}

redact.Redact([]string{"User", "Password"}, &req)
fmt.Println(req) // {abc {John }}
```

It also works on arrays and slices. In that case, *every* element of the array
or slice gets recursively redacted:

```go
users := []User{
  User{Name: "John", Password: "letmein"},
  User{Name: "Mary", Password: "123456"},
}

fmt.Println(users) // [{John letmein}, {Mary 123456}]

redact.Redact([]string{"Password"}, &users)
fmt.Println(users) // [{John }, {Mary }]
```

It also works on maps in a way similar to structs.

```go
user := map[string]string{
  "Name": "John",
  "Password": "letmein",
}
fmt.Println(user) // {Name: John, Password: letmein}

redact.Redact([]string{"Password"}, &user)
fmt.Println(user) // {Name: John, Password: }
```

Note, however, that Go does not support mutating map elements. That's not
something any package can work around. So if you want to mutate elements of a
map (in the example above, we *replaced* `Password`, not modified it). You'll
have to make the map elements be pointers:

```go
users := map[string]*User{
  "a": &User{Name: "John", Password: "letmein"},
}
fmt.Println(users["a"]) // &{John letmein}

redact.Redact([]string{"a", "Password"}, &users)
fmt.Println(users["a"]) // &{John }
```
