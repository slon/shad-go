package iprange_test

import (
	"fmt"
	"log"

	"gitlab.com/slon/shad-go/iprange"
)

func ExampleParseList() {
	list, err := iprange.ParseList("10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24")
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range list {
		fmt.Println(i)
	}

	// Output:
	// {10.0.0.1 10.0.0.1}
	// {10.0.0.5 10.0.0.10}
	// {192.168.1.0 192.168.1.255}
	// {192.168.10.0 192.168.10.255}
}
