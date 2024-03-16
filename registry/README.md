# gRPC Registry

![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg?style=flat-square)
![Golang](https://img.shields.io/badge/Golang-1.22.0-blue?style=flat-square)

## Presentation

The registry is a store that records the location of running ML models.  
The information recorded are:

- the model name
- the model version
- the hostname where the model is located
- the TCP port to communicate with the model

## Service endpoints

The gRPC server has the following endpoints:

- **SetService** : register a new model inside the registry
- **GetService** : retrieve the parameters for a previously recorded model

### SetService message

The message format to set a new service is:

``` go
ModelName       string
ModelVersion    string
Hostname        string
Port            uint32
```

### GetService message

The message format to query a service is:

``` go
ModelName       string
ModelVersion    string
```

## Running

The registry accepts the following flags on the command line:

``` console
    -tls        bool
        True uses TLS, otherwise use plain TCP (default: false)
    -cert_file  string
        The certificate file for TLS communication
    -key_file   string
        The key file for TLS communication
    - hostname  string
        The hostname / IP where to bind the server (default: localhost)
    - port      uint32
        The TCP port to listen (default: 50051)
```

To start the registry server

``` console
$ go run registry/main.go
2024/03/16 21:39:55 Starting registry server on localhost:50051
```
