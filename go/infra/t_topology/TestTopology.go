package t_topology

import (
	"fmt"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_servicepoints"
	"github.com/saichler/layer8/go/overlay/health"
	. "github.com/saichler/layer8/go/overlay/vnet"
	. "github.com/saichler/layer8/go/overlay/vnic"
	. "github.com/saichler/types/go/common"
	"strconv"
	"sync"
)

type TestTopology struct {
	vnets       map[string]*VNet
	vnetsOrder  []*VNet
	vnics       map[string]IVirtualNetworkInterface
	handlers    map[string]*TestServicePointHandler
	trHandlers  map[string]*TestServicePointTransactionHandler
	repHandlers map[string]*TestServicePointReplicationHandler
	mtx         *sync.RWMutex
}

func NewTestTopology(vnicCountPervNet int, vnetPorts []int, level LogLevel) *TestTopology {
	this := &TestTopology{}
	this.vnets = make(map[string]*VNet)
	this.vnics = make(map[string]IVirtualNetworkInterface)
	this.handlers = make(map[string]*TestServicePointHandler)
	this.trHandlers = make(map[string]*TestServicePointTransactionHandler)
	this.repHandlers = make(map[string]*TestServicePointReplicationHandler)
	this.vnetsOrder = make([]*VNet, 0)
	this.mtx = &sync.RWMutex{}

	for _, vNetPort := range vnetPorts {
		_vnet := createVnet(vNetPort, level)
		this.vnets[_vnet.Resources().SysConfig().LocalAlias] = _vnet
		this.vnetsOrder = append(this.vnetsOrder, _vnet)
	}

	Sleep()
	Sleep()

	for _, vNetPort := range vnetPorts {
		for i := 0; i < vnicCountPervNet; i++ {
			if i == vnicCountPervNet-1 {
				_vnic, _, _, _ := createVnic(vNetPort, i+1, -1, level)
				this.vnics[_vnic.Resources().SysConfig().LocalAlias] = _vnic
			} else {
				_vnic, handler, trHandler, repHandler := createVnic(vNetPort, i+1, 0, level)
				this.vnics[_vnic.Resources().SysConfig().LocalAlias] = _vnic
				this.handlers[_vnic.Resources().SysConfig().LocalAlias] = handler
				this.trHandlers[_vnic.Resources().SysConfig().LocalAlias] = trHandler
				this.repHandlers[_vnic.Resources().SysConfig().LocalAlias] = repHandler
			}
		}
	}

	if !WaitForCondition(this.areVnetsReady1, 2, nil, "Vnet are not ready 1") {
		panic("Vnet are not ready 1")
	}

	for i := 0; i < len(this.vnetsOrder)-1; i++ {
		for j := i + 1; j < len(this.vnetsOrder); j++ {
			connectVnets(this.vnetsOrder[i], this.vnetsOrder[j])
			Sleep()
			Sleep()
		}
	}

	if !WaitForCondition(this.areVnetsReady2, 3, nil, "Vnet are not ready 2") {
		for _, vnet := range this.vnets {
			hc := health.Health(vnet.Resources())
			fmt.Println(vnet.Resources().SysConfig().LocalAlias, " ", vnet.ExternalCount(), vnet.LocalCount(), len(hc.All()))
		}
		panic("Vnet are not ready 2")
	}

	if !WaitForCondition(this.areVnicReady, 2, nil, "Vnics are not ready!") {
		nic := this.VnicByVnetNum(1, 1)
		hc := health.Health(nic.Resources())
		all := hc.All()
		for _, hp := range all {
			fmt.Println("Vnic 1 -> ", hp.Alias)
		}
		panic("Vnics are not ready, it has only " + strconv.Itoa(len(all)) + " instead of 15")
	}
	return this
}

func (this *TestTopology) areVnicReady() bool {
	for vnetNum := 1; vnetNum <= 3; vnetNum++ {
		for vnicNum := 1; vnicNum <= 4; vnicNum++ {
			nic := this.VnicByVnetNum(vnetNum, vnicNum)
			hc := health.Health(nic.Resources())
			hp := hc.All()
			if len(hp) != 15 {
				return false
			}
		}
	}
	return true
}

func (this *TestTopology) areVnetsReady1() bool {
	for _, vnet := range this.vnets {
		hc := health.Health(vnet.Resources())
		all := hc.All()
		if len(all) != 5 {
			return false
		}
	}
	return true
}

func (this *TestTopology) areVnetsReady2() bool {
	for _, vnet := range this.vnets {
		if vnet.ExternalCount() != 2 {
			return false
		}
		if vnet.LocalCount() != 4 {
			return false
		}
		hc := health.Health(vnet.Resources())
		if len(hc.All()) != 15 {
			return false
		}
	}
	return true
}

func (this *TestTopology) Vnet(vnetPort int) *VNet {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, -1)
	return this.vnets[alias]
}

func (this *TestTopology) VnicByPort(vnetPort, vnicNum int) IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) VnicByVnetNum(vnetNum, vnicNum int) IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) HandlerByPort(vnetPort, vnicNum int) *TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) HandlerByVnetNum(vnetNum, vnicNum int) *TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) TrHandlerByVnetNum(vnetNum, vnicNum int) *TestServicePointTransactionHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.trHandlers[alias]
}

func (this *TestTopology) RepHandlerByVnetNum(vnetNum, vnicNum int) *TestServicePointReplicationHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.repHandlers[alias]
}

func (this *TestTopology) Shutdown() {
	for _, _vnic := range this.vnics {
		_vnic.Shutdown()
	}
	for _, _vnet := range this.vnets {
		_vnet.Shutdown()
	}
}

func (this *TestTopology) ResetHandlers() {
	for _, _handler := range this.handlers {
		_handler.Reset()
	}
	for _, _handler := range this.trHandlers {
		_handler.Reset()
	}
	for _, _handler := range this.repHandlers {
		_handler.Reset()
	}
}

func (this *TestTopology) AllHandlers() []*TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServicePointHandler, 0)
	for _, h := range this.handlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllTrHandlers() []*TestServicePointTransactionHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServicePointTransactionHandler, 0)
	for _, h := range this.trHandlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllRepHandlers() []*TestServicePointReplicationHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServicePointReplicationHandler, 0)
	for _, h := range this.repHandlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllVnics() []IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]IVirtualNetworkInterface, 0)
	for _, vnic := range this.vnics {
		result = append(result, vnic)
	}
	return result
}

func (this *TestTopology) RenewVnic(alias string) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	nic, ok := this.vnics[alias]
	if ok {
		nic.Shutdown()
		delete(this.vnics, alias)
		r := nic.Resources()
		r.SysConfig().LocalUuid = ""
		r.SysConfig().RemoteUuid = ""
		nic = NewVirtualNetworkInterface(nic.Resources(), nil)
		nic.Start()
		this.vnics[alias] = nic
	} else {
		Log.Error("Unable to find vnic ", alias)
	}
}

func (this *TestTopology) SetLogLevel(lvl LogLevel) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for _, net := range this.vnets {
		net.Resources().Logger().SetLogLevel(lvl)
	}
	for _, nic := range this.vnics {
		nic.Resources().Logger().SetLogLevel(lvl)
	}
}

func (this *TestTopology) ReActivateTestService(nic IVirtualNetworkInterface) {
	h, err := nic.Resources().ServicePoints().Activate(ServicePointType, ServiceName, 0, nic.Resources(), nil,
		nic.Resources().SysConfig().LocalAlias)
	if err != nil {
		panic(err)
	}
	handler := h.(*TestServicePointHandler)
	this.handlers[nic.Resources().SysConfig().LocalAlias] = handler
}
