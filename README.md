## To Do

A very simple CRUD API implementation in Golang

### Development Iterations

#### Iteration 1

1. ~~Database connection~~
2. ~~Task CRUD routes setup~~
3. ~~Database credentials in env file~~

#### Iteration 2

1. ~~Error handling~~
2. ~~Error feedback~~
3. ~~Logging~~

#### Iteration 3

1. Authentication
2. Authorization

#### Iteration 4

1. Enable HTTPS

### Learnings

1. Always use the request's context inside the route handlers. These contexts are tied to the lifecycle of the request so will correctly allocate and free resources when the request fails or completes.
2. Contexts are only meant to be used for timing out operations that are time consuming and/or error prone.
3. Contexts should be created inside the part of the code that decides the timeouts, cancellations and lifecycle. Contexts should be passed on to functions that need to communicate with the context for such properties