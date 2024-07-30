package faqdb

import (
	"futureEVChronJobs/pkg/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (*service) AddFaq(faq entity.FAQDB) (string, error) {
	faq.ID = primitive.NewObjectID()
	faq.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	return repo.InsertOne(faq)
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

	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
}

func (*service) DeleteFaq(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

func (*service) GetAllFaq() ([]entity.FAQDB, error) {
	return repo.Find(bson.M{}, bson.M{})
}
