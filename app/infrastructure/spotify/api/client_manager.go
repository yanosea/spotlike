package api

import (
	"context"
	"errors"
	"sync"

	"github.com/yanosea/spotlike/pkg/proxy"
)

var (
	// gcm is a global client manager.
	gcm ClientManager
	// gmutex is a global mutex.
	gmutex = &sync.Mutex{}
	// GetClientManagerFunc is a function to get the client manager.
	GetClientManagerFunc = getClientManager
)

// ClientManager is an interface that manages api clients.
type ClientManager interface {
	CloseClient() error
	GetClient() (Client, error)
	InitializeClient(ctx context.Context, config *ClientConfig) error
	IsClientInitialized() bool
}

// connectionManager is a struct that implements the ConnectionManager interface.
type clientManager struct {
	spotify proxy.Spotify
	client  Client
	http    proxy.Http
	randstr proxy.Randstr
	url     proxy.Url
	mutex   *sync.RWMutex
}

// NewClientManager initializes the client manager.
func NewClientManager(spotify proxy.Spotify, http proxy.Http, randstr proxy.Randstr, url proxy.Url) ClientManager {
	gmutex.Lock()
	defer gmutex.Unlock()

	if gcm == nil {
		gcm = &clientManager{
			spotify: spotify,
			client:  nil,
			http:    http,
			randstr: randstr,
			url:     url,
			mutex:   &sync.RWMutex{},
		}
	}

	return gcm
}

// GetClientManager gets the client manager.
func GetClientManager() ClientManager {
	return GetClientManagerFunc()
}

// getClientManager gets the client manager.
func getClientManager() ClientManager {
	gmutex.Lock()
	defer gmutex.Unlock()

	if gcm == nil {
		return nil
	}

	return gcm
}

// ResetClientManager resets the client manager.
func ResetClientManager() error {
	gmutex.Lock()
	defer gmutex.Unlock()

	if gcm == nil {
		return nil
	}

	if err := gcm.CloseClient(); err != nil {
		return err
	} else {
		gcm = nil
	}

	return nil
}

// IsClientInitialized checks if the client is initialized.
func (cm *clientManager) IsClientInitialized() bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	return cm.client != nil
}

// CloseClient closes the client.
func (cm *clientManager) CloseClient() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.client != nil {
		if err := cm.client.Close(); err != nil {
			return err
		}
		cm.client = nil
	}

	return nil
}

// GetClient gets the api client.
func (cm *clientManager) GetClient() (Client, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if cm.client == nil {
		return nil, errors.New("client not initialized")
	}

	return cm.client, nil
}

// InitializeClient initializes the api client.
func (cm *clientManager) InitializeClient(ctx context.Context, config *ClientConfig) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.client != nil {
		return errors.New("client already initialized")
	}

	cm.client = &client{
		spotify: cm.spotify,
		client:  nil,
		http:    cm.http,
		randstr: cm.randstr,
		url:     cm.url,
		config:  config,
		context: ctx,
		mutex:   &sync.RWMutex{},
	}

	return nil
}
