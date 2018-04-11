package node_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/stackint"
)

const (
	NumberOfBootstrapNodes = 5
	NumberOfOrders         = 10
)

var primeVal, _ = stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
var Prime = &primeVal
var trader, _ = identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
var epochDNR dnr.DarkNodeRegistry

var dnrOuterLock = new(sync.Mutex)
var dnrInnerLock = new(sync.Mutex)

type OrderBook struct {
	LastUpdateID int             `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}

// HeapInt creates a stackint on the heap - temporary convenience method
func heapInt(n uint) *stackint.Int1024 {
	tmp := stackint.FromUint(n)
	return &tmp
}

var _ = Describe("Dark nodes", func() {

	var err error
	epochDNR, err = dnr.TestnetDNR(nil)
	if err != nil {
		panic(err)
	}

	var mu = new(sync.Mutex)

	BeforeEach(func() {
		mu.Lock()
	})

	AfterEach(func() {
		mu.Unlock()
	})

	// Bootstrapping
	for _, numberOfNodes := range []int{5} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("when bootstrapping %d nodes", numberOfNodes), func() {

				var err error
				var nodes []*node.DarkNode

				BeforeEach(func() {
					By("generate nodes")
					nodes, err = generateNodes(numberOfNodes)
					Ω(err).ShouldNot(HaveOccurred())

					By("start node services")
					startNodeServices(nodes)
				})

				It("should reach a fault tolerant level of connectivity", func() {
					By("bootstrap nodes")
					bootstrapNodes(nodes)
					n := 0
					for i := range nodes {
						numberOfPeers := len(nodes[i].DHT.MultiAddresses())
						if numberOfPeers > numberOfNodes*2/3 {
							n++
						}
					}
					Ω(n).Should(BeNumerically(">", numberOfNodes*2/3))
				})

				AfterEach(func() {
					By("stop node services")
					stopNodes(nodes)
				})
			})
		}(numberOfNodes)
	}

	// Connectivity
	for _, numberOfNodes := range []int{5} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("when connecting %d nodes", numberOfNodes), func() {
				for _, connectivity := range []int{20, 40, 60, 80, 100} {
					func(connectivity int) {
						Context(fmt.Sprintf("with %d%% connectivity", connectivity), func() {

							var err error
							var nodes []*node.DarkNode

							BeforeEach(func() {
								By("generate nodes")
								nodes, err = generateNodes(numberOfNodes)
								Ω(err).ShouldNot(HaveOccurred())

								By("start node service")
								startNodeServices(nodes)

								By("bootstrap nodes")
								bootstrapNodes(nodes)
							})

							It("should succeed for the super majority", func() {
								By("ping connections")
								numberOfPings, numberOfErrors := connectNodes(nodes, connectivity)
								if (numberOfPings / 3) == 0 {
									Ω(numberOfErrors).Should(Equal(0))
								} else {
									Ω(numberOfErrors).Should(BeNumerically("<", numberOfPings/3))
								}
							})

							AfterEach(func() {
								stopNodes(nodes)
							})
						})
					}(connectivity)
				}
			})
		}(numberOfNodes)
	}

	// Order matching
	for _, numberOfNodes := range []int{5} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("when sending orders to %d nodes", numberOfNodes), func() {

				var err error
				var nodes []*node.DarkNode

				BeforeEach(func() {
					By("generate nodes")
					nodes, err = generateNodes(numberOfNodes)
					Ω(err).ShouldNot(HaveOccurred())
					err = registerNodes(nodes)
					Ω(err).ShouldNot(HaveOccurred())

					By("start node service")
					startNodeServices(nodes)

					By("start node background workers")
					startNodeBackgroundWorkers(nodes)

					By("bootstrap nodes")
					bootstrapNodes(nodes)

					By("watching for changes to the dark ocean")
					watchDarkOcean(nodes)
				})

				It("should succeed for the super majority", func() {
					By("verify configuration")
					for _, node := range nodes {
						// A node does not include itself in its pool and so
						// we account for this by including an extra +1
						k := node.DarkPool.Size()*2/3 + 2
						Ω(k).Should(Equal(len(nodes)*2/3 + 1))
					}

					By("send orders")
					err := sendOrders(nodes)
					Ω(err).ShouldNot(HaveOccurred())

					By("verify order matches")
					timer := time.NewTimer(time.Minute * time.Duration(len(nodes)))
					for _, node := range nodes {
						n := 0
						for i := 0; i < NumberOfOrders; i++ {
							select {
							case <-node.DeltaNotifications:
								n++
							case <-timer.C:
								i = NumberOfOrders
							}
						}
						Ω(n).Should(Equal(NumberOfOrders))
					}
				})

				AfterEach(func() {
					err := deregisterNodes(nodes)
					Ω(err).ShouldNot(HaveOccurred())
					stopNodes(nodes)
				})
			})
		}(numberOfNodes)
	}

	// Synchronization
	for _, numberOfNodes := range []int{32} {
		func(numberOfNodes int) {
			FContext(fmt.Sprintf("synchronizing with %d nodes", numberOfNodes), func() {

				var err error
				var nodes []*node.DarkNode

				BeforeEach(func() {
					By("generate nodes")
					nodes, err = generateNodes(numberOfNodes)
					Ω(err).ShouldNot(HaveOccurred())

					By("start node services")
					startNodeServices(nodes)
				})

				It("should reach a fault tolerant level of connectivity", func() {
					start := time.Now()
					By("bootstrap nodes")
					bootstrapNodes(nodes)
					log.Println("bootstrapping takes ", time.Since(start))

					By("send orders")
					err := sendOrders(nodes)
					Ω(err).ShouldNot(HaveOccurred())

					go func() {
						defer GinkgoRecover()
						for i := 0; i < 5; i++ {
							time.Sleep(10 * time.Second)
							log.Println("sending a new patch of orders")
							err := sendOrders(nodes)
							Ω(err).ShouldNot(HaveOccurred())
						}
					}()

					By("synchronization")
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
					defer cancel()

					syncBlocks, errs := nodes[0].ClientPool.Sync(ctx, nodes[1].NetworkOptions.MultiAddress)
					continuing := true
					for continuing {
						select {
						case err := <-errs:
							if err != nil {
								log.Println("cannot sync", err)
							}
							continuing = false
						case block, ok := <-syncBlocks:
							if !ok {
								continuing = false
								break
							}
							// Handle received blocks
							var ord order.Order
							switch block.OrderBlock.(type) {
							case *rpc.SyncBlock_Open:
								ord = rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Open).Open)
							case *rpc.SyncBlock_Confirmed:
								ord = rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Confirmed).Confirmed)
							case *rpc.SyncBlock_Unconfirmed:
								ord = rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Unconfirmed).Unconfirmed)
							case *rpc.SyncBlock_Canceled:
								ord = rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Canceled).Canceled)
							case *rpc.SyncBlock_Settled:
								ord = rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Settled).Settled)
							default:
								log.Printf("unknown order status, %t", block.OrderBlock)
							}

							if ord.Parity == order.ParityBuy {
								log.Println("buy  order received from synchronization, orderID : ", ord.ID.String())
							} else {
								log.Println("sell order received from synchronization, orderID : ", ord.ID.String())
							}

						}
					}
				})

				AfterEach(func() {
					By("stop node services")
					stopNodes(nodes)
				})
			})
		}(numberOfNodes)
	}
})

func generateNodes(numberOfNodes int) ([]*node.DarkNode, error) {
	// Generate nodes from the config files
	nodes := make([]*node.DarkNode, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		var err error
		var config *node.Config
		if i < NumberOfBootstrapNodes {
			config, err = node.LoadConfig(fmt.Sprintf("../test/config/bootstrap-node-%d.json", i+1))
		} else {
			config, err = node.LoadConfig(fmt.Sprintf("../test/config/node-%d.json", i-NumberOfBootstrapNodes+1))
		}
		if err != nil {
			return nil, err
		}
		auth := bind.NewKeyedTransactor(config.EthereumKey.PrivateKey)
		dnr, err := dnr.TestnetDNR(auth)
		if err != nil {
			return nil, err
		}
		node, err := node.NewDarkNode(*config, dnr)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

func registerNodes(nodes []*node.DarkNode) error {
	dnrOuterLock.Lock()
	dnrInnerLock.Lock()
	defer dnrInnerLock.Unlock()
	for _, node := range nodes {
		isRegistered, err := node.DarkNodeRegistry.IsRegistered(nodes[0].NetworkOptions.MultiAddress.ID())
		if isRegistered {
			return errors.New("already registered")
		}

		bond := stackint.FromUint(10)
		err = node.DarkNodeRegistry.ApproveRen(&bond)
		if err != nil {
			return err
		}

		node.DarkNodeRegistry.SetGasLimit(300000)
		_, err = node.DarkNodeRegistry.Register(node.ID, []byte{}, &bond)
		node.DarkNodeRegistry.SetGasLimit(0)
		if err != nil {
			return err
		}
	}
	_, err := epochDNR.WaitForEpoch()
	return err
}

func deregisterNodes(nodes []*node.DarkNode) error {
	defer dnrOuterLock.Unlock()
	dnrInnerLock.Lock()
	defer dnrInnerLock.Unlock()
	for _, node := range nodes {
		node.DarkNodeRegistry.SetGasLimit(300000)
		_, err := node.DarkNodeRegistry.Deregister(node.ID)
		node.DarkNodeRegistry.SetGasLimit(0)
		if err != nil {
			panic(err)
		}
	}
	epochDNR.SetGasLimit(300000)
	_, err := epochDNR.WaitForEpoch()
	epochDNR.SetGasLimit(0)
	if err != nil {
		panic(err)
	}
	for _, node := range nodes {
		node.DarkNodeRegistry.SetGasLimit(300000)
		_, err := node.DarkNodeRegistry.Refund(node.ID)
		node.DarkNodeRegistry.SetGasLimit(0)
		if err != nil {
			panic(err)
		}
	}
	_, err = epochDNR.WaitForEpoch()
	return err
}

func startNodeServices(nodes []*node.DarkNode) {
	for i := range nodes {
		go func(i int) {
			defer GinkgoRecover()
			nodes[i].StartServices()
		}(i)
	}
	time.Sleep(time.Millisecond * time.Duration(10*len(nodes)))
}

func startNodeBackgroundWorkers(nodes []*node.DarkNode) {
	for i := range nodes {
		nodes[i].StartBackgroundWorkers()
	}
	time.Sleep(time.Millisecond * time.Duration(10*len(nodes)))
}

func bootstrapNodes(nodes []*node.DarkNode) {
	do.CoForAll(nodes, func(i int) {
		nodes[i].Bootstrap()
	})
	do.CoForAll(nodes, func(i int) {
		nodes[i].Bootstrap()
	})
}

func watchDarkOcean(nodes []*node.DarkNode) {
	for i := range nodes {
		go func(i int) {
			defer GinkgoRecover()
			nodes[i].WatchDarkOcean()
		}(i)
	}

	_, err := epochDNR.WaitForEpoch()
	if err != nil {
		panic(err)
	}

	// time.Sleep(time.Minute)
	time.Sleep(2 * time.Second)
}

func stopNodes(nodes []*node.DarkNode) {
	for i := range nodes {
		nodes[i].Stop()
	}
}

func connectNodes(nodes []*node.DarkNode, connectivity int) (int, int) {
	numberOfPings := 0
	numberOfErrors := 0
	do.CoForAll(nodes, func(i int) {
		// Select nodes randomly
		from := nodes[i]
		isSelected := rand.Intn(100) < connectivity
		if !isSelected {
			return
		}
		for j, to := range nodes {
			if i == j {
				continue
			}
			// Connect nodes randomly
			isConnected := rand.Intn(100) < connectivity
			if isConnected {
				numberOfPings++
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := from.ClientPool.Ping(ctx, to.NetworkOptions.MultiAddress); err != nil {
					log.Printf("error pinging: %v", err)
					numberOfErrors++
				}
			}
		}
	})
	return numberOfPings, numberOfErrors
}

func sendOrders(nodes []*node.DarkNode) error {
	// Get order data from Binance
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%d", NumberOfOrders))
	if err != nil {
		return err
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	orderBook := new(OrderBook)
	if err := json.Unmarshal(response, orderBook); err != nil {
		log.Println(response)
		return err
	}

	// Generate order from the Binance data
	buyOrders := make([]*order.Order, len(orderBook.Asks))
	sellOrders := make([]*order.Order, len(orderBook.Asks))

	for i, j := range orderBook.Asks {
		price, err := strconv.ParseFloat(j[0].(string), 10)
		if err != nil {
			return errors.New("fail to parse the price into a float")
		}
		price = price * 1000000000000

		amount, err := strconv.ParseFloat(j[1].(string), 10)
		if err != nil {
			return errors.New("fail to parse the amount into a float")
		}
		amount = amount * 1000000000000
		sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
			order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(uint(price)), heapInt(uint(amount)),
			heapInt(uint(amount)), heapInt(1))
		sellOrders[i] = sellOrder

		buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour),
			order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(uint(price)), heapInt(uint(amount)),
			heapInt(uint(amount)), heapInt(1))
		buyOrders[i] = buyOrder
	}

	// Send order fragment to the nodes
	totalNodes := len(nodes)
	pool := rpc.NewClientPool(trader).WithTimeout(10 * time.Second).WithTimeoutBackoff(5 * time.Second)
	for i := range buyOrders {
		buyOrder, sellOrder := buyOrders[i], sellOrders[i]
		log.Printf("Sending matched order. [BUY] %s <---> [SELL] %s", buyOrder.ID, sellOrder.ID)
		buyShares, err := buyOrder.Split(int64(totalNodes), int64(totalNodes*2/3+1), Prime)
		if err != nil {
			return err
		}
		sellShares, err := sellOrder.Split(int64(totalNodes), int64(totalNodes*2/3+1), Prime)
		if err != nil {
			return err
		}

		do.CoForAll(buyShares, func(j int) {
			orderRequest := &rpc.OpenOrderRequest{
				From: &rpc.MultiAddress{
					Signature:    []byte{},
					MultiAddress: nodes[0].NetworkOptions.MultiAddress.String(),
				},
				OrderFragment: rpc.MarshalOrderFragment(buyShares[j]),
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := pool.OpenOrder(ctx, nodes[j].NetworkOptions.MultiAddress, orderRequest)
			if err != nil {
				log.Printf("Coudln't send order fragment to %s\n", nodes[j].NetworkOptions.MultiAddress.ID())
				log.Fatal(err)
			}
		})

		do.CoForAll(sellShares, func(j int) {
			orderRequest := &rpc.OpenOrderRequest{
				From: &rpc.MultiAddress{
					Signature:    []byte{},
					MultiAddress: nodes[0].NetworkOptions.MultiAddress.String(),
				},
				OrderFragment: rpc.MarshalOrderFragment(sellShares[j]),
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := pool.OpenOrder(ctx, nodes[j].NetworkOptions.MultiAddress, orderRequest)
			if err != nil {
				log.Printf("Coudln't send order fragment to %s\n", nodes[j].NetworkOptions.MultiAddress.ID())
				log.Fatal(err)
			}
		})
	}

	time.Sleep(time.Second)
	return nil
}
