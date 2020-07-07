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
Workflows are created in Clamp by making a POST request to its `/workflow` API endpoint. Here's a sample payload:
```
{
  "description": "Maintenance Request",
  "enabled": true, // not mandatory
  "name": "mtcReq",
  "steps": [
    ...
  ]
}
```
##### Workflow Metadata
There's some basic metadata that needs to be defined when a workflow is created. The following attributes are mandatory and must be present in the request:
- `name` is the unique identifier for any workflow. It is recommended that you keep this short and name every workflow in a consistent manner. Camel case is recommended, but not mandatory. You could hyphenate between words or use underscores, or choose any other convention, as long as you choose one and stick with it. This field **does not accept spaces**.
- `description` is typically a brief title describing what the workflow is for.
- `steps` are used to describe the workflow in terms of service calls and payload transformations. Steps support simple branching strategies, as well as rollback strategies for error scenarios. See more on [defining steps](#defining-steps) below.

##### Defining Steps
The following section covers how to define steps in your workflow specification. Here's a sample step:
```
{
    "canStepExecute": true,
    "enabled": true,
    "mode": "HTTP",
    "name": "create-mtcReq",
    "onFailure": [
        ... (steps)
    ],
    "requestTransform": {},
    "transform": true,
    "transformFormat": "string",
    "type": "string",
    "val": {},
    "when": "string"
}
```

#### View Workflow
Once a workflow is defined, you can view the structure and metadata of the workflow by performing a GET request to the `/workflow/{name}` endpoint for Clamp. For example, if your workflow was called `mtcReq` and Clamp was running on `localhost:8080`, you would make the following cURL request:
```
curl http://localhost:8080/workflow/mtcReq
```
This should return a response as below:
```
{
  "created_at": "string",
  "description": "Maintenance Request",
  "enabled": true,
  "id": "1",
  "name": "mtcReq",
  "steps": [
    {
      "canStepExecute": true,
      "enabled": true,
      "id": 0,
      "mode": "string",
      "name": "string",
      "onFailure": [
        null
      ],
      "requestTransform": {},
      "transform": true,
      "transformFormat": "string",
      "type": "string",
      "val": {},
      "when": "string"
    }
  ],
  "updated_at": "string"
}
```
### Service Requests
A service request essentially tells Clamp to execute a particular workflow. Depending upon the workflow in question, it may or may not require a request body to go along with it. Let us take the example of `mtcReq`, the workflow we created in the above sections. 
#### Creation
By making a POST request on `/serviceRequest/mtcReq`, we can instruct Clamp to start the maintenance request workflow. If our workflow requires a payload, for example, the address at which maintenance is requested, we can send it in the request body.

**Request**:
```
curl -X POST "http://localhost:8080/serviceRequest/mtcReq" -H "accept: application/json" -H "Content-Type: application/json" -d "{\"address\": \"H12, Cloud Town\"}"
```
The above request will trigger the maintenane request workflow with the following initial payload:
```
{
    "address": "H12, Cloud Town"
}
```

**Response**:
```
{
  "pollUrl": "string",
  "serviceRequestId": "string",
  "status": "string"
}
```
- The `pollUrl` will contain the `GET` endpoint which needs to be polled to monitor the status of every service request.
- The `serviceRequestId` is the unique identifier for a service request.
- The `status` field will contain the completion status for the service request.

#### Check Status
The status of a service request can be polled by making a GET request on the `/serviceRequest/{id}` endpoint, where the `{id}` parameter is the service request ID obtained during creation. Hence, the following request:
```
curl http://localhost:8080/serviceRequest/a1bc
```
should respond back with the status of service request "a1bc", which would look as follows:
```
{
  "reason": "string",
  "service_request_id": "string",
  "status": "string",
  "steps": [
    {
      "id": 0,
      "name": "string",
      "payload": {
        "request": {
          "additionalProp1": {}
        },
        "response": {
          "additionalProp1": {}
        }
      },
      "status": "string",
      "time_taken": 0
    }
  ],
  "total_time_in_ms": 0,
  "workflow_name": "string"
}
```