package provider

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethNode "github.com/ethereum/go-ethereum/node"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/status-im/status-go/geth/account"
	"github.com/status-im/status-go/geth/jail"
	"github.com/status-im/status-go/geth/node"
	"github.com/status-im/status-go/geth/transactions"
)

// ServiceProvider provides access to status and geth services
type ServiceProvider struct {
	fcmServerKey string

	nodeManager    *node.NodeManager
	accountManager *account.Manager
	jailManager    *jail.Jail
	txQueueManager *transactions.Manager

	account *accounts.Manager
	whisper *whisper.Whisper
}

// New builds a Serviceprovider based on a NodeManager and a fcmServerKey
func New(nodeManager *node.NodeManager) *ServiceProvider {
	return &ServiceProvider{
		nodeManager: nodeManager,
	}
}

// SetFCMServerKey sets the fcmServerKey for needed for notification manager
func (p *ServiceProvider) SetFCMServerKey(fcmServerKey string) {
	p.fcmServerKey = fcmServerKey
}

// Node gets the underlying geth Node of the given nodeManager
func (p *ServiceProvider) Node() (*gethNode.Node, error) {
	return p.nodeManager.Node()
}

// Account get the underlying accounts.Manager under account.Manager
func (p *ServiceProvider) Account() (*accounts.Manager, error) {
	if p.account != nil {
		return p.account, nil
	}

	node, err := p.nodeManager.Node()
	if err != nil {
		return nil, err
	}
	p.account = node.AccountManager()

	return p.account, nil
}

// AccountKeyStore exposes reference to accounts key store
func (p *ServiceProvider) AccountKeyStore() (*keystore.KeyStore, error) {
	am, err := p.Account()
	if err != nil {
		return nil, account.ErrInvalidAccountManager
	}

	backends := am.Backends(keystore.KeyStoreType)
	if len(backends) == 0 {
		return nil, account.ErrAccountKeyStoreMissing
	}

	keyStore, ok := backends[0].(*keystore.KeyStore)
	if !ok {
		return nil, account.ErrAccountKeyStoreMissing
	}

	return keyStore, nil
}

// Whisper gets a WhisperService
func (p *ServiceProvider) Whisper() (*whisper.Whisper, error) {
	var err error

	if p.whisper != nil {
		return p.whisper, nil
	}
	p.whisper, err = p.nodeManager.WhisperService()

	return p.whisper, err
}

// NodeManager get the related NodeManager
func (p *ServiceProvider) NodeManager() *node.NodeManager {
	return p.nodeManager
}

// AccountManager  get the AccountManager
func (p *ServiceProvider) AccountManager() (*account.Manager, error) {
	if p.accountManager != nil {
		return p.accountManager, nil
	}

	p.accountManager = account.NewManager(p)

	return p.accountManager, nil
}

// JailManager get the jail manager
func (p *ServiceProvider) JailManager() *jail.Jail {
	if p.jailManager != nil {
		return p.jailManager
	}

	return jail.New(p.nodeManager)
}

// TxQueueManager get transaction manager
func (p *ServiceProvider) TxQueueManager() *transactions.Manager {
	if p.txQueueManager != nil {
		return p.txQueueManager
	}
	am, err := p.AccountManager()
	if err != nil {
		return nil
	}
	return transactions.NewManager(p.nodeManager, am)
}

// Reset resets all managers
func (p *ServiceProvider) Reset() {
	p.accountManager = nil
	p.jailManager = nil
	p.txQueueManager = nil
	p.account = nil
	p.whisper = nil
}
