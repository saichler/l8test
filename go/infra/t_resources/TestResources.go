package t_resources

import (
	"bytes"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/servicepoints/go/points/service_points"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/testtypes"
	"github.com/saichler/types/go/types"
	"strconv"
	"time"
)

const (
	VNET_PREFIX = "vnet-"
	VNIC_PREFIX = "-vnic-"
)

var Log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func CreateResources(vnetPort, vnicNum int) (common.IResources, string) {
	alias := AliasOf(vnetPort, vnicNum)
	_log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	_registry := registry.NewRegistry()
	_security, err := common.LoadSecurityProvider("security.so", "../../../../")
	if err != nil {
		panic("Failed to load security provider " + err.Error())
	}
	_config := &types.SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  alias,
		VnetPort:    uint32(vnetPort)}
	_introspector := introspecting.NewIntrospect(_registry)
	_servicepoints := service_points.NewServicePoints(_introspector, _config)
	_resources := resources.NewResources(_registry, _security, _servicepoints, _log, nil, nil, _config, _introspector)
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

func CreateTestModelInstance(index int) *testtypes.TestProto {
	tag := strconv.Itoa(index)
	sub := &testtypes.TestProtoSub{
		MyString: "string-sub-" + tag,
		MyInt64:  time.Now().Unix(),
		MySubs:   make(map[string]*testtypes.TestProtoSubSub),
	}
	sub.MySubs["sub"] = &testtypes.TestProtoSubSub{MyString: "sub", Int32Map: make(map[int32]int32)}
	sub.MySubs["sub"].Int32Map[0] = 0
	sub.MySubs["sub"].Int32Map[1] = 0

	sub1 := &testtypes.TestProtoSub{
		MyString: "string-sub-1-" + tag,
		MyInt64:  time.Now().Unix(),
	}
	sub2 := &testtypes.TestProtoSub{
		MyString: "string-sub-2-" + tag,
		MyInt64:  time.Now().Unix(),
		MySubs:   make(map[string]*testtypes.TestProtoSubSub),
	}
	sub2.MySubs["sub2"] = &testtypes.TestProtoSubSub{MyString: "sub2-string-sub", Int32Map: make(map[int32]int32)}
	sub2.MySubs["sub2"].Int32Map[0] = 0
	sub2.MySubs["sub2"].Int32Map[1] = 0
	i := &testtypes.TestProto{
		MyString:           "string-" + tag,
		MyFloat64:          123456.123456,
		MyBool:             true,
		MyFloat32:          123.123,
		MyInt32:            int32(index),
		MyInt64:            int64(index * 10),
		MyInt32Slice:       []int32{1, 2, 3, int32(index)},
		MyStringSlice:      []string{"a", "b", "c", "d", tag},
		MyInt32ToInt64Map:  map[int32]int64{1: 11, 2: 22, 3: 33, 4: 44, int32(index): int64(index * 10)},
		MyString2StringMap: map[string]string{"a": "aa", "b": "bb", "c": "cc", tag: tag + tag},
		MySingle:           sub,
		MyModelSlice:       []*testtypes.TestProtoSub{sub1, sub2},
		MyString2ModelMap:  map[string]*testtypes.TestProtoSub{sub1.MyString: sub1, sub2.MyString: sub2},
		MyEnum:             testtypes.TestEnum_ValueOne,
	}
	return i
}
