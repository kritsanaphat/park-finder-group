package search

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s SearchServices) Search(ctx context.Context, req *models.SearchQueryRequest, keyword string) *[]models.ParkingArea {

	province1, province2 := utility.ProvinceCal(req.Latitude, req.Longitude)

	difference_end := 0
	if req.MinEnd > 0 {
		difference_end = 1
	}
	var parkings []models.ParkingArea

	if len(req.Date) == 1 {
		date, err := time.Parse("2006-01-02", req.Date[0])
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Println(province1, province2)
		parkings = s.FilterAreaOneDay(ctx, date, req, province1, province2, keyword, difference_end)
	} else if len(req.Date) > 1 {
		parkings = s.FilterAreaMoreThanOneDay(ctx, req.Date, req, province1, province2, keyword, difference_end)
	} else {
		return nil
	}

	return &parkings

}

func (s SearchServices) FilterAreaOneDay(ctx context.Context, date time.Time, req *models.SearchQueryRequest, province1, province2, keyword string, difference_end int) []models.ParkingArea {
	var parkings []models.ParkingArea
	daliy_open := "open_detail." + strings.ToLower(date.Weekday().String()) + ".open_time"
	daliy_close := "open_detail." + strings.ToLower(date.Weekday().String()) + ".close_time"
	search := models.ReserveLog{
		DateStart: req.Date[0],
		DateEnd:   req.Date[0],
		HourStart: req.HourStart,
		MinStart:  req.MinStart,
		HourEnd:   req.HourEnd,
		MinEnd:    req.MinEnd,
	}
	if keyword == "" {
		fmt.Println("Without Tag mode")
		area, err := s.ParkingAreaStorage.FindParkingAvailableOneDay(ctx, province1, province2, req.Date[0], daliy_open, daliy_close, req.MinPrice, req.MaxPrice, req.Review, req.HourStart, req.HourEnd, req.MinStart, difference_end)
		if err == mongo.ErrNoDocuments {
			return nil
		}

		parkings = s.FilterTimeAndSlot(ctx, search, area, req.Latitude, req.Longitude)

	} else {
		fmt.Println("Tag mode")
		area, err := s.ParkingAreaStorage.FindParkingAvailableWithTagOneDay(ctx, province1, province2, req.Date[0], daliy_open, daliy_close, keyword, req.MinPrice, req.MaxPrice, req.Review, req.HourStart, req.HourEnd, req.MinStart, difference_end)
		if err == mongo.ErrNoDocuments {
			return nil
		}

		parkings = s.FilterTimeAndSlot(ctx, search, area, req.Latitude, req.Longitude)
	}

	sort.Slice(parkings, func(i, j int) bool {
		return parkings[i].Distance < parkings[j].Distance
	})

	limit := 30
	if len(parkings) > limit {
		parkings = parkings[:limit]
	}
	return parkings
}

func (s SearchServices) FilterAreaMoreThanOneDay(ctx context.Context, date []string, req *models.SearchQueryRequest, province1, province2, keyword string, difference_end int) []models.ParkingArea {
	var parkings []models.ParkingArea
	date_names := utility.ConvertToDayNames(date[0], date[len(date)-1])
	search := models.ReserveLog{
		DateStart: req.Date[0],
		DateEnd:   req.Date[1],
		HourStart: req.HourStart,
		MinStart:  req.MinStart,
		HourEnd:   req.HourEnd,
		MinEnd:    req.MinEnd,
	}

	if keyword == "" {
		fmt.Println("Without Tag mode")
		area, err := s.ParkingAreaStorage.FindParkingAvailableMoreOneDay(ctx, province1, province2, date_names, date, req.MinPrice, req.MaxPrice, req.Review, req.HourStart, req.HourEnd, req.MinStart, difference_end)
		if err == mongo.ErrNoDocuments {
			return nil
		}
		parkings = s.FilterTimeAndSlot(ctx, search, area, req.Latitude, req.Longitude)

	} else {
		fmt.Println("Tag mode")
		area, err := s.ParkingAreaStorage.FindParkingAvailableWithTagMoreOneDay(ctx, keyword, date_names, date, req.MinPrice, req.MaxPrice, req.Review, req.HourStart, req.HourEnd, req.MinStart, difference_end)
		if err == mongo.ErrNoDocuments {
			return nil
		}
		parkings = s.FilterTimeAndSlot(ctx, search, area, req.Latitude, req.Longitude)
	}

	sort.Slice(parkings, func(i, j int) bool {
		return parkings[i].Distance < parkings[j].Distance
	})

	limit := 30
	if len(parkings) > limit {
		parkings = parkings[:limit]
	}
	return parkings
}

func (s SearchServices) FilterTimeAndSlot(ctx context.Context, search models.ReserveLog, area *mongo.Cursor, lat, long float64) []models.ParkingArea {
	var parkings []models.ParkingArea

	for area.Next(ctx) {
		var parking models.ParkingArea

		if err := area.Decode(&parking); err != nil {
			return nil
		}
		if checkReservation(search, parking.ReserveLog, parking.ToatalParkingCount) {
			distance := utility.DistanceCal(lat, long, parking.Address.Latitude, parking.Address.Longitude)
			parking.Distance = float32(distance)
			parkings = append(parkings, parking)
		}
	}
	return parkings
}

func checkReservation(desiredSlot models.ReserveLog, reservations []models.ReserveLog, slot int) bool {
	var count int
	desiredStart, err := time.Parse("2006-01-02 15:04", desiredSlot.DateStart+" "+fmt.Sprintf("%02d:%02d", desiredSlot.HourStart, desiredSlot.MinStart))
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return false
	}
	desiredEnd, err := time.Parse("2006-01-02 15:04", desiredSlot.DateEnd+" "+fmt.Sprintf("%02d:%02d", desiredSlot.HourEnd, desiredSlot.MinEnd))
	if err != nil {
		fmt.Println("Error parsing end time:", err)
		return false
	}

	for _, reservation := range reservations {
		start, err := time.Parse("2006-01-02 15:04", reservation.DateStart+" "+fmt.Sprintf("%02d:%02d", reservation.HourStart, reservation.MinStart))
		if err != nil {
			fmt.Println("Error parsing start time for reservation:", err)
			continue
		}
		end, err := time.Parse("2006-01-02 15:04", reservation.DateEnd+" "+fmt.Sprintf("%02d:%02d", reservation.HourEnd, reservation.MinEnd))
		if err != nil {
			fmt.Println("Error parsing end time for reservation:", err)
			continue
		}

		if !(desiredEnd.Before(start) || desiredStart.After(end)) {
			count++
		}
	}

	return slot > count
}
