# Hanger.io

- run: `go run main.go`
- client1: `curl "localhost:8080/pause/1"`
- client2: `curl "localhost:8080/continue/1?rampup=10"`

All flags are optional:
- `-p` port
  - Default value: 8080
- `-h` host
  - Default value: 127.0.0.1
- `-rp` redis port
  - Default value: 6379
- `-pass` Redis password
  - Default value: <empty_string>
- `-db` Redis DB index
  - Default value: 0