### Project structure

```
├── cmd                 # entry point to run service
│   ├── main.go
├── configs             # config for different env (for simple, I only add one config)
│   └── ...
├── internal
│   ├── domains         # define model for request-response usage
│   ├── entities        # define model for struct in database
│   ├── repositorires   # storage interface
│   │   └── ...         # storage implementation
│   └── service         # define and implement service
├── Makefile
├── mock            # mock/stub for testing
├── test            # integration tests
```

### Run
- Run `bash ./start.sh` to start the application
- In case you're using `docker compose v2`, please using run `bash ./startv2.sh` (not yet test this command)