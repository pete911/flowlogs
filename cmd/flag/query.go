package flag

import (
	"github.com/pete911/flowlogs/internal/aws/query"
	"github.com/spf13/cobra"
)

var Query QueryFlags

type QueryFlags struct {
	Pretty       bool
	limit        int
	sinceMinutes int
	niId         string
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

func (f QueryFlags) GetQuery() query.Query {
	q := query.NewQuery(f.limit, f.sinceMinutes)
	q = q.NoNoData().NoSkipData()
	if f.niId != "" {
		q = q.InterfaceId(f.niId)
	}
	if f.protocol != "" {
		q = q.Protocol(f.protocol)
	}
	if f.egress {
		q = q.Egress()
	}
	if f.ingress {
		q = q.Ingress()
	}
	if f.accept {
		q = q.Accept()
	}
	if f.reject {
		q = q.Reject()
	}
	if f.port > -1 {
		q = q.Port(f.port)
	}
	if f.addr != "" {
		q = q.Address(f.addr)
	}
	if f.srcPort > -1 {
		q = q.SourcePort(f.srcPort)
	}
	if f.srcAddr != "" {
		q = q.SourceAddress(f.srcAddr)
	}
	if f.pktSrcAddr != "" {
		q = q.PktSourceAddress(f.pktSrcAddr)
	}
	if f.dstPort > -1 {
		q = q.DestinationPort(f.dstPort)
	}
	if f.dstAddr != "" {
		q = q.DestinationAddress(f.dstAddr)
	}
	if f.pktDstAddr != "" {
		q = q.PktDestinationAddress(f.pktDstAddr)
	}
	return q.Sort()
}

func InitPersistentQueryFlags(cmd *cobra.Command, flags *QueryFlags) {
	cmd.PersistentFlags().BoolVar(
		&flags.Pretty,
		"pretty",
		getBoolEnv("PRETTY", false),
		"whether to enhance flow logs with names",
	)
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
		&flags.niId,
		"ni-id",
		getStringEnv("NI_ID", ""),
		"network interface id",
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
