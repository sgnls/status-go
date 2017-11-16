package api

import (
	"net"

	"github.com/status-im/status-go/geth/api"
	"github.com/status-im/status-go/geth/params"
)

// adminService exposes functions for administrative tasks.
type adminService struct{}

// newAdminService creates an instance of the administrative
// service to expose.
func newAdminService() *adminService {
	return &adminService{}
}

// GetAddresses returns the IP address of the client.
// nolint: unparam
func (svc *adminService) GetAddresses(args *NoArgs, reply *StringsReply) error {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	var addrs []string
	for _, ifcAddr := range ifcAddrs {
		addrs = append(addrs, ifcAddr.String())
	}
	reply.Strings = addrs
	return nil
}

// statusService exposes the functions of the StatusAPI via JSON-RPC.
type statusService struct {
	statusAPI *api.StatusAPI
}

// newStatusService creates an instance of the Status service to expose.
// If backend is nil a new StatusAPI will be created, elso one using
// the backend.
func newStatusService(backend *api.StatusBackend) *statusService {
	var statusAPI *api.StatusAPI
	if backend == nil {
		statusAPI = api.NewStatusAPI()
	} else {
		statusAPI = api.NewStatusAPIWithBackend(backend)
	}
	return &statusService{
		statusAPI: statusAPI,
	}
}

// StartNode loads the configuration out of the passed string and
// starts a node with it.
// nolint: unparam
func (svc *statusService) StartNode(args *ConfigArgs, reply *NoReply) error {
	config, err := params.LoadNodeConfig(args.Config)
	if err != nil {
		return err
	}

	return svc.statusAPI.StartNode(config)
}

// StopNode starts the stopped node.
// nolint: unparam
func (svc *statusService) StopNode(args *NoArgs, reply *NoReply) error {
	return svc.statusAPI.StopNode()
}

// CreateAccount creates an internal geth account.
// nolint: unparam
func (svc *statusService) CreateAccount(args *AccountArgs, reply *AccountReply) error {
	address, publicKey, mnemonic, err := svc.statusAPI.CreateAccount(args.Password)
	if err != nil {
		return err
	}
	reply.Address = address
	reply.PublicKey = publicKey
	reply.Mnemonic = mnemonic
	return nil
}

// SelectAccount selects the addressed account.
// nolint: unparam
func (svc *statusService) SelectAccount(args *AccountArgs, reply *NoReply) error {
	return svc.statusAPI.SelectAccount(args.Address, args.Password)
}

// Logout clears the Whisper identities.
// nolint: unparam
func (svc *statusService) Logout(args *NoArgs, reply *NoReply) error {
	return svc.statusAPI.Logout()
}
