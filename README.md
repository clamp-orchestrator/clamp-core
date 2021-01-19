# Overview - clamp-core
 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/clamp-orchestrator/clamp-core/blob/master/LICENSE)
 [![Maintainability](https://api.codeclimate.com/v1/badges/7dae82e6001dcd176930/maintainability)](https://codeclimate.com/repos/5f721f2b64cdeb01a0007ceb/maintainability)
 [![Test Coverage](https://api.codeclimate.com/v1/badges/7dae82e6001dcd176930/test_coverage)](https://codeclimate.com/repos/5f721f2b64cdeb01a0007ceb/test_coverage) [![Join the chat at https://gitter.im/ClampOrchestrator/community](https://badges.gitter.im/ClampOrchestrator/announcements.svg)](https://gitter.im/ClampOrchestrator/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Clamp is a workflow management and microservices orchestrator. It has been written in go-lang. The documentation on various aspects of the framework can be found [here](https://clamp-orchestrator.github.io/clamp-orchestrator/docs/about-docs). 

## Backlogs & Issues

- [Project dashboard](https://github.com/orgs/clamp-orchestrator/projects/1)

## Docker Compose

Clamp Core and its dependencies can be run through Docker Compose with the following command

```bash
$ docker-compose up
```

Optionally, you can run the following command if you want to inspect Clamp Core metrics through Prometheus annd Grafana.

```bash
$ docker-compose -f docker-compose.yml -f prometheus-grafana.yml u
```