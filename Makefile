build:
	env GOOS=linux go build -ldflags="-s -w" -o main cmd/datasets/main.go
	mkdir -p bin/datasets
	zip bin/datasets/main.zip main
	mv main bin/datasets

start-api:
	sam local start-api -t sam.yaml --skip-pull-image --warm-containers EAGER --parameter-overrides dockerhost=host.docker.internal

create-tables: create-dataset-table

create-dataset-table: 
	aws dynamodb create-table --table-name datasets --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://localhost:4566