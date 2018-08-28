# go-graphql-example

This is a sample project that implements Graphql server using Golang and firebase.


#### Prerequisite

1. [Install Golang](https://golang.org/doc/install)
2. Setup GOPATH [Link1](https://golang.org/doc/code.html#GOPATH) and [Link2](https://github.com/golang/go/wiki/GOPATH)
3. [Install Glide](https://github.com/Masterminds/glide)
4. Setup Firebase account and get your service account key [Firebase Docs](https://firebase.google.com/docs/admin/setup#add_firebase_to_your_app)

#### Getting Started
1. Clone the repo
2. Run `glide install`
3. Copy/replace your service key file under the root with the name `serviceAccountKey.json`
4. Update firebase database URL in `dal/firebase/firebase.go` file.
3. Run `go run server/server.go`
4. Open `http://localhost:8090/` for GraphQL Playground
6. Hit the queries

#### TODO
1. Config files
2. Authentication
3. Validations
4. CORS
5. Dataloaders
5. Small tweaks