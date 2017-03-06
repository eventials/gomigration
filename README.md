# migration
Migration lib for Golang projects.

## Preparing the environment

Install Docker:

* [Install steps](https://docs.docker.com/engine/installation/)

## Compiling and running the app

To compile and run the application, just run:

```
docker-compose up
```

## Docker Images

### Requirements

- [Docker](https://docs.docker.com/engine/installation/)
- [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/installing.html#install-with-pip)

### Docker login

Login to AWS ECR (Docker Registry):

```sh
./images.sh configure
```

### Running Tests

```sh
./images.sh test
```

### Building Images

```sh
./images.sh build
```

### Pushing images

```sh
./images.sh push
```

## Deploy

### Requirements

- [ecs-deploy-py](https://github.com/gfronza/ecs-deploy.py)
- [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/installing.html#install-with-pip)

### Staging

```sh
./images.sh deploy staging
```

### Production

```sh
./images.sh deploy production
```
