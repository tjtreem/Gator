package main

import (
	"fmt"
	"github.com/tjtreem/gator/internal/config"

)


func main() {
    
    cfg, err := config.Read()
    if err != nil {
	fmt.Println("Error:", err)
	return
    }

    err = cfg.SetUser("tjtreem")
    if err != nil {
	fmt.Println("Error:", err)
	return
    }

    cfg, err = config.Read()
    if err != nil {
	fmt.Println("Error:", err)
	return
    }

    fmt.Printf("%+v\n",cfg)
}
