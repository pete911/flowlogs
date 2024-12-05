package query

import (
	"fmt"
	"strings"
)

type FlowLogFields []string

func (f FlowLogFields) Format() string {
	var logFormat []string
	for _, v := range f {
		logFormat = append(logFormat, fmt.Sprintf("${%s}", v))
	}
	return strings.Join(logFormat, " ")
}

// FlowLogFieldsV2V5 used when creating flow logs - https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html
var FlowLogFieldsV2V5 = FlowLogFields{
	"interface-id", "srcaddr", "dstaddr", "srcport", "dstport", "protocol", "packets", "bytes", "start", "end", // version 2 fields
	"action", "log-status",
	"vpc-id", "subnet-id", "instance-id", "tcp-flags", "type", "pkt-srcaddr", "pkt-dstaddr", // version 3 fields
	"pkt-src-aws-service", "pkt-dst-aws-service", "flow-direction", "traffic-path", // version 5 fields
}

// FlowLogFieldsV7 V7 fields can only be created if there is at least one ECS cluster in VPC
// this is another crazy half-baked product by AWS, what if we want to create ECS cluster after?
var FlowLogFieldsV7 = FlowLogFields{
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
