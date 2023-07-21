package main

import (
	"flag"
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

// getTime Парсинг времени из строки
func getTime(body string) (hours int, minutes int, err error) {
	matches := regexp.MustCompile(`(\d+)(h|m)`).FindAllStringSubmatch(body, -1)

	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, 0, err
		}

		if match[2] == "h" {
			hours += value
		} else if match[2] == "m" {
			minutes += value
		}
	}

	return hours, minutes, nil
}

// getDateFromNote Парсинг даты из строки
func getDateFromNote(body string) (*time.Time, error) {
	regex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)

	match := regex.FindStringSubmatch(body)
	if len(match) == 0 {
		return nil, fmt.Errorf("дата не найдена в строке")
	}

	searchDate, err := time.Parse("2006-01-02", match[1])
	if err != nil {
		return nil, err
	}

	return &searchDate, nil
}

func main() {
	// Получение текущей даты и форматирование в строку
	currentDate := time.Now().Format("2006-01-02")

	// Определение флага -date с значением по умолчанию
	dateFlag := flag.String("date", currentDate, "Specify a date in the format 'yyyy-mm-dd'")
	flag.Parse()

	// Создание переменной даты из флага -date
	var searchDate time.Time

	// Парсинг значения аргумента -date
	parsedDate, err := time.Parse("2006-01-02", *dateFlag)
	if err != nil {
		fmt.Println("Date parsing error:", err)
		return
	}
	searchDate = parsedDate
	fmt.Println(searchDate)

	// Создание новой таблицы для вывода в stdout
	table := tablewriter.NewWriter(os.Stdout)
	// Установка заголовка таблицы
	table.SetHeader([]string{"Task", "Time"})
	// Отключение переноса текста в ячейках таблицы
	table.SetAutoWrapText(false)

	// Создание Gitlab API клиента
	client, err := gitlab_issues.GitlabAuth(gitlabToken, gitlabURL)
	if err != nil {
		log.Fatal(err)
	}

	// Получение имени активного спринта (Milestone)
	milTitle, _, err := gitlab_issues.GetActiveMilestone(client, projectID)
	if err != nil {
		log.Fatal(err)
	}

	// Получение списка тасок
	issues, err := gitlab_issues.GetMilestoneIssues(client, projectID, milTitle, userName)
	if err != nil {
		log.Fatal(err)
	}

	totalHours := 0
	totalMinutes := 0

	// Проходим по каждой таске
	for _, i := range issues {
		taskHours := 0
		taskMinutes := 0
		// Получаем заметки(note) таски
		note, _, err := client.Notes.ListIssueNotes(projectID, i.IID, &gitlab.ListIssueNotesOptions{
			ListOptions: gitlab.ListOptions{},
			OrderBy:     gitlab.String("updated_at"),
			Sort:        gitlab.String("desc")})
		if err != nil {
			log.Fatal(err)
		}

		// Смотрим каждую заметку
		for _, n := range note {
			if n.UpdatedAt.Format("2006-01-02") == searchDate.Format("2006-01-02") && strings.Contains(n.Body, "time spent") {
				// Если заметка не содержит at (была создана в searchDate)
				if !strings.Contains(n.Body, "at") {
					hours, minutes, err := getTime(n.Body)
					if err != nil {
						log.Fatal(err)
					}
					taskHours += hours
					taskMinutes += minutes
					totalHours += hours
					totalMinutes += minutes

					if taskMinutes >= 60 {
						taskHours++
						taskMinutes -= 60
					}
				}
			}

			// Если заметка была создана в другую дату (/spend 2h 15m 2023-07-01)
			if strings.Contains(n.Body, "at") {
				noteDate, err := getDateFromNote(n.Body)
				if err != nil {
					log.Fatal(err)
				}
				// Сравниваем дату из заметки с датой из аргумента -date
				if searchDate.Equal(*noteDate) {
					hours, minutes, err := getTime(n.Body)
					if err != nil {
						log.Fatal(err)
					}
					taskHours += hours
					taskMinutes += minutes
					totalHours += hours
					totalMinutes += minutes

					if taskMinutes >= 60 {
						taskHours++
						taskMinutes -= 60
					}
				} else {
					continue
				}
			}
			if totalMinutes >= 60 {
				totalHours++
				totalMinutes -= 60
			}
		}
		if taskHours > 0 || taskMinutes > 0 {
			table.Append([]string{i.Title, fmt.Sprintf("%dh %dm", taskHours, taskMinutes)})
		}
	}
	// Отрисовка таблицы в stdout
	table.Render()

	fmt.Printf("Total spend time today: %dh %dm\n", totalHours, totalMinutes)
}
