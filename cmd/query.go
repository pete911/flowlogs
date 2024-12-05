package cmd

import (
	"fmt"
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/out"
	"github.com/pete911/flowlogs/cmd/prompt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"github.com/pete911/flowlogs/internal/aws/query"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strings"
)

var (
	Query = &cobra.Command{
		Use:   "query",
		Short: "query AWS flow logs",
		Long:  "",
	}

	QueryInstance = &cobra.Command{
		Use:     "instance",
		Aliases: []string{"instances", "ec2", "ec2s"},
		Short:   "query flow logs for specific instance",
		Long:    "",
		Run:     runQuery,
	}

	QueryNat = &cobra.Command{
		Use:     "nat",
		Aliases: []string{"nats", "nat-gateway", "nat-gateways"},
		Short:   "query flow logs for specific nat gateway",
		Long:    "",
		Run:     runQuery,
	}

	QuerySG = &cobra.Command{
		Use:     "sg",
		Aliases: []string{"sgs", "security-group", "security-groups"},
		Short:   "query flow logs for specific security group",
		Long:    "",
		Run:     runQuery,
	}

	QuerySubnet = &cobra.Command{
		Use:     "subnet",
		Aliases: []string{"subnets"},
		Short:   "query flow logs for specific subnet",
		Long:    "",
		Run:     runQuery,
	}

	QueryVPC = &cobra.Command{
		Use:     "vpc",
		Aliases: []string{"vpcs"},
		Short:   "query flow logs for vpc",
		Long:    "",
		Run:     runQuery,
	}
)

func init() {
	flag.InitPersistentQueryFlags(Query, &flag.Query)
	Root.AddCommand(Query)
	Query.AddCommand(QueryInstance)
	Query.AddCommand(QueryNat)
	Query.AddCommand(QuerySG)
	Query.AddCommand(QuerySubnet)
	Query.AddCommand(QueryVPC)
}

func runQuery(cmd *cobra.Command, _ []string) {
	var flowLogType aws.FlowLogType
	switch cmd.Name() {
	case "instance":
		flowLogType = aws.FlowLogTypeInstance
	case "nat":
		flowLogType = aws.FlowLogTypeNatGateway
	case "sg":
		flowLogType = aws.FlowLogTypeSecurityGroup
	case "subnet":
		flowLogType = aws.FlowLogTypeSubnet
	case "vpc":
		flowLogType = aws.FlowLogTypeVPC
	}

	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedFlowLogs := prompt.SelectFlowLogs(prompt.ListFlowLogs(client, flowLogType), false)
	logs, err := client.QueryFlowLogs(selectedFlowLogs, flag.Query.GetQuery())
	if err != nil {
		fmt.Printf("query flow logs: %v\n", err)
		os.Exit(1)
	}

	if flag.Query.Pretty {
		interfaces, err := client.ListNetworkInterfaces()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		prettyPrintQuery(logger, logs, interfaces)
		return
	}
	printQuery(logger, logs)
}

func printQuery(logger *slog.Logger, logs []map[string]string) {
	table := out.NewTable(logger, os.Stdout)
	table.AddRow("TIME", "NI ID", "NI ADDRESS", "NI PORT", "FLOW", "ADDRESS", "PORT", "ACTION", "PACKETS", "BYTES", "PROTOCOL", "TCP FLAGS", "TRAFFIC PATH")
	for _, row := range logs {
		flow := ToFlow(row)
		table.AddRow(
			query.ToTime(row["@timestamp"]), row["interfaceId"], flow.NiAddr, flow.NiPort, flow.Flow, flow.Addr, flow.Port, row["action"],
			row["packets"], row["bytes"], query.ToProtocolKeyword(row["protocol"]),
			strings.Join(query.ToTcpFlagNames(row["tcpFlags"]), ", "),
			query.ToPathName(row["trafficPath"]),
		)
	}
	table.Print()
}

func prettyPrintQuery(logger *slog.Logger, logs []map[string]string, interfaces ec2.NetworkInterfaces) {
	table := out.NewTable(logger, os.Stdout)
	table.AddRow("TIME", "NI ID", "TYPE", "NAME", "NI ADDRESS", "NI PORT", "FLOW", "ADDRESS", "PORT", "ACTION", "PACKETS", "BYTES", "PROTOCOL", "TCP FLAGS", "TRAFFIC PATH")
	for _, row := range logs {
		flow := ToFlow(row)
		ni := interfaces.GetById(row["interfaceId"])
		name := ni.Name
		if ecsSvc := row["ecsServiceName"]; ecsSvc != "" {
			name = ecsSvc
		}

		table.AddRow(
			query.ToTime(row["@timestamp"]), row["interfaceId"], ni.Type, name, flow.NiAddr, flow.NiPort, flow.Flow, flow.Addr, flow.Port, row["action"],
			row["packets"], row["bytes"], query.ToProtocolKeyword(row["protocol"]),
			strings.Join(query.ToTcpFlagNames(row["tcpFlags"]), ", "),
			query.ToPathName(row["trafficPath"]),
		)
	}
	table.Print()
}

type Flow struct {
	Flow   string
	NiAddr string
	NiPort string
	Addr   string
	Port   string
}

func ToFlow(in map[string]string) Flow {
	if in["flowDirection"] == "ingress" {
		niAddr := in["dstAddr"]
		//if in["dstAddr"] != in["pktDstAddr"] {
		//	niAddr = fmt.Sprintf("%s <- %s", in["pktDstAddr"], in["dstAddr"])
		//}
		addr := in["srcAddr"]
		//if in["srcAddr"] != in["pktSrcAddr"] {
		//	addr = fmt.Sprintf("%s <- %s", in["srcAddr"], in["pktSrcAddr"])
		//}
		return Flow{
			Flow:   "<-ingress-",
			NiAddr: niAddr,
			NiPort: in["dstPort"],
			Addr:   addr,
			Port:   in["srcPort"],
		}
	}
	if in["flowDirection"] == "egress" {
		niAddr := in["srcAddr"]
		//if in["srcAddr"] != in["pktSrcAddr"] {
		//	niAddr = fmt.Sprintf("%s -> %s", in["srcAddr"], in["pktSrcAddr"])
		//}
		addr := in["dstAddr"]
		//if in["dstAddr"] != in["pktDstAddr"] {
		//	addr = fmt.Sprintf("%s -> %s", in["pktDstAddr"], in["dstAddr"])
		//}
		return Flow{
			Flow:   "-egress-->",
			NiAddr: niAddr,
			NiPort: in["srcPort"],
			Addr:   addr,
			Port:   in["dstPort"],
		}
	}
	return Flow{}
}
