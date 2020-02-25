# rabbitmq-cluster
RabbitMQ Cluster with HA Queue

## Nodes
- rab01
- rab02

## Access to Management UI
- rab01 - http://IPADDRESS:15672
- rab02 - http://IPADDRESS:15673

## Create nodes
```
docker-compose up -d
```

## Create cluster
https://www.rabbitmq.com/clustering.html
- Use file `.erlang.cookie` to authenticate nodes, copied during docker-compose
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

## Add HA policy
https://www.rabbitmq.com/ha.html
```
dockerexec1 rabbitmqctl set_policy ha-all "^ha\." '{"ha-mode":"all"}'
```

## Send data to queue
```
watch -n 1 go run helloworld/send/send.go
watch -n 1 go run helloworld/send/send2.go
```

