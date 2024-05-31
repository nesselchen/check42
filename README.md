# CHECK42
Simple REST-API for creating and managing todos

## Capabilities
On startup the application logs out all available paths and their respective methods. \
Anything under the API route is best accessed via a tool like Postman or Thunderclient.

### Todo endpoints 
Path: /api/todo
- GET: returns all todos for the logged in user.
- POST: create new todo specified in the JSON body. Any left out fields are zeroed as per the Go's JSON Unmarshalling rules.
```json
{
    "text": "My urgent task",
    "category": 1
}
```
- GET, DELETE with /id: perform the action on the specified todo. 
- PATCH with /id: change directly via `text` and `done` URL params.
- PUT with /id: change all fields via the JSON provided in the body. This works the same as creating a todo.

### Category endpoints
Path: /api/todo/category
- GET: returns all categories for the logged in user.
- POST: create a new category via the `name` URL parameter.
- DELETE with /id: deletes the category and all todos that are associated with it.
- PATCH with /id: change the name via the `name` URL parameter.  

### Login and authentication
- POST /auth/signin: Create a new user via the JSON body.
```json
{
    "name":     "Nessel",
    "email":    "test@test.com",
    "password": "password" 
}
```
- POST /auth/login: Create a JWT token cookie by sending a Basic authenticated request to this endpoint. The token will expire after one week or when you unset it. All other endpoints rely on this token as method of authorization.
- POST /auth/logout: Overrides the JWT with an expired cookie.

### Frontend
If you open the project at :2442 in a browser, a rudimentary frontend should be served up. This prompts you to log in with a previously created user (admin:password is the dummy user :D). The frontend is vanilla HTML and Javascript for ease of bundling and lets you create, check and delete todos in different categories.

## Setup

### Outside Docker
Run `go run .` to start the application and log out the available endpoints.
For running outside Docker the .env file should be adjusted like this:

> DB_HOST=localhost \
> SERVER_HOST=

Run the initialization script found in sql/initdb.sql to initialize the database scheme and insert some dummy values.

### Inside Docker
Run `docker compose up`. \
Here the default .env configuration should suffice. This will also run the DB initialization script.

### Demo
For demonstration purposes you can use the dummy user `admin` with password `password` which already has some todos registered.
