package geecache

import (
	// "bytes"
	"fmt"
	"geecache/consistenthash"
	pb "geecache/geecachepb"
	"google.golang.org/protobuf/proto"
	"io"
	// "io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type httpGetter struct {
	baseURL   string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v/%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)

	res,err := http.Get(u)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("reading resonse body: %v", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)

	err = proto.Unmarshal(bytes, out)
	
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}	

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)

const (
	defaultBasePath = "/_geecache"
	defaultReplicas = 50
)

type HTTPPool struct {
	self       string
	basePath   string
	mu         sync.Mutex
	peers      *consistenthash.Map
	httpGetter map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		basePath: defaultBasePath,

	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}  


func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)
    // /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath) + 1:], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return 
	}

	groupName := parts[0]
	key :=   parts[1]
	// log.Println(groupName)
	// log.Println(key)

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetter = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetter[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetter[peer], true
	}

	// p.Log("Pick self %s")

	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)