# set env value NAME
NAME = app
# set env value VERSION
VERSION = "1.0.0"

# set heroku token (already set in computer env)
# HEROKU_TOKEN = ""

# set heroku email
# HEROKU_EMAIL = ""

# set the app name
APP_NAME = "cfg-cendit"

# set the image sha
IMAGE_SHA = ""

all:
	make graph
	make build
	make start

build:
	go build -o bin/$(NAME) -ldflags "-X main.Version=$(VERSION)" ./cmd/$(NAME).go

run:
	make fmt
	make build
	./bin/$(NAME)

start:
	./bin/$(NAME)

# regenerate the graphql schema and resolver files
graph:
	go run github.com/99designs/gqlgen generate

# format the code
fmt:
	go fmt ./...


# this makes use of ~/.netrc file to login to heroku
# you can generate one by running the following command
# in your terminal (of course you have to have heroku cli installed and also substitute the HEROKU_EMAIL and HEROKU_TOKEN with your own values)

# export HEROKU_EMAIL="your heroku email"
# export HEROKU_TOKEN="your heroku token"

# cat >~/.netrc <<EOF
# machine api.heroku.com
#     login ${HEROKU_EMAIL}
#     password ${HEROKU_TOKEN}
# machine git.heroku.com
#     login ${HEROKU_EMAIL}
#     password ${HEROKU_TOKEN}
# EOF

deploy:
	heroku container:login && \
	heroku container:push web -a $(APP_NAME) && \
	heroku container:release web -a $(APP_NAME)

# sha:
# 	export SHASIZE=$$(export SHA=$$(cat imagesha.txt) && echo $$(SHA#*latest: digest: )) && echo $$(SHASIZE:0:71) > ./imageshaa.txt

# deploying:
# 	docker build -f Dockerfile --iidfile imageid.txt -t registry.heroku.com/cfg-cendit/web . && \
# 	docker login -u _ -p $(HEROKU_TOKEN) registry.heroku.com && \
# 	docker push registry.heroku.com/cfg-cendit/web
	
# deployed:
# 	curl -X PATCH https://api.heroku.com/apps/cfg-cendit/formation --header "Content-Type: application/json" --header "Accept: application/vnd.heroku+json; version=3.docker-releases" --header "Authorization: Bearer $(HEROKU_TOKEN)" --data '{ "updates": [ { "type": "web", "docker_image": $(IMAGE_SHA) } ] }'


# echo "$$(docker push registry.heroku.com/cfg-cendit/web)" > ./imagesha.txt
# export SHASIZE=$$(export SHA=$$(cat imagesha.txt) && echo $${SHA#*latest: digest: }) && echo $${SHASIZE:0:71} > ./imagesha.txt && \
# 	echo "Docker Image ID is :: $$(cat imagesha.txt)" && \

# deploy:
# 	docker build -f Dockerfile --iidfile imageid.txt -t registry.heroku.com/cfg-cendit/cendit . && \
# 	docker login -u _ -p $(HEROKU_TOKEN) registry.heroku.com && \
# 	echo "Docker Image ID is $$(cat imageid.txt)" && \
# 	heroku container:push web -a $(APP_NAME) && \
# 	heroku container:release web -a $(APP_NAME)

# cat >~/.netrc <<EOF
# machine api.heroku.com
#     login ${HEROKU_EMAIL}
#     password ${HEROKU_TOKEN}
# machine git.heroku.com
#     login ${HEROKU_EMAIL}
#     password ${HEROKU_TOKEN}
# EOF
