package query

var pathNameByNumber = map[string]string{
	"1": "vpc",
	"2": "internet gateway/vpc endpoint",
	"3": "virtual private gateway",
	"4": "intra-region vpc peering",
	"5": "inter-region vpc peering",
	"6": "local gateway",
	"7": "vpc endpoint",     // nitro-based instances only
	"8": "internet gateway", // nitro-based instances only
}

// ToPathName takes traffic-path flow log field and return name representation. This applies only to egress traffic
func ToPathName(in string) string {
	return pathNameByNumber[in]
}
