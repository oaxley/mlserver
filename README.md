# Machine Learning Server

![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg?style=flat-square)

## Goals

Run ML/AI models on GPU servers and expose them through an API running on a different server.  
The communication between the API Service and the ML models is done via RPC in asynchronous mode.  

Models are registered on a registry service, the API will re-route the requests from the user to their corresponding model.  
A rate limiter (via Redis) allows the administrator to ensure all users are served equaly without overloading the GPUs.

## Registry Service

This is a store where models register themselves so that other components can find them and initiate communication.

The documentation is available [here](./registry/README.md).


## License

This program is under the **Apache License 2.0**
A copy of the license is availabe [here](https://choosealicense.com/licenses/apache-2.0/).
