package t_topology

import (
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_servicepoints"
	. "github.com/saichler/layer8/go/overlay/vnet"
	. "github.com/saichler/types/go/common"
	"sync"
)

type TestTopology struct {
	vnets    map[string]*VNet
	vnics    map[string]IVirtualNetworkInterface
	handlers map[string]*TestServicePointHandler
	mtx      *sync.RWMutex
}

func NewTestTopology(vnicCountPervNet int, vnetPorts ...int) *TestTopology {
	this := &TestTopology{}
	this.vnets = make(map[string]*VNet)
	this.vnics = make(map[string]IVirtualNetworkInterface)
	this.handlers = make(map[string]*TestServicePointHandler)
	this.mtx = &sync.RWMutex{}

	for _, vNetPort := range vnetPorts {
		_vnet := createVnet(vNetPort)
		this.vnets[_vnet.Resources().Config().LocalAlias] = _vnet
	}
	Sleep()
	for i := 0; i < vnicCountPervNet; i++ {
		for _, vNetPort := range vnetPorts {
			if i < vnicCountPervNet-1 {
				_vnic, handler := createVnic(vNetPort, i+1, 0)
				this.vnics[_vnic.Resources().Config().LocalAlias] = _vnic
				this.handlers[_vnic.Resources().Config().LocalAlias] = handler
			} else {
				_vnic, _ := createVnic(vNetPort, i+1, -1)
				this.vnics[_vnic.Resources().Config().LocalAlias] = _vnic
			}
		}
	}
	Sleep()
	list := make([]*VNet, 0)
	for _, _vnet := range this.vnets {
		list = append(list, _vnet)
	}
	for i := 0; i < len(list)-1; i++ {
		connectVnets(list[i], list[i+1])
		Sleep()
	}
	return this
}

func (this *TestTopology) Vnet(vnetPort int) *VNet {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, -1)
	return this.vnets[alias]
}

func (this *TestTopology) Vnic(vnetPort, vnicNum int) IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) Handler(vnetPort, vnicNum int) *TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) Shutdown() {
	for _, _vnic := range this.vnics {
		_vnic.Shutdown()
	}
	for _, _vnet := range this.vnets {
		_vnet.Shutdown()
	}
}
