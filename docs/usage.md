# Usage Guide

The following sections cover detailed usage details and specifications for the JSON configuration used to describe workflows. Before we begin, let's go over a few basic concepts:

- **Workflow** A workflow is a series of steps executed in a set order to achieve a specific function, such as booking an order, renewing a license, scheduling an appointment, and so forth. In a microservices setup, each step might be taken up by a different µ-service catering to something specific. Clamp is able to understand and execute such workflows given a JSON configuration according to the given specification.

- **Service Request** A service request is a synchronous or asynchronous call made to Clamp, telling it to execute a specific workflow. Depending on the nature of the workflow, a service request might contain a payload or not. Upon creation of service requests, their state can be monitored via their service request IDs.

## Pre-requisites

### Data Store
Clamp makes use of a data store to keep track of service requests, workflows, and payloads while orchestrating between µ-services. Right now, the only data store supported is [PostgreSQL](https://www.postgresql.org/). You can connect Clamp to your own Postgres setup by configuring it in `config/env.go`.

### Message Broker
A message broker is required to facilitate communication between Clamp and all the µ-services in your environment if you choose to not communicate over HTTP. Clamp ships with integrations for RabbitMQ and Kafka. These are also configured the same way as above with the data store, by editing `config/env.go`.

## Integrating services

## Clamp API

The following section covers API documentation for Clamp's REST API, which ships with endpoints handling workflow creation, triggering a service request, and so forth.

### Swagger

If you're someone who likes to get their hands dirty immediately, maybe you would like to check out this [Swagger](http://34.222.166.218:8080/swagger/index.html) link to try out the APIs.

### Workflows

#### Creation

#### View Details

### Service Requests

#### Creation

#### Check Status

