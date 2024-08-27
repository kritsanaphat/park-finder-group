package utility

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateOrderID() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateOTP() string {

	minValue := int64(1)
	maxValue := int64(1)
	for i := 1; i < 6; i++ {
		minValue *= 10
		maxValue *= 10
	}
	maxValue = (maxValue * 10) - 1

	otp := rand.Int63n(maxValue-minValue+1) + minValue

	return fmt.Sprintf("%0*d", 6, otp)
}

func generateUUID() string {
	return uuid.New().String()
}
