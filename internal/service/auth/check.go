package auth

import "context"

func (s *service) Check(ctx context.Context, accessToken, endpointAddress string) (bool, error) {
	return false, nil
}
