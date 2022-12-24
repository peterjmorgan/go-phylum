package phylum

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"

	"github.com/coreos/go-oidc"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func CheckResponse(resp *resty.Response) *string {
	var jsonER JsonErrorResponse

	if resp.IsError() {
		err := json.Unmarshal(resp.Body(), &jsonER)
		if err != nil {
			fmt.Printf("failed to parse json: %v\n", err)
		}
		retString := fmt.Sprintf("%v - %v\n", jsonER.Error.Code, jsonER.Error.Description)
		return &retString
	}
	return nil
}

type ClientOptions struct {
	Token string
}

type PhylumClient struct {
	RefreshToken string
	OauthToken   oauth2.Token
	Ctx          context.Context
	Client       *resty.Client
	Groups       ListUserGroupsResponse
	AllProjects  []ProjectSummaryResponse
}

func NewClient(opts *ClientOptions) (*PhylumClient, error) {
	var PhylumToken string = ""
	var err error

	ctx := context.Background()
	client := resty.New()

	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		if opts.Token != "" {
			PhylumToken = opts.Token
		}
	}

	// Token wasn't set via options, try to extract from CLI
	if PhylumToken == "" {
		PhylumToken, err = GetTokenFromCLI()
		if err != nil {
			return nil, fmt.Errorf("Failed to get token from cli: %v\n", err)
		}
	}

	pClient := PhylumClient{
		RefreshToken: PhylumToken,
		Ctx:          ctx,
		Client:       client,
	}
	if err = pClient.GetAccessToken(); err != nil {
		return nil, fmt.Errorf("Failed to get access token: %v\n", err)
	}
	return &pClient, nil
}

func (p *PhylumClient) GetAccessToken() error {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://login.phylum.io/realms/phylum")
	if err != nil {
		fmt.Printf("failed to get oidc provider: %v\n", err)
		return err
	}

	oauth2Config := oauth2.Config{
		ClientID: "phylum_cli",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	ts := oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: p.RefreshToken})
	tok, err := ts.Token()
	if err != nil {
		fmt.Printf("failed to get access token: %v\n", err)
		return err
	}

	p.OauthToken = *tok

	return nil
}

type AuthStatus struct {
	Sub               string `json:"sub"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

func (p *PhylumClient) GetAuthStatus(token string) (bool, error) {
	var status AuthStatus
	var url string = "https://login.phylum.io/realms/phylum/protocol/openid-connect/userinfo"

	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(token).
		Get(url)
	if resp.StatusCode() != 200 {
		return false, nil
	}
	if err != nil {
		fmt.Printf("failed to GetAuthStatus: %v\n", err)
		return false, errors.New("failed to GetAuthStatus")
	}

	body := resp.Body()
	err = json.Unmarshal(body, &status)
	if err != nil {
		fmt.Printf("GetAuthStatus(): failed to parse response: %v\n", err)
		return false, err
	}

	if status.EmailVerified == true {
		return true, nil
	} else {
		return false, nil
	}
}

func GetTokenFromCLI() (string, error) {
	var stdErrBytes bytes.Buffer
	var phylumTokenArgs = []string{"auth", "token"}
	phylumTokenCmd := exec.Command("phylum", phylumTokenArgs...)
	phylumTokenCmd.Stderr = &stdErrBytes
	output, err := phylumTokenCmd.Output()
	stdErrString := stdErrBytes.String()
	if err != nil {
		fmt.Printf("error running phylum auth token: %v\n", err)
		fmt.Printf("stderr: %v\n", stdErrString)
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}

// GetUserGroups Get Phylum groups for which the user is a member or owner
// Write the result to the PhylumClient struct
func (p *PhylumClient) GetUserGroups() (*ListUserGroupsResponse, error) {
	userGroups := new(ListUserGroupsResponse)

	var url string = "https://api.phylum.io/api/v0/groups"

	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)

	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to get groups: %v\n", err)
		return nil, errors.New(*test)
	}

	body := resp.Body()
	err = json.Unmarshal(body, userGroups)
	if err != nil {
		fmt.Printf("GetGroups(): failed to parse response: %v\n", err)
		return nil, err
	}

	p.Groups = *userGroups

	return userGroups, nil
}

func (p *PhylumClient) GetHealth() (bool, error) {
	client := resty.New()
	token, err := GetTokenFromCLI()
	if err != nil {
		fmt.Printf("failed to get token from CLI")
		return false, err
	}

	resp, err := client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(token).
		Get("https://api.phylum.io/api/v0/health")
	if err != nil {
		fmt.Printf("failed to get health")
	}
	_ = resp
	if bytes.Contains(resp.Body(), []byte("alive")) {
		return true, nil
	}
	return false, fmt.Errorf("Health: couldn't find alive")
}

func (p *PhylumClient) ListProjects() ([]ProjectSummaryResponse, error) {
	var temp []ProjectSummaryResponse
	var url string = "https://api.phylum.io/api/v0/data/projects/overview"

	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)

	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to get projects: %v\n", err)
		return nil, errors.New(*test)
	}

	body := resp.Body()
	err = json.Unmarshal(body, &temp)
	if err != nil {
		fmt.Printf("GetProjects(): failed to parse response: %v\n", err)
	}

	return temp, nil
}

// TODO: abstract group elements in the optional struct
type ProjectOpts struct {
	GroupName string
}

func (p *PhylumClient) CreateProject(name string, opts *ProjectOpts) (*ProjectSummaryResponse, error) {
	var respPSR ProjectSummaryResponse
	var url string = "https://api.phylum.io/api/v0/data/projects"

	bodyMap := make(map[string]string, 0)
	bodyMap["name"] = name

	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		if opts.GroupName != "" {
			bodyMap["group_name"] = opts.GroupName
		}
	}

	resp, err := p.Client.R().
		SetAuthToken(p.OauthToken.AccessToken).
		SetBody(bodyMap).
		Post(url)
	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to create project\n")
		return nil, errors.New(*test)
	}
	err = json.Unmarshal(resp.Body(), &respPSR)
	if err != nil {
		fmt.Printf("CreateProject(): failed parse json: %v\n", err)
	}

	return &respPSR, nil
}

func CheckProjectId(projectId string) error {
	_, err := uuid.Parse(projectId)
	if err != nil {
		fmt.Printf("Error: must provide a valid GUID for project ID")
		return errors.New("ProjectID is not a guid")
	}
	return nil
}

func (p *PhylumClient) DeleteProject(projectId string) (*ProjectSummaryResponse, error) {
	var respPSR ProjectSummaryResponse

	if err := CheckProjectId(projectId); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.phylum.io/api/v0/data/projects/%v", projectId)

	resp, err := p.Client.R().
		SetAuthToken(p.OauthToken.AccessToken).
		Delete(url)
	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to delete project\n")
		return nil, errors.New(*test)
	}
	err = json.Unmarshal(resp.Body(), &respPSR)
	if err != nil {
		fmt.Printf("DeleteProject(): failed parse json: %v\n", err)
	}

	return &respPSR, nil
}

// GetProject Gets a project based on a Phylum project ID. It can get user or group projects.
func (p *PhylumClient) GetProject(projectID string) (*ProjectResponse, error) {
	var returnProject *ProjectResponse

	projects, err := p.ListAllProjects()
	if err != nil {
		return nil, err
	}

	// Find the project
	var targetProject ProjectSummaryResponse
	for _, proj := range projects {
		if proj.Id.String() == projectID {
			targetProject = proj
		}
	}
	if targetProject.Id.String() == "" {
		return nil, fmt.Errorf("GetProject: failed to find project with ID: %v\n", projectID)
	}

	if *targetProject.GroupName != "" {
		// group project
		returnProject, err = p.GetGroupProject(*targetProject.GroupName, targetProject.Id.String())
		if err != nil {
			return nil, err
		}
	} else {
		// user project
		returnProject, err = p.GetUserProject(targetProject.Id.String())
		if err != nil {
			return nil, err
		}
	}

	return returnProject, nil
}

// GetUserProject Gets a user project based on a Phylum project ID.
func (p *PhylumClient) GetUserProject(projectID string) (*ProjectResponse, error) {
	var result ProjectResponse

	url := fmt.Sprintf("https://api.phylum.io/api/v0/data/projects/%s", projectID)
	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)

	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("GetUserProject(): failed to get projects: %v\n", err)
		return nil, errors.New(*test)
	}

	body := resp.Body()
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("GetUserProject(): failed to parse response: %v\n", err)
	}

	return &result, nil
}

// TODO: this should be folded into GetProject() with GetProjectOpts struct
func (p *PhylumClient) GetGroupProject(groupName string, projectID string) (*ProjectResponse, error) {
	var result ProjectResponse

	url := fmt.Sprintf("https://api.phylum.io/api/v0/groups/%s/projects/%s", groupName, projectID)
	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)

	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to get projects: %v\n", err)
		return nil, errors.New(*test)
	}

	body := resp.Body()
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("GetProjects(): failed to parse response: %v\n", err)
	}

	return &result, nil
}

// TODO: this should be folded into ListProjects() with an optional struct
func (p *PhylumClient) ListGroupProjects(groupName string) ([]ProjectSummaryResponse, error) {
	var result []ProjectSummaryResponse
	url := fmt.Sprintf("https://api.phylum.io/api/v0/groups/%s/projects", groupName)

	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)

	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to ListGroupProjects: %v\n", err)
		return nil, errors.New(*test)
	}

	body := resp.Body()
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("ListGroupProjects(): failed to parse response: %v\n", err)
	}

	return result, nil
}

// type MetaProject struct {
// 	Project ProjectSummaryResponse
// 	Group string
// }

func (p *PhylumClient) ListAllProjects() ([]ProjectSummaryResponse, error) {
	var allProjects []ProjectSummaryResponse

	// Get all group projects into a slice
	groups, err := p.GetUserGroups()
	if err != nil {
		fmt.Printf("Failed to GetUserGroups(): %v\n", err)
		return nil, err
	}

	for _, group := range groups.Groups {
		groupProjectList, err := p.ListGroupProjects(group.GroupName)
		if err != nil {
			fmt.Printf("Failed to ListGroupProjects: %v\n", err)
			return nil, err
		}
		// for _, gp := range groupProjectList {
		// 	newMetaProject := MetaProject{
		// 		Project: gp,
		// 		Group:   group.GroupName,
		// 	}
		// 	allProjects = append(allProjects, newMetaProject)
		// }
		allProjects = append(allProjects, groupProjectList...)
	}

	// Add User Projects to slice
	projectList, err := p.ListProjects()
	if err != nil {
		fmt.Printf("Failed to ListProjects(): %v\n", err)
		return nil, err
	}

	// for _, up := range projectList {
	// 	newProject := MetaProject{
	// 		Project: up,
	// 	}
	// 	allProjects = append(allProjects, newProject)
	// }
	allProjects = append(allProjects, projectList...)

	//p.AllProjects = allProjects

	return allProjects, nil
}

// Default should get all projects in all groups
func (p *PhylumClient) GetAllProjects() ([]*ProjectResponse, error) {
	var result []*ProjectResponse
	var allProjectList []ProjectSummaryResponse

	// Get all group projects into a slice
	groups, err := p.GetUserGroups()
	if err != nil {
		fmt.Printf("Failed to GetUserGroups(): %v\n", err)
		return nil, err
	}

	for _, group := range groups.Groups {
		groupProjectList, err := p.ListGroupProjects(group.GroupName)
		if err != nil {
			fmt.Printf("Failed to ListGroupProjects: %v\n", err)
			return nil, err
		}
		allProjectList = append(allProjectList, groupProjectList...)
	}

	projectList, err := p.ListProjects()
	if err != nil {
		fmt.Printf("Failed to ListProjects(): %v\n", err)
		return nil, err
	}

	allProjectList = append(allProjectList, projectList...)

	chRecv := make(chan *ProjectResponse)
	var wg sync.WaitGroup

	for _, proj := range allProjectList {
		wg.Add(1)
		go func(inProj ProjectSummaryResponse) {
			defer wg.Done()
			var temp *ProjectResponse

			if inProj.GroupName != nil {
				temp, err = p.GetGroupProject(*inProj.GroupName, inProj.Id.String())
				if err != nil {
					fmt.Printf("Failed to GetGroupProject: %v\n", err)
					return
				}
			} else {
				temp, err = p.GetUserProject(inProj.Id.String())
				if err != nil {
					fmt.Printf("Failed to GetUserProject: %v\n", err)
					return
				}
			}
			chRecv <- temp
		}(proj)
	}

	go func() {
		for res := range chRecv {
			result = append(result, res)
		}
		close(chRecv)
	}()
	wg.Wait()

	return result, err
}

func (p *PhylumClient) GetAllGroupProjects(groupName string) ([]*ProjectResponse, error) {
	var result []*ProjectResponse

	groupProjectList, err := p.ListGroupProjects(groupName)
	if err != nil {
		fmt.Printf("Failed to ListGroupProjects: %v\n", err)
		return nil, err
	}

	chRecv := make(chan *ProjectResponse)
	var wg sync.WaitGroup

	for _, proj := range groupProjectList {
		wg.Add(1)
		go func(inProj ProjectSummaryResponse) {
			defer wg.Done()
			temp, err := p.GetGroupProject(groupName, inProj.Id.String())
			if err != nil {
				fmt.Printf("Failed to GetGroupProject: %v\n", err)
				return
			}
			chRecv <- temp

		}(proj)
	}

	go func() {
		for res := range chRecv {
			result = append(result, res)
		}
		close(chRecv)
	}()
	wg.Wait()

	return result, err
}

// TODO: should be removed
func (p *PhylumClient) GetAllGroupProjectsByEcosystem(groupName string, ecosystem string) ([]*ProjectResponse, error) {
	var result []*ProjectResponse
	var targetList []ProjectSummaryResponse
	var wg sync.WaitGroup

	groupProjectList, err := p.ListGroupProjects(groupName)
	if err != nil {
		fmt.Printf("Failed to ListGroupProjects: %v\n", err)
		return nil, err
	}

	for _, proj := range groupProjectList {
		if proj.Ecosystem != nil {
			projEcosystem := fmt.Sprintf("%v", *proj.Ecosystem)
			if projEcosystem == ecosystem {
				targetList = append(targetList, proj)
			}
		}
	}

	chRecv := make(chan *ProjectResponse)

	for _, proj := range targetList {
		wg.Add(1)
		go func(inProj ProjectSummaryResponse) {
			defer wg.Done()
			temp, err := p.GetGroupProject(groupName, inProj.Id.String())
			if err != nil {
				fmt.Printf("Failed to GetGroupProject: %v\n", err)
				return
			}
			chRecv <- temp
		}(proj)
	}

	go func() {
		for res := range chRecv {
			result = append(result, res)
		}
		close(chRecv)
	}()
	wg.Wait()

	return result, err
}

func (p *PhylumClient) AnalyzeParsedPackages(projectType string, projectID string, packages *[]PackageDescriptor) (string, error) {
	var respSPR SubmitPackageResponse
	var url string = "https://api.phylum.io/api/v0/data/jobs"

	submitPackageRequest := SubmitPackageRequest{
		GroupName: nil,
		IsUser:    true,
		Label:     "",
		Packages:  *packages,
		Project:   projectID,
		Type:      projectType,
	}

	resp, err := p.Client.R().
		SetAuthToken(p.OauthToken.AccessToken).
		SetBody(submitPackageRequest).
		Post(url)
	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to analyze packages\n")
		return "", errors.New(*test)
	}
	err = json.Unmarshal(resp.Body(), &respSPR)
	if err != nil {
		fmt.Printf("AnalyzeParsedPackages(): failed parse json: %v\n", err)
		return "", fmt.Errorf("AnalyzeParsedPackages(): failed parse json: %v\n", err)
	}
	if respSPR.JobId.String() == "" {
		return "", fmt.Errorf("AnalyzeParsedPackages(): failed to read JobID, submission may not have been successful")
	}

	return respSPR.JobId.String(), nil
}

// TODO: handle non verbose
// TODO: handle jobstatusresponsevariant
func (p *PhylumClient) GetJobVerbose(jobID string) (*JobStatusResponseForPackageStatusExtended, *[]byte, error) {
	var jobResponse JobStatusResponseForPackageStatusExtended
	url := fmt.Sprintf("https://api.phylum.io/api/v0/data/jobs/%s?verbose=true", jobID)

	resp, err := p.Client.R().
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)
	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to GetJob\n")
		return nil, nil, errors.New(*test)
	}
	err = json.Unmarshal(resp.Body(), &jobResponse)
	if err != nil {
		fmt.Printf("GetJob(): failed parse json: %v\n", err)
		return nil, nil, err
	}
	jsonData := resp.Body()

	return &jobResponse, &jsonData, nil
}

// ParseLockfile parses a lockfile into a struct that can be submitted for analysis.
// It takes the path to a lockfile as input, and returns a pointer to a slice of PackageDescriptors
func (p *PhylumClient) ParseLockfile(lockfilePath string) (*[]PackageDescriptor, error) {
	if _, err := os.Stat(lockfilePath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("lockfilePath: %v is not a file", lockfilePath)
	}
	var packages []PackageDescriptor
	url := "https://parse.phylum.io"

	resp, err := p.Client.R().
		SetAuthToken(p.OauthToken.AccessToken).
		SetFile("lockfile", lockfilePath).
		Post(url)
	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("ParseLockfile(): failed to ParseLockfile\n")
		return nil, errors.New(*test)
	}

	err = json.Unmarshal(resp.Body(), &packages)
	if err != nil {
		fmt.Printf("ParseLockfile(): failed parse json: %v\n", err)
		return nil, err
	}
	return &packages, nil
}

func (p *PhylumClient) GetProjectPreferences(projectID string) (*ProjectPreferencesResponse, error) {
	var result ProjectPreferencesResponse

	url := fmt.Sprintf("https://api.phylum.io/api/v0/preferences/project/%s", projectID)
	resp, err := p.Client.R().
		SetHeader("accept", "application/json").
		SetAuthToken(p.OauthToken.AccessToken).
		Get(url)

	test := CheckResponse(resp)
	if test != nil || err != nil {
		fmt.Printf("failed to get projects: %v\n", err)
		return nil, errors.New(*test)
	}

	body := resp.Body()
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("GetProjects(): failed to parse response: %v\n", err)
	}

	return &result, nil
}

// GetProjectIssues gets the issues for a project. Only returns issues that are not ignored/suppressed
func (p *PhylumClient) GetProjectIssues(projectId string) ([]IssuesListItem, error) {
	var issues []IssuesListItem
	if err := CheckProjectId(projectId); err != nil {
		return nil, err
	}

	projectResponse, err := p.GetUserProject(projectId)
	if err != nil {
		return nil, err
	}

	for _, elem := range projectResponse.Issues {
		if elem.Ignored == "false" {
			issues = append(issues, elem)
		}
	}

	return issues, nil
}
