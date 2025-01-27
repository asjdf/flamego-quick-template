package sonyflake

import (
	"github.com/sony/sonyflake"
	"github.com/wujunyi792/flamego-quick-template/pkg/logx"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenSonyFlakeId() (int64, error) {

	id, err := flake.NextID()
	if err != nil {
		logx.Warning.Println("flake NextID failed: ", err)
		return 0, err
	}

	return int64(id), nil
}
