package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/pure-dc-net/internal/nodes"
	"github.com/number571/pure-dc-net/internal/token"
	"github.com/number571/pure-dc-net/pkg/dc"
)

func DoRequest(ctx context.Context, nodes nodes.Nodes, tokenData []byte) {
	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))

	for _, node := range nodes {
		node := node
		go func() {
			defer wg.Done()

			token := token.GenerateToken([]byte(node.Pasw), tokenData)
			_ = doExternalRequest(ctx, node.Addr, token)
		}()
	}

	wg.Wait()
}

func doExternalRequest(ctx context.Context, addr string, token *token.Token) error {
	jsonRequest, _ := json.Marshal(token)
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s/dc", addr),
		bytes.NewBuffer(jsonRequest),
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := doHTTPRequest(req); err == nil {
				return nil
			}
			time.Sleep(time.Second)
		}
	}
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

func NewDCExternalServer(addr string, nodes nodes.Nodes, dcNet dc.IDCNet, totalizer dc.ITotalizer) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/dc", handleExternal(nodes, dcNet, totalizer))
	return &http.Server{Handler: mux, Addr: addr}
}

func handleExternal(nodes nodes.Nodes, dcNet dc.IDCNet, totalizer dc.ITotalizer) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		reqToken := &token.Token{}
		if err := json.NewDecoder(r.Body).Decode(reqToken); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenData, err := token.UnmarshalTokenData(reqToken.Data)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		node, ok := nodes[tokenData.Name]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := token.ValidateMAC([]byte(node.Pasw), reqToken); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if tokenData.Iter != dcNet.Iteration() {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if totalizer.Size() == len(nodes) {
			w.WriteHeader(http.StatusLengthRequired)
			return
		}

		totalizer.Store(tokenData.Byte)
	}
}
