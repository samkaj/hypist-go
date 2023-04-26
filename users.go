package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollectionName = "users"

func InsertUser(ctx context.Context, db *mongo.Database, user *User) (*User, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	collection := db.Collection(usersCollectionName)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	fmt.Printf("[hypist] inserted user with id %s\n", result.InsertedID)
	return user, nil
}

func DeleteUser(ctx context.Context, db *mongo.Database, name string) error {
	collection := db.Collection(usersCollectionName)
	filter := bson.D{{"name", name}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete user: %v", err))
	}

	fmt.Printf("[hypist] deleted user: %s\n", name)
	return nil
}

func (u *User) validate() error {
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !strings.ContainsAny(u.Password, "0123456789") {
		return errors.New("password must contain at least one special character")
	}

	if !strings.ContainsAny(u.Password, "!@#$%^&*()_+-=[]{}\\|;:'\",.<>/?") {
		return errors.New("password must contain at least one special character")
	}

	if strings.ToLower(u.Password) == u.Password {
		return errors.New("password must contain at least one uppercase character")
	}

	if len(u.Name) > 32 || len(u.Name) < 4 {
		return errors.New("name must be between 4 and 32 characters")
	}

	return nil
}
