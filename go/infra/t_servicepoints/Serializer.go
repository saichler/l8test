package t_servicepoints

import "github.com/saichler/types/go/common"

type TestSerializer struct {
}

func (ts *TestSerializer) Mode() common.SerializerMode {
	return common.BINARY
}
func (ts *TestSerializer) Marshal(interface{}, common.IRegistry) ([]byte, error) {
	return nil, nil
}
func (ts *TestSerializer) Unmarshal([]byte, string, common.IRegistry) (interface{}, error) {
	return nil, nil
}
