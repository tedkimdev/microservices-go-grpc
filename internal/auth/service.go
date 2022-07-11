package auth

//project
//├── cmd                     # contains microservices
//│   ├── '{microservice a}'  # "microservice a" project which contains "main" function to start app
//│   │   └── main.go
//│   └── '{microservice b}'  # "microservice b" project, also only contains "main" function to start app
//│       └── main.go
//├── docker                  # dockefile for each microservices
//│   ├── '{microservice a}'
//│   └── '{microservice b}'
//├── internal                # internal package of each microservice. codes inside internal package only visible to its own microservice and cannot be shared
//│   ├── '{microservice a}'
//│   │   ├── rest            # HTTP entry point for REST API
//│   │   ├── grpc            # gRPC entry
//│   │   ├── repo            # repository / connector to database
//│   │   └── service         # business logic layer which is handled by the microservice
//│   └── '{microservice b}'
//│       ├── rest
//│       ├── grpc
//│       ├── repo
//│       └── service
//├── pkg                     # public package directory which can be shared across project / microservices
//│   └── models              # data model which usually 1:1 relation with database
//│       ├── something.go
//│       └── another.go
//└── readme.md