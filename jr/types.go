package jr

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

type User struct {
	Guid      string  `json:"guid"`
	IsActive  bool    `json:"isActive"`
	Balance   float64 `json:"balance"`
	Age       int     `json:"age"`
	EyeColor  string  `json:"eyeColor"`
	Name      string  `json:"name"`
	Company   string  `json:"company"`
	Email     string  `json:"email"`
	Alt_email string  `json:"alt_email"`
	About     string  `json:"about"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
