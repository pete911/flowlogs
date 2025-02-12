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
	return Query{
		query:        append(q.query, `| filter (logStatus != "NODATA"`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) NoSkipData() Query {
	return Query{
		query:        append(q.query, `| filter (logStatus != "SKIPDATA"`),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
	}
}

func (q Query) InterfaceId(id string) Query {
	return Query{
		query:        append(q.query, fmt.Sprintf(`| filter (interfaceId == "%s"`, id)),
		limit:        q.limit,
		sinceMinutes: q.sinceMinutes,
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
	protoNumber := protocolFromKeywordToNumber(proto)
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
