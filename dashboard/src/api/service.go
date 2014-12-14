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
)

func NewService(model Model) *Service {
	return &Service{ForServe: model}
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
