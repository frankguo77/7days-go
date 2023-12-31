package server

import (
	"encoding/json"
	"errors"
	"geerpc"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)



type Server struct{
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}

		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {DefaultServer.Accept(lis)}
func Register(rcvr interface{}) error {return DefaultServer.Register(rcvr)}

func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}

	return nil
}



func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {conn.Close()}()

	var opt geerpc.Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error:", err)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return		
	}

	server.serveCodec(f(conn))
}

var invalidRequest = struct{}{}

func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)

	for {
		req, err := server.readRequest(cc)
		if err != nil {
			log.Println("Server: serveCodec", err)
			if req == nil {
				break
			}

		    req.h.Error = err.Error()
		    server.sendResponse(cc, req.h, invalidRequest, sending)
		    continue
		}

		wg.Add(1)
		go server.hanleRequest(cc, req, sending, wg)
	}

	wg.Wait()
	log.Println("Server: goint to close codec")
	cc.Close()
}

type request struct {
	h            *codec.Header
	argv,replyv  reflect.Value
	mtype        *methodType
	svc          *service
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h,err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
    req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()

	if req.argv.Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}

	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read argv err:", err)
	}

	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()

	if err := cc.Write(h, body); err != nil {
		log.Panicln("rpc server: write response error:", err)
	}
}

func (server *Server) hanleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	err := req.svc.call(req.mtype,req.argv, req.replyv)
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, req.h, invalidRequest, sending)
		return
	}

	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

func (server *Server) findService(serviceMethod string)(svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
	}

	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot + 1:]
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service" + serviceName)
	}

	svc = svci.(*service)

	mtype, ok = svc.method[methodName]
	if !ok {
		err = errors.New("rpc server: can't find method " + methodName)
	}

	return
}