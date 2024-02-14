# Python workers Variables

Python workers handle the real-time processing of information within the data pipeline. These configurations can assist in setting up additional nodes for workers. These nodes are organized into clusters known as worker groups, which can be allocated to particular pipelines or designated for an individual stage within a pipeline.

:::info Environments Workers and their associated worker groups are specific to an environment. Only pipelines in that environment can be run on worker groups in the same environment. The isolation of environments is an important concept in Dataplane to assist data operations in segregating access, projects and compute resources. :::

#### Environment variables common across Dataplane and python workers

| Environment variable    | Description                                                                                                                                                                                |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| secret\_db\_host        | Host of the Postgresql database                                                                                                                                                            |
| secret\_db\_user        | User for connection to Postgresql database                                                                                                                                                 |
| secret\_db\_pwd         | Password for connection to Postgresql database                                                                                                                                             |
| secret\_db\_ssl         | One of disable, allow, prefer, require, verify-ca, verify-full - https://www.postgresql.org/docs/current/libpq-ssl.html                                                                    |
| secret\_db\_port        | Database port                                                                                                                                                                              |
| secret\_db\_database    | Database name, default dataplane                                                                                                                                                           |
| secret\_jwt\_secret     | Generate a UUID secret for JWT. It is important that you keep this secret safe. To create a secret, you can use an online generator for example https://www.uuidgenerator.net/             |
| secret\_encryption\_key | Generate a 32 charater long random password. It is important you keep this password safe. You can use an online generator for example https://www.lastpass.com/features/password-generator |

#### Environment variables specifc to workers

| Environment variable           | Options              | Example                                 | Description                                                                                                                                                                 |
| ------------------------------ | -------------------- | --------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| DP\_NATS                       |                      | nats://nats:4222, nats://nats-r\_1:4222 | Connection string to NATS                                                                                                                                                   |
| DP\_DEBUG                      | "true", "false"      | false                                   | Print debug logs to console. Recommended to turn off in production.                                                                                                         |
| DP\_DB\_DEBUG                  | "true", "false"      | false                                   | Print database debug logs to console.                                                                                                                                       |
| DP\_MQ\_DEBUG                  | "true", "false"      | false                                   | Print message queue debug logs to console.                                                                                                                                  |
| DP\_METRIC\_DEBUG              | "true", "false"      | false                                   | Print CPU and memory metrics debug logs to console.                                                                                                                         |
| DP\_WORKER\_HEARTBEAT\_SECONDS |                      | 1                                       | The interval in seconds that the worker sends a heart beat to the main app.                                                                                                 |
| DP\_WORKER\_GROUP              |                      | python\_1                               | The worker group is the collection of worker nodes that have the same configuration. For example, a python worker group that runs the python scripts in the pipeline.       |
| DP\_WORKER\_CMD                |                      | /bin/bash                               | The shell command installed on the linux. This is useful for different linux installations.                                                                                 |
| DP\_WORKER\_TYPE               | "container", "other" | container                               | The worker type is for CPU and memory metrics collection. This can differ between a containerised or bare metal installation. If unsure, recommended to keep it to "other". |
| DP\_WORKER\_LB                 | "roundrobin"         | roundrobin                              | The load balancer strategy is how analytical workloads are distributed to worker nodes.                                                                                     |
| DP\_WORKER\_ENV                |                      | Development                             | This is the name of the environment the worker node belongs to. This must match environments set inside the main app.                                                       |
| DP\_WORKER\_PORT               |                      | 9005                                    | The port that the worker node runs on.                                                                                                                                      |
