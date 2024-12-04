# AWS Flow Logs

Dynamically create or delete aws flow logs for EC2 instances, security groups, subnet or VPC.

CLI creates AWS Flow Logs for specific group (EC2 instance(s) - grouped by the same name), security group, subnet or VPC).

Logs can be searched either via cli `flowlogs query <instance|sg|subnet|vpc|nat> <flags>` or in CloudWatch Logs Insights
by select log group with `/fl-cli/` prefix.

## usage

If you have multiple accounts you need to prefix command with `AWS_PROFILE=<your profile> flowlogs ...`.

- create `flowlogs create <instance|sg|subnet|vpc|nat>`
- list `flowlogs list` flowlogs created by this cli
- delete `flowlogs delete <instance|sg|subnet|vpc|nat|all>` (use all argument to clean up all flowlogs)
- query `flowlogs query <instance|sg|subnet|vpc|nat>`

```
flowlogs create vpc

flowlogs query vpc
TIME                     FLOW     ACTION  PACKETS  BYTES    PROTOCOL  SRC ADDR  PKT SRC ADDR  SRC PORT  DST ADDR  PKT DST ADDR  DST PORT  TCP FLAGS
2029-02-04 12:13:13.000  egress   ACCEPT  4        216      TCP       10.0.0.0  10.0.0.0      30572     10.0.0.1  10.0.0.1      6379      FIN, SYN
2029-02-04 12:13:06.000  egress   ACCEPT  4        216      TCP       10.0.0.0  10.0.0.0      25023     10.0.0.1  10.0.0.1      6379      FIN, SYN
2029-02-04 12:12:47.000  ingress  ACCEPT  2        160      TCP       10.0.0.1  10.0.0.1      6379      10.0.0.0  10.0.0.0      24444     FIN, SYN, ACK
2029-02-04 12:12:47.000  ingress  ACCEPT  2        160      TCP       10.0.0.1  10.0.0.1      6379      10.0.0.0  10.0.0.0      25259     FIN, SYN, ACK
2029-02-04 12:12:29.000  ingress  ACCEPT  4        216      TCP       10.0.0.0  10.0.0.0      11308     10.0.0.1  10.0.0.1      6379      FIN, SYN
2029-02-04 12:12:29.000  egress   ACCEPT  329      2368184  TCP       10.0.0.1  10.0.0.1      6379      10.0.0.0  10.0.0.0      36493     ACK
...
```

**Available query flags**
 ```
--accept                accepted traffic
--addr string           address - source, destination or packet
--dst-addr string       destination address
--dst-port int          destination port, negative value means all ports (default -1)
--egress                egress flow logs
--ingress               ingress flow logs
--limit int             number of returned results (default 100)
--minutes int           minutes 'ago' to search logs (default 60)
--pkt-dst-addr string   packet destination address
--pkt-src-addr string   packet source address
--port int              port - source or destination, negative value means all ports (default -1)
--protocol string       protocol
--reject                rejected traffic
--src-addr string       source address
--src-port int          source port, negative value means all ports (default -1)
```

## install

### brew

- add tap `brew tap pete911/tap`
- install `brew install flowlogs`

### binary

Download binary from [releases page](https://github.com/pete911/flowlogs/releases). Unzip and move the binary to your PATH.

## release

Releases are published when the new tag is created e.g.
`git tag -m "<message>" v1.0.0 && git push --follow-tags`

## design/architecture

CLI creates CloudWatch log group in the `/fl-cli/<id>` format. It also creates IAM role and flow log either per VPC, 
subnet or ENI (when instance or sg argument is used). 

### aws flow logs

Flow logs are grouped by ENI. If the flow direction is ingress, destination address and destination port belong to the
ENI that produced the logs. If the flow direction is egress, ENI would be source address and source port.

```
+---- eni xyz ----+
|                 |
| +--------------------------------------+
| | +- ingress -+          +-----------+ |
| | | dst Addr  |<---------| src Addr  | |
| | | dst Port  |          | src Port  | |
| | +-----------+          +-----------+ |
| +--------------------------------------+
|                 |
| +--------------------------------------+
| | +- egress --+          +-----------+ |
| | | src Addr  |--------->| dst Addr  | |
| | | src Port  |          | dst Port  | |
| | +-----------+          +-----------+ |
| +--------------------------------------+
+-----------------|
```