package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	TYPE "url-shortner/model/type"
	"url-shortner/repository"
)

type UrlService struct {
	Repo repository.Repository
}

// TODO: who should take care of the go routines?
// Ideally the service to repo transaction should be wrapped in a go routine

func (s *UrlService) MakeShortUrl(ctx context.Context) error {
	return s.Repo.Transaction(ctx, func(repo repository.Repository) error {
		return nil // business logic goes here
	})
}

func (s *UrlService) GetLongUrl(ctx context.Context, url *TYPE.Url) error {
	return s.Repo.Transaction(ctx, func(repo repository.Repository) error {
		if err := repo.Query(ctx, &url); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("url not found")
			} else {
				return errors.New("something went wrong, Please try again later")
			}
		}

		if len(url.LongUrl) == 0 {
			return errors.New("requested url is empty")
		}

		return nil
	})
}
