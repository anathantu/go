package main

import(
	"github.com/hashicorp/consul/api"
	"github.com/dmitriyGarden/consul-leader-election"
	"fmt"
	"time"
)
type notify struct {
	T  string
}
func (n *notify)EventLeader(f bool)  {
	if f {
		fmt.Println(n.T, "I'm the leader!")
	} else {
		fmt.Println(n.T, "I'm no longer the leader!")
	}
}
func main(){
	conf := api.DefaultConfig()
	consul, _ := api.NewClient(conf)
	n := &notify{
		T: "My service",
	}

	elconf := &ElectionConfig{
		CheckTimeout: 5 * time.Second,
		Client: consul,
		Checks: []string{"healthID"},
		Key: "service/test-election/leader",
		LogLevel: election.LogDebug
		Event: n,
	}
	e := election.NewElection(elconf)
	// start election
	go  e.Init()
	if e.IsLeader() {
		fmt.Println("I'm a leader!")
	}
	// to do something ....
	if e.IsLeader() {
		fmt.Println("I'm a leader!")
	}
	// re-election
	e.ReElection()
	// to do something ....
	if e.IsLeader() {
		fmt.Println("I'm a leader!")
	}
	e.Stop()
}

