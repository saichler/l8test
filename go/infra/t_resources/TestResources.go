package t_resources

import (
	"bytes"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/servicepoints/go/points/service_points"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"strconv"
)

const (
	VNET_PREFIX = "vnet-"
	VNIC_PREFIX = "-vnic-"
)

var Log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func CreateResources(vnetPort, vnicNum int) (common.IResources, string) {
	alias := AliasOf(vnetPort, vnicNum)
	_registry := registry.NewRegistry()
	_security, err := common.LoadSecurityProvider("security.so")
	if err != nil {
		panic("Failed to load security provider")
	}
	_config := &types.VNicConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  alias,
		VnetPort:    uint32(vnetPort)}
	_introspector := introspecting.NewIntrospect(_registry)
	_servicepoints := service_points.NewServicePoints(_introspector, _config)
	_resources := resources.NewResources(_registry, _security, _servicepoints, Log, nil, nil, _config, _introspector)
	return _resources, alias
}

func AliasOf(vnetPort, vnicNum int) string {
	alias := bytes.Buffer{}
	alias.WriteString(VNET_PREFIX)
	alias.WriteString(strconv.Itoa(vnetPort))
	if vnicNum != -1 {
		alias.WriteString(VNIC_PREFIX)
		alias.WriteString(strconv.Itoa(vnicNum))
	}
	return alias.String()
}
