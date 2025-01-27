package awgconfig

import (
	"3x-api/internal"
	"3x-api/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"slices"
	"strconv"
)

const (
	messageInitiationSize = 148
	messageResponseSize   = 92
	interfaceConfig       = "sadf/wg0.conf"
	jsonConfig            = "conf.json"
	permissions           = 0o600
	client                = `[Interface]
Address = %s
DNS = %s
PrivateKey = %s

%s

[Peer]
PublicKey = %s
PresharedKey = %s
AllowedIPs = %s
Endpoint = %s
PersistentKeepalive = %s
`
	addClient = `[Peer]
PublicKey = %s
PresharedKey = %s
AllowedIPs = %s
`
)

// Some obfuscating values, part of main config and every client config
type parameters struct {
	Jc, Jmin, Jmax, S1, S2, H1, H2, H3, H4 int
}

// Using reflection to get names of the fields
func (p *parameters) toByte() []byte {
	t := reflect.TypeOf(*p)
	v := reflect.ValueOf(*p)
	b := &bytes.Buffer{}
	for i := range t.NumField() {
		fmt.Fprintf(b, "%s = %s\n", t.Field(i).Name, strconv.Itoa(int(v.Field(i).Int())))
	}
	return b.Bytes()
}

// Header of the main config
type wgConfig struct {
	PrivateKey, PublicKey, Address, ListenPort string
}

func (c *wgConfig) toByte() []byte {
	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)
	b := &bytes.Buffer{}
	for i := range t.NumField() {
		name := t.Field(i).Name
		// Can use struct tags
		if name == "PublicKey" || name == "PresharedKey" {
			continue
		}
		fmt.Fprintf(b, "%s = %s\n", t.Field(i).Name, v.Field(i).String())
	}
	return b.Bytes()
}

// This is the main config
// TODO: Add fields persistent across clients: AllowedIPs, Endpoint, DNS
// PersistentKeepalive
type baseConf struct {
	wgConfig
	parameters
}

func (c *baseConf) toByte() []byte {
	return slices.Concat(c.wgConfig.toByte(), c.parameters.toByte())
}

type configuration interface {
	toByte() []byte
}

// Client config
// TODO: Add fields like name etc
type clientConf struct {
	PrivateKey, PublicKey, PresharedKey, AllowedIPs, Endpoint,
	PersistentKeepalive, DNS string
}

func generateBaseConfig(cfg *config.Config) error {
	slog.Debug("Generating values for config")
	configu := new(baseConf)

	configu.Jmin = cfg.WgConfig.JcMin
	configu.Jmax = cfg.WgConfig.JcMax
	privKey, err := internal.GeneratePrivateKey()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	configu.PrivateKey = privKey.String()
	configu.PublicKey = privKey.PublicKey().String()
	// TODO: Validate IP and subnet mask
	// TODO: Generate subnet IP and subnet mask instead of using env
	configu.Address = cfg.WgConfig.SubNetIP + "/" + cfg.WgConfig.SubNetMask
	configu.ListenPort = cfg.Port
	configu.Jc = internal.RandIntBound(cfg.WgConfig.JPMin, cfg.WgConfig.JPMax)
	s1 := internal.RandIntBound(15, 150)
	s2 := internal.RandIntBound(15, 150)
	for s1+messageInitiationSize == s2+messageResponseSize {
		s2 = internal.RandIntBound(15, 150)
	}
	configu.S1 = s1
	configu.S2 = s2
	configu.H1 = internal.RandIntBound(5, math.MaxInt32)
	configu.H2 = internal.RandIntBound(5, math.MaxInt32)
	configu.H3 = internal.RandIntBound(5, math.MaxInt32)
	configu.H4 = internal.RandIntBound(5, math.MaxInt32)

	f1, _ := json.Marshal(configu)

	return internal.WriteFiles(internal.File{interfaceConfig, configu.toByte(), permissions},
		internal.File{jsonConfig, f1, permissions})
}

func Boot(cfg *config.Config) error {
	return generateBaseConfig(cfg)
}
