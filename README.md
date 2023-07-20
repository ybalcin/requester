# Requester

A cli tool to provides sending requests in multi thread

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://go.dev/doc/install) to
install Go on your computer.

## Usage
```shell
git clone https://github.com/ybalcin/requester.git

cd requester

# Run with go run directly
go run main.go -parallel=3 https://google.com facebook.com

# Or build first
go build
./requester -parallel=3 https://google.com facebook.com
```

### Flags
* `-parallel`: int number of worker which runs at the background and listens queue of requests to send,
if not set it will be 10 as default