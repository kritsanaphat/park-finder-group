package utility

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ConvertThaiTimeToUTC(input string) (time.Time, error) {
	year, _ := strconv.Atoi(input[0:4])
	month, _ := strconv.Atoi(input[5:7])
	day, _ := strconv.Atoi(input[8:10])
	hour, _ := strconv.Atoi(input[11:13])
	minute, _ := strconv.Atoi(input[14:16])

	t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	result := t.Add(-7 * time.Hour)
	return result, nil

}

func ConvertToDayName(dateString string) string {
	parsedDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "Null"
	}

	dayOfWeek := parsedDate.Weekday()
	return dayOfWeek.String()
}

func ConvertToBangkokTimeAndHours(dateTimeString time.Time) int {
	bangkokLocation, _ := time.LoadLocation("Asia/Bangkok")
	bangkokTime := dateTimeString.In(bangkokLocation)
	return bangkokTime.Hour()
}

func ConvertToDayNames(startDate, endDate string) []string {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return nil
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println("Error parsing end date:", err)
		return nil
	}

	if startTime.After(endTime) {
		fmt.Println("Start date is after end date")
		return nil
	}

	var result []string

	currentTime := startTime
	for !currentTime.After(endTime) {
		result = append(result, "open_detail."+strings.ToLower(currentTime.Weekday().String())+".open_time")
		result = append(result, "open_detail."+strings.ToLower(currentTime.Weekday().String())+".close_time")

		currentTime = currentTime.AddDate(0, 0, 1)
	}

	return result
}

func ConvertTimeToHourAndMin(currentHour, currentMin int) (int, int) {
	if currentMin > 30 {
		return currentHour + 1, 0
	}
	return currentHour, 30
}

func FormatThaiDateTime(t time.Time) string {
	thaiMonths := map[time.Month]string{
		1:  "มกราคม",
		2:  "กุมภาพันธ์",
		3:  "มีนาคม",
		4:  "เมษายน",
		5:  "พฤษภาคม",
		6:  "มิถุนายน",
		7:  "กรกฎาคม",
		8:  "สิงหาคม",
		9:  "กันยายน",
		10: "ตุลาคม",
		11: "พฤศจิกายน",
		12: "ธันวาคม",
	}

	dateFormat := fmt.Sprintf("%02d %s %d เวลา %02d:%02d", t.Day(), thaiMonths[t.Month()], t.Year(), t.Hour(), t.Minute())

	return dateFormat
}

func FormatThaiDateTimeFromString(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ""
	}

	thaiMonths := map[time.Month]string{
		1:  "มกราคม",
		2:  "กุมภาพันธ์",
		3:  "มีนาคม",
		4:  "เมษายน",
		5:  "พฤษภาคม",
		6:  "มิถุนายน",
		7:  "กรกฎาคม",
		8:  "สิงหาคม",
		9:  "กันยายน",
		10: "ตุลาคม",
		11: "พฤศจิกายน",
		12: "ธันวาคม",
	}

	dateFormat := fmt.Sprintf("%02d %s %d", t.Day(), thaiMonths[t.Month()], t.Year())

	return dateFormat
}

func ParseMonth(monthStr string) time.Month {
	months := map[string]time.Month{
		"January":   time.January,
		"February":  time.February,
		"March":     time.March,
		"April":     time.April,
		"May":       time.May,
		"June":      time.June,
		"July":      time.July,
		"August":    time.August,
		"September": time.September,
		"October":   time.October,
		"November":  time.November,
		"December":  time.December,
	}
	return months[monthStr]
}
