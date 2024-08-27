package services

import (
	"context"
	"errors"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (as AdminServices) AdminRegister(ctx context.Context, user *models.RegisterAccountRequest) error {

	pHash, err := utility.HashPassword([]byte(user.Password))
	if err != nil {
		return err
	}
	data := &models.AdminAccount{
		ID:        primitive.NewObjectID(),
		FirstName: user.FristName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  pHash,
		TimeStamp: time.Now(),
	}

	result, err := as.AdminAccoutStorage.InsertAccount(ctx, data)
	if err != nil {
		return err
	}

	if result.InsertedID == nil {
		return errors.New("no documents were inserted")
	}

	return nil
}

func (cs AdminServices) CheckExistEmail(ctx context.Context, email string) *models.AdminAccount {
	filter := bson.M{"email": email}

	user := new(models.AdminAccount)
	err := cs.AdminAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return user
}

func (cs AdminServices) CheckVerifyEmail(ctx context.Context, email string) *models.AdminAccount {
	filter := bson.M{"email": email, "verify": true}

	user := new(models.AdminAccount)
	err := cs.AdminAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err != nil {
		return nil
	}
	return user
}

func (cs AdminServices) CheckPassword(user models.AdminAccount, password string) bool {
	err := utility.CheckPasswordHash([]byte(user.Password), []byte(password))
	return err == nil
}

func (cs AdminServices) AddToken(ctx context.Context, user *models.AdminAccount, token string) error {
	user_id := user.IDToString()
	data := &models.Token{
		UserID: user_id,
		Token:  token,
		Valid:  true,
		Role:   "admin",
	}
	cs.TokenStorage.InsertToken(ctx, data)

	return nil
}

func (cs AdminServices) CheckExistToken(ctx context.Context, tk string) *models.Token {
	token := cs.TokenStorage.FindToken(ctx, tk)
	return token
}

func (cs AdminServices) RevokeToken(ctx context.Context, user *models.AdminAccount, token string) error {

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

func (cs AdminServices) RevokeExpireToken(ctx context.Context, token string, expireDate time.Time) error {
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
