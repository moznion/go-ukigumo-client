package ukigumo

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type IkasanNotifier struct {
	URL           string
	Channel       string
	IgnoreSuccess bool
	IgnoreSkip    bool
	Method        string
}

func NewIkasanNotifier(url string, channel string) *IkasanNotifier {
	notifier := new(IkasanNotifier)

	notifier.URL = url
	notifier.Channel = channel
	notifier.IgnoreSuccess = true
	notifier.IgnoreSkip = true
	notifier.Method = "notice"

	return notifier
}

func (notifier *IkasanNotifier) send(c Client, status string, lastStatus string, reportURL string, currentRev string, reposOwner string, reposName string) error {
	ignoreSkip := notifier.IgnoreSkip

	if notifier.IgnoreSuccess && status != "" && lastStatus != "" && (lastStatus == StatusSuccess || lastStatus == StatusSkip) {
		// Don't notify if status represents successful and last status is not failed
		log.Printf("The test was succeeded. There is no reason to notify(%s, %s).", status, lastStatus)
		return nil
	}

	if ignoreSkip && status == StatusSkip {
		// Don't notify if status represents skipping
		log.Printf("The test was skiped. There is no reason to notify.")
		return nil
	}

	ikasanURL := fmt.Sprintf("%s/%s", strings.TrimRight(notifier.URL, "/"), notifier.Method)

	// TODO support to show status and colorize
	message := fmt.Sprintf("%s %s [%s] %s %s", reportURL, c.Project, c.VC.GetVCInfo().Branch, status, currentRev[:10])
	log.Printf("Sending message to irc server: %s", message)

	res, err := http.PostForm(ikasanURL, url.Values{
		"channel": {notifier.Channel},
		"message": {message},
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	resStatusCode := res.StatusCode
	if !(200 <= resStatusCode && resStatusCode <= 299) {
		errMsg := fmt.Sprintf("Cannot send ikachan notification: %s %s %s %d",
			notifier.Method, notifier.URL, notifier.Channel, resStatusCode)
		err := errors.New(errMsg)
		log.Fatal(err)
		return err
	}

	log.Printf("Sent notification for %s @ %s", notifier.URL, notifier.Channel)
	return nil
}
