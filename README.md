# grpc-hellow

Start servers with `docker-compose up -d`.

Make a go -> go call with `docker-compose exec go /hellow client world 42`

Make a go -> node call with `docker-compose exec go /hellow --host node client world 42`

Make a node -> node call with `docker-compose exec node yarn start client world 42`.

Make a node -> go call with `docker-compose exec node yarn start client --host go world 42`.
