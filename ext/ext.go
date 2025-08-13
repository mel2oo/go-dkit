package ext

import (
	"strings"
	"sync"
)

const (
	KEY_TID           = "tid"
	KEY_SID           = "sid"
	KEY_QID           = "qid"
	KEY_ORG           = "org"
	KEY_UID           = "uid"
	KEY_AID           = "aid"
	KEY_XUID          = "x-uid"
	KEY_XORG          = "x-org"
	KEY_XORGRiskLevel = "x-org-risklevel"
	KEY_ORDER         = "order"
	KEY_Principle     = "principle"
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
	id := e.GetValue(KEY_XUID)
	if len(id) == 0 {
		id = e.GetValue(KEY_XORG)
	}
	return id
}
