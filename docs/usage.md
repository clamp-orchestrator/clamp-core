# Usage Guide
The following sections cover detailed usage details and specifications for the JSON configuration used to describe workflows. Before we begin, let's go over a few basic concepts:
- **Workflow**: A workflow is a series of steps executed in a set order to achieve a specific function, such as booking an order, renewing a license, scheduling an appointment, and so forth. In a microservices setup, each step will be taken up by a different µ-service catering to something specific. Clamp is able to understand and execute such workflows given a JSON configuration according to the given specification.
- **Step**: A step is the smallest unit in a workflow. A step can be used to execute a specific function. A step can execute in 2 modes SYNC & ASYNC. Step can be integrated with HTTP API / Kafka / Rabbit MQ. Step can be conditionally executed and the request for the step can be transformed in runtime.
    - **Step Modes**: Sync step is a type of step where clamp expects a response after executing it. Async step is a step where the response can be deferred, the downstream service can send back the response latter in time.
- **Service Request**: A service request is an asynchronous call made to Clamp, telling it to execute a specific workflow. Depending on the nature of the workflow, a service request might contain a payload or not. Upon creation of service requests, their state can be monitored via their service request IDs.
- **Context Object**: In the process of executing a workflow, a context object will store the request and respose of each step. This request context can be used in step for conditional executing it or transforming the request to the step.

## Pre-requisites
### Data Store
Clamp makes use of a data store to keep track of workflows, service requests and payloads while orchestrating between µ-services. Right now, the only data store supported is [PostgreSQL](https://www.postgresql.org/). You can connect Clamp to your own Postgres setup by configuring it in `config/env.go`.
```
DBConnectionStr string `env:"CLAMP_DB_CONNECTION_STR" envDefault:"host=<ip_address>:<port> user=<user_name> dbname=<db_name> password=<password>"`
```
### Message Broker
A message broker is required to facilitate async communication between Clamp and all the µ-services in your environment if you choose to not communicate over HTTP. Clamp ships with integrations for RabbitMQ and Kafka. These are also configured the same way as above with the data store, by editing `config/env.go`.
- AMQP
```
QueueConnectionStr string `env:"CLAMP_QUEUE_CONNECTION_STR" envDefault:"amqp://<user_name>:<password>@<ip_address>:<port>/"`
```
- Kafka
```
KafkaConnectionStr string `env:"CLAMP_KAFKA_CONNECTION_STR" envDefault:"<ip_address>:<port>"`
```
## Clamp API
The following section covers API documentation for Clamp's REST API, which ships with endpoints handling workflow creation, triggering a service request, and so forth.
### Swagger
If you're someone who likes to get their hands dirty immediately, maybe you would like to check out this [Swagger](http://54.149.76.62:8642/swagger/index.html) link to try out the APIs.
### Workflows
#### Creation
Workflows are created in Clamp by making a **POST** request to its **`/workflow`** API endpoint. 
<details>
  <summary> Here's a sample payload: (Click to expand)</summary>
  
    ```
    {
      "name": "process_claim",
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
    "name": "submit motor claim",
    "when": "user_authentication.request.claimDetails.claimType == 'MOTOR'",
    "mode": "AMQP",
    "transform" : true,
    "requestTransform": {
        "spec":{
            "claimId": "create_claim.response.claimId",
            "userId": "user_authentication.response.id",
            "claimStatus": "create_claim.response.claimStatus",
            "claimType": "user_authentication.request.claimDetails.claimType",
            "claimDate": "create_claim.response.claimDate",
            "policyId": "create_claim.response.policyId",
            "garageId":"create_claim.response.garageId",
            "inspectorDetails":"create_claim.response.inspectorDetails"
        }
    },
    "val": {
        "connection_url": "amqp://clamp:clampdev!@172.31.0.152:5672/",
        "queue_name": "clamp_queue",
        "content_type": "text/plain"
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
                "url": "https://run.mocky.io/v3/f4ee8258-49b1-4579-a8c2-5881a0c65206",
                "headers": "Content-Type:application/json"
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
- `when` is used to specify the condition based on which the step execution depends. The possible options to use for comparision are [here](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md#comparison-operators)
- `requestTransform` is used for transforming the request object to the step. `transform` need to be enabled to apply the transformations.
- `Context object` can be used in both `when` and `requestTransform`. The context object can be accessed by directly specifying the `step_name` and then specify whether `request` or `response` and specify the key to access it. Ex:`step_name.response.key`. It can be nested to any level like `step_name.response.key1.key1a`

#### View Workflow
Once a workflow is defined, you can view the structure and metadata of the workflow by performing a GET request to the `/workflow/{name}` endpoint for Clamp. For example, if your workflow was called `mtcReq` and Clamp was running on `http://54.149.76.62:8642`, you would make the following cURL request:
```
curl http://54.149.76.62:8642/workflow/process_claim
```

<details>
  <summary>This should return a response as below:(Click here to expand)</summary>
  
```
{
    "id": "348",
    "name": "process_claim",
    "description": "processing of medical claim",
    "enabled": true,
    "created_at": "2020-09-14T12:08:37.989947Z",
    "updated_at": "0001-01-01T00:00:00Z",
    "steps": [
        {
            "id": 1,
            "name": "user_authentication",
            "type": "SYNC",
            "mode": "HTTP",
            "val": {
                "method": "POST",
                "url": "https://run.mocky.io/v3/f4ee8258-49b1-4579-a8c2-5881a0c65206",
                "headers": ""
            },
            "transform": false,
            "enabled": false,
            "when": "",
            "transformFormat": "",
            "requestTransform": null,
            "onFailure": null
        },
        {
            "id": 2,
            "name": "user_authorization",
            "type": "SYNC",
            "mode": "HTTP",
            "val": {
                "method": "POST",
                "url": "https://run.mocky.io/v3/d9f2e6d1-3100-4ffb-88a4-633e89e1b99c",
                "headers": ""
            },
            "transform": true,
            "enabled": false,
            "when": "",
            "transformFormat": "",
            "requestTransform": {
                "spec": {
                    "userId": "user_authentication.response.id",
                    "username": "user_authentication.response.username"
                }
            },
            "onFailure": null
        },
        ...
        ...
        {
            "id": 8,
            "name": "update_reject_claim",
            "type": "SYNC",
            "mode": "HTTP",
            "val": {
                "method": "POST",
                "url": "https://run.mocky.io/v3/b0ab4d1c-263b-41f5-9888-c8913160c20f",
                "headers": ""
            },
            "transform": false,
            "enabled": false,
            "when": "update_reject_claim.request.claimStatus == 'REJECTED'",
            "transformFormat": "",
            "requestTransform": null,
            "onFailure": null
        },
        {
            "id": 9,
            "name": "process_disbursement",
            "type": "SYNC",
            "mode": "HTTP",
            "val": {
                "method": "POST",
                "url": "https://run.mocky.io/v3/a2a9bb05-f043-4a6e-b513-0377902bd85d",
                "headers": ""
            },
            "transform": true,
            "enabled": false,
            "when": "update_approved_claim.request.claimStatus == 'APPROVED'",
            "transformFormat": "",
            "requestTransform": {
                "spec": {
                    "approvedAmount": "process_disbursement.request.reviewerDetails.approvedAmount",
                    "claimId": "create_claim.response.claimId",
                    "claimStatus": "process_disbursement.request.claimStatus",
                    "reviewerDate": "process_disbursement.request.reviewerDetails.reviewDate",
                    "reviewerId": "process_disbursement.request.reviewerDetails.reviewerId",
                    "userId": "user_authentication.response.id"
                }
            },
            "onFailure": null
        }
    ]
}
```
</details>

### Service Requests
A service request essentially tells Clamp to execute a particular workflow. Depending upon the workflow, it may or may not require a request body to go along with it. Let us take the example of `process_claim`, the workflow we created in the above sections. 
#### Creation
By making a POST request on `/serviceRequest/process_claim`, we can instruct Clamp to start the `process_claim` workflow. If our workflow requires an initial payload, we can send it in the request body.

**Request**:
```
curl -X POST 'http://54.149.76.62:8642/serviceRequest/process_claim' \
--header 'Content-Type: application/json' \
--data-raw '{
    "userDetails" :{
    	"username": "xyz",
		"password": "***",
		"channel": "web"
    },
	"claimDetails":{
		"claimType":"MOTOR",
		"claimDate":"23/06/2020",
		"policyId":"908",
		"garageId":"5000",
		"supportingDocuments":""
	}
}'
```
The above request will trigger the `process_claim` workflow with the following initial payload:
```
{
    "userDetails" :{
    	"username": "xyz",
		"password": "***",
		"channel": "web"
    },
	"claimDetails":{
		"claimType":"MOTOR",
		"claimDate":"23/06/2020",
		"policyId":"908",
		"garageId":"5000",
		"supportingDocuments":""
	}
}
```

**Response**:
```
{
    "pollUrl": "/serviceRequest/6102fa39-d209-4b98-8c75-d9f2ef9aa791",
    "status": "NEW",
    "serviceRequestId": "6102fa39-d209-4b98-8c75-d9f2ef9aa791"
}}
```
- The `pollUrl` will contain the `GET` endpoint which needs to be polled to monitor the status of every service request.
- The `serviceRequestId` is the unique identifier for a service request.
- The `status` field will contain the completion status for the service request.

#### Check Status
The status of a service request can be polled by making a GET request on the `/serviceRequest/{id}` endpoint, where the `{id}` parameter is the service request ID obtained during creation. Hence, the following request:
```
curl http://54.149.76.62:8642/serviceRequest/6102fa39-d209-4b98-8c75-d9f2ef9aa791
```
should respond back with the status of service request "6102fa39-d209-4b98-8c75-d9f2ef9aa791", which would look as follows:
```

    "service_request_id": "6102fa39-d209-4b98-8c75-d9f2ef9aa791",
    "workflow_name": "process_claim",
    "status": "COMPLETED",
    "total_time_in_ms": 170382,
    "steps": [
        {
            "id": 1,
            "name": "user_authentication",
            "status": "STARTED",
            "time_taken": 0,
            "payload": {
                "request": {
                    "claimDetails": {
                        "claimDate": "23/06/2020",
                        "claimType": "MOTOR",
                        "garageId": "5000",
                        "policyId": "908",
                        "supportingDocuments": ""
                    },
                    "userDetails": {
                        "channel": "web",
                        "password": "jungle-green-t0p!",
                        "username": "shambhu.shikari"
                    }
                },
                "response": null
            }
        },
        {
            "id": 1,
            "name": "user_authentication",
            "status": "COMPLETED",
            "time_taken": 798,
            "payload": {
                "request": {
                    "claimDetails": {
                        "claimDate": "23/06/2020",
                        "claimType": "MOTOR",
                        "garageId": "5000",
                        "policyId": "908",
                        "supportingDocuments": ""
                    },
                    "userDetails": {
                        "channel": "web",
                        "password": "jungle-green-t0p!",
                        "username": "shambhu.shikari"
                    }
                },
                "response": {
                    "id": "1234567890",
                    "name": "Shambhu Shikari",
                    "username": "shambhu.shikari"
                }
            }
        },
        ...
        ...
        {
            "id": 9,
            "name": "process_disbursement",
            "status": "STARTED",
            "time_taken": 0,
            "payload": {
                "request": {
                    "approvedAmount": "5000",
                    "claimId": "90990908324",
                    "claimStatus": "APPROVED",
                    "reviewerDate": "2020 Jun 23 00:00:00.000 IST",
                    "reviewerId": "12924",
                    "userId": "1234567890"
                },
                "response": null
            }
        },
        {
            "id": 9,
            "name": "process_disbursement",
            "status": "COMPLETED",
            "time_taken": 514,
            "payload": {
                "request": {
                    "approvedAmount": "5000",
                    "claimId": "90990908324",
                    "claimStatus": "APPROVED",
                    "reviewerDate": "2020 Jun 23 00:00:00.000 IST",
                    "reviewerId": "12924",
                    "userId": "1234567890"
                },
                "response": {
                    "claimId": "90990908324",
                    "disbursedAmount": "5000",
                    "disbursementDate": "2020 Jun 23 00:00:00.000 IST",
                    "disbursementRefId": "234234434",
                    "partyDetails": {
                        "partyId": "23432431",
                        "partyName": "Apple Auto"
                    },
                    "paymentInstrumentId": "CHEQUE",
                    "userId": "1234567890"
                }
            }
        }
    ],
    "reason": ""
}
```
- The `status` field will contain the completion status for the service request. It will contain `IN_PROGRESS` / `COMPLETED`
- `total_time_in_ms` will contain the time taken in ms to execute the complete workflow
- `steps` will contain step level status, it contains both the request/response that is sent/received for each step.
- In each step in `steps` the `status` defines the state of each step it went through. The possible values are `STARTED` / `COMPLETED` / `SKIPPED` / `FAILED`
- The step status will be `SKIPPED` if the `when` condition is not met.

#### Send Response
When an async step gets executed, the response for the step needs to be sent explicitly to clamp. The response can be sent back using an HTTP API or through AMQP / Kafka queue. 
- HTTP

    By making a POST request on `/stepResponse`, we can send response to Clamp. The response can be sent in below format.
    
    **Request**:
    ```
    curl -X POST 'http://54.149.76.62:8642/stepResponse' \
    --header 'Content-Type: application/json' \
    --data-raw '{
                    "serviceRequestId": "{{serviceRequestId}}",
                    "stepId": 5,
                    "response": {
                        "claimId": "90990908324",
                        ...
                        "notes": "Inspection not required approved based on documentation. CASHLESS"
                    }
                }
    }'
    ```
    The above request will trigger the workflow to resume execution.
    
    - When an async step gets executed, **clamp** sends the current id of the service request in the request body. The downstream service needs to send back the same **serviceRequestId** back in the response to continue the same service request. 
    - The `stepId` is the next step id which needs to be executed. The value of this will also be sent as request to downstream services.
    - The `response` field should contain the respective json response that needs to be sent back.
    
- AMQO / KAFKA

    The reponse can be sent back through AMQP / KAFKA. Clamp listens to specific topic in both AMQP / KAFKA. `clamp_steps_response` is the topic name to which the below response can be sent.
    ```
  {
      "serviceRequestId": "{{serviceRequestId}}",
      "stepId": 5,
      "response": {
          "claimId": "90990908324",
          ...
          "notes": "Inspection not required approved based on documentation. CASHLESS"
      }
  }
  ```