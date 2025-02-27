package utility

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
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

func GenerateRewardCODE() string {
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

func GenerateTransactionID() int {

	minValue := int64(1)
	maxValue := int64(1)
	for i := 1; i < 6; i++ {
		minValue *= 10
		maxValue *= 10
	}
	maxValue = (maxValue * 10) - 1

	transaction_id := rand.Int63n(maxValue-minValue+1) + minValue
	value := fmt.Sprintf("%0*d", 10, transaction_id)
	value_int, err := strconv.Atoi(value)

	if err != nil {
		fmt.Println(err)
	}
	return value_int
}

func GenHeaderLinePay(URL string, req interface{}) map[string]string {
	nonce := generateUUID()
	ChannelID := os.Getenv("CHANNEL_ID")
	ChannelSecret := os.Getenv("CHANNEL_SECRET")

	var RequestDetailStr string
	switch v := req.(type) {
	case string:
		RequestDetailStr = v
	default:
		jsonData, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		RequestDetailStr = string(jsonData)
	}

	data := ChannelSecret + URL + RequestDetailStr + nonce

	keyBytes := []byte(ChannelSecret)
	hmacSha256 := hmac.New(sha256.New, keyBytes)
	hmacSha256.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(hmacSha256.Sum(nil))

	header := map[string]string{
		"Content-Type":               "application/json",
		"X-LINE-ChannelId":           ChannelID,
		"X-LINE-Authorization-Nonce": nonce,
		"X-LINE-Authorization":       signature,
	}
	return header
}

func generateUUID() string {
	return uuid.New().String()
}
