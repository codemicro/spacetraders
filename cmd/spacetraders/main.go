package main

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
)

func main() {
	user, err := stapi.GetUserInfo("akpa")
	fmt.Println(err)
	fmt.Printf("%#v\n", user)
	fmt.Printf("%#v\n", user.Ships[0])
	fmt.Printf("%#v %s\n", user.Loans[0], user.Loans[0].Due)

	locations, err := stapi.GetSystemLocations("OE")
	fmt.Println(err)
	for _, x := range locations {
		fmt.Printf("%#v\n", x)
	}

}
