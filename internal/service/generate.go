package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i AuditLogService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i HashService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i CacheUserService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i ValidatorService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i ProducerService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i AuthConfig -o ./mocks/ -s "_minimock.go"
