package dreck

import (
	"fmt"
	"strings"

	"github.com/miekg/dreck/types"
)

func (d Dreck) isAutosubmit(req types.PullRequestOuter, conf *types.DreckConfig) (bool, error) {
	client, ctx, err := d.newClient(req.Installation.ID)
	if err != nil {
		return false, err
	}

	pull, _, err := client.PullRequests.Get(ctx, req.Repository.Owner.Login, req.Repository.Name, req.PullRequest.Number)
	if err != nil {
		return false, fmt.Errorf("getting PR %d: %s", req.PullRequest.Number, err.Error())
	}

	permitted := permittedUserFeature(featureAutosubmit, conf, pull.User.GetLogin())
	if !permitted {
		return false, nil
	}

	return isautosubmit(pull.GetBody()), nil
}

func isautosubmit(msg string) bool { return strings.Contains(msg, Trigger+autosubmitConst) }

// PullRequestAutosubmit will kick off autosubmit, by calling d.autosubmit.
func (d Dreck) pullRequestAutosubmit(req types.PullRequestOuter) error {
	reqComment := types.PullRequestToIssueComment(req)

	return d.autosubmit(reqComment)
}
