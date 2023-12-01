package types

import (
	"time"
)

type ContentReportSearchQuery struct {
	ContentType    string `schema:"content_type"`
	ReportState    string `schema:"report_state"`
	ReportedBy     string `schema:"reported_by"`
	ReportedUser   string `schema:"reported_user"`
	ResolvedBy     string `schema:"resolved_by"`
	Page           int64  `schema:"page"`
	PageSize       int64  `schema:"page_size"`
	OrderBy        string `schema:"order_by"`
	OrderDirection string `schema:"order_direction"`
	IncludeTotal   bool   `schema:"include_total"`
}

type ContentReport struct {
	ID                int64        `json:"id"`
	ContentRef        string       `json:"content_ref"`
	ReportState       string       `json:"report_state"`
	ReportedBy        *UserProfile `json:"reported_by"`
	ReportReason      string       `json:"report_reason"`
	AdditionalContext string       `json:"additional_context"`
	ReportedUser      *UserProfile `json:"reported_user"`
	ResolvedBy        *UserProfile `json:"resolved_by"`
	ResolvedAt        *time.Time   `json:"resolved_at"`
	ActionTaken       string       `json:"action_taken"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type ContentReportSearchResponse struct {
	Reports []*ContentReport `json:"reports"`
	Total   int64            `json:"total"`
}

type SubmittedContentReport struct {
	ContentRef        string `json:"content_ref"`
	ReportReason      string `json:"reason"`
	AdditionalContext string `json:"context"`
}
