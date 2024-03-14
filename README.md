# Machine Learning Server


## Goals

Run ML/AI models on GPU servers and expose them through an API running on a different server.  
The communication between the API Service and the ML models is done via RPC in asynchronous mode.  

Models are registered on a registry service, the API will re-route the requests from the user to their corresponding model.  
A rate limiter (via Redis) allows the administrator to ensure all users are served equaly without overloading the GPUs.

