package query

// FlowLogFields used when creating flow logs - https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html
var FlowLogFields = []string{
	"interface-id", "srcaddr", "dstaddr", "srcport", "dstport", "protocol", "packets", "bytes", "start", "end", // version 2 fields
	"action", "log-status",
	"vpc-id", "subnet-id", "instance-id", "tcp-flags", "type", "pkt-srcaddr", "pkt-dstaddr", // version 3 fields
	"pkt-src-aws-service", "pkt-dst-aws-service", "flow-direction", "traffic-path", // version 5 fields
	"ecs-cluster-arn", "ecs-cluster-name", "ecs-container-instance-arn", "ecs-container-instance-id", "ecs-container-id", // version 7 fields
	"ecs-second-container-id", "ecs-service-name", "ecs-task-definition-arn", "ecs-task-arn", "ecs-task-id",
}

// Fields used when querying flow logs (unsurprisingly naming convention is different from the above fields)
var Fields = []string{
	"@timestamp", "interfaceId", "srcAddr", "dstAddr", "srcPort", "dstPort", "protocol", "packets", "bytes",
	"action",
	"tcpFlags", "pktSrcAddr", "pktDstAddr",
	"flowDirection", "trafficPath",
	"ecsServiceName",
}
