package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDate(date interface{}, formats ...string) time.Time {
	if date == nil {
		return time.Time{}
	}
	
	switch d := date.(type) {
	case time.Time:
		return d
	case int64:
		return time.UnixMilli(d)
	case string:
		return parseDateString(d, formats...)
	case float64:
		return time.UnixMilli(int64(d))
	default:
		return time.Time{}
	}
}

func parseDateString(dateStr string, formats ...string) time.Time {
	if dateStr == "" || dateStr == "刚刚" {
		return time.Now()
	}

	dateStr = strings.TrimSpace(dateStr)
	
	dateStr = strings.ToLower(dateStr)
	dateStr = regexp.MustCompile(`^(an?\s)|(\san?\s)`).ReplaceAllString(dateStr, "1")
	dateStr = regexp.MustCompile(`几|幾`).ReplaceAllString(dateStr, "3")
	dateStr = regexp.MustCompile(`[\s,]`).ReplaceAllString(dateStr, "")

	now := time.Now()
	
	// 处理相对时间
	if matches := regexp.MustCompile(`^(\d+)(秒|分钟|分|小时|时|天|日|周|星期|月|年)`).FindStringSubmatch(dateStr); len(matches) > 0 {
		num, _ := strconv.Atoi(matches[1])
		unit := matches[2]
		switch unit {
		case "秒":
			return now.Add(-time.Duration(num) * time.Second)
		case "分钟", "分":
			return now.Add(-time.Duration(num) * time.Minute)
		case "小时", "时":
			return now.Add(-time.Duration(num) * time.Hour)
		case "天", "日":
			return now.AddDate(0, 0, -num)
		case "周":
			return now.AddDate(0, 0, -num*7)
		case "月":
			return now.AddDate(0, -num, 0)
		case "年":
			return now.AddDate(-num, 0, 0)
		}
	}

	// 处理 "X分钟前", "X小时前", "X天前"
	if matches := regexp.MustCompile(`^(\d+)(分钟|分|小时|时|天|日|周)前$`).FindStringSubmatch(dateStr); len(matches) > 0 {
		num, _ := strconv.Atoi(matches[1])
		unit := matches[2]
		switch unit {
		case "分钟", "分":
			return now.Add(-time.Duration(num) * time.Minute)
		case "小时", "时":
			return now.Add(-time.Duration(num) * time.Hour)
		case "天", "日":
			return now.AddDate(0, 0, -num)
		case "周":
			return now.AddDate(0, 0, -num*7)
		}
	}

	// 处理中文日期 "今天", "昨天", "前天"
	datePatterns := map[string]time.Time{
		"今天":     todayStart(),
		"今日":     todayStart(),
		"昨天":     todayStart().AddDate(0, 0, -1),
		"昨日":     todayStart().AddDate(0, 0, -1),
		"前天":     todayStart().AddDate(0, 0, -2),
		"明天":     todayStart().AddDate(0, 0, 1),
		"明日":     todayStart().AddDate(0, 0, 1),
		"后天":     todayStart().AddDate(0, 0, 2),
	}

	for pattern, baseTime := range datePatterns {
		if strings.HasPrefix(dateStr, pattern) {
			timePart := strings.TrimPrefix(dateStr, pattern)
			if t, err := parseTimePart(timePart); err == nil {
				return baseTime.Add(t)
			}
		}
	}

	// 处理周几
	weekdayMap := map[string]time.Weekday{
		"周一": time.Monday,
		"周二": time.Tuesday,
		"周三": time.Wednesday,
		"周四": time.Thursday,
		"周五": time.Friday,
		"周六": time.Saturday,
		"周日": time.Sunday,
		"星期日": time.Sunday,
	}

	for prefix, weekday := range weekdayMap {
		if strings.HasPrefix(dateStr, prefix) {
			timePart := strings.TrimPrefix(dateStr, prefix)
			weekStart := getLastWeekday(weekday)
			if t, err := parseTimePart(timePart); err == nil {
				return weekStart.Add(t)
			}
		}
	}

	// 尝试标准格式
	formatsToTry := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"01/02/2006",
		"2006/01/02",
	}

	if len(formats) > 0 {
		formatsToTry = append(formats, formatsToTry...)
	}

	for _, format := range formatsToTry {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	// 最后尝试 Unix 时间戳
	if ts, err := strconv.ParseInt(dateStr, 10, 64); err == nil {
		if ts > 1e11 {
			return time.UnixMilli(ts)
		}
		return time.Unix(ts, 0)
	}

	return time.Time{}
}

func todayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func getLastWeekday(weekday time.Weekday) time.Time {
	now := time.Now()
	for i := 0; i <= 7; i++ {
		d := now.AddDate(0, 0, -i)
		if d.Weekday() == weekday {
			return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
		}
	}
	return now
}

func parseTimePart(timePart string) (time.Duration, error) {
	timePart = strings.TrimSpace(timePart)
	if timePart == "" {
		return 0, nil
	}
	
	// 匹配 "10:30", "10:30:00"
	if matches := regexp.MustCompile(`^(\d{1,2}):(\d{2})(?::(\d{2}))?$`).FindStringSubmatch(timePart); len(matches) >= 3 {
		hour, _ := strconv.Atoi(matches[1])
		minute, _ := strconv.Atoi(matches[2])
		second := 0
		if len(matches) > 3 && matches[3] != "" {
			second, _ = strconv.Atoi(matches[3])
		}
		return time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(second)*time.Second, nil
	}

	return 0, nil
}

func ParseRelativeDate(date string) time.Time {
	return parseDateString(date)
}

func ParseRelativeDateToUnix(date string) int64 {
	return ParseRelativeDate(date).Unix()
}
