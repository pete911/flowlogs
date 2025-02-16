package query

import "strconv"

func ToTcpFlagNames(in string) []string {
	// tcp flags do not have to be set, do not return error
	if in == "-" || in == "" {
		return nil
	}
	intVal, err := strconv.Atoi(in)
	if err != nil {
		return nil
	}
	if intVal > 255 { // tcp flags are 8 bits
		return nil
	}

	// When a flow log entry consists of only ACK packets, the flag value is 0, not 16 (see AWS docs)
	// But it can be 0, also if no flags are sent, or flag is not supported by AWS flow logs. To be on the safe
	// side, return 0.
	if intVal == 0 {
		return []string{"0"}
	}

	// order is important, it is used to check TCP Flag bits
	var tcpFlags = []string{
		"FIN",
		"SYN",
		"RST",
		"PSH",
		"ACK",
		"URG",
		"RESERVED",
		"RESERVED",
	}

	var out []string
	for i, v := range tcpFlags {
		if byte(intVal)&(1<<i) != 0 {
			out = append(out, v)
		}
	}
	return out
}
