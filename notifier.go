package ukigumo

type Notifier interface {
	send(c Client, status string, lastStatus string, reportURL string, currentRev string, reposOwner string, reposName string) error
}
