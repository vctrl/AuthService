# AuthService

Simple autorization server using OAUTH 2.0 including implementation for Google service.

## Adding new authorization servers
To add new autorization servers, we have to extend config.yml file and add new endpoint to constants in auth.go. All settings could be moved to config file in the future.

## Endpoints
We have two endpoints:
GET /login?site={siteName} which returns link to the authorization page.

GET /callback?site={siteName} to which the authorization service redirects.

siteName is the required parameter. Currently service including implementation for Google server.

## Test coverage
ok      github.com/vctrl/authService/delivery   0.017s  coverage: 94.3% of statements
ok      github.com/vctrl/authService/usecase    0.016s  coverage: 76.5% of statements

## Requirements
- go 1.15
- docker & docker-compose

## Run Project
Use ```make run``` to build and run docker containers with application
Use http://localhost:8080/login?site=google to login with Google.