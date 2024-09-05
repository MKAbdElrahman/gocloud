## Architecture

- For this project, I will take the problem domain the as the pricipal axes of modularity. This is architecture is similar to approaches like modular monolith or vertical slice architecture which discourages technical partitioning based on layers (presentaion, business, data access). 

- Each module will have own  its handlers, business and infrastrucuture objects. 

- For this architecture to work well, communcation and sharing between modules should be close two zero.

## Dependancy Injection 

- Wiring dependancies and all object creation is centralized to main.go
- Each object should be passed only the essentail dependancies it needs to do its job, it should never have access to a functionality it doesn't need. 


## Logging

- Logging is dones only at top level handlers, error comming from below could be wrapped and given context, but the handlers are responsible for logging ther errors


- I usually couple server errors  with logging

## Configuration

- All configuration variables are read at main.go
