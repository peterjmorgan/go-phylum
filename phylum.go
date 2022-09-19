package phylum

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type PhylumClient struct {
	RefreshToken string
	OauthToken   oauth2.Token
	Ctx          context.Context
	Client       *resty.Client
	Groups       ListUserGroupsResponse
}

func (p *PhylumClient) GetAccessToken() error {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://login.phylum.io/auth/realms/phylum")
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

func NewClient() *PhylumClient {
	ctx := context.Background()
	client := resty.New()

	token, err := GetTokenFromCLI()
	if err != nil {
		fmt.Printf("Failed to get token from cli: %v\n", err)
		return nil
	}

	pClient := PhylumClient{
		RefreshToken: token,
		Ctx:          ctx,
		Client:       client,
	}
	if err = pClient.GetAccessToken(); err != nil {
		fmt.Printf("Failed to get access token: %v\n", err)
		return nil
	}
	return &pClient
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

func (p *PhylumClient) GetUserGroups() (*ListUserGroupsResponse, error) {
	var userGroups *ListUserGroupsResponse

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

func checkProjectId(projectId string) error {
	_, err := uuid.Parse(projectId)
	if err != nil {
		fmt.Printf("Error: must provide a valid GUID for project ID")
		return errors.New("not a guid")
	}
	return nil
}

func (p *PhylumClient) DeleteProject(projectId string) (*ProjectSummaryResponse, error) {
	var respPSR ProjectSummaryResponse

	if err := checkProjectId(projectId); err != nil {
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

// TODO: add group and label handling
func (p *PhylumClient) AnalyzeParsedPackages(projectType string, projectID string, packages *[]PackageDescriptor) (*interface{}, error) {
	var respPSR SubmitPackageResponse
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
		return nil, errors.New(*test)
	}
	err = json.Unmarshal(resp.Body(), &respPSR)
	if err != nil {
		fmt.Printf("AnalyzeParsedPackages(): failed parse json: %v\n", err)
		return nil, err
	}

	return nil, nil
}

// TODO: handle non verbose
// TODO: handle jobstatusresponsevariant
func (p *PhylumClient) GetJob(jobID string) (*JobStatusResponseForPackageStatusExtended, *[]byte, error) {
	var jobResponse JobStatusResponseForPackageStatusExtended
	url := fmt.Sprintf("https://api.phylum.io/api/v0/data/jobs/%s", jobID)

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
