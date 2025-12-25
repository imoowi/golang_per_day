package services

import "golang_per_day_24/internal/repos"

type Services func()

var services = []Services{}

func InitServices() {
	repos.InitRepos()
	for _, service := range services {
		service()
	}
}

func RegisterServices(r ...Services) {
	services = append(services, r...)
}
