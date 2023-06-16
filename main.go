package main

import (
	"fmt"
	"github.com/alexey-ash/gitlab-spendtime-calc/pkg/gitlab_issues"
	"github.com/olekukonko/tablewriter"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const gitlabURL = ""
const gitlabToken = ""
const projectID = ""
const userName = ""

func main() {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Task", "Hours"})
	table.SetAutoWrapText(false)

	client, err := gitlab_issues.GitlabAuth(gitlabToken, gitlabURL)
	if err != nil {
		log.Fatal(err)
	}

	milTitle, _, err := gitlab_issues.GetActiveMilestone(client, projectID)
	if err != nil {
		log.Fatal(err)
	}

	issues, err := gitlab_issues.GetMilestoneIssues(client, projectID, milTitle, userName)
	if err != nil {
		log.Fatal(err)
	}

	totalHours := 0

	for _, i := range issues {
		taskHours := 0
		note, _, err := client.Notes.ListIssueNotes(projectID, i.IID, &gitlab.ListIssueNotesOptions{
			ListOptions: gitlab.ListOptions{},
			OrderBy:     gitlab.String("updated_at"),
			Sort:        gitlab.String("desc")})
		if err != nil {
			log.Fatal(err)
		}
		// Смотрим каждую заметку
		for _, n := range note {

			if n.UpdatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") && strings.Contains(n.Body, "time spent") {
				re := regexp.MustCompile(`(\d+)h`)
				match := re.FindStringSubmatch(n.Body)
				if len(match) > 1 {
					hour, err := strconv.Atoi(match[1])
					if err != nil {
						log.Fatal(err)
					}
					totalHours += hour
					taskHours += hour
				}
			}
		}
		if taskHours > 0 {
			table.Append([]string{i.Title, strconv.Itoa(taskHours)})
		}
	}
	table.Render()
	fmt.Printf("Total spend time today: %dh\n", totalHours)
}
