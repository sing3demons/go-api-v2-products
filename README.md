|       Method      |     RelativePath  |
| ------------- | ------------- |
|   GET      | /healthz                  |
|   GET      | /x                        |
|   GET      | /api/v1/auth/profile      |
|   POST     | /api/v1/auth/register     |
|   POST     | /api/v1/auth/login        |
|   PATCH    | /api/v1/auth/profile/:id  |
|   PUT      | /api/v1/auth/profile/:id  |
|   GET      | /api/v1/categories        |
|   GET      | /api/v1/categories/:id    |
|   POST     | /api/v1/categories        |
|   PUT      | /api/v1/categories/:id    |
|   DELETE   | /api/v1/categories/:id    |
|   GET      | /api/v1/products          |
|   GET      | /api/v1/products/:id      |
|   POST     | /api/v1/products          |
|   PUT      | /api/v1/products/:id      |
|   DELETE   | /api/v1/products/:id      |
|   GET      | /api/v1/users             |
|   POST     | /api/v1/users             |
|   GET      | /api/v1/users/:id         |
|   PUT      | /api/v1/users/:id         |
|   DELETE   | /api/v1/users/:id         |
|   PATCH    | /api/v1/users/:id/promote |
|   PATCH    | /api/v1/users/:id/demote  |


<hr>

## Usage
### Start using it
#### Add comments to your API source code, [See Declarative Comments Format](https://swaggo.github.io/swaggo.io/declarative_comments_format/).
#### Download Swag for Go by using:
```install swagger
$ go get -u github.com/swaggo/swag/cmd/swag

# 1.16 or newer
$ go install github.com/swaggo/swag/cmd/swag@latest
```

#### Download gin-swagger by using:
```gin-swagger
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

#### Download fiber-swagger by using:
```fiber-swagger
go get -u github.com/arsmn/fiber-swagger/v2
```

#### Import following in your code:
```import
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files
```

## Run the Swag at your Go project root path(for instance ~/root/go-peoject-name), Swag will parse comments and generate required files(docs folder and docs/doc.go) at ~/root/go-peoject-name/docs.
```swagger
swag init
```
```go
go run main.go
```

```go to
http://localhost:8080/swagger/swagger/index.html
```

```tree
.
├── cache
│   ├── cacher.go
│   └── config.go
├── config
│   ├── acl_model.conf
│   ├── policy.csv
│   └── redis_cfg
│       └── redis.conf
├── controllers
│   ├── auth.go
│   ├── category.go
│   ├── controller.go
│   ├── product.go
│   └── user.go
├── database
│   └── db.go
├── dev.env
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── k8s
│   ├── go-product
│   │   ├── 00-namespace.yml
│   │   ├── 01-secret.yml
│   │   ├── 02-pvc.yml
│   │   ├── 03-deployment.yml
│   │   ├── 04-service.yml
│   │   └── 05-ingress.yml
│   ├── postgres
│   │   ├── 00-namespace.yml
│   │   ├── 01-secret.yml
│   │   ├── 02-pv-pvc.yml
│   │   ├── 03-deployment.yml
│   │   └── 04-service.yml
│   └── redis
│       ├── 00-namespace.yml
│       ├── 01-deployment.yml
│       └── 02-service.yml
├── main.go
├── middleware
│   ├── authentication.go
│   └── authorization.go
├── models
│   ├── category.go
│   ├── product.go
│   └── user.go
├── routes
│   ├── gin.go
│   └── routes.go
├── seeds
│   └── seed.go
├── store
│   └── store.go
├── test
│   ├── auth.http
│   ├── next.png
│   └── products
│       ├── create_product.http
│       └── find_product.http
└── uploads
    ├── products
    └── users
```

```chmod +x k8s.sh
./k8s.sh 
kubectl get po -n go-product  
```
#### output
```output
NAME                       READY   STATUS    RESTARTS   AGE
go-server-5d4689dc-jdkjl   1/1     Running   0          114s
redis-5c4c454dd-94xmj      1/1     Running   0          114s
```

#### open browser
```
http://kubernetes.docker.internal/swagger/index.html
```