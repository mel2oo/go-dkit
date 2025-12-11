package ext

import (
	"context"
	"net/http"
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
	return New("")
}

func InjectHeader(ctx context.Context, header http.Header) {
	extv := FromContextValue(ctx)
	exts := extv.ToString()
	if exts != "" {
		header.Set("EXT", exts)
	}

	if val := extv.GetValue(KeySID); len(val) > 0 {
		header.Add("SID", val)
	}

	if val := extv.GetValue(KeyTID); len(val) > 0 {
		header.Add("X-TID", val)
	}

	if val := extv.GetValue(KeyXUID); len(val) > 0 {
		header.Add("X-UID", val)
	}

	if val := extv.GetValue(KeyXORG); len(val) > 0 {
		header.Add("X-ORG", val)
	}

	if val := extv.GetValue(KeyXORGRiskLevel); len(val) > 0 {
		header.Add("X-ORG-RISKLEVEL", val)
	}

	if val := extv.GetValue(KeyOrder); len(val) > 0 {
		header.Add("ORDER", val)
	}

	if val := extv.GetValue(KeyPrinciple); len(val) > 0 {
		header.Add("PLATFORMPRINCIPLE", val)
	}
}

func ExtractHeader(ctx context.Context, header http.Header) context.Context {
	exts := header.Get("EXT")
	extv := New(exts)

	sid := header.Get("SID")
	if len(sid) > 0 {
		extv.SetValue(KeySID, sid)
	}
	xtid := header.Get("X-TID")
	if len(xtid) > 0 {
		extv.SetValue(KeyTID, xtid)
	}
	xuid := header.Get("X-UID")
	if len(xuid) > 0 {
		extv.SetValue(KeyXUID, xuid)
	}
	xorg := header.Get("X-ORG")
	if len(xorg) > 0 {
		extv.SetValue(KeyXORG, xorg)
	}
	xorgRiskLevel := header.Get("X-ORG-RISKLEVEL")
	if len(xorgRiskLevel) > 0 {
		extv.SetValue(KeyXORGRiskLevel, xorgRiskLevel)
	}
	order := header.Get("ORDER")
	if len(order) > 0 {
		extv.SetValue(KeyOrder, order)
	}
	platformPrinciple := header.Get("PLATFORMPRINCIPLE")
	if len(platformPrinciple) > 0 {
		extv.SetValue(KeyPrinciple, platformPrinciple)
	}

	return WithContextValue(ctx, extv)
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
		e.kvs[strings.ToLower(kv[0])] = kv[1]
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
