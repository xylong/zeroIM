package snowflake

import (
	"github.com/sony/sonyflake/v2"
	"strconv"
	"time"
)

var sf *sonyflake.Sonyflake

func init() {
	var err error
	sf, err = sonyflake.New(sonyflake.Settings{
		StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (int, error) {
			return 1, nil
		},
	})
	if err != nil {
		panic(err)
	}
	if sf == nil {
		panic("failed to initialize snowflake")
	}
}

func GenID() (int64, error) {
	return sf.NextID()
}

func GenIDStr() string {
	id, err := sf.NextID()
	if err != nil {
		return ""
	}
	return strconv.FormatInt(id, 10)
}
