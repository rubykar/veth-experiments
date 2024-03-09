package network

import (
	"fmt"
	"os/exec"
)


func runCommand(command string , args ...string) error{
	cmd := exec.Command(command,args...)
	
	output , err := cmd.CombinedOutput()
	
	if err != nil{
		return fmt.Errorf("Error running command '%s %s' : %s",command,args,output)
	}
	return  nil
}


func SetupNetwork() error {
	if err := createNetworkNamespace("net1") ; err != nil {
		return err
	}
	
	if err := createNetworkNamespace("net2") ; err != nil {
		return err
	}
	
	if err := createVethPair("veth1","veth2"); err != nil {
		return err
	}
	
	if err := configureIP("net1","veth1","10.100.0.1/16"); err != nil{
		return err
	}
	
	if err := configureIP("net2","veth2","10.100.0.2/16"); err != nil {
		return err
	}
	if err := setLinkUp("net1","veth1") ; err != nil {
		return err
	}
		
	if err := setLinkUp("net2","veth2") ; err != nil {
		return err
	}
	
	if err := testConnectivity("net1","net2"); err != nil {
		return err;
	}
	
	
	return nil
	
}

func TeardownNetwork() error {
	if err := deleteNetworkNamespace("net1"); err != nil {
		return err
	}
	
	if err := deleteNetworkNamespace("net2"); err != nil {
		return err
	}
	
	return nil
}

func createNetworkNamespace(namespace string) error {
	return runCommand("ip","netns","add",namespace)
}

func deleteNetworkNamespace(namespace string) error {
	return runCommand("ip","netns","delete",namespace)
}

func createVethPair(name1 , name2 string) error {
	return runCommand("ip","link","add",name1,"type","veth","peer","name",name2)
}

func configureIP(namespace , interfaceName, ipAddress string) error {
	return runCommand("ip","netns","exec",namespace, "ip","addr" , "add" , ipAddress,"dev",interfaceName)
}

func setLinkUp(namespace , interfaceName string) error{
	return runCommand("ip","netns", "exec",namespace, "ip","link" , "set" , interfaceName,"up")
}

func testConnectivity(namespace1, namespace2 string) error {
	return runCommand("ip","netns","exec",namespace1, "ping","-c","3","10.100.0.2")
}


