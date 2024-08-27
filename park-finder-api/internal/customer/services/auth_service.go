package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var redisOTPKeyFormat string = "OTP:Email:%v"

func (cs CustomerServices) CustomerRegister(ctx context.Context, user *models.RegisterAccountRequest) error {

	pHash, err := utility.HashPassword([]byte(user.Password))
	if err != nil {
		return err
	}
	data := &models.CustomerAccount{
		ID:           primitive.NewObjectID(),
		FirstName:    user.FristName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Email:        user.Email,
		Password:     pHash,
		TimeStamp:    time.Now(),
		Address:      []models.CustomerAddress{},
		Cashback:     0,
		Point:        0,
		Reward:       []models.CustomerReward{},
		FavoriteArea: []string{},
		Fine:         &models.CustomerFine{},
	}

	result, err := cs.CustomerAccoutStorage.InsertAccount(ctx, data)
	if err != nil {
		return err
	}

	if result.InsertedID == nil {
		return errors.New("no documents were inserted")
	}

	return nil
}

func (cs CustomerServices) UpdateCustomerPassword(ctx context.Context, user *models.ResetPasswordRequest, email string) error {
	pHash, err := utility.HashPassword([]byte(user.NewPassword))
	if err != nil {
		return err
	}
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": models.UpdatePassword{
			Password:  pHash,
			TimeStamp: time.Now(),
		}}
	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) UpdateCustomerPasswordByForgot(ctx context.Context, password string, email string) error {
	pHash, err := utility.HashPassword([]byte(password))
	if err != nil {
		return err
	}
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": models.UpdatePassword{
			Password:  pHash,
			TimeStamp: time.Now(),
		}}
	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) UpdateCustomerVerify(ctx context.Context, email string) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": models.VerifyCustomer{
			Verify:    true,
			TimeStamp: time.Now().In(location),
		}}
	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) SaveOTP(email string, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	// Cache for 2 minute
	result := cs.Redis.Set(key, otp, 120)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (cs CustomerServices) SaveOTPForgot(email string, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	// Cache for 2 minute
	result := cs.Redis.Set(key+"_forgot", otp, 120)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (cs CustomerServices) CheckOTP(email, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	result := cs.Redis.Get(key)
	storedOTP, err := result.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("OTP not found or expired")
		}
		return err
	}

	if storedOTP != otp {
		return errors.New("Invalid OTP")
	}

	fmt.Println("OTP is valid for email:", email)
	return nil
}

func (cs CustomerServices) CheckOTPForgot(email, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	result := cs.Redis.Get(key + "_forgot")
	storedOTP, err := result.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("OTP not found or expired")
		}
		return err
	}

	if storedOTP != otp {
		return errors.New("Invalid OTP")
	}

	fmt.Println("OTP is valid for email:", email)
	return nil
}

func (cs CustomerServices) RemoveOTP(email string, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf(redisOTPKeyFormat, email) + "_forgot"

	// Cache for 2 minute
	result := cs.Redis.Client.Del(ctx, key)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (cs CustomerServices) RemoveOTPForgot(email string, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf(redisOTPKeyFormat, email) + "_forgot"

	// Cache for 2 minute
	result := cs.Redis.Client.Del(ctx, key)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (cs CustomerServices) CheckExistEmail(ctx context.Context, email string) *models.CustomerAccount {
	filter := bson.M{"email": email}

	user := new(models.CustomerAccount)
	err := cs.CustomerAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return user
}

func (cs CustomerServices) CheckVerifyEmail(ctx context.Context, email string) *models.CustomerAccount {
	filter := bson.M{"email": email, "verify": true}

	user := new(models.CustomerAccount)
	err := cs.CustomerAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err != nil {
		return nil
	}
	return user
}

func (cs CustomerServices) CheckPassword(user models.CustomerAccount, password string) bool {
	err := utility.CheckPasswordHash([]byte(user.Password), []byte(password))
	return err == nil
}

func (cs CustomerServices) AddToken(ctx context.Context, user *models.CustomerAccount, token string) error {
	user_id := user.IDToString()
	data := &models.Token{
		UserID: user_id,
		Token:  token,
		Valid:  true,
		Role:   "customer",
	}
	cs.TokenStorage.InsertToken(ctx, data)

	return nil
}

func (cs CustomerServices) CheckExistToken(ctx context.Context, tk string) *models.Token {
	token := cs.TokenStorage.FindToken(ctx, tk)
	return token
}

func (cs CustomerServices) RevokeToken(ctx context.Context, user *models.CustomerAccount, token string) error {

	user_id := user.IDToString()
	filter := bson.M{
		"user_id": user_id,
		"token":   token,
	}
	update := bson.M{
		"$set": bson.M{
			"valid":       false,
			"revoke_date": time.Now(),
		},
	}
	cs.TokenStorage.UpdateToken(ctx, filter, update)
	return nil
}

func (cs CustomerServices) RevokeExpireToken(ctx context.Context, token string, expireDate time.Time) error {
	filter := bson.M{
		"token": token,
	}
	update := bson.M{
		"$set": bson.M{
			"valid":       false,
			"revoke_date": expireDate,
		},
	}
	cs.TokenStorage.UpdateExpireToken(ctx, filter, update)
	return nil
}
