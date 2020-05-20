# clamp-core

The clamp core that can be deployed as a binary and does service orchestration

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