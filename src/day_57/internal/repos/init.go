package repos

import "golang.org/x/sync/singleflight"

type Repos func()

var repos = []Repos{}

func InitRepos() {
	for _, repo := range repos {
		repo()
	}
}

func RegisterRepos(r ...Repos) {
	repos = append(repos, r...)
}

var (
	sfGroup       singleflight.Group
	cacheNilValue = "{}"
)
