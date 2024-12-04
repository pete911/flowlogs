package fields

import (
	"fmt"
	"strings"
)

// FlowLogFields - https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html
var FlowLogFields = []string{
	"interface-id", "srcaddr", "dstaddr", "srcport", "dstport", "protocol", "packets", "bytes", "start", "end", // version 2 fields
	"action", "log-status",
	"vpc-id", "subnet-id", "instance-id", "tcp-flags", "type", "pkt-srcaddr", "pkt-dstaddr", // version 3 fields
	"pkt-src-aws-service", "pkt-dst-aws-service", "flow-direction", "traffic-path", // version 5 fields
	"ecs-cluster-arn", "ecs-cluster-name", "ecs-container-instance-arn", "ecs-container-instance-id", "ecs-container-id", // version 7 fields
	"ecs-second-container-id", "ecs-service-name", "ecs-task-definition-arn", "ecs-task-arn", "ecs-task-id",
}

var QueryFields = QueryFieldItems{
	{name: "@timestamp", header: "TIME"},
	{name: "flowDirection", header: "FLOW"},
	{name: "action", header: "ACTION"},
	{name: "packets", header: "PACKETS"},
	{name: "bytes", header: "BYTES"},
	{name: "protocol", header: "PROTOCOL"},
	{name: "srcAddr", header: "SRC ADDR"},
	{name: "pktSrcAddr", header: "PKT SRC ADDR"},
	{name: "srcPort", header: "SRC PORT"},
	{name: "dstAddr", header: "DST ADDR"},
	{name: "pktDstAddr", header: "PKT DST ADDR"},
	{name: "dstPort", header: "DST PORT"},
	{name: "tcpFlags", header: "TCP FLAGS"},
}

type QueryFieldItem struct {
	name   string
	header string
}

type QueryFieldItems []QueryFieldItem

func (f QueryFieldItems) names() []string {
	var out []string
	for _, v := range f {
		out = append(out, v.name)
	}
	return out
}

func (f QueryFieldItems) Header() []string {
	var out []string
	for _, v := range f {
		out = append(out, v.header)
	}
	return out
}

func (f QueryFieldItems) Values(in []map[string]string) [][]string {
	var rows [][]string
	for _, row := range in {
		var columns []string
		for _, fieldName := range f.names() {
			column := row[fieldName]
			if fieldName == "tcpFlags" {
				column = strings.Join(toTcpFlagNames(column), ", ")
			}
			if fieldName == "protocol" {
				column = toProtocol(column).keyword
			}
			if fieldName == "trafficPath" {
				column = toPathName(column)
			}
			columns = append(columns, column)
		}
		rows = append(rows, columns)
	}
	return rows
}

type Query struct {
	query        []string
	limit        int
	sinceMinutes int
}

func NewQuery(limit, sinceMinutes int) Query {
	return Query{
		query:        []string{fmt.Sprintf("fields %s", strings.Join(QueryFields.names(), ", "))},
		limit:        limit,
		sinceMinutes: sinceMinutes,
	}
}

func (q Query) Ingress() Query {
	return Query{
		query:        append(q.query, `| filter (flowDirection == "ingress"`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Egress() Query {
	return Query{
		query:        append(q.query, `| filter (flowDirection == "egress"`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Accept() Query {
	return Query{
		query:        append(q.query, `| filter (action == "ACCEPT"`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Reject() Query {
	return Query{
		query:        append(q.query, `| filter (action == "REJECT"`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Protocol(proto string) Query {
	protoNumber := fromProtocol(proto)
	// not found, query all protocols
	if protoNumber < 0 {
		return q
	}
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (protocol == "%d"`, protoNumber)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Port(port int) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (srcPort == "%d" or dstPort == "%d"`, port, port)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) SourcePort(port int) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (srcPort == "%d"`, port)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) DestinationPort(port int) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (dstPort == "%d"`, port)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Address(addr string) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (srcAddr == "%s" or pktSrcAddr == "%s" or dstAddr == "%s" or pktDstAddr == "%s"`, addr, addr, addr, addr)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) SourceAddress(addr string) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (srcAddr == "%s"`, addr)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) PktSourceAddress(addr string) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (pktSrcAddr == "%s"`, addr)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) DestinationAddress(addr string) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (dstAddr == "%s"`, addr)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) PktDestinationAddress(addr string) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (pktDstAddr == "%s"`, addr)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) Sort() Query {
	return Query{
		query:        append(q.query, `| sort @timestamp desc`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) GetQuery() string {
	return strings.Join(q.query, "\n")
}

func (q Query) GetLimit() int {
	return q.limit
}

func (q Query) GetSinceMinutes() int {
	return q.sinceMinutes
}
