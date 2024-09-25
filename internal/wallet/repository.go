package wallet

import (
	"context"
	"go-api/pkg/db"

	usr "go-api/internal/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var walletCollection = db.GetCollection("wallets")

func DbGetWalletID(id primitive.ObjectID) (Wallet, error) {
	var wallet Wallet
	err := walletCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&wallet)
	return wallet, err
}

func DbGetWalletByUserID(userID primitive.ObjectID) (Wallet, error) {
	var wallet Wallet
	err := walletCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&wallet)
	return wallet, err
}

func DbUserWalletHasBalance(userID primitive.ObjectID, amount int) (bool, error) {
	var wallet Wallet
	err := walletCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		return false, err
	}
	return wallet.Balance >= amount, nil
}

func DbGetUserWallerByNickName(nickname string) (Wallet, error) {
	var user usr.User
	user, err := usr.DbGetUserByUsername(nickname)
	if err != nil {
		return Wallet{}, err
	}
	return DbGetWalletByUserID(user.ID)
}

func DbGetUserWalletBalanceByNickName(nickname string) (int, error) {
	var user usr.User
	user, err := usr.DbGetUserByUsername(nickname)
	if err != nil {
		return 0, err
	}
	return DbGetUserWalletBalance(user.ID)

}

func DbGetUserWalletBalance(userID primitive.ObjectID) (int, error) {
	var wallet Wallet
	err := walletCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func DbWalletHasBalance(id primitive.ObjectID, amount int) (bool, error) {
	var wallet Wallet
	err := walletCollection.FindOne(context.Background(), bson.M{"_id": id }).Decode(&wallet)
	if err != nil {
		return false, err
	}
	return wallet.Balance >= amount, nil
}