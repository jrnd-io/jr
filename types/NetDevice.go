//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

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
