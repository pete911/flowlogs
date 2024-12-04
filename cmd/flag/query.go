package flag

import (
	"github.com/pete911/flowlogs/internal/aws/fields"
	"github.com/spf13/cobra"
)

var Query QueryFlags

type QueryFlags struct {
	limit        int
	sinceMinutes int
	protocol     string
	ingress      bool
	egress       bool
	accept       bool
	reject       bool
	port         int
	addr         string
	srcPort      int
	srcAddr      string
	pktSrcAddr   string
	dstPort      int
	dstAddr      string
	pktDstAddr   string
}

func (q QueryFlags) GetQuery() fields.Query {
	query := fields.NewQuery(q.limit, q.sinceMinutes)
	if q.protocol != "" {
		query = query.Protocol(q.protocol)
	}
	if q.egress {
		query = query.Egress()
	}
	if q.ingress {
		query = query.Ingress()
	}
	if q.accept {
		query = query.Accept()
	}
	if q.reject {
		query = query.Reject()
	}
	if q.port > -1 {
		query = query.Port(q.port)
	}
	if q.addr != "" {
		query = query.Address(q.addr)
	}
	if q.srcPort > -1 {
		query = query.SourcePort(q.srcPort)
	}
	if q.srcAddr != "" {
		query = query.SourceAddress(q.srcAddr)
	}
	if q.pktSrcAddr != "" {
		query = query.PktSourceAddress(q.pktSrcAddr)
	}
	if q.dstPort > -1 {
		query = query.DestinationPort(q.dstPort)
	}
	if q.dstAddr != "" {
		query = query.DestinationAddress(q.dstAddr)
	}
	if q.pktDstAddr != "" {
		query = query.PktDestinationAddress(q.pktDstAddr)
	}
	return query.Sort()
}

func InitPersistentQueryFlags(cmd *cobra.Command, flags *QueryFlags) {
	cmd.PersistentFlags().IntVar(
		&flags.limit,
		"limit",
		getIntEnv("LIMIT", 100),
		"number of returned results",
	)
	cmd.PersistentFlags().IntVar(
		&flags.sinceMinutes,
		"minutes",
		getIntEnv("MINUTES", 60),
		"minutes 'ago' to search logs",
	)
	cmd.PersistentFlags().StringVar(
		&flags.protocol,
		"protocol",
		getStringEnv("PROTOCOL", ""),
		"protocol",
	)
	cmd.PersistentFlags().BoolVar(
		&flags.ingress,
		"ingress",
		getBoolEnv("INGRESS", false),
		"ingress flow logs",
	)
	cmd.PersistentFlags().BoolVar(
		&flags.egress,
		"egress",
		getBoolEnv("EGRESS", false),
		"egress flow logs",
	)
	cmd.PersistentFlags().BoolVar(
		&flags.accept,
		"accept",
		getBoolEnv("ACCEPT", false),
		"accepted traffic",
	)
	cmd.PersistentFlags().BoolVar(
		&flags.reject,
		"reject",
		getBoolEnv("REJECT", false),
		"rejected traffic",
	)
	cmd.PersistentFlags().IntVar(
		&flags.port,
		"port",
		getIntEnv("PORT", -1),
		"port - source or destination, negative value means all ports",
	)
	cmd.PersistentFlags().StringVar(
		&flags.addr,
		"addr",
		getStringEnv("ADDR", ""),
		"address - source, destination or packet",
	)
	cmd.PersistentFlags().IntVar(
		&flags.srcPort,
		"src-port",
		getIntEnv("SRC_PORT", -1),
		"source port, negative value means all ports",
	)
	cmd.PersistentFlags().StringVar(
		&flags.srcAddr,
		"src-addr",
		getStringEnv("SRC_ADDR", ""),
		"source address",
	)
	cmd.PersistentFlags().StringVar(
		&flags.pktSrcAddr,
		"pkt-src-addr",
		getStringEnv("PKT_SRC_ADDR", ""),
		"packet source address",
	)
	cmd.PersistentFlags().IntVar(
		&flags.dstPort,
		"dst-port",
		getIntEnv("DST_PORT", -1),
		"destination port, negative value means all ports",
	)
	cmd.PersistentFlags().StringVar(
		&flags.dstAddr,
		"dst-addr",
		getStringEnv("DST_ADDR", ""),
		"destination address",
	)
	cmd.PersistentFlags().StringVar(
		&flags.pktDstAddr,
		"pkt-dst-addr",
		getStringEnv("PKT_DST_ADDR", ""),
		"packet destination address",
	)
}
