// Package main provides primitives to interact with the openapi HTTP API.
package phylum

import (
	"encoding/json"
	"fmt"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	Phylum_CLI_TokenScopes = "Phylum_CLI_Token.Scopes"
)

// Defines values for Action.
const (
	ActionBreak Action = "break"
	ActionNone  Action = "none"
	ActionWarn  Action = "warn"
)

// Defines values for IgnoredReason.
const (
	False         IgnoredReason = "false"
	FalsePositive IgnoredReason = "falsePositive"
	NotRelevant   IgnoredReason = "notRelevant"
	Other         IgnoredReason = "other"
)

// Defines values for IngestionSource.
const (
	CLI               IngestionSource = "CLI"
	GitHubIntegration IngestionSource = "GitHub Integration"
	GitLabIntegration IngestionSource = "GitLab Integration"
)

// Defines values for PackageType.
const (
	Maven    PackageType = "maven"
	Npm      PackageType = "npm"
	Nuget    PackageType = "nuget"
	Pypi     PackageType = "pypi"
	Rubygems PackageType = "rubygems"
)

// Defines values for PaginateDirection.
const (
	Backward PaginateDirection = "Backward"
	Forward  PaginateDirection = "Forward"
)

// Defines values for ProjectField.
const (
	CreatedAt ProjectField = "CreatedAt"
	Ecosystem ProjectField = "Ecosystem"
	Group     ProjectField = "Group"
	Name      ProjectField = "Name"
	UpdatedAt ProjectField = "UpdatedAt"
)

// Defines values for RiskDomain.
const (
	RiskDomainAuthor        RiskDomain = "author"
	RiskDomainEngineering   RiskDomain = "engineering"
	RiskDomainLicense       RiskDomain = "license"
	RiskDomainMaliciousCode RiskDomain = "malicious_code"
	RiskDomainVulnerability RiskDomain = "vulnerability"
)

// Defines values for RiskLevel.
const (
	Critical RiskLevel = "critical"
	High     RiskLevel = "high"
	Info     RiskLevel = "info"
	Low      RiskLevel = "low"
	Medium   RiskLevel = "medium"
)

// Defines values for RiskType.
const (
	AuthorsRisk       RiskType = "authorsRisk"
	EngineeringRisk   RiskType = "engineeringRisk"
	LicenseRisk       RiskType = "licenseRisk"
	MaliciousCodeRisk RiskType = "maliciousCodeRisk"
	TotalRisk         RiskType = "totalRisk"
	Vulnerabilities   RiskType = "vulnerabilities"
)

// Defines values for SortDirection.
const (
	Ascending  SortDirection = "Ascending"
	Descending SortDirection = "Descending"
)

// Defines values for Status.
const (
	Complete   Status = "complete"
	Incomplete Status = "incomplete"
)

// Defines values for ThresholdViolationAction.
const (
	ThresholdViolationActionBreak ThresholdViolationAction = "break"
	ThresholdViolationActionNone  ThresholdViolationAction = "none"
	ThresholdViolationActionWarn  ThresholdViolationAction = "warn"
)

// Defines values for ValidatedGroupNameError.
const (
	Invalid ValidatedGroupNameError = "Invalid"
)

// When a job is completed, and some requirement is not met ( such as quality level ), what action should be taken? In the case of the CLI, the value of this result is used to determine if the CLI should print a warning, or exit with a non-zero exit code.
type Action string

// Represents a response that summarizes the output of all current jobs
type AllJobsStatusResponse struct {
	Count uint32 `json:"count"`

	// A description of the latest jobs
	Jobs []JobDescriptor `json:"jobs"`

	// Total jobs run
	TotalJobs uint32 `json:"total_jobs"`
}

// Author information
type Author struct {
	AvatarUrl  string `json:"avatarUrl"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	ProfileUrl string `json:"profileUrl"`
}

// CreateGroupResponse defines model for CreateGroupResponse.
type CreateGroupResponse struct {
	GroupName  string `json:"group_name"`
	OwnerEmail string `json:"owner_email"`
}

// Rquest to create a project
type CreateProjectRequest struct {
	GroupName *string `json:"group_name"`
	Name      string  `json:"name"`
}

// DependenciesCounts defines model for DependenciesCounts.
type DependenciesCounts struct {
	AboveThreshold uint32 `json:"aboveThreshold"`
	BelowThreshold uint32 `json:"belowThreshold"`
	NumIncomplete  uint32 `json:"numIncomplete"`
	Total          uint32 `json:"total"`
}

// DependenciesStatsBlock defines model for DependenciesStatsBlock.
type DependenciesStatsBlock struct {
	Counts DependenciesCounts `json:"counts"`
}

// Responsiveness of developers
type DeveloperResponsiveness struct {
	OpenIssueAvgDuration       *uint32 `json:"open_issue_avg_duration"`
	OpenIssueCount             *uint   `json:"open_issue_count"`
	OpenPullRequestAvgDuration *uint32 `json:"open_pull_request_avg_duration"`
	OpenPullRequestCount       *uint   `json:"open_pull_request_count"`
	TotalIssueCount            *uint   `json:"total_issue_count"`
	TotalPullRequestCount      *uint   `json:"total_pull_request_count"`
}

// Information about a downloadable file.
type DownloadInfo struct {
	// The canonical name of the file.
	Filename string `json:"filename"`

	// The URL at which the file is available.
	Url string `json:"url"`
}

// The kind of error encountered
type ErrorKind interface{}

// Extended information about a package.
//
// This contains additional information and can only be requested for one package at a time.
type ExtendedPackage struct {
	// An empty list.
	//
	// This is maintained for compatibility. Author data is retrieved via a separate API call.
	Authors []interface{} `json:"authors"`

	// True if the package has been processed.
	Complete bool `json:"complete"`

	// The declared dependencies of the package.
	DepSpecs []PackageSpecifier `json:"depSpecs"`

	// Details about the declared dependencies of the package.
	Dependencies []FullPackage `json:"dependencies"`

	// The description of the package, as provided by the package author.
	Description *string `json:"description"`

	// Metrics about developer reponsiveness to issues and pull requests.
	DeveloperResponsiveness *struct {
		OpenIssueAvgDuration       *uint32 `json:"open_issue_avg_duration"`
		OpenIssueCount             *uint   `json:"open_issue_count"`
		OpenPullRequestAvgDuration *uint32 `json:"open_pull_request_avg_duration"`
		OpenPullRequestCount       *uint   `json:"open_pull_request_count"`
		TotalIssueCount            *uint   `json:"total_issue_count"`
		TotalPullRequestCount      *uint   `json:"total_pull_request_count"`
	} `json:"developerResponsiveness"`

	// The number of times the package has been downloaded from its package registry.
	DownloadCount uint32 `json:"downloadCount"`

	// The Phylum-specific ID of the package version.
	Id string `json:"id"`

	// Whether our heuristics consider this package to be considered abandonware.
	IsAbandonware *bool `json:"is_abandonware"`

	// Statistics about how many issues exist at each impact level.
	IssueImpacts struct {
		Critical *uint32 `json:"critical,omitempty"`
		High     *uint32 `json:"high,omitempty"`
		Low      *uint32 `json:"low,omitempty"`
		Medium   *uint32 `json:"medium,omitempty"`
	} `json:"issueImpacts"`

	// The list of known issues for this package.
	Issues []IssuesListItem `json:"issues"`

	// The list of known issues for this package.
	IssuesDetails []Issue `json:"issuesDetails"`

	// The latest version of the package, as reported by the package registry.
	LatestVersion *string `json:"latestVersion"`

	// The license of the package.
	License *string `json:"license"`

	// Whether our heuristics consider this package to have recently changed maintainers.
	MaintainersRecentlyChanged *bool `json:"maintainers_recently_changed"`

	// The name of the package.
	Name string `json:"name"`

	// The date this version of the package was published.
	PublishedDate string `json:"publishedDate"`

	// The registry where the package is located.
	Registry string `json:"registry"`

	// Information about this package's different releases.
	ReleaseData *struct {
		FirstReleaseDate time.Time `json:"first_release_date"`
		LastReleaseDate  time.Time `json:"last_release_date"`
	} `json:"releaseData"`

	// The URL for the source code repository.
	RepoUrl *string `json:"repoUrl"`

	// The risk scores for this version of the package.
	RiskScores struct {
		Author        float32 `json:"author"`
		Engineering   float32 `json:"engineering"`
		License       float32 `json:"license"`
		MaliciousCode float32 `json:"malicious_code"`
		Total         float32 `json:"total"`
		Vulnerability float32 `json:"vulnerability"`
	} `json:"riskScores"`

	// The version of the package.
	Version string `json:"version"`

	// The known versions of this package.
	Versions []ScoredVersion `json:"versions"`
}

// Information about a package.
//
// This is available only to authenticated users.
type FullPackage struct {
	// An empty list.
	//
	// This is maintained for compatibility. Author data is retrieved via a separate API call.
	Authors []interface{} `json:"authors"`

	// True if the package has been processed.
	Complete bool `json:"complete"`

	// The declared dependencies of the package.
	DepSpecs []PackageSpecifier `json:"depSpecs"`

	// The description of the package, as provided by the package author.
	Description *string `json:"description"`

	// Zero.
	//
	// This is left here for compatibility.
	DownloadCount uint32 `json:"downloadCount"`

	// The Phylum-specific ID of the package version.
	Id string `json:"id"`

	// Statistics about how many issues exist at each impact level.
	IssueImpacts struct {
		Critical *uint32 `json:"critical,omitempty"`
		High     *uint32 `json:"high,omitempty"`
		Low      *uint32 `json:"low,omitempty"`
		Medium   *uint32 `json:"medium,omitempty"`
	} `json:"issueImpacts"`

	// The list of known issues for this package.
	Issues []IssuesListItem `json:"issues"`

	// The list of known issues for this package.
	IssuesDetails []Issue `json:"issuesDetails"`

	// The license of the package.
	License *string `json:"license"`

	// The name of the package.
	Name string `json:"name"`

	// The date this version of the package was published.
	PublishedDate string `json:"publishedDate"`

	// The registry where the package is located.
	Registry string `json:"registry"`

	// The URL for the source code repository.
	RepoUrl *string `json:"repoUrl"`

	// The risk scores for this version of the package.
	RiskScores struct {
		Author        float32 `json:"author"`
		Engineering   float32 `json:"engineering"`
		License       float32 `json:"license"`
		MaliciousCode float32 `json:"malicious_code"`
		Total         float32 `json:"total"`
		Vulnerability float32 `json:"vulnerability"`
	} `json:"riskScores"`

	// The version of the package.
	Version string `json:"version"`

	// An empty list.
	//
	// This is maintained for compatibility. Version data is only available on ExtendedPackage.
	Versions []interface{} `json:"versions"`
}

// Information about a package.
//
// This is available only to authenticated users.
//
// This internal version is used inside the server. Use `FullPackage::from` to convert it into a `Package` compatible format for sending over the network.
//
// ```mermaid classDiagram FullPackageInternal <|-- FullPackage FullPackageInternal <|-- ExtendedPackageInternal link FullPackage "struct.FullPackage.html" link ExtendedPackageInternal "struct.ExtendedPackageInternal.html" ```
type FullPackageInternal struct {
	// True if the package has been processed.
	Complete bool `json:"complete"`

	// The declared dependencies of the package.
	DepSpecs []PackageSpecifier `json:"depSpecs"`

	// The description of the package, as provided by the package author.
	Description *string `json:"description"`

	// The Phylum-specific ID of the package version.
	Id string `json:"id"`

	// Statistics about how many issues exist at each impact level.
	IssueImpacts struct {
		Critical *uint32 `json:"critical,omitempty"`
		High     *uint32 `json:"high,omitempty"`
		Low      *uint32 `json:"low,omitempty"`
		Medium   *uint32 `json:"medium,omitempty"`
	} `json:"issueImpacts"`

	// The list of known issues for this package.
	Issues []IssuesListItem `json:"issues"`

	// The list of known issues for this package.
	IssuesDetails []Issue `json:"issuesDetails"`

	// The license of the package.
	License *string `json:"license"`

	// The name of the package.
	Name string `json:"name"`

	// The date this version of the package was published.
	PublishedDate string `json:"publishedDate"`

	// The registry where the package is located.
	Registry string `json:"registry"`

	// The URL for the source code repository.
	RepoUrl *string `json:"repoUrl"`

	// The risk scores for this version of the package.
	RiskScores struct {
		Author        float32 `json:"author"`
		Engineering   float32 `json:"engineering"`
		License       float32 `json:"license"`
		MaliciousCode float32 `json:"malicious_code"`
		Total         float32 `json:"total"`
		Vulnerability float32 `json:"vulnerability"`
	} `json:"riskScores"`

	// The version of the package.
	Version string `json:"version"`
}

// GroupMember defines model for GroupMember.
type GroupMember struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserEmail string `json:"user_email"`
}

// GroupPreferences defines model for GroupPreferences.
type GroupPreferences struct {
	// The default label to use when none is supplied (can be overridden per project).
	DefaultLabel *string `json:"defaultLabel"`

	// Group specific ignored issues (in addition to any project specific ignored issues).
	IgnoredIssues        *[]IgnoredIssue        `json:"ignoredIssues"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// The preferences for a given group which may contain several projects.
type GroupPreferencesResponse struct {
	// The id of the group these preferences apply to.
	GroupId openapi_types.UUID `json:"groupId"`

	// The preference settings
	Preferences struct {
		// The default label to use when none is supplied (can be overridden per project).
		DefaultLabel *string `json:"defaultLabel"`

		// Group specific ignored issues (in addition to any project specific ignored issues).
		IgnoredIssues *[]IgnoredIssue `json:"ignoredIssues"`
	} `json:"preferences"`
}

// Health defines model for Health.
type Health struct {
	Response string `json:"response"`
}

// Issues ignored from package score
type IgnoredIssue struct {
	Id     string        `json:"id"`
	Reason IgnoredReason `json:"reason"`
	Tag    string        `json:"tag"`
}

// IgnoredReason defines model for IgnoredReason.
type IgnoredReason string

// How / where was the package ingested?
type IngestionSource string

// A single package issue.
type Issue struct {
	Description string `json:"description"`

	// Risk domains.
	Domain RiskDomain `json:"domain"`
	Id     *string    `json:"id"`

	// Issue severity.
	Severity RiskLevel `json:"severity"`
	Tag      *string   `json:"tag"`
	Title    string    `json:"title"`
}

// Count of issues for each severity.
type IssueImpacts struct {
	Critical *uint32 `json:"critical,omitempty"`
	High     *uint32 `json:"high,omitempty"`
	Low      *uint32 `json:"low,omitempty"`
	Medium   *uint32 `json:"medium,omitempty"`
}

// IssueStatusCounts defines model for IssueStatusCounts.
type IssueStatusCounts struct {
	Accept      uint32 `json:"accept"`
	NotRelevant uint32 `json:"notRelevant"`
	Untagged    uint32 `json:"untagged"`
	WillFix     uint32 `json:"willFix"`
}

// IssueStatusesStatsBlock defines model for IssueStatusesStatsBlock.
type IssueStatusesStatsBlock struct {
	Counts IssueStatusCounts `json:"counts"`
}

// Issue description.
type IssuesListItem struct {
	Description string        `json:"description"`
	Id          *string       `json:"id"`
	Ignored     IgnoredReason `json:"ignored"`

	// Issue severity.
	Impact   RiskLevel `json:"impact"`
	RiskType RiskType  `json:"riskType"`
	Score    float32   `json:"score"`
	Tag      *string   `json:"tag"`
	Title    string    `json:"title"`
}

// Metadata about a job
type JobDescriptor struct {
	Date            string              `json:"date"`
	Ecosystem       string              `json:"ecosystem"`
	JobId           openapi_types.UUID  `json:"job_id"`
	Label           string              `json:"label"`
	Msg             string              `json:"msg"`
	NumDependencies uint32              `json:"num_dependencies"`
	NumIncomplete   *uint32             `json:"num_incomplete,omitempty"`
	Packages        []PackageDescriptor `json:"packages"`
	Pass            bool                `json:"pass"`
	Project         string              `json:"project"`
	Score           float64             `json:"score"`
}

// JobScore defines model for JobScore.
type JobScore struct {
	// Whether or not all of the underlying job's packages have completed processing. Note that there is a stop-gap where packages with nonstandard versions are allowed to be missing from Redis and still be considered complete.
	Complete bool    `json:"complete"`
	Value    float32 `json:"value"`
}

// JobStatusResponseVariant defines model for JobStatusResponseVariant.
type JobStatusResponseVariant interface{}

// Data returned when querying the job status endpoint
type JobStatusResponseForPackageStatus struct {
	// The action to take if the job fails
	Action interface{} `json:"action"`

	// The time the job started in epoch seconds
	CreatedAt int64 `json:"created_at"`

	// The language ecosystem TODO: How is this different than package type ( npm, etc ) or language?
	Ecosystem string `json:"ecosystem"`

	// The id of the job processing the top level package
	JobId openapi_types.UUID `json:"job_id"`

	// A label associated with this job, most often a branch name
	Label *string `json:"label"`

	// The last time the job metadata was updated
	LastUpdated uint64 `json:"last_updated"`
	Msg         string `json:"msg"`

	// Dependencies that have not completed processing
	NumIncomplete *uint32 `json:"num_incomplete,omitempty"`

	// The packages that are a part of this job
	Packages []PackageStatus `json:"packages"`
	Pass     bool            `json:"pass"`

	// The id of the project associated with this job
	Project string `json:"project"`

	// The project name
	ProjectName string `json:"project_name"`

	// The current score
	Score float64 `json:"score"`

	// The job status
	Status interface{} `json:"status"`

	// The currently configured threshholds for this job. If the scores fall below these thresholds, then the client should undertake the action spelled out by the action field.
	Thresholds struct {
		Author        float32 `json:"author"`
		Engineering   float32 `json:"engineering"`
		License       float32 `json:"license"`
		Malicious     float32 `json:"malicious"`
		Total         float32 `json:"total"`
		Vulnerability float32 `json:"vulnerability"`
	} `json:"thresholds"`

	// The user email
	UserEmail string `json:"user_email"`

	// The id of the user submitting the job
	UserId openapi_types.UUID `json:"user_id"`
}

// Data returned when querying the job status endpoint
type JobStatusResponseForPackageStatusExtended struct {
	// The action to take if the job fails
	Action interface{} `json:"action"`

	// The time the job started in epoch seconds
	CreatedAt int64 `json:"created_at"`

	// The language ecosystem TODO: How is this different than package type ( npm, etc ) or language?
	Ecosystem string `json:"ecosystem"`

	// The id of the job processing the top level package
	JobId openapi_types.UUID `json:"job_id"`

	// A label associated with this job, most often a branch name
	Label *string `json:"label"`

	// The last time the job metadata was updated
	LastUpdated uint64 `json:"last_updated"`
	Msg         string `json:"msg"`

	// Dependencies that have not completed processing
	NumIncomplete *uint32 `json:"num_incomplete,omitempty"`

	// The packages that are a part of this job
	Packages []PackageStatusExtended `json:"packages"`
	Pass     bool                    `json:"pass"`

	// The id of the project associated with this job
	Project string `json:"project"`

	// The project name
	ProjectName string `json:"project_name"`

	// The current score
	Score float64 `json:"score"`

	// The job status
	Status interface{} `json:"status"`

	// The currently configured threshholds for this job. If the scores fall below these thresholds, then the client should undertake the action spelled out by the action field.
	Thresholds struct {
		Author        float32 `json:"author"`
		Engineering   float32 `json:"engineering"`
		License       float32 `json:"license"`
		Malicious     float32 `json:"malicious"`
		Total         float32 `json:"total"`
		Vulnerability float32 `json:"vulnerability"`
	} `json:"thresholds"`

	// The user email
	UserEmail string `json:"user_email"`

	// The id of the user submitting the job
	UserId openapi_types.UUID `json:"user_id"`
}

// A document describing an error.
type JsonErrorResponse struct {
	// Information about an error that occurred while servicing a request.
	Error struct {
		// The class of error.
		ApiError interface{} `json:"apiError"`

		// The HTTP error code.
		Code uint16 `json:"code"`

		// A general description of this class of error.
		Description string `json:"description"`

		// A unique ID for this error.
		ErrorId openapi_types.UUID `json:"error_id"`

		// A reason for the error.
		Reason string `json:"reason"`
	} `json:"error"`
}

// A description of an error that occurred while servicing a request.
type JsonErrorResponseBody struct {
	// The class of error.
	ApiError interface{} `json:"apiError"`

	// The HTTP error code.
	Code uint16 `json:"code"`

	// A general description of this class of error.
	Description string `json:"description"`

	// A unique ID for this error.
	ErrorId openapi_types.UUID `json:"error_id"`

	// A reason for the error.
	Reason string `json:"reason"`
}

// LicensesStatsBlock defines model for LicensesStatsBlock.
type LicensesStatsBlock struct {
	Counts LicensesStatsBlock_Counts `json:"counts"`
}

// LicensesStatsBlock_Counts defines model for LicensesStatsBlock.Counts.
type LicensesStatsBlock_Counts struct {
	AdditionalProperties map[string]uint32 `json:"-"`
}

// ListGroupMembersResponse defines model for ListGroupMembersResponse.
type ListGroupMembersResponse struct {
	Members []GroupMember `json:"members"`
}

// ListUserGroupsResponse defines model for ListUserGroupsResponse.
type ListUserGroupsResponse struct {
	Groups []UserGroup `json:"groups"`
}

// The maintainer information is (currently) sparser than the contributor data
type Maintainer struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Metadata related to a particular pagination request. The values are ultimately influenced by the `FilterProjects` value (sorting / paging won't affect these values).
type Metadata struct {
	// The distinct ecosystems across all projects satisfying the given query.
	Ecosystems []PackageType `json:"ecosystems"`

	// The groups that have existing projects satisfying the given query.
	Groups Metadata_Groups `json:"groups"`
}

// The groups that have existing projects satisfying the given query.
type Metadata_Groups struct {
	AdditionalProperties map[string]string `json:"-"`
}

// PackageAuthorsResponse defines model for PackageAuthorsResponse.
type PackageAuthorsResponse struct {
	Contributors []Author     `json:"contributors"`
	Maintainers  []Maintainer `json:"maintainers"`
}

// Describes a package in the system
type PackageDescriptor struct {
	Name string `json:"name"`

	// The package ecosystem
	Type    PackageType `json:"type"`
	Version string      `json:"version"`
}

// PackageReleaseData defines model for PackageReleaseData.
type PackageReleaseData struct {
	FirstReleaseDate time.Time `json:"first_release_date"`
	LastReleaseDate  time.Time `json:"last_release_date"`
}

// PackageSearchListing defines model for PackageSearchListing.
type PackageSearchListing struct {
	Description *string `json:"description"`

	// Count of issues for each severity.
	IssueImpacts  IssueImpacts `json:"issueImpacts"`
	Name          string       `json:"name"`
	PublishedDate string       `json:"publishedDate"`
	Registry      string       `json:"registry"`

	// Risk scores by domain.
	RiskScores RiskScores `json:"riskScores"`
	Version    string     `json:"version"`
}

// PackageSpecifier defines model for PackageSpecifier.
type PackageSpecifier struct {
	Name     string `json:"name"`
	Registry string `json:"registry"`
	Version  string `json:"version"`
}

// Basic core package meta data
type PackageStatus struct {
	// Last updates, as epoch seconds
	LastUpdated uint64 `json:"last_updated"`

	// Package license
	License *string `json:"license"`

	// Name of the package
	Name string `json:"name"`

	// Number of dependencies
	NumDependencies uint32 `json:"num_dependencies"`

	// Number of vulnerabilities found in this package and all transitive dependencies
	NumVulnerabilities uint32 `json:"num_vulnerabilities"`

	// The overall quality score of the package
	PackageScore *float64 `json:"package_score"`

	// Package processing status
	Status interface{} `json:"status"`

	// Package version
	Version string `json:"version"`
}

// Package metadata with extended info info
type PackageStatusExtended struct {
	// Dependencies of this package
	Dependencies PackageStatusExtended_Dependencies `json:"dependencies"`

	// Any issues found that may need action, but aren't in and of themselves vulnerabilities
	Issues []Issue `json:"issues"`

	// Last updates, as epoch seconds
	LastUpdated uint64 `json:"last_updated"`

	// Package license
	License *string `json:"license"`

	// Name of the package
	Name string `json:"name"`

	// Number of dependencies
	NumDependencies uint32 `json:"num_dependencies"`

	// Number of vulnerabilities found in this package and all transitive dependencies
	NumVulnerabilities uint32 `json:"num_vulnerabilities"`

	// The overall quality score of the package
	PackageScore *float64                          `json:"package_score"`
	RiskVectors  PackageStatusExtended_RiskVectors `json:"riskVectors"`

	// Package processing status
	Status interface{} `json:"status"`

	// The package_type, npm, etc.
	Type interface{} `json:"type"`

	// Package version
	Version string `json:"version"`
}

// Dependencies of this package
type PackageStatusExtended_Dependencies struct {
	AdditionalProperties map[string]string `json:"-"`
}

// PackageStatusExtended_RiskVectors defines model for PackageStatusExtended.RiskVectors.
type PackageStatusExtended_RiskVectors struct {
	AdditionalProperties map[string]float64 `json:"-"`
}

// The package ecosystem
type PackageType string

// Requests
type PaginateDirection string

// Responses
type PaginatedForProjectListEntryAndMetadata struct {
	// Indication of whether the current query has more values past the last element in `values`.
	HasMore bool `json:"has_more"`

	// Optional metadata that can be sent about the result set, such as a total count estimate.
	Metadata *struct {
		// The distinct ecosystems across all projects satisfying the given query.
		Ecosystems []PackageType `json:"ecosystems"`

		// The groups that have existing projects satisfying the given query.
		Groups PaginatedForProjectListEntryAndMetadata_Metadata_Groups `json:"groups"`
	} `json:"metadata"`

	// The curent page of values.
	Values []ProjectListEntry `json:"values"`
}

// The groups that have existing projects satisfying the given query.
type PaginatedForProjectListEntryAndMetadata_Metadata_Groups struct {
	AdditionalProperties map[string]string `json:"-"`
}

// The project fields on which users can sort.
type ProjectField string

// Project folder response format/data
type ProjectFolderResponse struct {
	CreatedAt time.Time                `json:"createdAt"`
	Id        openapi_types.UUID       `json:"id"`
	IsDefault bool                     `json:"isDefault"`
	Name      string                   `json:"name"`
	Projects  []ProjectListingResponse `json:"projects"`
	UpdatedAt time.Time                `json:"updatedAt"`
	UserId    openapi_types.UUID       `json:"userId"`
}

// ProjectListEntry defines model for ProjectListEntry.
type ProjectListEntry struct {
	// Project created time
	CreatedAt time.Time `json:"created_at"`

	// Project ecosystem (as of latest job run)
	Ecosystem *interface{} `json:"ecosystem"`

	// If this is a project belonging to the group, then the group name, otherwise None.
	GroupName *string `json:"group_name"`

	// Project id
	Id openapi_types.UUID `json:"id"`

	// Project name
	Name string `json:"name"`

	// The total risk score of the latest analysis job (if there is one).
	TotalRiskScore *struct {
		// Whether or not all of the underlying job's packages have completed processing. Note that there is a stop-gap where packages with nonstandard versions are allowed to be missing from Redis and still be considered complete.
		Complete bool    `json:"complete"`
		Value    float32 `json:"value"`
	} `json:"total_risk_score"`

	// Project updated time
	UpdatedAt time.Time `json:"updated_at"`
}

// ProjectListingResponse defines model for ProjectListingResponse.
type ProjectListingResponse struct {
	CreatedAt time.Time          `json:"createdAt"`
	GroupName *string            `json:"groupName"`
	Id        openapi_types.UUID `json:"id"`

	// How / where was the package ingested?
	IngestionSource IngestionSource `json:"ingestionSource"`

	// Count of issues for each severity.
	IssueImpacts IssueImpacts     `json:"issueImpacts"`
	Issues       []IssuesListItem `json:"issues"`
	Label        *string          `json:"label"`
	Name         string           `json:"name"`
	Registry     *string          `json:"registry"`
	Stats        *struct {
		Dependencies  DependenciesStatsBlock  `json:"dependencies"`
		IssueStatuses IssueStatusesStatsBlock `json:"issueStatuses"`
		Licenses      LicensesStatsBlock      `json:"licenses"`
	} `json:"stats"`
	TotalRiskScore float32    `json:"totalRiskScore"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

// ProjectPreferences defines model for ProjectPreferences.
type ProjectPreferences struct {
	// The default label to use when none is supplied.
	DefaultLabel *string `json:"defaultLabel"`

	// Project specific ignored issues.
	IgnoredIssues *[]IgnoredIssue `json:"ignoredIssues"`

	// The risk thresholds to apply.
	Thresholds struct {
		// Capture the user threshold settings
		Author ThresholdDescriptor `json:"author"`

		// Capture the user threshold settings
		Engineering ThresholdDescriptor `json:"engineering"`

		// Capture the user threshold settings
		License ThresholdDescriptor `json:"license"`

		// Capture the user threshold settings
		MaliciousCode ThresholdDescriptor `json:"maliciousCode"`

		// Capture the user threshold settings
		Total ThresholdDescriptor `json:"total"`

		// Capture the user threshold settings
		Vulnerability ThresholdDescriptor `json:"vulnerability"`
	} `json:"thresholds"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// The preferences for a given project.
type ProjectPreferencesResponse struct {
	// The preference settings
	Preferences struct {
		// The default label to use when none is supplied.
		DefaultLabel *string `json:"defaultLabel"`

		// Project specific ignored issues.
		IgnoredIssues *[]IgnoredIssue `json:"ignoredIssues"`

		// The risk thresholds to apply.
		Thresholds struct {
			// Capture the user threshold settings
			Author ThresholdDescriptor `json:"author"`

			// Capture the user threshold settings
			Engineering ThresholdDescriptor `json:"engineering"`

			// Capture the user threshold settings
			License ThresholdDescriptor `json:"license"`

			// Capture the user threshold settings
			MaliciousCode ThresholdDescriptor `json:"maliciousCode"`

			// Capture the user threshold settings
			Total ThresholdDescriptor `json:"total"`

			// Capture the user threshold settings
			Vulnerability ThresholdDescriptor `json:"vulnerability"`
		} `json:"thresholds"`
	} `json:"preferences"`

	// The id of the project these preferences apply to.
	ProjectId openapi_types.UUID `json:"projectId"`
}

// ProjectResponse defines model for ProjectResponse.
type ProjectResponse struct {
	CreatedAt    time.Time             `json:"createdAt"`
	Dependencies []FullPackageInternal `json:"dependencies"`
	Id           openapi_types.UUID    `json:"id"`

	// How / where was the package ingested?
	IngestionSource IngestionSource `json:"ingestionSource"`

	// Count of issues for each severity.
	IssueImpacts IssueImpacts     `json:"issueImpacts"`
	Issues       []IssuesListItem `json:"issues"`
	Label        *string          `json:"label"`

	// The created time of the latest job, if at least one job has been created for this project.
	LatestJobCreatedAt *time.Time `json:"latestJobCreatedAt"`
	Name               string     `json:"name"`
	Registry           *string    `json:"registry"`

	// Risk scores by domain.
	RiskScores           RiskScores `json:"riskScores"`
	RiskThresholdActions *struct {
		// Capture the user threshold settings
		Author ThresholdDescriptor `json:"author"`

		// Capture the user threshold settings
		Engineering ThresholdDescriptor `json:"engineering"`

		// Capture the user threshold settings
		License ThresholdDescriptor `json:"license"`

		// Capture the user threshold settings
		MaliciousCode ThresholdDescriptor `json:"maliciousCode"`

		// Capture the user threshold settings
		Total ThresholdDescriptor `json:"total"`

		// Capture the user threshold settings
		Vulnerability ThresholdDescriptor `json:"vulnerability"`
	} `json:"riskThresholdActions"`
	Stats          ProjectStatsBlock `json:"stats"`
	TotalRiskScore float32           `json:"totalRiskScore"`
	UpdatedAt      *time.Time        `json:"updatedAt"`
}

// ProjectStatsBlock defines model for ProjectStatsBlock.
type ProjectStatsBlock struct {
	Dependencies  DependenciesStatsBlock  `json:"dependencies"`
	IssueStatuses IssueStatusesStatsBlock `json:"issueStatuses"`
	Licenses      LicensesStatsBlock      `json:"licenses"`
}

// Summary response for a project
type ProjectSummaryResponse struct {
	// When the project was created
	CreatedAt time.Time `json:"created_at"`

	// The ecosystem of the project; determined by its latest job
	Ecosystem *string `json:"ecosystem"`

	// The project's group's name, if this is a group project
	GroupName *string `json:"group_name"`

	// The project id
	Id openapi_types.UUID `json:"id"`

	// The project name
	Name string `json:"name"`

	// When the project was updated
	UpdatedAt time.Time `json:"updated_at"`
}

// Rick cut off thresholds for a project
type ProjectThresholds struct {
	Author        float32 `json:"author"`
	Engineering   float32 `json:"engineering"`
	License       float32 `json:"license"`
	Malicious     float32 `json:"malicious"`
	Total         float32 `json:"total"`
	Vulnerability float32 `json:"vulnerability"`
}

// RequestManagerProjectHistoryJobResponse defines model for RequestManagerProjectHistoryJobResponse.
type RequestManagerProjectHistoryJobResponse struct {
	Date            time.Time          `json:"date"`
	JobId           openapi_types.UUID `json:"job_id"`
	Label           *string            `json:"label"`
	NumDependencies int32              `json:"num_dependencies"`
}

// ResultOfValidatedGroupNamePathOrValidatedGroupNameError defines model for Result_of_ValidatedGroupNamePath_or_ValidatedGroupNameError.
type ResultOfValidatedGroupNamePathOrValidatedGroupNameError interface{}

// Risk domains.
type RiskDomain string

// Issue severity.
type RiskLevel string

// Risk scores by domain.
type RiskScores struct {
	Author        float32 `json:"author"`
	Engineering   float32 `json:"engineering"`
	License       float32 `json:"license"`
	MaliciousCode float32 `json:"malicious_code"`
	Total         float32 `json:"total"`
	Vulnerability float32 `json:"vulnerability"`
}

// Capture the user threshold settings
type RiskThresholds struct {
	// Capture the user threshold settings
	Author ThresholdDescriptor `json:"author"`

	// Capture the user threshold settings
	Engineering ThresholdDescriptor `json:"engineering"`

	// Capture the user threshold settings
	License ThresholdDescriptor `json:"license"`

	// Capture the user threshold settings
	MaliciousCode ThresholdDescriptor `json:"maliciousCode"`

	// Capture the user threshold settings
	Total ThresholdDescriptor `json:"total"`

	// Capture the user threshold settings
	Vulnerability ThresholdDescriptor `json:"vulnerability"`
}

// RiskType defines model for RiskType.
type RiskType string

// ScoredVersion defines model for ScoredVersion.
type ScoredVersion struct {
	TotalRiskScore *float32 `json:"total_risk_score"`
	Version        string   `json:"version"`
}

// SortDirection defines model for SortDirection.
type SortDirection string

// Did the processing of the Package or Job complete successfully
type Status string

// Submit Package for analysis
type SubmitPackageRequest struct {
	// The group that owns the project, if applicable
	GroupName *string `json:"group_name"`

	// Was this submitted by a user interactively and not a CI?
	IsUser bool `json:"is_user"`

	// A label for this package. Often it's the branch.
	Label string `json:"label"`

	// The subpackage dependencies of this package
	Packages []PackageDescriptor `json:"packages"`

	// The id of the project this top level package should be associated with
	Project string `json:"project"`

	// The 'type' of package, NPM, RubyGem, etc
	//TODO: this probably needs to be a *string
	Type interface{} `json:"type"`
}

// Initial response after package has been submitted
type SubmitPackageResponse struct {
	// The id of the job processing the package
	JobId openapi_types.UUID `json:"job_id"`
}

// Capture the user threshold settings
type ThresholdDescriptor struct {
	// When phylum is integrated with CI, what action should be taken when the quality threshold of a package falls below the limit.
	Action    ThresholdViolationAction `json:"action"`
	Active    bool                     `json:"active"`
	Threshold float32                  `json:"threshold"`
}

// When phylum is integrated with CI, what action should be taken when the quality threshold of a package falls below the limit.
type ThresholdViolationAction string

// An explanation of why an action was rejected for account tier reasons.
type TierExceeded interface{}

// UserGroup defines model for UserGroup.
type UserGroup struct {
	CreatedAt    time.Time `json:"created_at"`
	GroupName    string    `json:"group_name"`
	IsAdmin      *bool     `json:"is_admin,omitempty"`
	IsOwner      *bool     `json:"is_owner,omitempty"`
	LastModified time.Time `json:"last_modified"`
	OwnerEmail   string    `json:"owner_email"`
}

// Preferences for a given user
type UserPreferencesResponse struct {
	Preferences UserPreferencesResponse_Preferences `json:"preferences"`

	// The id of the user the preferences apply to
	UserId openapi_types.UUID `json:"userId"`
}

// UserPreferencesResponse_Preferences defines model for UserPreferencesResponse.Preferences.
type UserPreferencesResponse_Preferences struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// A group name that has been checked and validated.
type ValidatedGroupName struct {
	GroupName string `json:"group_name"`
}

// An error that occured during validation of the group name
type ValidatedGroupNameError string

// A uri path component that has also been validated as a group name.
type ValidatedGroupNamePath struct {
	GroupName string `json:"group_name"`
}

// The package ecosystem
type WrapForPackageType = PackageType

// GetUserJobsParams defines parameters for GetUserJobs.
type GetUserJobsParams struct {
	Limit   uint16 `form:"limit" json:"limit"`
	Verbose *bool  `form:"verbose,omitempty" json:"verbose,omitempty"`
}

// StartJobJSONBody defines parameters for StartJob.
type StartJobJSONBody = SubmitPackageRequest

// StartJobPutJSONBody defines parameters for StartJobPut.
type StartJobPutJSONBody = SubmitPackageRequest

// GetJobStatusParams defines parameters for GetJobStatus.
type GetJobStatusParams struct {
	Verbose *bool `form:"verbose,omitempty" json:"verbose,omitempty"`
}

// PackagesSearchEndpointParams defines parameters for PackagesSearchEndpoint.
type PackagesSearchEndpointParams struct {
	Search string `form:"search" json:"search"`
}

// PackagesVersionSearchEndpointParams defines parameters for PackagesVersionSearchEndpoint.
type PackagesVersionSearchEndpointParams struct {
	Search *string `form:"search,omitempty" json:"search,omitempty"`
}

// ProjectsCreateProjectJSONBody defines parameters for ProjectsCreateProject.
type ProjectsCreateProjectJSONBody = CreateProjectRequest

// ProjectsGetEndpointParams defines parameters for ProjectsGetEndpoint.
type ProjectsGetEndpointParams struct {
	Label *string `form:"label,omitempty" json:"label,omitempty"`
}

// ProjectsUpdateProjectJSONBody defines parameters for ProjectsUpdateProject.
type ProjectsUpdateProjectJSONBody = CreateProjectRequest

// ProjectsHistoryEndpointParams defines parameters for ProjectsHistoryEndpoint.
type ProjectsHistoryEndpointParams struct {
	Label *string `form:"label,omitempty" json:"label,omitempty"`
}

// GroupsPostCreateGroupJSONBody defines parameters for GroupsPostCreateGroup.
type GroupsPostCreateGroupJSONBody = ValidatedGroupName

// GroupsGetProjectParams defines parameters for GroupsGetProject.
type GroupsGetProjectParams struct {
	Label *string `form:"label,omitempty" json:"label,omitempty"`
}

// GroupsGetProjectHistoryParams defines parameters for GroupsGetProjectHistory.
type GroupsGetProjectHistoryParams struct {
	Label *string `form:"label,omitempty" json:"label,omitempty"`
}

// UpdateUserPreferencesEndpointJSONBody defines parameters for UpdateUserPreferencesEndpoint.
type UpdateUserPreferencesEndpointJSONBody struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// UpdateGroupPreferencesEndpointJSONBody defines parameters for UpdateGroupPreferencesEndpoint.
type UpdateGroupPreferencesEndpointJSONBody = GroupPreferences

// UpdateProjectPreferencesEndpointJSONBody defines parameters for UpdateProjectPreferencesEndpoint.
type UpdateProjectPreferencesEndpointJSONBody = ProjectPreferences

// ProjectsListProjectsParams defines parameters for ProjectsListProjects.
type ProjectsListProjectsParams struct {
	// A limit on the number of objects to be returned, between 1 and 100. Default is Paginate::<Id>::LIMIT_DEFAULT
	Limit uint32 `form:"limit" json:"limit"`

	// The direction of pagination, i.e. given some sorting, get the next or previous page. Default is PaginateDirection::Forward
	Direction PaginateDirection `form:"direction" json:"direction"`

	// A cursor for use in pagination. This is an object ID that defines your place in the list. For instance, if you make a request and receive a list of objects ending with <end_obj_uuid>, your subsequent call can include cursor=<end_obj_uuid> in order to fetch the next page of the list. Similarly, cursor=<first_obj_uuid>&direction=backward will fetch the previous page.
	//
	// Not specifying a cursor will get you the _first_ page if direction is forward, and the _last_ page if direction is backward.
	//
	// Default is None
	Cursor *openapi_types.UUID `form:"cursor,omitempty" json:"cursor,omitempty"`

	// A flag for requesting arbitrary metadata that may be offered depending on the endpoint. This may include an estimate for the total count, or a range of distinct values for relevant filtering.
	Metadata      bool          `form:"metadata" json:"metadata"`
	Field         ProjectField  `form:"field" json:"field"`
	SortDirection SortDirection `form:"direction" json:"direction"`

	// Only include projects which contain the given string (case insensitively).
	NameContains *string `form:"name_contains,omitempty" json:"name_contains,omitempty"`

	// Only include projects in this ecosystem.
	Ecosystem *WrapForPackageType `form:"ecosystem,omitempty" json:"ecosystem,omitempty"`

	// Only include projects under this group name. No auth errors are thrown here; if the user doesn't exist in a group with this name, just return no projects.
	Group *string `form:"group,omitempty" json:"group,omitempty"`
}

// StartJobJSONRequestBody defines body for StartJob for application/json ContentType.
type StartJobJSONRequestBody = StartJobJSONBody

// StartJobPutJSONRequestBody defines body for StartJobPut for application/json ContentType.
type StartJobPutJSONRequestBody = StartJobPutJSONBody

// ProjectsCreateProjectJSONRequestBody defines body for ProjectsCreateProject for application/json ContentType.
type ProjectsCreateProjectJSONRequestBody = ProjectsCreateProjectJSONBody

// ProjectsUpdateProjectJSONRequestBody defines body for ProjectsUpdateProject for application/json ContentType.
type ProjectsUpdateProjectJSONRequestBody = ProjectsUpdateProjectJSONBody

// GroupsPostCreateGroupJSONRequestBody defines body for GroupsPostCreateGroup for application/json ContentType.
type GroupsPostCreateGroupJSONRequestBody = GroupsPostCreateGroupJSONBody

// UpdateUserPreferencesEndpointJSONRequestBody defines body for UpdateUserPreferencesEndpoint for application/json ContentType.
type UpdateUserPreferencesEndpointJSONRequestBody UpdateUserPreferencesEndpointJSONBody

// UpdateGroupPreferencesEndpointJSONRequestBody defines body for UpdateGroupPreferencesEndpoint for application/json ContentType.
type UpdateGroupPreferencesEndpointJSONRequestBody = UpdateGroupPreferencesEndpointJSONBody

// UpdateProjectPreferencesEndpointJSONRequestBody defines body for UpdateProjectPreferencesEndpoint for application/json ContentType.
type UpdateProjectPreferencesEndpointJSONRequestBody = UpdateProjectPreferencesEndpointJSONBody

// Getter for additional properties for UpdateUserPreferencesEndpointJSONBody. Returns the specified
// element and whether it was found
func (a UpdateUserPreferencesEndpointJSONBody) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for UpdateUserPreferencesEndpointJSONBody
func (a *UpdateUserPreferencesEndpointJSONBody) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for UpdateUserPreferencesEndpointJSONBody to handle AdditionalProperties
func (a *UpdateUserPreferencesEndpointJSONBody) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for UpdateUserPreferencesEndpointJSONBody to handle AdditionalProperties
func (a UpdateUserPreferencesEndpointJSONBody) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for GroupPreferences. Returns the specified
// element and whether it was found
func (a GroupPreferences) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for GroupPreferences
func (a *GroupPreferences) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for GroupPreferences to handle AdditionalProperties
func (a *GroupPreferences) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["defaultLabel"]; found {
		err = json.Unmarshal(raw, &a.DefaultLabel)
		if err != nil {
			return fmt.Errorf("error reading 'defaultLabel': %w", err)
		}
		delete(object, "defaultLabel")
	}

	if raw, found := object["ignoredIssues"]; found {
		err = json.Unmarshal(raw, &a.IgnoredIssues)
		if err != nil {
			return fmt.Errorf("error reading 'ignoredIssues': %w", err)
		}
		delete(object, "ignoredIssues")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for GroupPreferences to handle AdditionalProperties
func (a GroupPreferences) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	if a.DefaultLabel != nil {
		object["defaultLabel"], err = json.Marshal(a.DefaultLabel)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'defaultLabel': %w", err)
		}
	}

	if a.IgnoredIssues != nil {
		object["ignoredIssues"], err = json.Marshal(a.IgnoredIssues)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'ignoredIssues': %w", err)
		}
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for LicensesStatsBlock_Counts. Returns the specified
// element and whether it was found
func (a LicensesStatsBlock_Counts) Get(fieldName string) (value uint32, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for LicensesStatsBlock_Counts
func (a *LicensesStatsBlock_Counts) Set(fieldName string, value uint32) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]uint32)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for LicensesStatsBlock_Counts to handle AdditionalProperties
func (a *LicensesStatsBlock_Counts) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]uint32)
		for fieldName, fieldBuf := range object {
			var fieldVal uint32
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for LicensesStatsBlock_Counts to handle AdditionalProperties
func (a LicensesStatsBlock_Counts) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for Metadata_Groups. Returns the specified
// element and whether it was found
func (a Metadata_Groups) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for Metadata_Groups
func (a *Metadata_Groups) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for Metadata_Groups to handle AdditionalProperties
func (a *Metadata_Groups) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for Metadata_Groups to handle AdditionalProperties
func (a Metadata_Groups) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for PackageStatusExtended_Dependencies. Returns the specified
// element and whether it was found
func (a PackageStatusExtended_Dependencies) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for PackageStatusExtended_Dependencies
func (a *PackageStatusExtended_Dependencies) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for PackageStatusExtended_Dependencies to handle AdditionalProperties
func (a *PackageStatusExtended_Dependencies) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for PackageStatusExtended_Dependencies to handle AdditionalProperties
func (a PackageStatusExtended_Dependencies) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for PackageStatusExtended_RiskVectors. Returns the specified
// element and whether it was found
func (a PackageStatusExtended_RiskVectors) Get(fieldName string) (value float64, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for PackageStatusExtended_RiskVectors
func (a *PackageStatusExtended_RiskVectors) Set(fieldName string, value float64) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]float64)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for PackageStatusExtended_RiskVectors to handle AdditionalProperties
func (a *PackageStatusExtended_RiskVectors) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]float64)
		for fieldName, fieldBuf := range object {
			var fieldVal float64
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for PackageStatusExtended_RiskVectors to handle AdditionalProperties
func (a PackageStatusExtended_RiskVectors) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for PaginatedForProjectListEntryAndMetadata_Metadata_Groups. Returns the specified
// element and whether it was found
func (a PaginatedForProjectListEntryAndMetadata_Metadata_Groups) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for PaginatedForProjectListEntryAndMetadata_Metadata_Groups
func (a *PaginatedForProjectListEntryAndMetadata_Metadata_Groups) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for PaginatedForProjectListEntryAndMetadata_Metadata_Groups to handle AdditionalProperties
func (a *PaginatedForProjectListEntryAndMetadata_Metadata_Groups) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for PaginatedForProjectListEntryAndMetadata_Metadata_Groups to handle AdditionalProperties
func (a PaginatedForProjectListEntryAndMetadata_Metadata_Groups) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ProjectPreferences. Returns the specified
// element and whether it was found
func (a ProjectPreferences) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ProjectPreferences
func (a *ProjectPreferences) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ProjectPreferences to handle AdditionalProperties
func (a *ProjectPreferences) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["defaultLabel"]; found {
		err = json.Unmarshal(raw, &a.DefaultLabel)
		if err != nil {
			return fmt.Errorf("error reading 'defaultLabel': %w", err)
		}
		delete(object, "defaultLabel")
	}

	if raw, found := object["ignoredIssues"]; found {
		err = json.Unmarshal(raw, &a.IgnoredIssues)
		if err != nil {
			return fmt.Errorf("error reading 'ignoredIssues': %w", err)
		}
		delete(object, "ignoredIssues")
	}

	if raw, found := object["thresholds"]; found {
		err = json.Unmarshal(raw, &a.Thresholds)
		if err != nil {
			return fmt.Errorf("error reading 'thresholds': %w", err)
		}
		delete(object, "thresholds")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ProjectPreferences to handle AdditionalProperties
func (a ProjectPreferences) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	if a.DefaultLabel != nil {
		object["defaultLabel"], err = json.Marshal(a.DefaultLabel)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'defaultLabel': %w", err)
		}
	}

	if a.IgnoredIssues != nil {
		object["ignoredIssues"], err = json.Marshal(a.IgnoredIssues)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'ignoredIssues': %w", err)
		}
	}

	object["thresholds"], err = json.Marshal(a.Thresholds)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'thresholds': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for UserPreferencesResponse_Preferences. Returns the specified
// element and whether it was found
func (a UserPreferencesResponse_Preferences) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for UserPreferencesResponse_Preferences
func (a *UserPreferencesResponse_Preferences) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for UserPreferencesResponse_Preferences to handle AdditionalProperties
func (a *UserPreferencesResponse_Preferences) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for UserPreferencesResponse_Preferences to handle AdditionalProperties
func (a UserPreferencesResponse_Preferences) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}
