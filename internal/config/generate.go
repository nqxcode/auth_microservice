package config

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i AuthConfig -o ./mocks/ -s "_minimock.go"
