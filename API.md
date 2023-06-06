# API Documentation

The distillery comes with an API served under `/api/`.
It is still a work in progress, and will be polished and properly implemented at a later point.
The API is currently disabled by default, and needs to be enabled in `distillery.yaml`. 

- `/api/v1/auth`: Returns user information
- `/api/v1/systems`: Returns a (publically visible) list of systems 
- `/api/v1/news`: Returns JSON containing all news items