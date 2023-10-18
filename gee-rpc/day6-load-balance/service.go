package geerpc

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

type methodType struct {
	method        reflect.Method
	ArgType       reflect.Type
	ReplyType     reflect.Type
	numCalls      uint64
}

func (m *methodType) NumCalls() uint64  {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value

	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}

	return argv
}

func (m *methodType) newReplyv() reflect.Value {
	replyv := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}

	return replyv
}

type service struct {
	name    string
	typ     reflect.Type
	rcvr    reflect.Value
	method  map[string]*methodType
}

func newService(rcvr interface{}) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	// s.name = s.typ.Name()
	// log.Printf("name: %s", s.typ)
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server :%s is not a valid service name", s.name)
	}
	
	s.registerMethods()
	
	return s
}

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)

	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		// method := s.rcvr.Method(i)
		mtype := method.Type
		// log.Println(method)
		// log.Println(mtype)
		if mtype.NumIn() != 3 || mtype.NumOut() != 1 {
			continue
		}

		if mtype.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		argType, replyType := mtype.In(1), mtype.In(2)
		if !isExportedOrBuitinType(argType) || !isExportedOrBuitinType(replyType) {
			continue
		}

		// log.Printf("rpc server: register %s.%s", s.name, mtype.Name())

		s.method[method.Name] = &methodType{
			method:  method,
			ArgType: argType,
			ReplyType: replyType,
		}

		log.Printf("rpc server %p: register %s.%s", s, s.name, method.Name)
	}
}

func isExportedOrBuitinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)

	f := m.method.Func
	returnv := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInf := returnv[0].Interface(); errInf != nil {
		return errInf.(error)
	}

	return nil
}