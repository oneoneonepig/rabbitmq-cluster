`send.go`
```
Send a single message

Usage of send.go:
  -h, --host string       host address and port (default "localhost:5672")
  -n, --name string       queue name (default "two.hello")
  -p, --password string   password (default "admin")
  -u, --username string   username (default "admin")
  
Example:
  send -h 10.20.131.53 -n tw.hello -u user -p pass
```

`rapidsend.go`
```
Send a single message rapidly

Usage of rapidsend.go:
  -h, --host string       host address and port (default "localhost:5672")
  -i, --interval string   interval between messages (default "1s")
  -n, --name string       queue name (default "two.hello")
  -p, --password string   password (default "admin")
  -u, --username string   username (default "admin")
  
 Example:
   rapidsend -h 10.20.131.53 -n tw.hello -u user -p pass -i 100ms
```

`receive.go`
```
Receive message

Usage of receive.go:
  -h, --host string       host address and port (default "localhost:5672")
  -n, --name string       queue name (default "two.hello")
  -p, --password string   password (default "admin")
  -u, --username string   username (default "admin")
  
Example:
  receive -h 10.20.131.53 -n tw.hello -u user -p pass
```
