rest-api

Test user API for Passio.

USAGE:
    * To run the "production" version run:
        * `docker-compose up --build`
        * use either postman, curl, or any other library of your choosing to create users in the database using the examples at the bottom of this file
        * the APIs routes will be:
            * http://localhost:8080/api/user - GET - get all users 
            * http://localhost:8080/api/user/1 - GET - get user by id
            * http://localhost:8080/api/user?username=test - GET - get user by username "test"
            * http://localhost:8080/api/user - POST - create a user using the JSON body format described below
            * http://localhost:8080/api/user/1 - PUT - update user by id using the JSON body format (or partial) described below
            * http://localhost:8080/api/user/1 - DELETE - delete a user by ID
            * http://localhost:8080/api/auth/user - GET - **This is a strange one:** This GET request also takes a body. It is JSON with 2 fields: username and password. This endpoint will return the user object as JSON if the given password, once hashed, matches the stored password. Returns BadRequest(400) if the password doesnt match.

    * To run the tests (since this is not CI) run:
        * `docker-compose -f docker-compose.test.yml up --remove-orphans --force-recreate --build`
        * Then, from another cmd in this directory, run `go test --tags=e2e -v ./...`
        * To clean up, kill the process (Ctrl+C) and run `docker-compose -f docker-compose.test.yml down`
        * This will test all the endpoints (in reasonably basic ways for now) to make sure everything is working



Example JSON for post request to create a user:
------------------------------------------------------------------------------------------------------------------------------
**POST**
*http://localhost:8080/api/user*
```json
{
    "Username": "testuser",
    "Password": "password",
    "FirstName": "Test",
    "LastName": "User",
    "Email": "testuser@example.ca",
    "Telephone": "5555555555"
}
```

**or cURL**
```shell
curl --location --request POST 'http://localhost:8080/api/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Username": "testuser",
    "Password": "branton98",
    "FirstName": "Test",
    "LastName": "User",
    "Email": "testuser@example.ca",
    "Telephone": "9054959271"
}'
```
------------------------------------------------------------------------------------------------------------------------------

Example JSON to make sure the user password hash works
------------------------------------------------------------------------------------------------------------------------------
**GET**
*http://localhost:8080/api/auth/user*
```json
{
    "username": "testuser",
    "password": "password"
}
```

**or cURL**
```shell
curl --location --request GET 'http://localhost:8080/api/auth/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "testuser",
    "password": "branton98"
}'
```
------------------------------------------------------------------------------------------------------------------------------