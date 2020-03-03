# rabbitmq-cluster
```
                Host
+----------------------------------+
|                                  |
|              Docker              |
|  +---------------------------+   |
|  |                           |   |
|  |  +-----------+  5672      |   |
|  |  |           +<-------------------  5672
|  |  |           |            |   |      AMQP
|  |  |   rab01   |            |   |
|  |  |           | 15672      |   |
|  |  |           +<------------------- 15672
|  |  +----+-+----+            |   |      MGMT UI
|  |       | ^                 |   |
|  |       | |                 |   |
|  |       | |                 |   |
|  |       v |                 |   |
|  |  +----+-+----+  5672      |   |
|  |  |           +<-------------------  5673
|  |  |           |            |   |      AMQP
|  |  |   rab02   |            |   |
|  |  |           | 15672      |   |
|  |  |           +<------------------- 15673
|  |  +-----------+            |   |      MGMT UI
|  |                           |   |
|  +---------------------------+   |
|                                  |
+----------------------------------+
```

## [RabbitMQ ports](https://www.rabbitmq.com/clustering.html#ports)
- TCP 5672 - AMQP
- TCP 15672 - Management UI

## Docker exposed ports
- rab01
  - TCP 5672 - AMQP
  - TCP 15672 - Management UI
- rab02
  - TCP 5673 - AMQP
  - TCP 15673 - Management UI

## Create nodes
```
docker-compose up -d
```

## [Create cluster](https://www.rabbitmq.com/clustering.html)

- Use [.erlang.cookie](https://www.rabbitmq.com/clustering.html#erlang-cookie) to authenticate nodes; Copied during docker-compose

```
alias dockerexec1='docker exec rabbitmq_rab01_1'
alias dockerexec2='docker exec rabbitmq_rab02_1'

dockerexec2 rabbitmqctl stop_app
dockerexec2 rabbitmqctl reset
dockerexec2 rabbitmqctl join_cluster rabbit@rab01
dockerexec2 rabbitmqctl start_app

dockerexec1 rabbitmqctl cluster_status
dockerexec2 rabbitmqctl cluster_status
```

## [Add HA policy](https://www.rabbitmq.com/ha.html)

- Queue name starts with "ha." will be configured with HA

```
dockerexec1 rabbitmqctl set_policy ha-all "^ha\." '{"ha-mode":"all"}'
```

## Send data to queue

```
watch -n 1 go run helloworld/send/send.go
watch -n 1 go run helloworld/send/send2.go
```

