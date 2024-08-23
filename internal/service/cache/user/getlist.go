package user

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"runtime"
	"sync"

	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/pagination"
)

func (s *service) GetList(ctx context.Context, limit pagination.Limit) ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, err := redis.Int64s(s.redisClient.LRange(ctx, buildListCacheKeyByLimit(limit), 0, -1))
	if err != nil {
		return nil, err
	}

	idsCh := make(chan int64)
	go func() {
		for _, id := range ids {
			idsCh <- id
		}
		close(idsCh)
	}()

	usersCh := make(chan *model.User, len(ids))
	errCh := make(chan error, len(ids))

	var wg sync.WaitGroup
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idsCh {
				user, getErr := s.userRepository.Get(ctx, id)
				if getErr != nil {
					errCh <- getErr
					continue
				}

				if user == nil {
					errCh <- model.ErrorNoteNotFound
					continue
				}

				usersCh <- user
			}
		}()
	}

	wg.Wait()
	close(errCh)
	close(usersCh)

	for err = range errCh {
		if err != nil {
			return nil, err
		}
	}

	userMap := make(map[int64]model.User, len(ids))
	for u := range usersCh {
		userMap[u.ID] = *u
	}

	users := make([]model.User, 0, len(ids))
	for _, id := range ids {
		users = append(users, userMap[id])
	}

	return users, nil
}
