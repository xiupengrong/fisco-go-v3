package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	c := GetClient()
	if c == nil {
		fmt.Println("get client failed")
		return
	}
	defer c.Close()
	// cannot use big.NewInt to construct json request
	// TODO: analysis the ethereum's big.NewInt
	bn, err := c.GetBlockNumber(context.Background())
	if err != nil {
		fmt.Printf("block number not found: %v", err)
		return
	}
	bz, err := c.GetSealerList(context.Background())
	if err != nil {
		fmt.Printf("block number not found: %v", err)
		return
	}
	fmt.Printf("latest block number: \n%v", bn)
	fmt.Printf("bz: \n%v", bz)
	time.Sleep(10 * time.Second)
}
