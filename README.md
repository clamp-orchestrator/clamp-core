# clamp-core

The clamp core that can be deployed as a binary and does service orchestration

## Installation & Configuration

1. Install golang
2. Set GOPATH
3. Go to Application root path
4. Run `go build main.go`
5. `Postgres DB` and `RabbitMq` can be configured in `config/env.go` file
6. Presently we are using common postgres 
   ```
   hostname : 34.222.238.234:5432 
   user : clamp 
   dbname : clampdev
   password : clamppass
   ```
7. Instead recommend setup `Postgres locally` for development
8. Once everything is configured migration can be run using below command
   `./main migrate` if migration is not required then run `./main`
9. Finally application will be running locally on 8080 port [Swagger Link](http://localhost:8080/swagger/index.html)


## Docker Alternative

Build a dev image

```
docker build -t clamp-docker .

```

Run an instance

```
docker run -d -p 9090:8080 clamp-docker
```

The command above will utilize port `8080` of your host.
You can change it to any other port via `-p ANYOTHERPORT:8080`

## Documentation

- [Clamp Swagger Documentation](http://54.190.89.41:8080/swagger/index.html)

## Monitoring & Metrics

``
Grafana Credentials: admin/admin
``

- [Prometheus Dashboard](http://54.190.89.41:9090/graph)

- [Grafana Clamp Dashboard](http://54.190.89.41:3000/d/Ai6xpCgGz/clamp-dashboard?orgId=1)

- [Grafana System Dashboard](http://54.190.89.41:3000/d/rYdddlPWk/node-exporter-full?orgId=1&from=now-3h&to=now)

## Backlogs & Issues

- [Trello Dashboard](https://trello.com/b/oFb5UxvS/clamp)
