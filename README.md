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

# Usage

## First run

1) Run the binary: `./rocketchat-cli`
<img width="600" alt="image" src="https://user-images.githubusercontent.com/575138/209550209-4f4510bd-b8f6-45a0-8862-afb33b00760a.png">
2) Configure your install:
<img width="600" alt="image" src="https://user-images.githubusercontent.com/575138/209550277-780035f6-207c-4bc3-b94d-938a7f0d61d9.png">
3) Fill your hostname (without the protocol), your email (to generate the certificates) and click save
4) Quit the cli
5) Run with `docker compose up`

## Subsequent runs

1) Run the binary: `./rocketchat-cli`;
<img width="600" alt="image" src="https://user-images.githubusercontent.com/575138/209555510-d93429a1-3802-4592-b4ae-6ceef6992527.png">
2) You can reconfigure your install or quit;
3) Quit the cli
4) Run with `docker compose up`;
