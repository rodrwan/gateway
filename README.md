# Gateway

Service that connect GraphQL queries with a microservice architecture.

# Database

```
$ createdb graph_test
```

Install goose

```
$ go get -tags nosqlite3 github.com/steinbacher/goose/cmd/goose
```

Run migrations
```
goose -path postgres -env development up
```

# Run server

```
$ go run cmd/server/main.go
```

# Request server with cURL

```sh
$ curl -X POST \
  http://localhost:3000/users \
  -H 'Content-Type: application/graphql' \
  -d 'query {
    users {
				id,
        first_name,
        last_name,
        email,
        phone,
        birthdate,
   			address {
            address_line,
            city,
            locality,
            administrative_area_level_1,
            country,
            postal_code
        }
    }
}'
```
# queries

### To get all users

```
query {
    users {
        id,
        first_name,
        last_name,
        email,
        phone,
        birthdate
    }
}
```

### To get an specific user

```
query {
    user(email: "<email>") {
        id,
        first_name,
        last_name,
        email,
        phone,
        birthdate,
        address {
            address_line,
            city,
            locality,
            administrative_area_level_1,
            country,
            postal_code
        }
    }
}
```