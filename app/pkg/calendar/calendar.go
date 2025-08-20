package calendar

import (
	"fmt"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"time"
)

// GetChineseNewYearDate 接受一个公历年份，返回该年的春节（正月初一）的 time.Time 对象
func GetChineseNewYearDate(year int) (*time.Time, error) {
	// 创建一个农历正月初一的日历对象
	// calendar.ByLunar 的参数分别为：年, 月, 日, 时, 分, 秒, 是否闰月
	// 我们取指定年份的正月初一，时间设为0时0分0秒
	c := calendar.ByLunar(int64(year), 1, 1, 0, 0, 0, false)

	// 将农历日期对象转换为公历的 time.Time 对象
	// 这个库的 ToTime 方法返回的是一个 calendar.SolarTime 对象，需要进一步获取 time.Time
	solarTime, err := c.ToTime()
	if err != nil {
		return nil, fmt.Errorf("将农历转换为公历时出错: %w", err)
	}

	return &solarTime, nil
}
