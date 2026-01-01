package services

import "codee_jun/internal/repos"

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
