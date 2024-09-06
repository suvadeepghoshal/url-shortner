package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
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

func (s *UrlService) MakeShortUrl(c context.Context, url *TYPE.Url) error {
	ch := make(chan error, 1)
	ctx, cancel := context.WithTimeout(c, time.Millisecond*100) // 100 mili-seconds to perform the DB transaction
	defer cancel()

	start := time.Now()

	go s.Repo.Transaction(ctx, ch, func(repo repository.Repository) error {
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

	select {
	case <-ctx.Done():
		return fmt.Errorf("the connection timed out as it took %s", time.Since(start))
	case result := <-ch:
		if result != nil {
			return result
		}
	}
	return nil
}

func (s *UrlService) GetLongUrl(c context.Context, url *TYPE.Url) error {
	ch := make(chan error, 1)
	ctx, cancel := context.WithTimeout(c, time.Millisecond*100) // 100 mili-seconds to perform the DB transaction
	defer cancel()

	// start := time.Now()

	go s.Repo.Transaction(ctx, ch, func(repo repository.Repository) error {
		time.Sleep(time.Second * 2)

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

	/* select {
	case <-ctx.Done():
		return fmt.Errorf("the connection timed out as it took %s", time.Since(start))
	case result := <-ch:
		if result != nil {
			return result
		}
	} */

	handleChannelResponse(ctx, ch)

	return nil
}

// TODO: need to find out how the deadline works
func handleChannelResponse(ctx context.Context, ch chan error) error {
	select {
	case <-ctx.Done():
		deadline, ok := ctx.Deadline()
		slog.Debug("handleChannelResponse", "ok", ok)
		return fmt.Errorf("the connection timed out as it took %s", deadline)
	case result := <-ch:
		if result != nil {
			return result
		}
	}
	return nil
}
