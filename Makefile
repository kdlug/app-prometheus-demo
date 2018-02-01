.PHONY: install test build serve clean pack deploy ship
OWNER=casi
SERVICE=app-prometheus-demo
# it will take short hash (first 7 symbols) of last git commit.
# Then, we export this variable, so it’s available in commands run by make.
TAG=$(git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

install:
	go get .

test:
	go test ./...

build: install
	# build a binary
	go build -ldflags "-X main.version=$(TAG)" -o ${SERVICE} .

serve: build
	./${SERVICE}

clean:
	# remove binary file
	rm ./${SERVICE}

pack:
	# build docker image
	GOOS=linux make build
	docker build -t ${OWNER}/${SERVICE}:$(TAG) .

upload:
	# push docker image to registry
	gcloud docker -- push docker pull ${OWNER}/${SERVICE}:$(TAG)

deploy:
	# The envsubst program substitutes the values of environment variables.
	envsubst < k8s/deployment.yml | kubectl apply -f -

ship: test pack upload deploy clean
