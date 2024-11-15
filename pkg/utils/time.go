package utils

import (
	"strconv"
	"time"
)

func ConvertMillisecondsToTime(timestamp string) (time.Time, error) {
	// 将字符串转换为整数类型的毫秒数
	ms, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// 将毫秒转换为秒和纳秒，然后转为time.Time
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * int64(time.Millisecond)

	return time.Unix(seconds, nanoseconds), nil
}
