package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/xanzy/go-gitlab"
)

var  licenseFiles = []string{"LICENSE", "LICENSE.txt", "license", "license.txt"}


func main() {

	var wg sync.WaitGroup
	git := gitlab.NewClient(nil, "CHANGEME")
	git.SetBaseURL("https://gitlab.com/api/v4")
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}
	for {
		projects, resp , err := git.Projects.ListProjects(opt)
		if err != nil {
			log.Fatal(err)
		}

		for _, project := range projects {
			wg.Add(1)
			 checkLicenseFile(project, git, &wg)
		}
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opt.Page = resp.NextPage
	}
	wg.Wait()
}

func checkLicenseFile(project *gitlab.Project,   git *gitlab.Client, wg *sync.WaitGroup) {
	hasLicense := false
	fileName := ""
	typeLicense := "Unknown"
	gf := &gitlab.GetRawFileOptions{Ref: gitlab.String("master"),}
	for _,file := range licenseFiles {
		f, _, err := git.RepositoryFiles.GetRawFile(project.ID, file, gf)
		fileContent := string(f)
		if err == nil {
			hasLicense = true
			fileName = file


			if (strings.Contains(fileContent, "Apche")) {
				typeLicense = "Apache"
			}
			break
		}
	}
	if hasLicense {
		fmt.Printf("project \"%s\" has license \"%s\" in %s file  \n", project.Name,  typeLicense, fileName)
	} else {
		fmt.Printf("project \"%s\" hasn't license file \n", project.Name)
	}

	wg.Done()

}
