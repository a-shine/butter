package butter

import (
	"fmt"
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/persist"
	"github.com/a-shine/butter/retrieve"
	"github.com/a-shine/butter/tracker"
	"github.com/a-shine/butter/traverse"
	"os"
	"os/signal"
	"syscall"
)

// Spawn node into the network (the node serves as an entry-point to the butter network). You can also do this manually
// to have more control over the specific protocols used in your dapp. This function presents a simple abstraction with
// the included default butter protocols.
func Spawn(overlay node.Overlay, public bool, track bool) {
	n := overlay.Node()
	setupLeaveHandler(n)
	go discover.Discover(overlay)
	if track {
		go tracker.Track(overlay)
	}
	if public {
		go traverse.Traverse(n)
	}
	n.Start(overlay)
}

func SpawnDefaultOverlay(node *node.Node, public bool, track bool) {
	overlay := persist.NewOverlay(node) // Creates a new overlay network
	retrieve.AppendRetrieveBehaviour(overlay.Node())
	Spawn(&overlay, public, track)
}

// setupLeaveHandler creates a listener on a new goroutine which will notify the program if it receives an interrupt
// from the OS and then handles the node leaving the network gracefully.
func setupLeaveHandler(node *node.Node) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rLeaving the butter network...")
		node.Shutdown()
		os.Exit(0)
	}()
}

// SpawnAmbassador node which is a special community node with added ambassadorial behaviours that help it bridge
// connections across subnetworks. To be an ambassador a node inherently needs to be available publicly (must port
// forward either manually or via UPNP and have a public IP address). The added ambassadorial behaviours allows the node
// to share the public addresses of other traversed (i.e. public) nodes between each other.
//func SpawnAmbassador(node *node.Node) {
//	go discover.Discover(node)
//	go traverse.Traverse(node)
//	//go traverse.AppendAmbassadorBehaviour(node) // the node keeps track of ambassador so if someone needs an ambassador they can find them dynamically (improvement on bootstrapping)
//	node.Start()
//}
