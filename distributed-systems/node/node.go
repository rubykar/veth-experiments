package node

import (
	"fmt"
	"net/http"
)

type Node struct {
	Port int
	Stop chan struct{}
}

func (n *Node) StartServer() {
	http.HandleFunc(fmt.Sprintf("/node%d", n.Port), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Node on Port %d!", n.Port)
	})

	addr := fmt.Sprintf(":%d", n.Port)
	fmt.Printf("Node Listening on Port %d....\n", n.Port)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()
	<-n.Stop 
}

func (n *Node) StopServer() {
	close(n.Stop) 
}

