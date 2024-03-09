package main

import(
	"fmt"
	"os"
	"github.com/rubykar/virtual-ethernet/network"
)


func main(){
	action := os.Args[1]
	switch action {
		case "up" : 
			err := network.SetupNetwork()
			if err != nil {
				fmt.Println("Error setting up network" , err)
				os.Exit(1)
			}
		case "down" : 
			err := network.TeardownNetwork()
			if err != nil {
				fmt.Println("Error tearing down network" , err)
				os.Exit(1)
			}
			
		default : 
			fmt.Println("Invalid action . Use 'up' or 'down' . ")
			os.Exit(1)
	}
	
	fmt.Println("Network operation completed successfully . ")
}
