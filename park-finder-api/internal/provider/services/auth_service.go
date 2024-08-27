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

func (ps ProviderServices) ProviderRegister(ctx context.Context, user *models.RegisterAccountRequest) error {

	pHash, err := utility.HashPassword([]byte(user.Password))
	if err != nil {
		return err
	}
	var bank models.BankAccount
	data := &models.ProviderAccount{
		ID:          primitive.NewObjectID(),
		FirstName:   user.FristName,
		LastName:    user.LastName,
		Phone:       user.Phone,
		Email:       user.Email,
		Password:    pHash,
		TimeStamp:   time.Now(),
		BankAccount: bank,
	}

	ps.ProviderAccoutStorage.InsertAccount(ctx, data)

	return nil
}

func (ps ProviderServices) UpdateProviderPassword(ctx context.Context, user *models.ResetPasswordRequest, email string) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
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
			TimeStamp: time.Now().In(location),
		}}
	result, err := ps.ProviderAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (ps ProviderServices) UpdateProviderVerify(ctx context.Context, email string) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": models.VerifyProvider{
			Verify:    true,
			TimeStamp: time.Now().In(location),
		}}

	result, err := ps.ProviderAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (ps ProviderServices) SaveOTP(email, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	// Cache for 2 minute
	result := ps.Redis.Set(key, otp, 120)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (ps ProviderServices) CheckOTP(email, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	result := ps.Redis.Get(key)
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

func (cs ProviderServices) RemoveOTP(email string, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf(redisOTPKeyFormat, email) + "_forgot"

	result := cs.Redis.Client.Del(ctx, key)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (cs ProviderServices) RemoveOTPForgot(email string, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf(redisOTPKeyFormat, email) + "_forgot"

	result := cs.Redis.Client.Del(ctx, key)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (ps ProviderServices) CheckExistEmail(ctx context.Context, email string) *models.ProviderAccount {
	filter := bson.M{"email": email}

	user := new(models.ProviderAccount)
	err := ps.ProviderAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return user
}

func (ps ProviderServices) CheckVerifyEmail(ctx context.Context, email string) *models.ProviderAccount {
	filter := bson.M{"email": email, "verify": true}

	user := new(models.ProviderAccount)
	err := ps.ProviderAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err != nil {
		return nil
	}
	return user
}

func (ps ProviderServices) CheckPassword(user models.ProviderAccount, password string) bool {
	err := utility.CheckPasswordHash([]byte(user.Password), []byte(password))
	return err == nil
}

func (ps ProviderServices) UpdateCustomerPasswordByForgot(ctx context.Context, password string, email string) error {
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
	result, err := ps.ProviderAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (ps ProviderServices) SaveOTPForgot(email string, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	// Cache for 2 minute
	result := ps.Redis.Set(key+"_forgot", otp, 120)
	err := result.Err()
	if err != nil {
		return err
	}
	fmt.Println("Cache user otp for email:", email)

	return nil
}

func (ps ProviderServices) CheckOTPForgot(email, otp string) error {
	key := fmt.Sprintf(redisOTPKeyFormat, email)

	result := ps.Redis.Get(key + "_forgot")
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

func (ps ProviderServices) AddToken(ctx context.Context, user *models.ProviderAccount, token string) error {
	user_id := user.IDToString()
	data := &models.Token{
		UserID: user_id,
		Token:  token,
		Valid:  true,
		Role:   "provider",
	}
	ps.TokenStorage.InsertToken(ctx, data)

	return nil
}

func (ps ProviderServices) CheckExistToken(ctx context.Context, tk string) *models.Token {
	token := ps.TokenStorage.FindToken(ctx, tk)
	return token
}

func (ps ProviderServices) RevokeToken(ctx context.Context, user *models.ProviderAccount, token string) error {

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
	ps.TokenStorage.UpdateToken(ctx, filter, update)
	return nil
}

func (ps ProviderServices) RevokeExpireToken(ctx context.Context, token string, expireDate time.Time) error {
	filter := bson.M{
		"token": token,
	}
	update := bson.M{
		"$set": bson.M{
			"valid":       false,
			"revoke_date": expireDate,
		},
	}
	ps.TokenStorage.UpdateExpireToken(ctx, filter, update)
	return nil
}
