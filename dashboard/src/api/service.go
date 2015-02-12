package api

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Service struct {
	ForServe Model

	prefixTri map[string]prefixMethodEntry
}

type prefixMethodEntry struct {
	method     reflect.Value
	httpMethod string
}

type Model interface {
	Access(appKey string, token string) bool
}

var _ = http.Get
var _ = log.Println
var _ = (http.Handler)(&Service{})

const (
	kErrBadMethod      = "This method not support"
	kErrPermDenied     = "Permission denied"
	kErrMethodNotFound = "Method not found"
	kErrInternal       = "Server internal error"

	kPrefix = "Prefix"
)

// Prefix-Get|Post-Path-To
func NewService(model Model) *Service {
	self := &Service{
		ForServe:  model,
		prefixTri: make(map[string]prefixMethodEntry),
	}

	return self.registerAll()
}

func (self *Service) registerAll() *Service {
	stub := reflect.ValueOf(self.ForServe)
	stubTy := stub.Type()

	for i := 0; i < stubTy.NumMethod(); i++ {
		method := stubTy.Method(i)

		var entry prefixMethodEntry
		if !strings.HasPrefix(method.Name, kPrefix) {
			continue
		}

		partial := method.Name[len(kPrefix):]

		if strings.HasPrefix(partial, "Get") {
			entry.httpMethod = "Get"
		} else if strings.HasPrefix(partial, "Post") {
			entry.httpMethod = "Post"
		} else {
			log.Fatal("bad method.name", method.Name)
		}

		partial = titleToPath(partial[len(entry.httpMethod):])

		entry.method = stub.MethodByName(method.Name)

		self.prefixTri[partial] = entry
	}

	return self
}

func titleToPath(title string) string {
	path := ""
	var i, p int
	var ch rune

	for i, ch = range title {
		if i == 0 {
			continue
		}
		switch ch {
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			path += "/" + strings.ToLower(title[p:i])
			p = i
		}
	}
	path += "/" + strings.ToLower(title[p:])

	return path
}

func (self *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	values := url.Query()
	if len(values) > 0 {
		appKey := ""
		if value, found := values["appkey"]; found {
			appKey = value[0]
		}

		token := ""
		if value, found := values["token"]; found {
			token = value[0]
		}
		if !self.ForServe.Access(appKey, token) {
			http.Error(w, kErrPermDenied, 503)
			return
		}
	}

	for prefix, entry := range self.prefixTri {
		if strings.HasPrefix(url.Path, prefix) {
			dispatchPrefixMethod(w, r, url.Path[len(prefix):], entry.method)
			return
		}
	}

	methodName := ""
	if r.Method == "GET" {
		methodName = "Get"
	} else if r.Method == "POST" {
		methodName = "Post"
	} else {
		http.Error(w, kErrBadMethod, 503)
		return
	}

	sgements := strings.Split(url.Path, "/")
	for _, sgement := range sgements {
		methodName += strings.Title(sgement)
	}
	//log.Println("request function: ", methodName)

	stub := reflect.ValueOf(self.ForServe)
	method := stub.MethodByName(methodName)
	if !method.IsValid() {
		http.Error(w, kErrInternal, 503)

		log.Printf("method: %s not found\n", methodName)
		return
	}

	dispatchMethod(w, r, method)

}

func dispatchPrefixMethod(w http.ResponseWriter, r *http.Request, path string, method reflect.Value) {
	methodTy := method.Type()

	if methodTy.NumIn() != 3 || methodTy.NumOut() != 1 {

		log.Fatal("bad prefix method prototype.")
		return
	}

	inputs := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r), reflect.ValueOf(path)}
	rv := method.Call(inputs)

	if len(rv) >= 1 {
		if !rv[0].IsNil() && rv[0].Kind() == reflect.Interface {
			if err, ok := rv[0].Interface().(error); ok {
				http.Error(w, kErrInternal, 503)

				log.Printf("call fail(%v): %v", method.Type(), err)
				return
			}
		}
	}
}

func dispatchMethod(w http.ResponseWriter, r *http.Request, method reflect.Value) {
	methodTy := method.Type()

	if r.Method == "GET" {
		if methodTy.NumIn() != 1 || methodTy.NumOut() != 1 {
			http.Error(w, kErrInternal, 503)

			//log.Printf("method: %s not found\n", methodName)
			return
		}

		responsePtrTy := methodTy.In(0)
		responseStub := reflect.New(responsePtrTy.Elem())

		inputs := []reflect.Value{responseStub}
		call(w, method, inputs[:])
	} else if r.Method == "POST" {
		if methodTy.NumIn() != 2 || methodTy.NumOut() != 1 {
			http.Error(w, kErrInternal, 503)

			//log.Printf("method: %s not found\n", methodName)
			return
		}

		responsePtrTy := methodTy.In(0)
		requestPtrTy := methodTy.In(1)

		responseStub := reflect.New(responsePtrTy.Elem())
		requestStub := reflect.New(requestPtrTy.Elem())

		decoder := json.NewDecoder(r.Body)
		decoder.Decode(requestStub.Interface())

		inputs := []reflect.Value{responseStub, requestStub}
		call(w, method, inputs[:])
	}
}

func call(w http.ResponseWriter, method reflect.Value, inputs []reflect.Value) {
	rv := method.Call(inputs)
	if len(rv) >= 1 {
		if !rv[0].IsNil() && rv[0].Kind() == reflect.Interface {
			if err, ok := rv[0].Interface().(error); ok {
				http.Error(w, kErrInternal, 503)

				log.Printf("call fail(%v): %v", method.Type(), err)
				return
			}
		}
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(inputs[0].Interface()); err != nil {
		log.Printf("write back fail: %v", err)
	}
}

func EncodeJson(w http.ResponseWriter, input interface{}) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(input); err != nil {
		log.Printf("write back fail: %v", err)
	}
}
