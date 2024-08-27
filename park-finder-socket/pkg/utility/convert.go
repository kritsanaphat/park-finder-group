package utility

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertThaiTimeToUTC(input string) (time.Time, error) {
	year, _ := strconv.Atoi(input[0:4])
	month, _ := strconv.Atoi(input[5:7])
	day, _ := strconv.Atoi(input[8:10])
	hour, _ := strconv.Atoi(input[11:13])
	minute, _ := strconv.Atoi(input[14:16])

	// Create a time.Time instance using time.Date()
	t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	// Subtract 7 hours
	result := t.Add(-7 * time.Hour)
	return result, nil

}

func ConvertToDayName(dateString string) string {
	parsedDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "Null"
	}

	// Get the day of the week
	dayOfWeek := parsedDate.Weekday()
	return dayOfWeek.String()
}

func ConvertToBangkokTimeAndHours(dateTimeString time.Time) int {
	// โหลดตำแหน่งเวลาสำหรับ Asia/Bangkok
	bangkokLocation, _ := time.LoadLocation("Asia/Bangkok")

	// แปลงเวลาไปยัง timezone ที่กำหนด
	bangkokTime := dateTimeString.In(bangkokLocation)

	// คืนค่าชั่วโมงเป็นตัวเลข
	return bangkokTime.Hour()
}
