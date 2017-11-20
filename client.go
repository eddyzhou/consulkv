package consulkv

import (
	"errors"
	"time"

	"github.com/hashicorp/consul/api"
)

var (
	ErrNotExist = errors.New("consul-kv: key not exist")
)

type ConsulClient struct {
	client *api.KV
}

func NewClient(nodes []string, scheme string) (*ConsulClient, error) {
	conf := api.DefaultConfig()

	conf.Scheme = scheme
	if len(nodes) > 0 {
		conf.Address = nodes[0]
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &ConsulClient{client.KV()}, nil
}

func (c *ConsulClient) Get(key string, fn ValueMapper) (*ConfKV, error) {
	v, _, err := c.client.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNotExist
	}

	ret := &ConfKV{
		kv:          v,
		valueMapper: fn,
	}
	return ret, nil
}

func (c *ConsulClient) WatchPrefix(prefix string, waitIndex uint64) (uint64, error) {
	opts := api.QueryOptions{
		WaitIndex: waitIndex,
	}
	_, meta, err := c.client.List(prefix, &opts)
	if err != nil {
		return waitIndex, err
	}

	return meta.LastIndex, nil
}

type Watcher struct {
	client        *ConsulClient
	prefixIndices map[string]uint64
	onUpdate      func(prefix string)
}

func NewWatcher(cli *ConsulClient, onUpdate func(prefix string)) *Watcher {
	return &Watcher{
		client:        cli,
		prefixIndices: make(map[string]uint64),
		onUpdate:      onUpdate,
	}
}

func (p *Watcher) Process() {
	for prefix, idx := range p.prefixIndices {
		go p.monitorPrefix(prefix, idx)
	}
}

func (p *Watcher) monitorPrefix(prefix string, lastIndex uint64) {
	for {
		index, err := p.client.WatchPrefix(prefix, lastIndex)
		if err != nil {
			// TODO: add log
			time.Sleep(time.Second * 2)
			continue
		}
		p.prefixIndices[prefix] = index
		if index > lastIndex {
			p.onUpdate(prefix)
		}
	}
}
