# REST API using Go

## Overview

This is just a simple RESTful API written in Golang using Go-Fiber Framework. One of the main goals was to make this production ready and so it is.

All the secrets like database username/password being passed using [docker secrets](https://docs.docker.com/engine/swarm/secrets/). Those are working with Docker Swarm only. So to run this code, it is important to create a simple Docker Swarm cluster. I'll tell you how to do it later on.

## Table of Contents

- [Tech Stack](#techstack)
- [Installation](#instalation)
  - [Creating VMs for Docker Swarm](#creating-vms-for-docker-swarm)
  - [Creating Docker Swarm Cluster](#creating-docker-swarm-cluster)
  - [Add docker secrets](#add-docker-secrets)
  - [Cloning repo and building an image](#cloning-repo-and-building-an-image)
- [Usage](#usage)

## TechStack

- go version 1.22.2 linux/amd64
- Docker version 27.0.3
- Docker Compose version v2.3.3
- Docker Swarm
- Multipass 1.13.1
- PostgreSQL 15.6
- GNU Make 4.3

## Instalation

### Creating VMs for Docker Swarm

In order to create a Docker Swarm cluster, we should create two VMs for that. One VM is a manager and another one is a worker.

This repository has two script files with neccessary commands to run to create two VMs. You can just execute them:

```shell
sh -x init-instance.sh manager
sh -x init-instance.sh worker
```

## Creating Docker Swarm Cluster

Now we can create a Docker Swarm cluster with our VMs.

To initialize the Docker Swarm on manager node execute next command:

```shell
multipass exec manager -- docker swarm init
```

You will get this message:

```shell
Swarm initialized: current node (793322xn7el4jb0ujnkqgeqqg) is now a manager.

To add a worker to this swarm, run the following command:

    docker swarm join --token {TOKEN IS HERE} 10.223.24.185:2377

To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
```

Execute this command on your worker node like here:

```shell
 multipass exec worker -- docker swarm join --token {TOKEN IS HERE} 10.223.24.185:2377
```

You will get this message:

```shell
This node joined a swarm as a worker.
```

## Add docker secrets

To add docker secrets we should enter our manager VM with the next command

```shell
multipass shell manager
```

Now we should add our docker secrets `pg_pass` and `pg_user` that will be used for our database

```shell
echo "username" | docker secret create pg_user
echo "password" | docker secret create pg_pass
```

You can check this secrets by runing

```shell
docker secret ls
```

## Cloning repo and building an image

Now we can clone our repository. In the home directory execute this command

```shell
git clone https://github.com/sshaparenko/restApiOnGo.git
```

In the project directory execute the next comand

```shell
docker stack deploy -c docker-compose.yml my_stack
```

This will create two docker services. One for application itself and one for the database. You can check status of those services with next comand:

```shell
docker service ls
```

```shell
ID             NAME                MODE         REPLICAS IMAGE                             PORTS
1r44gi1367a6   my_stack_app        replicated   1/1        sshaparenkos/restapiongo:latest   *:8080->8080/tcp
rr6rs3nc8oa7   my_stack_postgres   replicated   1/1        bitnami/postgresql:15             *:5432->5432/tcp
```

When all services are up and runing you can send requests to the app itself. For that you should know IP address of the worker node. You can get it with the next command:

```shell
multipass ls
```

Now you are ready to go

## Usage

API has seven endopints:

- POST /signup
- POST /login
- GET, POST /items
- GET, PUT, DELETE /items/:id

### POST `/signup`

This endpoint will register a user with specified `email` and `password`. Returns JWT token in responce.

```shell
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"email": "someemail@gmail.com, "password": "1231234"}' \
	http://<WORKER_IP>:8080/api/v1/signup
```

### POST `/login`

This endpoint is for login a user with specified `email` and `password`. Returns JWT toke in respose.

```shell
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"email": "someemail@gmail.com, "password": "1231234"}' \
	http://<WORKER_IP>:8080/api/v1/login
```

### GET `/items`

This endpoint will return a list of all items.

```shell
curl -X GET \
	http://<WORKER_IP>:8080/api/v1/items
```

### POST `/items`

This endpint will create new item in a database. Authorization is required.

```shell
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "item_name", "price": 100, "quantity": 10}' \
  http://<WORKER_IP>:8080/api/v1/items
```

### GET `/items/:id`

This endpoint will return information about item by id.

```shell
curl -X GET \
	http://<WORKER_IP>:8080/api/v1/items:1
```

### PUT `/items/:id`

Will update item data by id.

```shell
curl -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "item_name", "price": 100,"quantity": 10}' \
  http://<WORKER_IP>:8080/api/v1/items/1
```

### DELETE `/items/:id`

Will remove item by id.

```shell
curl -X DELETE \
	-H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://<WORKER_IP>:8080/api/v1/items/1
```
