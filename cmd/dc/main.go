package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/number571/pure-dc-net/internal/nodes"
	"github.com/number571/pure-dc-net/internal/service"
	"github.com/number571/pure-dc-net/internal/token"
	"github.com/number571/pure-dc-net/pkg/dc"
)

var (
	servicePath = os.Getenv("SERVICE_PATH")
	consumeAddr = os.Getenv("CONSUME_ADDR")
	produceAddr = os.Getenv("PRODUCE_ADDR")
)

var (
	currIter = loadDCIter()
	nodeName = loadDCName()
	nodesMap = loadDCNodesMap()
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var (
		dcNet  = dc.NewDCNet(currIter, nodes.NodesKeysToGenerators(nodesMap)...)
		ttlzr  = dc.NewTotalizer()
		bqueue = make(chan byte, 512)
	)

	log.Println("service is listening...")

	go runGenerator(ctx, dcNet, ttlzr, bqueue)

	go func() {
		server := service.NewDCProducerServer(produceAddr, bqueue)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	server := service.NewDCConsumerServer(consumeAddr, nodesMap, dcNet, ttlzr)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func runGenerator(ctx context.Context, dcNet dc.IDCNet, ttlzr dc.ITotalizer, bqueue chan byte) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b := dcNet.Generate()
			select {
			case x := <-bqueue:
				b ^= x // send message
			default:
			}
			doRequest(ctx, dcNet, nodesMap, nodeName, b)
			storeDCIter(dcNet.Iteration())
			for {
				if ttlzr.Size() == len(nodesMap) {
					break
				}
				time.Sleep(time.Second)
			}
			ttlzr.Store(b)
			if r := ttlzr.Sum(); r != 0 {
				log.Println(string(r))
			}
		}
	}
}

func doRequest(ctx context.Context, dcNet dc.IDCNet, nodes nodes.Nodes, name string, b byte) {
	tokenData := token.MarshalTokenData(&token.TokenData{
		Name: name,
		Iter: dcNet.Iteration(),
		Byte: b,
	})

	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))

	for _, node := range nodes {
		node := node
		go func() {
			defer wg.Done()

			token := token.GenerateToken(node.SKey, tokenData)
			_ = service.DoConsumeRequest(ctx, node.Addr, token)
		}()
	}

	wg.Wait()
}
