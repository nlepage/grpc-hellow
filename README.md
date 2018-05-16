gRPC Hello World in different languages with server side streaming

## Usage

Start servers with `docker-compose up -d`

Call the go client with `docker-compose exec go /hellow client [options] <name> <count>`

Call the node client with with `docker-compose exec node yarn start client [options] <name> <count>`

Client supports the following options:
```
    --host string   Host to call (default "localhost")
-s, --stream        Ask for a streamed response
```

Examples:
```sh
$ docker-compose exec node yarn start client world 42              # Make a simple node to node call
$ docker-compose exec go /hellow client -s world 42                # Make a streamed go to go call
$ docker-compose exec node yarn start client --host go -s world 42 # Make a streamed node to go call
```
