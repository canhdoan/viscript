package nodemanager

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skywire/src/apptracker"
	"github.com/skycoin/skywire/src/messages"
	"github.com/skycoin/skywire/src/node"
)

func TestDomainNameValidation(t *testing.T) {
	return
	messages.SetDebugLogLevel()

	config := &NodeManagerConfig{
		CtrlAddr: "127.0.0.1:5999",
	}

	config.Domain = "wrong_name_without_dot"
	_, e0 := newNodeManager(config)
	assert.NotNil(t, e0)

	config.Domain = "ok_with_dots.but#there@are-wrong,symbols!"
	_, e1 := newNodeManager(config)
	assert.NotNil(t, e1)

	config.Domain = "correct0.domain_name"
	nm, e2 := newNodeManager(config)
	assert.Nil(t, e2)

	defer nm.Shutdown()

	assert.Equal(t, nm.dnsServer.domain, "correct0.domain_name")
}

func TestMessagingServer(t *testing.T) {
	return
	messages.SetDebugLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	msgSrv := nm.nodeMsgServer
	assert.NotNil(t, msgSrv)

	host := net.ParseIP("127.0.0.1")
	port := 5999
	msgSrvAddr := net.UDPAddr{IP: host, Port: port}
	assert.Equal(t, msgSrvAddr.String(), msgSrv.conn.LocalAddr().String())
}

func TestRegisterNode(t *testing.T) {
	return

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	assert.Len(t, nm.nodeList, 0)

	n, err := node.CreateNode(&node.NodeConfig{"127.0.0.1:5992", []string{"127.0.0.1:5999"}, 4999, ""})
	assert.Nil(t, err)
	defer n.Shutdown()

	assert.Len(t, nm.nodeList, 1)
	assert.Equal(t, n.Id(), nm.nodeIdList[0])
}

func TestConnectNodes(t *testing.T) {
	return
	fmt.Println("")
	messages.SetDebugLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	n0, err := node.CreateNode(&node.NodeConfig{"127.0.0.1:5992", []string{"127.0.0.1:5999"}, 4990, ""})
	assert.Nil(t, err)
	defer n0.Shutdown()

	n1, err := node.CreateNode(&node.NodeConfig{"127.0.0.1:5993", []string{"127.0.0.1:5999"}, 4991, "node_one"})
	assert.Nil(t, err)
	defer n1.Shutdown()

	assert.Len(t, nm.nodeList, 2)

	err = n0.ConnectDirectly("node_one.demo.meshnet")
	assert.Nil(t, err)

	assert.True(t, n0.(*node.Node).ConnectedTo(n1.Id()))
	assert.True(t, n1.(*node.Node).ConnectedTo(n0.Id()))

	tf := nm.transportFactoryList[0]
	t0, t1 := tf.getTransports()
	assert.Equal(t, t0.id, t1.pair.id)
	assert.Equal(t, t1.id, t0.pair.id)
}

func TestNetwork(t *testing.T) {
	return
	fmt.Println("TestNetwork")
	messages.SetDebugLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	q := 20

	nodes := node.CreateNodeList(q, 14000)
	assert.Len(t, nodes, q, fmt.Sprintf("Should be %d nodes", q))
	assert.Len(t, nm.nodeIdList, q, fmt.Sprintf("Should be %d nodes", q))
	initRoute, err := nm.connectAllAndBuildRoute()
	assert.Nil(t, err)

	node0 := nodes[0].(*node.Node)

	inRouteMessage := messages.InRouteMessage{messages.NIL_TRANSPORT, initRoute, []byte{'t', 'e', 's', 't'}}
	node0.InjectTransportMessage(&inRouteMessage)
	time.Sleep(10 * time.Second)
	for i := 0; i < q-1; i++ {
		n0 := nodes[i]
		n1 := nodes[i+1]
		t0, err := n0.(*node.Node).GetTransportToNode(n1.Id())
		assert.Nil(t, err)
		t1, err := n1.(*node.Node).GetTransportToNode(n0.Id())
		assert.Nil(t, err)
		assert.Equal(t, uint32(1), t0.PacketsSent())
		assert.Equal(t, uint32(1), t0.PacketsConfirmed())
		assert.Equal(t, uint32(0), t1.PacketsSent())
		assert.Equal(t, uint32(0), t1.PacketsConfirmed())
	}

	node.ShutdownAll(nodes)

	fmt.Println("TestNetwork end")
}

func TestBuildRoute(t *testing.T) {
	return
	fmt.Println("TestBuildRoute")
	messages.SetInfoLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	n := 100
	m := 5

	allNodes := node.CreateNodeList(n, 15000)

	nodes := []cipher.PubKey{}

	for i := 0; i < m; i++ {
		nodenum := rand.Intn(n)
		node := allNodes[nodenum]
		nodes = append(nodes, node.Id())
	}

	for i := 0; i < m-1; i++ {
		_, err := nm.connectNodeToNode(nodes[i], nodes[i+1])
		assert.Nil(t, err)
	}

	routes, err := nm.buildRouteForward(nodes)
	assert.Nil(t, err)
	assert.Len(t, routes, m)

	node.ShutdownAll(allNodes)
	fmt.Println("TestBuildRoute end")
}

func TestFindRoute(t *testing.T) {
	return
	fmt.Println("TestFindRoute")
	messages.SetDebugLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	nodes := node.CreateNodeList(10, 16000)

	nodeList := []cipher.PubKey{}
	for _, n := range nodes {
		nodeList = append(nodeList, n.Id())
	}

	/*
		  1-2-3-4   long route
		 /	 \
		0---5-----9 short route, which should be selected
		 \ /     /
		  6_7_8_/   medium route
	*/
	nodes[0].ConnectDirectly(nodeList[1].Hex()) // making long route
	nodes[1].ConnectDirectly(nodeList[2].Hex())
	nodes[2].ConnectDirectly(nodeList[3].Hex())
	nodes[3].ConnectDirectly(nodeList[4].Hex())
	nodes[4].ConnectDirectly(nodeList[9].Hex())
	nodes[0].ConnectDirectly(nodeList[5].Hex()) // making short route
	nodes[5].ConnectDirectly(nodeList[9].Hex())
	nodes[0].ConnectDirectly(nodeList[6].Hex()) // make medium route, then findRoute should select the short one
	nodes[6].ConnectDirectly(nodeList[7].Hex())
	nodes[7].ConnectDirectly(nodeList[8].Hex())
	nodes[8].ConnectDirectly(nodeList[9].Hex())
	nodes[5].ConnectDirectly(nodeList[6].Hex())

	nm.rebuildRoutes()

	nodeFrom, nodeTo := nodeList[0], nodeList[9]
	routes, err := nm.findRouteForward(nodeFrom, nodeTo)
	assert.Nil(t, err)
	assert.Len(t, routes, 3, "Should be 3 routes")

	node.ShutdownAll(nodes)
	fmt.Println("TestFindRoute end")
}

func TestAddAndConnect2Nodes(t *testing.T) {
	return
	fmt.Println("TestAddAndConnect")
	messages.SetDebugLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	n0, err := node.CreateAndConnectNode(&node.NodeConfig{"127.0.0.1:5992", []string{"127.0.0.1:5999"}, 3990, ""})
	assert.Nil(t, err)
	defer n0.Shutdown()

	n1, err := node.CreateAndConnectNode(&node.NodeConfig{"127.0.0.1:5993", []string{"127.0.0.1:5999"}, 3991, ""})
	assert.Nil(t, err)
	defer n1.Shutdown()

	assert.Len(t, nm.nodeIdList, 2)
	assert.True(t, nm.connected(n0.Id(), n1.Id()))

	fmt.Println("TestAddAndConnect end")
}

func TestRandomNetwork100Nodes(t *testing.T) {
	return
	fmt.Println("TestRandomNetwork100Nodes")
	messages.SetInfoLogLevel()

	nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})
	defer nm.Shutdown()

	n := 100

	nodes := nm.CreateRandomNetwork(n, 17000)

	nodeIds := []cipher.PubKey{}

	for _, node := range nodes {
		nodeIds = append(nodeIds, node.Id())
	}

	assert.Len(t, nm.nodeIdList, n)
	assert.Equal(t, nm.nodeIdList, nodeIds)
	assert.True(t, nm.routeExists(nodeIds[0], nodeIds[n-1]))

	node.ShutdownAll(nodes)
	fmt.Println("TestRandomNetwork100Nodes end")
}

func TestSendThroughRandomNetworks(t *testing.T) {
	return
	fmt.Println("TestSendThroughRandomNetworks")
	messages.SetDebugLogLevel()

	lens := []int{2, 5, 10} // sizes of different networks which will be tested

	for _, n := range lens {

		nm, _ := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999"})

		nodes := nm.CreateRandomNetwork(n, 18000)

		n0 := nodes[0].(*node.Node)
		n1 := nodes[len(nodes)-1].(*node.Node)
		conn0, err := n0.Dial(n1.Id().Hex(), messages.AppId([]byte{}), messages.AppId([]byte{}))
		connId := conn0.Id()
		if err != nil {
			panic(err)
		}
		conn1 := n1.GetConnection(connId)
		assert.Equal(t, conn0.Status(), CONNECTED)
		assert.Equal(t, conn1.Status(), CONNECTED)
		fmt.Println(conn0.Id(), conn1.Id())
		msg := []byte{'t', 'e', 's', 't'}
		err = conn0.Send(msg)
		assert.Nil(t, err)
		time.Sleep(time.Duration(n) * time.Second)

		node.ShutdownAll(nodes)
		nm.Shutdown()
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
	fmt.Println("TestSendThroughRandomNetworks end")
}

func TestAppTrackerMsgServer(t *testing.T) {
	messages.SetDebugLogLevel()

	appTracker := apptracker.NewAppTracker("127.0.0.1:9000")
	defer appTracker.Shutdown()

	nm, err := newNodeManager(&NodeManagerConfig{Domain: "demo.meshnet", CtrlAddr: "127.0.0.1:5999", AppTrackerAddr: "127.0.0.1:9000"})
	assert.Nil(t, err)
	defer nm.Shutdown()

	assert.NotNil(t, nm.appTrackerMsgServer)
}
