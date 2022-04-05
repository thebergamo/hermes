build: build-notifications build-datasets build-campaings

build-campaings:
	env GOOS=linux go build -ldflags="-s -w" -o main cmd/campaings/main.go
	mkdir -p bin/campaings
	zip bin/campaings/main.zip main
	mv main bin/campaings

build-datasets:
	env GOOS=linux go build -ldflags="-s -w" -o main cmd/datasets/main.go
	mkdir -p bin/datasets
	zip bin/datasets/main.zip main
	mv main bin/datasets

build-notifications:
	env GOOS=linux go build -ldflags="-s -w" -o main cmd/notifications/main.go
	mkdir -p bin/notifications
	zip bin/notifications/main.zip main
	mv main bin/notifications


start-api:
	sam local start-api -t sam.yaml --skip-pull-image --warm-containers EAGER --parameter-overrides dockerhost=host.docker.internal

create-tables: create-dataset-table create-notification-table

create-campaing-table: 
	aws dynamodb create-table --table-name campaing --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://localhost:4566

create-notification-table: 
	aws dynamodb create-table --table-name notification --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://localhost:4566

create-dataset-table: 
	aws dynamodb create-table --table-name datasets --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://localhost:4566