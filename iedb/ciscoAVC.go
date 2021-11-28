package iedb

import (
	"bytes"
	"encoding/binary"
)

// RFC 6313 BasicList

// CiscoAppVarStringElement : Handle Cisco App
//       This type is used to handle IE: 9:12235
//		 based on different config,this IE could be observed
//		 multiple times in single record:
//
//		 Multiple collect type share this IE with different SubAppID
//
//		 http: url,useragent,refer
//		 dns:  domain name
//		 ssl:  common name
type CiscoAppVarStringKey struct {
	AppID     uint32
	SubAppTag uint32
}

// ParseCiscoAppVarString : Parse 9:12235 Field, return structure
func ParseCiscoAppVarString(b []byte) map[CiscoAppVarStringKey]string {
	result := make(map[CiscoAppVarStringKey]string)
	var k CiscoAppVarStringKey
	k.AppID = binary.BigEndian.Uint32(b[0:4])
	k.SubAppTag = uint32(binary.BigEndian.Uint16(b[4:6]))
	result[k] = string(bytes.TrimRight(b[6:], "\x00"))
	return result
}

// CiscoURLHitElement :  URL Hit Field
//This field use '/0' as delimeter
type CiscoURLHitItem struct {
	Name   string
	Number uint16
}

func ParseCiscoURLHitString(b []byte) []CiscoURLHitItem {
	var result []CiscoURLHitItem
	t := bytes.Split(b, make([]byte, 1, 1))
	//fmt.Println("DEBUG:---->", len(t))
	//fmt.Printf("DEBUG:::%+v\n", t)
	length := len(t)
	for idx := 0; idx < length-1; idx = idx + 2 {
		var uh CiscoURLHitItem
		uh.Name = string(t[idx])
		//fmt.Println("DEBUG::::%s\t\tLen:", uh.Name, len(t[idx+1]))
		if len(t[idx+1]) == 0 {
			tmp := make([]byte, 2)
			tmp[1] = t[idx+2][0]
			uh.Number = binary.BigEndian.Uint16(tmp)
			idx = idx + 1
		} else {
			uh.Number = binary.BigEndian.Uint16(t[idx+1])
		}
		result = append(result, uh)
		//fmt.Printf("DEBUG::::%s\t%d\n", uh.Name, uh.Number)
	}

	return result
}
