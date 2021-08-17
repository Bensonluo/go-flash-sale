package common

import (
	"net/http"
	"strings"
)

// FilterHandle 声明数据类型
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	filterMap map[string]FilterHandle
}

//构造函数
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// WebHandle 声明数据类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) {
				//执行拦截
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		webHandle(rw, r)
	}
}
