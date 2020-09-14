# Usage Guide
The following sections cover detailed usage details and specifications for the JSON configuration used to describe workflows. Before we begin, let's go over a few basic concepts:
- **Workflow** A workflow is a series of steps executed in a set order to achieve a specific function, such as booking an order, renewing a license, scheduling an appointment, and so forth. In a microservices setup, each step will be taken up by a different µ-service catering to something specific. Clamp is able to understand and execute such workflows given a JSON configuration according to the given specification.
- **Step**  A step is the smallest unit in a workflow. A step can be used to execute a specific function. A step can execute in 2 modes SYNC & ASYNC. Step can be integrated with HTTP API / Kafka / Rabbit MQ. Step can be conditionally executed and the request for the step can be transformed in runtime.
    - **Step Modes** Sync step is a type of step where clamp expects a response after executing it. Async step is a step where the response can be deferred, the downstream service can send back the response latter in time.
- **Service Request** A service request is an asynchronous call made to Clamp, telling it to execute a specific workflow. Depending on the nature of the workflow, a service request might contain a payload or not. Upon creation of service requests, their state can be monitored via their service request IDs.

## Pre-requisites
### Data Store
Clamp makes use of a data store to keep track of workflows, service requests and payloads while orchestrating between µ-services. Right now, the only data store supported is [PostgreSQL](https://www.postgresql.org/). You can connect Clamp to your own Postgres setup by configuring it in `config/env.go`.
### Message Broker
A message broker is required to facilitate communication between Clamp and all the µ-services in your environment if you choose to not communicate over HTTP. Clamp ships with integrations for RabbitMQ and Kafka. These are also configured the same way as above with the data store, by editing `config/env.go`.
## Clamp API
The following section covers API documentation for Clamp's REST API, which ships with endpoints handling workflow creation, triggering a service request, and so forth.
### Swagger
If you're someone who likes to get their hands dirty immediately, maybe you would like to check out this [Swagger](http://54.149.76.62:8642/swagger/index.html) link to try out the APIs.
### Workflows
#### Creation
Workflows are created in Clamp by making a **POST** request to its **`/workflow`** API endpoint. 
##### Here's a sample payload:
<details>
  <summary>Click here to expand!</summary>
  
```
{
  "name": "process_medical_claim",
  "description": "processing of medical claim",
  "steps": [
    {
      "name": "user_authentication",
      "mode": "HTTP",
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/f4ee8258-49b1-4579-a8c2-5881a0c65206"
      }
    },
    {
      "name": "user_authorization",
      "mode": "HTTP",
      "transform": true,
      "requestTransform": {
        "spec": {
          "username": "user_authentication.response.username",
          "userId": "user_authentication.response.id"
        }
      },
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/d9f2e6d1-3100-4ffb-88a4-633e89e1b99c"
      }
    },
    {
      "name": "get_user_details",
      "mode": "HTTP",
      "transform": true,
      "requestTransform": {
        "spec": {
          "username": "user_authentication.response.username",
          "userId": "user_authentication.response.id",
          "roles": "user_authorization.response.roles"
        }
      },
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/0df407a1-d4ea-41b3-bf2d-31f3c0fe03b5"
      }
    },
    {
      "name": "create_claim",
      "mode": "HTTP",
      "transform": true,
      "requestTransform": {
        "spec": {
          "claimDetails": "user_authentication.request.claimDetails",
          "userId": "user_authentication.response.id",
          "existingPolicies": "get_user_details.response.policyDetails"
        }
      },
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/c73e40b4-a044-44bd-931a-d0f08d58f0d3"
      }
    },
    {
      "name": "submit_motor_claim",
      "when": "user_authentication.request.claimDetails.claimType == 'MOTOR'",
      "mode": "AMQP",
      "transform": true,
      "requestTransform": {
        "spec": {
          "claimId": "create_claim.response.claimId",
          "userId": "user_authentication.response.id",
          "claimStatus": "create_claim.response.claimStatus",
          "claimType": "user_authentication.request.claimDetails.claimType",
          "claimDate": "create_claim.response.claimDate",
          "policyId": "create_claim.response.policyId",
          "garageId": "create_claim.response.garageId",
          "inspectorDetails": "create_claim.response.inspectorDetails"
        }
      },
      "val": {
        "connection_url": "amqp://clamp:clampdev!@172.31.0.152:5672/",
        "queue_name": "clamp_queue",
        "content_type": "text/plain"
      }
    },
    {
      "name": "submit_medical_claim",
      "when": "user_authentication.request.claimDetails.claimType == 'MEDICAL'",
      "mode": "KAFKA",
      "transform": true,
      "requestTransform": {
        "spec": {
          "claimId": "create_claim.response.claimId",
          "userId": "user_authentication.response.id",
          "claimStatus": "submit_medical_claim.request.claimStatus",
          "claimType": "user_authentication.request.claimDetails.claimType",
          "claimDate": "create_claim.response.claimDate",
          "policyId": "create_claim.response.policyId",
          "garageId": "create_claim.response.garageId",
          "inspectorDetails": "create_claim.response.inspectorDetails"
        }
      },
      "val": {
        "connection_url": "172.31.0.152:9092",
        "topic_name": "clamp_topic"
      }
    },
    {
      "name": "update_approved_claim",
      "when": "update_approved_claim.request.claimStatus == 'APPROVED'",
      "mode": "HTTP",
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/39528702-f29f-4a87-98e7-55b43c81fed3"
      }
    },
    {
      "name": "update_reject_claim",
      "when": "update_reject_claim.request.claimStatus == 'REJECTED'",
      "mode": "HTTP",
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/b0ab4d1c-263b-41f5-9888-c8913160c20f"
      }
    },
    {
      "name": "process_disbursement",
      "when": "update_approved_claim.request.claimStatus == 'APPROVED'",
      "mode": "HTTP",
      "transform": true,
      "requestTransform": {
        "spec": {
          "claimId": "create_claim.response.claimId",
          "userId": "user_authentication.response.id",
          "claimStatus": "process_disbursement.request.claimStatus",
          "approvedAmount": "process_disbursement.request.reviewerDetails.approvedAmount",
          "reviewerId": "process_disbursement.request.reviewerDetails.reviewerId",
          "reviewerDate": "process_disbursement.request.reviewerDetails.reviewDate"
        }
      },
      "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/a2a9bb05-f043-4a6e-b513-0377902bd85d"
      }
    }
  ]
}
```
</details>

##### Workflow Metadata
There's some basic metadata that needs to be defined when a workflow is created. The following attributes are mandatory and must be present in the request:
- `name` is the unique identifier for any workflow. It is recommended that you keep this short and name every workflow in a consistent manner. Camel case is recommended, but not mandatory. You could hyphenate between words or use underscores, or choose any other convention, as long as you choose one and stick with it. This field **does not accept spaces**.
- `description` is typically a brief title describing what the workflow is for.
- `steps` are used to describe the workflow in terms of service calls and payload transformations. Steps support simple branching strategies, as well as rollback strategies for error scenarios. See more on [defining steps](#defining-steps) below.

##### Defining Steps
The following section covers how to define steps in your workflow specification. Here's a sample step:
```
{
    "name": "user_authentication",
    "mode": "HTTP",
    "val": {
        "method": "POST",
        "url": "https://run.mocky.io/v3/f4ee8258-49b1-4579-a8c2-5881a0c65206"
    }
}
```
There's some basic metadata that needs to be defined for a step.
- `name` is the unique identifier for a step in a workflow. It is recommended that you keep this short. Camel case is recommended, but not mandatory. You could hyphenate between words or use underscores, or choose any other convention, as long as you choose one and stick with it. This field **does not accept spaces**.
- `mode` specifies what communication mode it needs to use. It supports HTTP / KAFKA / AMQP.
- `val` is the connection config to use for communication. Below are configs specific to each modes.
    - `HTTP`
        ```
        "val": {
                "method": "POST",
                "url": "https://run.mocky.io/v3/f4ee8258-49b1-4579-a8c2-5881a0c65206"
            }
        ```
    - `AMQP`
        ```
        "val": {
                "connection_url": "amqp://clamp:clampdev!@172.31.0.152:5672/",
                "queue_name": "clamp_queue",
                "content_type": "text/plain"
            }
      ```
    - `KAFKA`
        ```
        "val": {
                "connection_url": "172.31.0.152:9092",
                "topic_name": "clamp_topic"
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