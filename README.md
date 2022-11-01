### Project structure

```
├── cmd                 # entry point to run service
│   ├── main.go
├── configs             # config for different env (for simple, I only add one config)
│   └── ...
├── internal
│   ├── models          # define model for request-response usage
│   ├── entities        # define model for struct in database
│   ├── repositorires   # storage interface
│   │   └── ...         # storage implementation
│   └── service         # define and implement service
├── libs                # internal libs
├── start.sh            # using it to run the app.
├── mock                # mock/stub for testing
├── integration_test    # integration tests
```

### Run
- pre-condition: `go mod tidy`
- Run `bash ./start.sh` to start the application
- In case you're using `docker compose v2`, please using run `bash ./startv2.sh` (not yet test this command)
- I disable expose port `5432`, you can uncomment in `docker-compose.yaml`, to avoid conflict in your machine
- You can add ENV variable also to avoid leak secret.
- Unit test: `make unit-test`
- Integration test: `make integration-test` __ I'm using dockertest to write the integration test, so make sure your machine installed docker.
### Manual test
- Pre-condition: must run `bash ./start.sh`
- Test PlaceWager, example:
 ```
    curl --location --request POST 'localhost:8080/wagers' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "total_wager_value": 20,
        "odds": 30,
        "selling_percentage": 30,
        "selling_price": 50
    }'
```
- Test BuyWager (must call several times the PlaceWager, example:
```
    curl --location --request POST 'localhost:8080/buy/1' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "buying_price": 6
    }'
```
- Test ListWager (must call several times the PlaceWager, example:
```
    curl --location --request GET 'localhost:8080/wagers?page=:1&limit=:4'
```

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