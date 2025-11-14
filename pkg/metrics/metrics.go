package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
		[]string{"path"},
	)

	PRCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pr_created_total",
			Help: "Total number of pull requests created",
		},
		[]string{"team"},
	)

	PRWithReviewersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pr_with_reviewers_total",
			Help: "Total number of PRs by reviewers count",
		},
		[]string{"team", "reviewers_count"},
	)

	ReviewAssignmentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_assignments_total",
			Help: "Total number of review assignments",
		},
		[]string{"reviewer_id", "team"},
	)

	ReviewReassignTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_reassign_total",
			Help: "Total number of review reassignments",
		},
		[]string{"old_reviewer_id", "new_reviewer_id", "team"},
	)

	PRNeedMoreReviewersSetTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pr_need_more_reviewers_set_total",
			Help: "Total number of times needMoreReviewers flag was set",
		},
		[]string{"team"},
	)

	PRNeedMoreReviewersClearedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pr_need_more_reviewers_cleared_total",
			Help: "Total number of times needMoreReviewers flag was cleared",
		},
		[]string{"team"},
	)

	PRNeedMoreReviewersActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pr_need_more_reviewers_active",
			Help: "Current number of PRs with needMoreReviewers flag set",
		},
		[]string{"team"},
	)

	AssignmentNoCandidatesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "assignment_no_candidates_total",
			Help: "Total number of times no candidates were found for assignment",
		},
		[]string{"event"},
	)

	AssignmentCandidatesCount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "assignment_candidates_count",
			Help:    "Number of candidates available for assignment",
			Buckets: []float64{0, 1, 2, 3, 5, 10, 20},
		},
		[]string{"event"},
	)

	InvariantPRReviewersCountViolations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "invariant_pr_reviewers_count_violations_total",
			Help: "Total number of PRs with >2 reviewers",
		},
	)

	InvariantPRMergedMutatedReviewersViolations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "invariant_pr_merged_mutated_reviewers_violations_total",
			Help: "Total number of attempts to mutate reviewers on merged PRs",
		},
	)

	InvariantPRInactiveReviewersAssigned = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "invariant_pr_inactive_reviewers_assigned_total",
			Help: "Total number of PRs with inactive reviewers assigned",
		},
	)

	InvariantPRReviewerNotInAuthorTeam = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "invariant_pr_reviewer_not_in_author_team_total",
			Help: "Total number of PRs where reviewer is not in author's team",
		},
	)

	MergeIdempotentRetries = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "merge_idempotent_retries_total",
			Help: "Total number of idempotent merge retries",
		},
	)

	MergeIdempotentErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "merge_idempotent_errors_total",
			Help: "Total number of errors on idempotent merge",
		},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	DBQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_errors_total",
			Help: "Total number of database query errors",
		},
		[]string{"operation", "table"},
	)

	ActiveUsersGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users_gauge",
			Help: "Current number of active users",
		},
	)

	InactiveUsersGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "inactive_users_gauge",
			Help: "Current number of inactive users",
		},
	)

	TeamsCountGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "teams_count_gauge",
			Help: "Current number of teams",
		},
	)

	PRsOpenGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "prs_open_gauge",
			Help: "Current number of open PRs",
		},
	)

	PRsMergedGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "prs_merged_gauge",
			Help: "Current number of merged PRs",
		},
	)
)
