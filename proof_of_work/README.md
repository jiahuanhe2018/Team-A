### Deployment steps:
- `go run main.go`
- open a web browser and visit `http://localhost:8080/`
- to write new blocks, send a `POST` request (I like to use [Postman](https://www.getpostman.com/apps)) to `http://localhost:8080/` with a JSON payload with `BPM` as the key and an integer as the value. For example:
```
{"BPM":50}
```
- **watch your Terminal to see Proof of Work in action!**