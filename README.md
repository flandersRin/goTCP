# goTCP
A simple go language implements TCP connections

## QuickStart

1. Run

```go
// default port is 2222
// your can set the port inn config.conf
go run main.go
```

2. Sent message by terminal

- mac 

  ```shell
  # sent the message end with '\n'
  echo "some messages\n" |nc localhost 2222
  ```
