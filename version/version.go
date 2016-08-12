package version

var (
	version   string
	gitCommit string
)

func GetVersion() string {
	return version
}

func GetGitCommit() string {
	return gitCommit
}
