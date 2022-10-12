package phylum

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func Test_getTokenFromCLI(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"one"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTokenFromCLI()
			if err != nil {
				t.Errorf("Failed to get token")
			}
			_ = got
		})
	}
}

func Test_NewClient_WithOptions(t *testing.T) {
	opts := ClientOptions{
		Token: "",
	}

	pc, err := NewClient(&opts)
	if err != nil {
		t.Errorf("NewClient() failed with error: %v\n", err)
	}

	if reflect.TypeOf(pc) != reflect.TypeOf(&PhylumClient{}) {
		t.Errorf("NewClient() got = %v, want %v", reflect.TypeOf(pc), reflect.TypeOf(&PhylumClient{}))
	}

	p, _ := NewClient(&ClientOptions{})
	if reflect.TypeOf(p) != reflect.TypeOf(&PhylumClient{}) {
		t.Errorf("NewClient() got = %v, want %v", reflect.TypeOf(p), reflect.TypeOf(&PhylumClient{}))
	}

	fmt.Println("Done")
}

// func Test_api_getHealth(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 		{"one"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			api_getHealth()
// 		})
// 	}
// }

func Test_PhylumClient_ListProjects(t *testing.T) {
	pc, _ := NewClient(&ClientOptions{})
	tests := []struct {
		name    string
		want    []ProjectSummaryResponse
		wantErr bool
	}{
		{"one", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pc.ListProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("api_GetProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("ListProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_CreateProject(t *testing.T) {
	pc, _ := NewClient(&ClientOptions{})
	type args struct {
		name string
		opts *ProjectOpts
	}
	tests := []struct {
		name    string
		args    args
		want    *ProjectSummaryResponse
		wantErr bool
	}{
		{"w/o groups", args{"ABCD-testProject12", &ProjectOpts{}}, nil, false},
		{"with groups", args{"ABCD-testProject12", &ProjectOpts{GroupName: "test2"}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pc.CreateProject(tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("CreateProject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_GetGroupProject(t *testing.T) {
	pc, _ := NewClient(&ClientOptions{})

	type args struct {
		groupName string
		projectID string
	}
	tests := []struct {
		name    string
		args    args
		want    *ProjectResponse
		wantErr bool
	}{
		{"one", args{"test2", "85e3142f-efc9-41fc-b004-ca570df89af8"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pc.GetGroupProject(tt.args.groupName, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetGroupProject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_ListGroupProjects(t *testing.T) {
	p, _ := NewClient(&ClientOptions{})

	type args struct {
		groupName string
	}
	tests := []struct {
		name    string
		args    args
		want    []ProjectSummaryResponse
		wantErr bool
	}{
		{"one", args{"test2"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.ListGroupProjects(tt.args.groupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListGroupProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("ListGroupProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_GetAllGroupProjects(t *testing.T) {
	p, _ := NewClient(&ClientOptions{})
	// p.Client.SetProxy("http://localhost:8080")
	// p.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	type args struct {
		groupName string
	}
	tests := []struct {
		name    string
		args    args
		want    []*ProjectResponse
		wantErr bool
	}{
		{"one", args{"spring-test2"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetAllGroupProjects(tt.args.groupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllGroupProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetAllGroupProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_GetAllGroupProjectsByEcosystem(t *testing.T) {
	p, _ := NewClient(&ClientOptions{})

	type args struct {
		groupName string
		ecosystem string
	}
	tests := []struct {
		name    string
		args    args
		want    []*ProjectResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{"one", args{"test2", "maven"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.GetAllGroupProjectsByEcosystem(tt.args.groupName, tt.args.ecosystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllGroupProjectsByEcosystem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetAllGroupProjectsByEcosystem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_AnalyzePackages(t *testing.T) {
	p, _ := NewClient(&ClientOptions{})
	// p.Client.SetProxy("http://localhost:8080")
	// p.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	var ParseCmdArgs = []string{"parse", "package-lock.json"}
	parseCmd := exec.Command("phylum", ParseCmdArgs...)
	output, err := parseCmd.Output()
	if err != nil {
		fmt.Printf("Failed to exec 'phylum %v': %v\n", strings.Join(ParseCmdArgs, " "), err)
		return
	}
	packages := make([]PackageDescriptor, 0)
	err = json.Unmarshal(output, &packages)

	type args struct {
		projectType string
		projectID   string
		packages    *[]PackageDescriptor
	}
	tests := []struct {
		name    string
		args    args
		want    *interface{}
		wantErr bool
	}{
		{"one", args{"npm", "42b07f68-cf1c-42c8-a217-ce2d903c22b5", &packages}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.AnalyzeParsedPackages(tt.args.projectType, tt.args.projectID, tt.args.packages)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeParsedPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzeParsedPackages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhylumClient_GetJob(t *testing.T) {
	p, _ := NewClient(&ClientOptions{})

	type args struct {
		jobID   string
		verbose bool
	}
	tests := []struct {
		name    string
		args    args
		want    *JobStatusResponseForPackageStatusExtended
		wantErr bool
	}{
		{"one", args{"e1dc76a9-02c2-43c9-bd03-e654fb569ab9", true}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, jsonData, err := p.GetJobVerbose(tt.args.jobID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetJob() got = %v, want %v", got, tt.want)
			}
			blah := *jsonData
			_ = blah
		})
	}
}

//func TestPhylumClient_ParseLockfile(t *testing.T) {
//	p, _ := NewClient(&ClientOptions{})
//
//	type args struct {
//		lockfilePath string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *[]PackageDescriptor
//		wantLen int
//		wantErr bool
//	}{
//		{"package-lock", args{"test_lockfiles/package-lock.json"}, nil, 52, false},
//		// This lockfile fails parsing at the moment
//		{"package-lock-v6", args{"test_lockfiles/package-lock-v6.json"}, nil, 91, true},
//		{"requirements.txt", args{"test_lockfiles/requirements.txt"}, nil, 131, false},
//		{"poetry.lock", args{"test_lockfiles/poetry.lock"}, nil, 45, false},
//		{"Gemfile.lock", args{"test_lockfiles/Gemfile.lock"}, nil, 214, false},
//		{"Yarn", args{"test_lockfiles/yarn.lock"}, nil, 53, false},
//		{"pipfile.lock", args{"test_lockfiles/Pipfile.lock"}, nil, 27, false},
//		{"gradle.lockfile", args{"test_lockfiles/gradle.lockfile"}, nil, 6, false},
//		{"effective-pom", args{"test_lockfiles/effective-pom.xml"}, nil, 16, false},
//		// workspace-effective-pom Fails currently
//		{"workspace-effective-pom", args{"test_lockfiles/workspace-effective-pom.xml"}, nil, 16, true},
//		{"Calculator csproj", args{"test_lockfiles/Calculator.csproj"}, nil, 2, false},
//		{"sample csproj", args{"test_lockfiles/sample.csproj"}, nil, 5, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			var ParseCmdArgs = []string{"parse", tt.args.lockfilePath}
//			parseCmd := exec.Command("phylum", ParseCmdArgs...)
//			output, err := parseCmd.Output()
//			if err != nil {
//				t.Errorf("Failed to exec 'phylum %v': %v\n", strings.Join(ParseCmdArgs, " "), err)
//				return
//			}
//			phylumPackages := make([]PackageDescriptor, 0)
//			err = json.Unmarshal(output, &phylumPackages)
//			if err != nil {
//				t.Errorf("Failed to unmarshall json: %v\n", err)
//				return
//			}
//
//			got, err := p.ParseLockfile(tt.args.lockfilePath)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ParseLockfile() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			// temporary - write out the file
//			gotOutput, err := json.Marshal(got)
//			if err != nil {
//				log.Fatalf("failed to marshall json: %v\n", err)
//				return
//			}
//			err = ioutil.WriteFile("gotOUtput.json", gotOutput, 0644)
//			if err != nil {
//				log.Fatalf("failed to write file: %v\n", err)
//				return
//			}
//			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
//				t.Errorf("ParseLockfile() got = %v, want %v", got, tt.want)
//			}
//			if len(*got) != len(phylumPackages) {
//				t.Errorf("ParseLockfile(): Phylum parse and ParseLockfile returned slices of differing sizes")
//			}
//			if len(*got) != tt.wantLen {
//				t.Errorf("ParseLockfile(): len of parsed packages: got = %v, want %v", len(*got), tt.wantLen)
//			}
//		})
//	}
//}

func TestPhylumClient_ParseLockfile1(t *testing.T) {
	p, _ := NewClient(&ClientOptions{})

	type args struct {
		lockfilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]PackageDescriptor
		wantLen int
		wantErr bool
	}{
		{"package-lock", args{"test_lockfiles/package-lock.json"}, nil, 52, false},
		// This lockfile fails parsing at the moment
		{"package-lock-v6", args{"test_lockfiles/package-lock-v6.json"}, nil, 17, false},
		{"requirements.txt", args{"test_lockfiles/requirements.txt"}, nil, 131, false},
		{"poetry.lock", args{"test_lockfiles/poetry.lock"}, nil, 45, false},
		{"Gemfile.lock", args{"test_lockfiles/Gemfile.lock"}, nil, 214, false},
		{"Yarn", args{"test_lockfiles/yarn.lock"}, nil, 53, false},
		{"pipfile.lock", args{"test_lockfiles/Pipfile.lock"}, nil, 27, false},
		{"gradle.lockfile", args{"test_lockfiles/gradle.lockfile"}, nil, 6, false},
		{"effective-pom", args{"test_lockfiles/effective-pom.xml"}, nil, 16, false},
		// workspace-effective-pom Fails currently
		{"workspace-effective-pom", args{"test_lockfiles/workspace-effective-pom.xml"}, nil, 88, false},
		{"Calculator csproj", args{"test_lockfiles/Calculator.csproj"}, nil, 2, false},
		{"sample csproj", args{"test_lockfiles/sample.csproj"}, nil, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.ParseLockfile(tt.args.lockfilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLockfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(*got) != tt.wantLen {
				t.Errorf("ParseLockfile(): len of parsed packages: got = %v, want %v", len(*got), tt.wantLen)
			}
		})
	}
}
