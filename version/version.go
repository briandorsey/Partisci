package version

import (
	"encoding/json"
	"strings"
	"time"
)

type Version struct {
	Id         string `json:"id"`
	App        string `json:"app"`
	Version    string `json:"version,omitempty"`
	Host       string `json:"host,omitempty"`
	Instance   uint16 `json:"instance,omitempty"`
	HostIP     string `json:"host_ip,omitempty"`
	LastUpdate int64  `json:"last_update,omitempty"`
}

func safeRunes(r rune) rune {
	if '0' <= r && r <= '9' {
		return r
	}
	if 'a' <= r && r <= 'z' {
		return r
	}
	return '_'
}

func AppToID(app string) (id string) {
	id = strings.ToLower(app)
	id = strings.Map(safeRunes, id)
	return
}

func ParsePacket(host string, b []byte) (v Version, err error) {
	v.HostIP = host
	v.LastUpdate = time.Now().Unix()
	err = json.Unmarshal(b[:len(b)], &v)
	if err != nil {
		return v, err
	}
	v.Id = AppToID(v.App)
	return
}
