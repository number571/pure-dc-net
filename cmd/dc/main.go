package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/number571/pure-dc-net/pkg/dc"
	"github.com/number571/pure-dc-net/pkg/syncmap"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	myAddr := loadDCAddr(dcAddrFile)
	keys := loadDCKeys(dcKeysFile)
	dcNet := dc.NewDCNet(
		loadDCIter(dcIterFile),
		keysToGenerators(keys)...,
	)

	result := syncmap.NewSyncMap()
	bqueue := make(chan byte, 512)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				btext := []byte(text)
				for _, b := range btext {
					bqueue <- b
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				b := byte(0)
				select {
				case x := <-bqueue:
					b = x
				default:
				}
				doDCRequest(dcNet, keys, myAddr, b)
				for {
					if result.Size() == len(keys) {
						break
					}
					time.Sleep(time.Second)
				}
				if r := result.Sum(); r != 0 {
					fmt.Print(string(r))
				}
				result.Clear()
			}
		}
	}()

	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req := &dcRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := validateMAC(keys, req.ReqToken); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if req.Iteration != dcNet.Iteration() {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if result.Size() == len(keys) {
			w.WriteHeader(http.StatusLengthRequired)
			return
		}

		result.Store(req.ReqToken.Addr, req.Generated)
	})

	if err := http.ListenAndServe(myAddr, nil); err != nil {
		log.Println(err)
	}
}

func doDCRequest(dcNet dc.IDCNet, keys map[string][]byte, myAddr string, b byte) {
	var (
		generated = b ^ dcNet.Generate()
		iteration = dcNet.Iteration()
	)
	dcRequest := &dcRequest{
		Iteration: iteration,
		Generated: generated,
	}
	for addr := range keys {
		dcRequest.ReqToken = generateToken(keys, myAddr, addr)
		jsonRequest, _ := json.Marshal(dcRequest)
		req, _ := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("http://%s/push", addr),
			bytes.NewBuffer(jsonRequest),
		)
		for {
			if err := doHTTPRequest(req); err != nil {
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
	storeDCIter(dcIterFile, iteration)
}

func doHTTPRequest(req *http.Request) error {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	rsp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}
	return nil
}
