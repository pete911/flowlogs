package query

import (
	"fmt"
	"strings"
)

// Query is request to query cloud watch flow logs
type Query struct {
	query        []string
	limit        int
	sinceMinutes int
}

func NewQuery(limit, sinceMinutes int) Query {
	return Query{
		query:        []string{fmt.Sprintf("fields %s", strings.Join(Fields, ", "))},
		limit:        limit,
		sinceMinutes: sinceMinutes,
	}
}

func (q Query) NoNoData() Query {
	return q.add(`| filter logStatus != "NODATA"`)
}

func (q Query) NoSkipData() Query {
	return q.add(`| filter logStatus != "SKIPDATA"`)
}

func (q Query) InterfaceId(id string) Query {
	return q.add(fmt.Sprintf(`| filter interfaceId == "%s"`, id))
}

func (q Query) Ingress() Query {
	return q.add(`| filter flowDirection == "ingress"`)
}

func (q Query) Egress() Query {
	return q.add(`| filter flowDirection == "egress"`)
}

func (q Query) Accept() Query {
	return q.add(`| filter action == "ACCEPT"`)
}

func (q Query) Reject() Query {
	return q.add(`| filter action == "REJECT"`)
}

func (q Query) Protocol(proto string) Query {
	protoNumber := protocolFromKeywordToNumber(proto)
	// not found, query all protocols
	if protoNumber < 0 {
		return q
	}
	return q.add(fmt.Sprintf(`| filter protocol == "%d"`, protoNumber))
}

func (q Query) Port(port int) Query {
	return q.add(fmt.Sprintf(`| filter srcPort == "%d" or dstPort == "%d"`, port, port))
}

func (q Query) SourcePort(port int) Query {
	return q.add(fmt.Sprintf(`| filter srcPort == "%d"`, port))
}

func (q Query) DestinationPort(port int) Query {
	return q.add(fmt.Sprintf(`| filter dstPort == "%d"`, port))
}

func (q Query) Address(addr string) Query {
	return q.add(fmt.Sprintf(`| filter srcAddr == "%s" or pktSrcAddr == "%s" or dstAddr == "%s" or pktDstAddr == "%s"`, addr, addr, addr, addr))
}

func (q Query) SourceAddress(addr string) Query {
	return q.add(fmt.Sprintf(`| filter srcAddr == "%s"`, addr))
}

func (q Query) PktSourceAddress(addr string) Query {
	return q.add(fmt.Sprintf(`| filter pktSrcAddr == "%s"`, addr))
}

func (q Query) DestinationAddress(addr string) Query {
	return q.add(fmt.Sprintf(`| filter dstAddr == "%s"`, addr))
}

func (q Query) PktDestinationAddress(addr string) Query {
	return q.add(fmt.Sprintf(`| filter pktDstAddr == "%s"`, addr))
}

func (q Query) Sort() Query {
	return q.add(`| sort @timestamp desc`)
}

func (q Query) add(in string) Query {
	next := make([]string, len(q.query)+1)
	copy(next, q.query)
	next[len(q.query)] = in
	return Query{
		query:        next,
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
