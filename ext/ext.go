package ext

import (
	"context"
	"strings"
	"sync"
)

type ctxKey string

func WithContextValue(ctx context.Context, v ExtType) context.Context {
	return context.WithValue(ctx, ctxKey("ext"), v)
}

func FromContextValue(ctx context.Context) ExtType {
	val, ok := ctx.Value(ctxKey("ext")).(ExtType)
	if ok {
		return val
	}
	return ExtType{}
}

const (
	KeyTID           = "tid"
	KeySID           = "sid"
	KeyQID           = "qid"
	KeyORG           = "org"
	KeyUID           = "uid"
	KeyAID           = "aid"
	KeyXUID          = "x-uid"
	KeyXORG          = "x-org"
	KeyXORGRiskLevel = "x-org-risklevel"
	KeyOrder         = "order"
	KeyPrinciple     = "principle"
)

type ExtType struct {
	mux *sync.RWMutex
	kvs map[string]string
}

func New(ext string) ExtType {
	e := ExtType{
		mux: &sync.RWMutex{},
		kvs: map[string]string{},
	}
	for _, v := range strings.Split(ext, ";") {
		kv := strings.Split(v, ":")
		if len(kv) != 2 {
			continue
		}
		e.kvs[kv[0]] = kv[1]
	}
	return e
}

func (e *ExtType) ToString() string {
	e.mux.RLock()
	defer e.mux.RUnlock()

	ext := ""
	for k, v := range e.kvs {
		ext += k + ":" + v + ";"
	}
	return ext
}

func (e *ExtType) GetValue(key string) string {
	e.mux.RLock()
	defer e.mux.RUnlock()

	return e.kvs[key]
}

func (e *ExtType) SetValue(key, value string) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.kvs[key] = value
}

func (e *ExtType) GetIdentityID() string {
	id := e.GetValue(KeyXUID)
	if len(id) == 0 {
		id = e.GetValue(KeyXORG)
	}
	return id
}
