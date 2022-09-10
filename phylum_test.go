package phylum

import (
	"reflect"
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
	pc := NewClient()
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
	pc := NewClient()
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
	pc := NewClient()

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
	p := NewClient()

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
	p := NewClient()
	//p.Client.SetProxy("http://localhost:8080")
	//p.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

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
	p := NewClient()

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
