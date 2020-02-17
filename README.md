# DAPR with GRPC

Small example to try using Dapr to implement services in kubernetes, using gRPC for comms.

## Architecture

Simple two applications in _go_, a _server_ app that exposes a service over gRPC and a _client_ app that calls the server.

_client_ also exposes a http endpoint that you can _curl_ to trigger the backend _server_ and see the output.

_server_ exposes the services through _Dapr_ and _client_ uses _Dapr_ to call _server_.

```
             /-----------------------------------\   /-----------------------------------\
 O           |    +--------+           +-----+   |   |    +------+           +-------+   |   
/|\ - http - | -> | client | - grpc -> | dapr| - | - | -> | dapr | - grpc -> | server|   |
/ \          |    +--------+           +-----+   |   |    +------+           +-------+   |
             |             CLIENT POD            |   |             SERVER POD            |
             \-----------------------------------/   \-----------------------------------/
```
## Install locally

Uses Kind for local cluster and Ko for deployments.

1. Run the scripts under `./deployments/kind` to create a local kind cluser and install dapr.

2. Run `.apply.sh` to deploy the application locally.