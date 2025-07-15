package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/v3/client"
)

func GetClient() *client.Client {
	privateKey, err := hex.DecodeString("145e247e170ba3afd6ae97e88f00dbc976c2345d521b0f6713355d19d8b80b58")
	if err != nil {
		fmt.Printf("decode hex failed of %v", err)
		return nil
	}
	config := &client.Config{IsSMCrypto: false, GroupID: "group0", DisableSsl: false,
		PrivateKey: privateKey, Host: "36.139.226.161", Port: 20200,
		TLSCaFile:   "./ca1.crt",
		TLSKeyFile:  "./sdk1.key",
		TLSCertFile: "./sdk1.crt"}
	c, err := client.DialContext(context.Background(), config)
	if err != nil {
		fmt.Printf("Dial to %s:%d failed of %v", config.Host, config.Port, err)
		return nil
	}
	return c
}
