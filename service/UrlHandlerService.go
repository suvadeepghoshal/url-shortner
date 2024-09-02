package service

import (
	"context"
	"errors"
	"log/slog"
	"url-shortner/handlers/util"
	TYPE "url-shortner/model/type"
	"url-shortner/repository"

	"gorm.io/gorm"
)

type UrlService struct {
	Repo repository.Repository
}

// TODO: who should take care of the go routines?
// Ideally the service to repo transaction should be wrapped in a go routine

func (s *UrlService) MakeShortUrl(ctx context.Context, url *TYPE.Url) error {
	return s.Repo.Transaction(ctx, func(repo repository.Repository) error {
		short, e := util.CreateMd5Hash(url.LongUrl)
		if e != nil {
			return e
		}

		// taking first 7 chars of the md5 hash generated
		url.ShortUrl = short[0:7]
		if err := repo.Create(ctx, url); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errors.New("the url already exists")
			} else {
				// logging the generic error for monitoring the app
				slog.Error("MakeShortUrl", "req_id", ctx.Value("req_id"), "err", err)
				return errors.New("unable to create a short url")
			}
			// err.Error(*my)
		}

		// parsing from short hash to URL
		url.ShortUrl = util.ParseShortUrl(url.ShortUrl, ctx.Value("hostname").(string))

		return nil
	})
}

func (s *UrlService) GetLongUrl(ctx context.Context, url *TYPE.Url) error {
	return s.Repo.Transaction(ctx, func(repo repository.Repository) error {
		if err := repo.Query(ctx, &url); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("url not found")
			} else {
				// logging the generic error for monitoring the app
				slog.Error("MakeShortUrl", "req_id", ctx.Value("req_id"), "err", err)
				return errors.New("something went wrong, Please try again later")
			}
		}

		if len(url.LongUrl) == 0 {
			return errors.New("requested url is empty")
		}

		return nil
	})
}
