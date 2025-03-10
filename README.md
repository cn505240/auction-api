# auction-api
Sample Golang REST API to create and bid on auctions


Instructions:

1. Install the needed tooling:
```
  asdf plugin add golang
  asdf plugin add postgres
  asdf plugin add mockery
  asdf install
```
2. Run `docker compose up` to start the database
3. Run `go run server.go` to run the migrations and start the server
4. Run `go test -v ./...` to run the tests

Comments:
- There is lots missing here. There is no read functionality, and no user management, so it's still a kind of toy example.
- There is no `max bid` functionality, where a user can set a max bid and the system will automatically bid for them up to that amount.
- There are many quality and production-readiness issues that we can discuss in review.