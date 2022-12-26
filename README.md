# Rocket.Chat CLI

## Development

1) Clone the repository
2) Install packages: `go get`
3) Run: `go run .`

## Building

1) Clone the repository
2) Install packages: `go get`
3) Build: `go build .`
4) Run: `./rocketchat-cli`

## Building to different platforms

Instead of `go build .`, you can run:

`GOOS=linux GOARCH=arm go build .`

This will build the arm version for linux.

The list of possible values can be obtained with: `go tool dist list`

