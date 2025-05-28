package t_resources

import (
	"bytes"
	"github.com/saichler/l8services/go/services/manager"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"strconv"
	"testing"
	"time"
)

const (
	VNET_PREFIX = "vnet-"
	VNIC_PREFIX = "-vnic-"
)

var Log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func CreateResources(vnetPort, vnicNum int, level ifs.LogLevel) (ifs.IResources, string) {
	alias := AliasOf(vnetPort, vnicNum)
	_log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	_log.SetLogLevel(level)
	_resources := resources.NewResources(_log)
	_resources.Set(registry.NewRegistry())
	_security, err := ifs.LoadSecurityProvider()
	if err != nil {
		panic("Failed to load security provider " + err.Error())
	}
	_resources.Set(_security)
	_config := &types.SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  alias,
		VnetPort:    uint32(vnetPort)}
	_resources.Set(_config)
	_introspector := introspecting.NewIntrospect(_resources.Registry())
	_resources.Set(_introspector)
	_servicepoints := manager.NewServices(_resources)
	_resources.Set(_servicepoints)
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

func CreateSubSubModelInstance(index1, index2, index3, k int, v ...int32) *testtypes.TestProtoSubSub {
	key := bytes.Buffer{}
	key.WriteString("string-")
	key.WriteString(strconv.Itoa(index1))
	key.WriteString("-")
	key.WriteString(strconv.Itoa(index2))
	key.WriteString("-")
	key.WriteString(strconv.Itoa(index3))
	subsub := &testtypes.TestProtoSubSub{}

	subsub.MyString = key.String()
	subsub.MyInt64 = time.Now().UnixNano()
	subsub.Int32Map = make(map[int32]int32)
	for i := 0; i < k; i++ {
		subsub.Int32Map[int32(i)] = v[i]
	}
	return subsub
}

func CreateSubTestModelInstance(index1, index2 int, n, k int, v ...int32) *testtypes.TestProtoSub {
	key := bytes.Buffer{}
	key.WriteString("string-")
	key.WriteString(strconv.Itoa(index1))
	key.WriteString("-")
	key.WriteString(strconv.Itoa(index2))
	sub := &testtypes.TestProtoSub{}
	sub.MyString = key.String()
	sub.MyInt64 = time.Now().UnixNano()
	sub.MySubs = make(map[string]*testtypes.TestProtoSubSub)

	for i := 0; i < n; i++ {
		subsub := CreateSubSubModelInstance(index1, index2, i, k, v...)
		sub.MySubs[subsub.MyString] = subsub
	}

	return sub
}

func CreateTestModelInstance(index int) *testtypes.TestProto {
	key := bytes.Buffer{}
	key.WriteString("string-")
	key.WriteString(strconv.Itoa(index))

	sub1m := CreateSubTestModelInstance(index, 1, 2, 2, 0, 0)
	sub2m := CreateSubTestModelInstance(index, 2, 2, 2, 0, 1)
	sub3m := CreateSubTestModelInstance(index, 3, 2, 4, 0, 1, 2, 3)

	i := &testtypes.TestProto{
		MyString:           key.String(),
		MyFloat64:          123456.123456,
		MyBool:             true,
		MyFloat32:          123.123,
		MyInt32:            int32(index),
		MyInt64:            int64(index * 10),
		MyInt32Slice:       []int32{1, 2, 3, int32(index)},
		MyStringSlice:      []string{"a", "b", "c", "d", key.String()},
		MyInt32ToInt64Map:  map[int32]int64{1: 11, 2: 22, 3: 33, 4: 44, int32(index): int64(index * 10)},
		MyString2StringMap: map[string]string{"a": "aa", "b": "bb", "c": "cc", key.String(): key.String() + key.String()},
		MySingle:           CreateSubTestModelInstance(index, 0, 2, 2, 0, 0),
		MyModelSlice: []*testtypes.TestProtoSub{CreateSubTestModelInstance(index, 1, 2, 2, 0, 0),
			CreateSubTestModelInstance(index, 2, 2, 2, 0, 1)},
		MyString2ModelMap: map[string]*testtypes.TestProtoSub{sub1m.MyString: sub1m,
			sub2m.MyString: sub2m, sub3m.MyString: sub3m},
		MyEnum: testtypes.TestEnum_ValueOne,
	}
	return i
}

func CloneTestModel(a *testtypes.TestProto) *testtypes.TestProto {
	cloner := cloning.NewCloner()
	return cloner.Clone(a).(*testtypes.TestProto)
}

func WaitForCondition(cond func() bool, timeoutInSeconds int64, t *testing.T, failMessage string) bool {
	start := time.Now().UnixMilli()
	end := start + timeoutInSeconds*1000
	for start < end {
		if cond() {
			return true
		}
		time.Sleep(time.Millisecond * 100)
		start += 100
	}
	if t != nil {
		Log.Fail(t, failMessage)
	} else {
		Log.Error(failMessage)
	}
	return false
}
