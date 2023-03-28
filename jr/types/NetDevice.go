package types

type NetDevice struct {
	Vlan            string `json:"VLAN"`
	IPv4SrcAddress  string `json:"IPV4_SRC_ADDR"`
	IPv4DstAddress  string `json:"IPV4_DST_ADDR"`
	InBytes         int    `json:"IN_BYTES"`
	FirstSwitched   int64  `json:"FIRST_SWITCHED"`
	LastSwitched    int64  `json:"LAST_SWITCHED"`
	L4SrcPort       int64  `json:"L4_SRC_PORT"`
	L4DstPort       int64  `json:"L4_DST_PORT"`
	Protocol        int    `json:"PROTOCOL"`
	SrcTos          int    `json:"SRC_TOS"`
	SrcAs           int    `json:"SRC_AS"`
	DstAs           int    `json:"DST_AS"`
	L7Proto         int    `json:"L7_PROTO"`
	L7ProtoName     string `json:"L7_PROTO_NAME"`
	L7ProtoCategory string `json:"L7_PROTO_CATEGORY"`
}
