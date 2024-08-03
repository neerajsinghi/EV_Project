package faqdb

import (
	"bikeRental/pkg/entity"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (*service) AddFaq(faq entity.FAQDB) (string, error) {
	faq.ID = primitive.NewObjectID()
	faq.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	data, err := repo.InsertOne(faq)
	if err != nil {
		return "", errors.New("error in inserting faq")
	}
	return data, nil
}

func (*service) UpdateFaq(id string, faq entity.FAQDB) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	set := bson.M{}

	if faq.Question != "" {
		set["question"] = faq.Question
	}
	if faq.Answer != "" {
		set["answer"] = faq.Answer
	}
	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	data, err := repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
	if err != nil {
		return "", errors.New("error in updating faq")
	}
	return data, nil
}

func (*service) DeleteFaq(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	err := repo.DeleteOne(bson.M{"_id": idObject})
	if err != nil {
		return errors.New("error in deleting faq")
	}
	return nil
}

func (*service) GetAllFaq() ([]entity.FAQDB, error) {
	data, err := repo.Find(bson.M{}, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding faq")
	}
	return data, nil
}
