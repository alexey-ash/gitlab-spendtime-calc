package gitlab_issues

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"time"
)

const ISO8601 = "2006-01-02"

func GitlabAuth(gitlabAPIToken string, gitlabAPIUrl string) (*gitlab.Client, error) {
	return gitlab.NewClient(gitlabAPIToken, gitlab.WithBaseURL(gitlabAPIUrl))
}

func GetActiveMilestone(git *gitlab.Client, projectID string) (milestoneTitle string, statusCode int, err error) {
	currentDate := time.Now().UTC()

	milestones, resp, err := git.Milestones.ListMilestones(projectID, &gitlab.ListMilestonesOptions{State: gitlab.String("active")})
	if err != nil {
		return "", resp.StatusCode, err
	}

	for _, m := range milestones {
		startDate, err := time.Parse(ISO8601, m.StartDate.String())
		if err != nil {
			log.Fatal(err)
		}
		endDate, err := time.Parse(ISO8601, m.DueDate.String())

		if currentDate.After(startDate) && currentDate.Before(endDate) {
			return m.Title, resp.StatusCode, nil
		}
	}
	return "", 0, nil
}

func GetMilestoneIssues(git *gitlab.Client, projectID string, milestoneTitle string, username string) (issues []*gitlab.Issue, err error) {
	opt := &gitlab.ListProjectIssuesOptions{
		AssigneeUsername: gitlab.String(username),
		Milestone:        gitlab.String(milestoneTitle),
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	for {
		// Get issues from milestone
		issuesBatch, resp, err := git.Issues.ListProjectIssues(projectID, opt)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issuesBatch...)

		// Pagination
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return issues, nil
}
