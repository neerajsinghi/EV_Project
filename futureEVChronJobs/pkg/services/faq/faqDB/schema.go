package faqdb

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/faq"
)

var repo = faq.NewRepository("faq")

type FaqI interface {
	AddFaq(faq entity.FAQDB) (string, error)
	UpdateFaq(id string, faq entity.FAQDB) (string, error)
	DeleteFaq(id string) error
	GetAllFaq() ([]entity.FAQDB, error)
}
type service struct{}

func NewService() FaqI {
	return &service{}
}
