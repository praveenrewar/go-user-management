# go-user-management
This is a simple Web Application written in Go and uses MongoDB as database and can be used as a boilerplate for writing other such applications. It also includes writing unit tests for the APIs built and dockerizing and running the whole app in docker containers. There are 5 APIs in the app `/get_users`, `/login`, `/signup`, `/update_password` and `/delete_user/{user_id}`.

#### Directory structure
``` 
go-user-management

├── config
|       └── config.json
├── server
|       ├── logger
|       |     └── logger.go
|       ├── middleware
|       |     ├── jwtAuthenticate
|       |     |    ├── jwtAuthenticate.go
|       |     |    ├── model.go
|       |     |    └── repository.go
|       |     └── userMiddleware
|       |          └── userMiddleware.go
|       ├── modules
|       |     └── users
|       |           ├── controller.go
|       |           ├── model.go
|       |           ├── repository.go
|       |           └── router.go
|       └── sharedVariables
|               ├── sharedVariables.go
|               └── structures.go
├── .gitignore
├── docker-compose.yaml
├── Dockerfile
├── main.go
├── README.md
└── users_test.go
```

#### Installation
You can clone from this repository and use master
```
git clone https://github.com/praveenrewar/go-user-management.git
cd go-user-management
```
For package installation
```
go mod init
go mod vendor
```
To run unit tests and get code coverage percentage
```
WorkEnv=test go test -coverpkg ./...
```
To run the app
```
go run main.go
```
To run as a docker container (make sure you have [docker](https://docs.docker.com/install) and [docker-compose](https://docs.docker.com/compose/install/) installed)
```
sudo docker-compose up --build api
```
Currently, I have not mocked the db operations, hence the docker-compose file includes pulling the latest MongoDB image so that the tests run successfully in docker.

