package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/workshops/wallet/internal/repository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	Conn *mongo.Client
}

var ctx = context.TODO()

func NewRepository(conn *mongo.Client) *Repository {
	return &Repository{Conn: conn}
}

func NewMongoDB(dsn string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return client, nil
}

func (r *Repository) CreateUser(token string) error {
	collection := r.Conn.Database("wallet").Collection("users")
	user := &models.User{
		ID:    primitive.NewObjectID().String(),
		Token: &token,
	}

	_, err := collection.InsertOne(ctx, user)

	return err
}

func (r *Repository) GetUsers() ([]*models.User, error) {
	collection := r.Conn.Database("wallet").Collection("users")
	users := make([]*models.User, 0)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		user := new(models.User)
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return users, nil
}

func (r *Repository) CreateWallet(wallet *models.Wallet) error {
	collection := r.Conn.Database("wallet").Collection("wallets")

	wallet.ID = primitive.NewObjectID().String()

	_, err := collection.InsertOne(ctx, wallet)

	return err
}

func (r *Repository) GetWalletByID(id string) (*models.Wallet, error) {
	collection := r.Conn.Database("wallet").Collection("wallets")

	wallet := new(models.Wallet)

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&wallet)
	if err != nil {
		return nil, errors.Wrap(err, "Error from db")
	}

	return wallet, nil
}

func (r *Repository) GetWalletTransactionsByID(id string) ([]*models.Transaction, error) {
	collection := r.Conn.Database("wallet").Collection("transactions")

	transactions := make([]*models.Transaction, 0)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		transaction := new(models.Transaction)
		err := cur.Decode(&transaction)
		if err != nil {
			log.Fatal(err)
		}

		transactions = append(transactions, transaction)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return transactions, nil
}

func (r *Repository) GetTransactions() ([]*models.Transaction, error) {
	collection := r.Conn.Database("wallet").Collection("transactions")
	transactions := make([]*models.Transaction, 0)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		transaction := new(models.Transaction)
		err := cur.Decode(&transaction)
		if err != nil {
			log.Fatal(err)
		}

		transactions = append(transactions, transaction)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return transactions, nil
}

func (r *Repository) CreateTransaction(transaction *models.Transaction) error {
	collectionWallet := r.Conn.Database("wallet").Collection("wallets")
	collectionTransactions := r.Conn.Database("wallet").Collection("transactions")

	//var session mongo.Session

	transaction.ID = primitive.NewObjectID().String()
	t := time.Now()
	transaction.Date = t.String()

	_, err := collectionWallet.UpdateOne(ctx, bson.M{"_id": transaction.CreditWalletID}, bson.D{{"$inc", bson.D{{"balance", -transaction.Amount}}}}, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	_, err = collectionWallet.UpdateOne(ctx, bson.M{"_id": transaction.DebitWalletID}, bson.D{{"$inc", bson.D{{"balance", transaction.Amount}}}}, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	_, err = collectionTransactions.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}

	//session, err := r.Conn.StartSession()
	//if err != nil {
	//	return err
	//}
	//
	//defer session.EndSession(ctx)
	//
	//err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
	//	err = session.StartTransaction()
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = collectionWallet.UpdateOne(sc, bson.M{"_id": transaction.CreditWalletID}, bson.D{{"$inc", bson.D{{"balance", -transaction.Amount}}}}, options.Update().SetUpsert(false))
	//	if err != nil {
	//		session.AbortTransaction(sc)
	//		return err
	//	}
	//
	//	_, err = collectionWallet.UpdateOne(sc, bson.M{"_id": transaction.DebitWalletID}, bson.D{{"$inc", bson.D{{"balance", transaction.Amount}}}}, options.Update().SetUpsert(false))
	//	if err != nil {
	//		session.AbortTransaction(sc)
	//		return err
	//	}
	//
	//	_, err = collectionTransactions.InsertOne(ctx, transaction)
	//	if err != nil {
	//		session.AbortTransaction(sc)
	//		return err
	//	}
	//
	//	if err = session.CommitTransaction(sc); err != nil {
	//		session.AbortTransaction(sc)
	//		return err
	//	}
	//
	//	return nil
	//})

	return nil
}

func (r *Repository) GetWalletAmountDayByID(id string, week models.Week) ([]*models.Day, error) {
	collectionTransactions := r.Conn.Database("wallet").Collection("transactions")

	matchStage := bson.D{{"$match", bson.D{{"debitwalletid", id}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"$dateFromString", bson.D{
		{"dateString", "$b"},
		{"format", "%d/%m/%Y"}}}}}, {"total", bson.D{{"$sum", "$amount"}}}}}}

	showInfoCursor, err := collectionTransactions.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		panic(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println(showsWithInfo)
	return nil, nil
}

func (r *Repository) GetWalletAmountWeekByID(id string, week models.Week) ([]*models.Day, error) {
	return nil, nil
}
