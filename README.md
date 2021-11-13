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

### How to

- Find information related to the Makefile: make or make help
- Format code: `make check.fmt`
- Format import: `make check.imports`
- Code linter : `make check.lint`
- Execute TU: `make check.test`
- Create container for lint/format: `make build.tools`

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

## Design and assumptions

The app has a clean architecture design. It has two layers: transport layer (websocket, output `repo` module) and usecase. 
The idea is to decouple the _transport layer_ from _usecase_ to be able to scale up easily if a large number of pair trading are to be computed.

**Transport layer**

Only one connection is made to ws server although a connection per trading pair can be setup. On each message arrival, the client find the type of the message. Then the message is parsed into a corresponding `entity`.
This process of parsing the type of message and then the whole message was done in order to have a loose coupling of the low level read method `readWs` and the `receive` method of the client.
When the parsing is done, the _entity_ is written into a channel which is consumed by the _usecase_. The use of channel between the layers allows the _usecase_ to consume the message at its pace.

**Usecase**

The usecase has a central component (`AvgManager` the name could be better I admit) which consume messages from input channel and, for each trading pair, calls the `TradingPairAvgCalculator` for each _ticker_ of _heartbeat_ message.
Each trading pair has his own `TradingPairAvgCalculator` stored in a map.

The job of `TradingPairAvgCalculator` is to make sure that the sequence of the _ticker_ is equal or superior of the sequence of the last _hearbeat_. 
Internally, `TradingPairAvgCalculator` has an average calculator. 

The math is done in the `Calculator`.The calculator is very flexible. The maximum number of data points is set as a parameter. 

The data points are saved in a _LILO_ stack. The complexity of the is O(1). To achieve O(1), at each insert the sum of products `price * volume` is computed. If the maximum number of points is reached, the fall off data point is
substracted from the sum of products. Therefore, at each insert after the maximum number of data points has been reached, the sum is:
```
    sum = sum - price_falloff * vol_falloff + price_newpoint * vol_newpoint
```
Same idea is applied for total volume. Therefore, when average is compute it's just a simple matter of a division. 
 


