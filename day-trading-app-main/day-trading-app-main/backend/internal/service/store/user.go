package store

import (
	"context"
	"day-trading-app/backend/internal/service/models"

	"go.mongodb.org/mongo-driver/bson"
)

// Tested
func (mh *MongoHandler) RegisterUser(userName, password, name string) error {
	//create user in db
	collection := mh.client.Database("day-trading-app").Collection("users")
	// Insert the user into the database
	var user models.User = models.User{
		UserName:     userName,
		PasswordHash: password,
		Name:         name,
		Balance:      0,
		Stocks:       []models.PortfolioItem{},
		WalletTxns:   []models.WalletTransaction{},
	}
	_, err := collection.InsertOne(context.Background(), user)
	//_, err := collection.InsertOne(context.Background(), models.User{UserName: userName, PasswordHash: password, Name: name})
	if err != nil {
		return err
	}
	return nil
}

// Tested
func (mh *MongoHandler) GetUserByUserName(userName string) (models.User, error) {
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	//Test USE ONLY
	//fmt.Println(user)
	return user, nil
}

// Tested
func (mh *MongoHandler) GetWalletTransactions(userName string) ([]models.WalletTransaction, error) {
	//For testing purposes only:
	//userName = "VanguardETF"
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("wallet_transactions")

	// return every transaction in the wallet_transactions collection
	cursor, err := collection.Find(context.Background(), bson.M{"user_name": userName})
	if err != nil {
		return nil, err
	}
	var walletTransactions []models.WalletTransaction
	if err = cursor.All(context.Background(), &walletTransactions); err != nil {
		return nil, err
	}
	return walletTransactions, nil
}

// Tested
func (mh *MongoHandler) GetWalletBalance(userName string) (int, error) {
	//For testing purposes only:
	//userName = "VanguardETF"
	//access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	//find the user by their username
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

// Tested
func (mh *MongoHandler) SetWalletBalance(userName string, newBalance int) error {
	//For testing purposes only:
	//userName = "VanguardETF"
	//newBalance = 100000
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return err
	}

	// Update the user's balance
	_, err = collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$set": bson.M{"balance": newBalance}})
	if err != nil {
		return err
	}
	return nil
}

// Tested
func (mh *MongoHandler) ManageUserWalletBalance(userName string, amountToAdd int) error {
	// Access the collection where user data is stored
	collection := mh.client.Database("day-trading-app").Collection("users")

	// Define the update operation to increment the existing balance
	update := bson.M{"$inc": bson.M{"balance": amountToAdd}}

	// Perform the update operation atomically and get the updated document
	var updatedUser models.User
	err := collection.FindOneAndUpdate(context.Background(), bson.M{"user_name": userName}, update).Decode(&updatedUser)
	if err != nil {
		return err
	}

	// Return the updated balance
	return nil
}

// Tested.
func (mh *MongoHandler) AddWalletTransaction(userName string, walletTxID string, stockTxID string, is_debit bool, amount int, timeStamp int64) error {

	var walletTx models.WalletTransaction = models.WalletTransaction{
		UserName:   userName,
		WalletTxID: walletTxID,
		StockTxID:  stockTxID,
		Is_debit:   is_debit,
		Amount:     amount,
		TimeStamp:  timeStamp,
	}

	// Add to 'wallet_transactions' collection
	collection := mh.client.Database("day-trading-app").Collection("wallet_transactions")
	_, err := collection.InsertOne(context.Background(), walletTx)
	if err != nil {
		return err
	}

	// * Add to user's entry in 'users' collection *

	// Access the collection where user data is stored
	collection = mh.client.Database("day-trading-app").Collection("users")

	// Find the user by their username
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		return err
	}

	// update the user's wallet_txns
	_, err = collection.UpdateOne(context.Background(), bson.M{"user_name": user.UserName}, bson.M{"$push": bson.M{"wallet_txns": walletTx}})
	if err != nil {
		return err
	}
	return nil
}

// Tested.
func (mh *MongoHandler) DeleteWalletTransaction(userName string, walletTxID string) error {
	// Remove from 'wallet_transactions' collection using the walletTxID
	collection := mh.client.Database("day-trading-app").Collection("wallet_transactions")
	_, err := collection.DeleteOne(context.Background(), bson.M{"wallet_tx_id": walletTxID})
	if err != nil {
		return err
	}

	// Remove from user's entry in 'users' collection

	// Access the collection where user data is stored
	collection = mh.client.Database("day-trading-app").Collection("users")

	// Remove the transaction
	_, err = collection.UpdateOne(context.Background(), bson.M{"user_name": userName}, bson.M{"$pull": bson.M{"wallet_txns": bson.M{"wallet_tx_id": walletTxID}}})
	if err != nil {
		return err
	}
	return nil
}
