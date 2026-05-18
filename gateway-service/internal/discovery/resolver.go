package discovery

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"gateway-service/internal/model"

	frameworkconsul "github.com/calmlax/aevons-framework/core/consul"
)

type Resolver struct {
	registry *frameworkconsul.Registry
	rr       sync.Map
	random   *rand.Rand
	mu       sync.Mutex
}

func NewResolver(registry *frameworkconsul.Registry) *Resolver {
	return &Resolver{
		registry: registry,
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Resolver) Resolve(rule *model.ServiceRule) (*frameworkconsul.Instance, error) {
	if rule == nil {
		return nil, errors.New("service rule is nil")
	}
	if rule.Discovery != "consul" {
		return nil, errors.New("unsupported discovery type: " + rule.Discovery)
	}
	if r.registry == nil {
		return nil, errors.New("consul registry is not initialized")
	}

	instances, err := r.registry.Discover(rule.Name)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, errors.New("no healthy instances for service " + rule.Name)
	}

	index := 0
	switch rule.LoadBalance {
	case "random":
		r.mu.Lock()
		index = r.random.Intn(len(instances))
		r.mu.Unlock()
	default:
		counterAny, _ := r.rr.LoadOrStore(rule.Name, &atomic.Uint64{})
		counter := counterAny.(*atomic.Uint64)
		index = int(counter.Add(1)-1) % len(instances)
	}

	instance := instances[index]
	return &instance, nil
}
