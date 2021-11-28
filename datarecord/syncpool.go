package datarecord

import "sync"

var DATAFRAME_SYNC_POOL *sync.Pool

func init() {
	DATAFRAME_SYNC_POOL = &sync.Pool{New: func() interface{} {
		return &DataFrame{
			AgentID:    "",
			DomainID:   0,
			ExportTime: 0,
			SetID:      0,
			Type:       "",
			Record:     make([]map[string]interface{}, 0, 1),
		}
	}}
}

func NewDataFrame() *DataFrame {
	v := DATAFRAME_SYNC_POOL.Get().(*DataFrame)
	return v
}

func Free(d *DataFrame) {
	DATAFRAME_SYNC_POOL.Put(d)
}
