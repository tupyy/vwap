# vwap

## How to run

### Shell
If you have `go` installed:

```shell
make build.vendor
make build run
```

### Docker
```shell
make build.docker
make run.docker
```

## Configuration

The app can be configured by flags or with a configuration file:

```shell
./target/run --help
Usage of ./target/run:
  -config string
    	path of the configuration file
  -endpoint string
    	endpoint ws address
  -log_level string
    	log level (default "info")
  -max_data_points int
    	maximum number of data points used to compute the average (default 200)
  -output string
    	path of the output file
  -pairs string
    	comma separated trading pairs
```

If the configuration file is set, the other flags are ignored.
Configuration file:
```json
{
    "log_level": "info",
    "endpoint": "wss://ws-feed.exchange.coinbase.com",
    "trading_pairs": [
        "BTC-USD", "ETH-USD", "ETH-BTC"
    ],
    "max_data_points": 200
}
```
