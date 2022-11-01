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
├── libs                # internal libs
├── start.sh            # using it to run the app.
├── mock            # mock/stub for testing
├── test            # integration tests
```

### Run
- Run `bash ./start.sh` to start the application
- In case you're using `docker compose v2`, please using run `bash ./startv2.sh` (not yet test this command)
- Unit test: `make unit-test`

### Cool items:
- In postgres the `transaction_level default = read commited`, using lock row to lock the `wager record` when calling `buy wager` to avoid race condition. Using this way, we can easy scale when need improve throughput.
- Implement middleware to make the API more simple
- Change to use chi-go router. Why chi-go? Because it lightweight, idiomatic, and composable router for building Go HTTP services. Especially, chi's router is based on Radix trie, so it'll handle the request as fast as possible if we have a lot of handlers in the future.
- Add config (with file - just easy for testing) - after load config file, it will overwrite the variable environment, so this project still abide by 12factor https://12factor.net/ .P/s: Again the config file just save the infomation for quick run, when in production, use variable environment.
- Add adaptive checking solution  to validate `selling_price` 
- Add adative solution paging to save cost, not full scan DB. (if can, use the lastID from previous request and sort by desc (lastest first.. but out of scope so I omit it))
### If I have more time:
- Write more test to cover, especially egde cases
- Customize the error message to DRY
- Research deep in biz logic of betting