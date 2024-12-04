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
# wait couple of minute for AWS to aggregate flow logs
flowlogs query vpc
TIME      NI ID                  NI ADDRESS  NI PORT  FLOW        ADDRESS          PORT   ACTION  PACKETS  BYTES  PROTOCOL  TCP FLAGS  TRAFFIC PATH
21:43:55  eni-xxxxxxxxxxxxxxxxx  10.0.0.1    8075     <-ingress-  147.185.133.190  55053  REJECT  1        44     TCP       SYN        
21:43:55  eni-xxxxxxxxxxxxxxxxx  10.0.0.1    22       -egress-->  103.55.49.10     41360  ACCEPT  4        240    TCP       SYN, ACK   internet gateway
21:42:54  eni-xxxxxxxxxxxxxxxxx  10.0.0.1    23       <-ingress-  211.143.253.166  29207  REJECT  1        40     TCP       SYN        
21:42:54  eni-xxxxxxxxxxxxxxxxx  10.0.0.1    17933    <-ingress-  83.222.191.42    61000  REJECT  1        40     TCP       SYN        
...
```

Use `--pretty` flag to add network interface type and name columns.

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
--ni-id string          network interface id
--pkt-dst-addr string   packet destination address
--pkt-src-addr string   packet source address
--port int              port - source or destination, negative value means all ports (default -1)
--pretty                whether to enhance flow logs with names
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
