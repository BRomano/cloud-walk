# cloud-walk


### Implementation detail

I've created a cli and http as a solution.

I decide to develop a http server to explore more about design and development.

The endpoints receive the log file as multipart parse, it is not the best design for something to handle files who can be "big", 
but I understand it is to show some development abilities.
If it were a real program to parse log files, I would upload the file on a S3, and when it is done, I would process the whole log on an assyncronous thread and add statistics on a database.
And when the user request those statistics the endpoints would group and show the information requested on the endpoint.

1. DDD: The application is divided in infra/domain
   1. The service/handler objects implement their respective interface, helping dependency injection and low bound to objects. It can be seen on domain/service/log_parser.go and domain/repository/log_parser.go (Those names could be better)
   2. I design the solution to be capable to parse different kind of logs. To do so, it is needed to implement the repository interface and add the gameID on the factory, so the service can parse any implemented logs. It can be seen on infra/repository/factory.go
2. Scalability: 
   1. I have added the solution on a docker, so it can be scalable on ECS/Kubernets
   2. I have added a structured logging, so it is possible to send the output of container to an application to help visualise the logs
   3. It can be seen like so on output ![img.png](docs/img.png)
3. Tests: I have added some tests helping the development process and quality of code
   1. I was adding the logger on context, and it could be accessible in different point of code when needed
4. Environments: I added environment variables support, just using it on http port and some other information, but if the program need "feature flag", secrets or easy changeable feature, it could be done by changing some environment variable. It can be seen on cmd/settings.go

### To run
> docker compose up

> With the curl above you can call the death endpoint

```http request
curl -X POST --location "http://localhost:8080/log-parser/game/1/deaths" \
    -H "Content-Type: multipart/form-data" \
    -F "file=@.requests/qgames.log;filename=qgames.log;type=*/*" \
    -F "advertType=;type=*/*"
```

> With the curl above you can call the statistics endpoint
```http request
curl -X POST --location "http://localhost:8080/log-parser/game/1/matches" \
    -H "Content-Type: multipart/form-data" \
    -F "file=@.requests/qgames.log;filename=qgames.log;type=*/*" \
    -F "advertType=;type=*/*"
```


### Or you can build cli and run
> go build ./cmd/cli

For matches statistics 
> ./cli s .requests/qgames.log

For deaths statistics
> ./cli d .requests/qgames.log
