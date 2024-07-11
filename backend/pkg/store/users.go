package store

import (
	"backend/pkg/model"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserStore handles database operations for users.
type UserStore struct {
	db *mongo.Collection
}

// NewUserStore creates a new instance of UserStore.
func NewUserStore(client *mongo.Client, dbName, collectionName string) *UserStore {
	db := client.Database(dbName).Collection(collectionName)
	return &UserStore{db: db}
}

// CreateUser adds a new user to the database.
func (store *UserStore) CreateUser(ctx context.Context, user model.Users) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Status = "Pending"
	user.LastUpdatedAt = fmt.Sprintf("%v", time.Now())

	// Insert user into db
	_, err = store.db.InsertOne(ctx, user)
	return err
}

// FindUserByEmail retrieves a user by their email.
func (store *UserStore) FindUserByEmail(ctx context.Context, email string) (*model.Users, error) {
	var user model.Users
	filter := bson.M{"email": email}
	err := store.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// update's user session
func (store *UserStore) UpdateSession(ctx context.Context, userID string, sess string) error {
	filter := bson.M{"email": userID}
	update := bson.M{"$set": bson.M{"session": sess}}
	_, err := store.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// fetch users all
func (store *UserStore) FetchUsers(ctx context.Context) ([]model.Users, error) {
	var users []model.Users
	cursor, err := store.db.Find(ctx, bson.D{{}})
	if err != nil {
		log.Println("error in find users ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.Users
		if err := cursor.Decode(&user); err != nil {
			log.Println("error in decode users  ", err)
			return nil, err
		}
		time := strings.Split(user.LastUpdatedAt, ".")[0]
		user.LastUpdatedAt = time
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Println("cursor error  ", err)
		return nil, err
	}
	return users, nil
}

// update status of user's ACTIVE or INACTIVE
func (store *UserStore) UpdateStatus(ctx context.Context, userID string, status string) error {
	filter := bson.M{"email": userID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := store.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
