package archive

type GithubArchive struct {
	AccessToken string
}

func (a *GithubArchive) CreateArchive(u string) (location string, err error) {
	return "", nil
}
