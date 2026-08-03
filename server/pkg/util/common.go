package util

import (
	"fmt"
	"strconv"
	"time"
)

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func IntPtr(i int) *int {
	return &i
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func ParseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func ConvertStringOperatorTag(tag string) int {
	switch tag {
	case "自运营":
		return 1
	case "代运营":
		return 2
	case "走量":
		return 3
	case "收量":
		return 4
	case "非收":
		return 5
	case "无标签":
		return 6
	default:
		return 0
	}
}

func ConvertIntOperatorTagText(tag int) string {
	switch tag {
	case 1:
		return "自运营"
	case 2:
		return "代运营"
	case 3:
		return "走量"
	case 4:
		return "收量"
	case 5:
		return "非收"
	case 6:
		return "无标签"
	default:
		return ""
	}
}

func BuildYearMonthPatternByQuarter(year int, quarter int) []string {
	var monthPatterns []string
	switch quarter {
	case 1:
		monthPatterns = []string{
			fmt.Sprintf("%d-01", year),
			fmt.Sprintf("%d-02", year),
			fmt.Sprintf("%d-03", year),
		}
	case 2:
		monthPatterns = []string{
			fmt.Sprintf("%d-04", year),
			fmt.Sprintf("%d-05", year),
			fmt.Sprintf("%d-06", year),
		}
	case 3:
		monthPatterns = []string{
			fmt.Sprintf("%d-07", year),
			fmt.Sprintf("%d-08", year),
			fmt.Sprintf("%d-09", year),
		}
	case 4:
		monthPatterns = []string{
			fmt.Sprintf("%d-10", year),
			fmt.Sprintf("%d-11", year),
			fmt.Sprintf("%d-12", year),
		}
	}
	return monthPatterns
}
func BuildMonthPatternByQuarter(year int, quarter int) []string {
	var monthPatterns []string
	switch quarter {
	case 1:
		monthPatterns = []string{"1", "2", "3"}
	case 2:
		monthPatterns = []string{"4", "5", "6"}
	case 3:
		monthPatterns = []string{"7", "8", "9"}
	case 4:
		monthPatterns = []string{"10", "11", "12"}
	}
	return monthPatterns
}

func BuildMonthPeriodsByQuarter(year, quarter int) []string {
	startMonth := (quarter-1)*3 + 1
	periods := make([]string, 0, 3)
	for month := startMonth; month < startMonth+3; month++ {
		periods = append(periods, fmt.Sprintf("%d-%02d", year, month))
	}
	return periods
}

func GetYearByPeriod(period string) int {
	if len(period) >= 4 {
		if t, err := time.Parse("2006-01", period); err == nil {
			return t.Year()
		}
	}
	return time.Now().Year()
}

func GetQuarterByPeriod(period string) int {
	if t, err := time.Parse("2006-01", period); err == nil {
		return (int(t.Month())-1)/3 + 1
	}
	return 0
}
