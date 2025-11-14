package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kxddry/avito-backend-internship-2025/internal/api/generated"
	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

// Server is the HTTP server for the assignment service.
type Server struct {
	service domain.AssignmentService
}

// NewServer creates a new Server.
func NewServer(service domain.AssignmentService) *Server {
	return &Server{
		service: service,
	}
}

// PostPullRequestCreate creates a new pull request.
func (s *Server) PostPullRequestCreate(ctx context.Context, request generated.PostPullRequestCreateRequestObject) (
	generated.PostPullRequestCreateResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	input := &domain.CreatePullRequestInput{
		PullRequestID:   request.Body.PullRequestId,
		PullRequestName: request.Body.PullRequestName,
		AuthorID:        request.Body.AuthorId,
	}

	pr, err := s.service.CreatePullRequest(ctx, input)
	if err != nil {
		if resp, err := s.handlePostPullRequestCreateError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	prPayload := toGeneratedPullRequest(pr)
	return generated.PostPullRequestCreate201JSONResponse{Pr: &prPayload}, nil
}

// PostPullRequestMerge merges a pull request.
func (s *Server) PostPullRequestMerge(ctx context.Context, request generated.PostPullRequestMergeRequestObject) (
	generated.PostPullRequestMergeResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	input := domain.MergePullRequestInput{
		PullRequestID: request.Body.PullRequestId,
	}

	pr, err := s.service.MergePullRequest(ctx, &input)
	if err != nil {
		if resp, err := s.handlePostPullRequestMergeError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	prPayload := toGeneratedPullRequest(pr)
	return generated.PostPullRequestMerge200JSONResponse{Pr: &prPayload}, nil
}

// PostPullRequestReassign reassigns a pull request.
func (s *Server) PostPullRequestReassign(ctx context.Context, request generated.PostPullRequestReassignRequestObject) (
	generated.PostPullRequestReassignResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	input := domain.ReassignPullRequestInput{
		PullRequestID: request.Body.PullRequestId,
		OldUserID:     request.Body.OldUserId,
	}

	result, err := s.service.ReassignPullRequest(ctx, &input)
	if err != nil {
		if resp, err := s.handlePostPullRequestReassignError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	prPayload := toGeneratedPullRequest(&result.PullRequest)
	return generated.PostPullRequestReassign200JSONResponse{
		Pr:         prPayload,
		ReplacedBy: result.ReplacedBy,
	}, nil
}

// PostTeamAdd adds a new team.
func (s *Server) PostTeamAdd(ctx context.Context, request generated.PostTeamAddRequestObject) (
	generated.PostTeamAddResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	team := fromGeneratedTeam(*request.Body)

	created, err := s.service.CreateTeam(ctx, &team)
	if err != nil {
		if resp, err := s.handlePostTeamAddError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	teamPayload := toGeneratedTeam(created)
	return generated.PostTeamAdd201JSONResponse{
		Team: &teamPayload,
	}, nil
}

// GetTeamGet gets a team by name.
func (s *Server) GetTeamGet(ctx context.Context, request generated.GetTeamGetRequestObject) (
	generated.GetTeamGetResponseObject, error) {
	team, err := s.service.GetTeam(ctx, request.Params.TeamName)
	if err != nil {
		if resp, err := s.handleGetTeamGetError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	teamPayload := toGeneratedTeam(team)
	return generated.GetTeamGet200JSONResponse(teamPayload), nil
}

// GetUsersGetReview gets the reviewer assignments for a user.
func (s *Server) GetUsersGetReview(ctx context.Context, request generated.GetUsersGetReviewRequestObject) (
	generated.GetUsersGetReviewResponseObject, error) {
	result, err := s.service.GetReviewerAssignments(ctx, request.Params.UserId)
	if err != nil {
		if err := s.handleGetUsersGetReviewError(err); err != nil {
			return nil, err
		}
		return nil, err
	}

	return generated.GetUsersGetReview200JSONResponse{
		UserId:       result.UserID,
		PullRequests: toGeneratedPullRequestShortList(result.PullRequests),
	}, nil
}

// PostUsersSetIsActive sets the active status of a user.
func (s *Server) PostUsersSetIsActive(ctx context.Context, request generated.PostUsersSetIsActiveRequestObject) (
	generated.PostUsersSetIsActiveResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	input := domain.SetUserIsActiveInput{
		UserID:   request.Body.UserId,
		IsActive: request.Body.IsActive,
	}

	user, err := s.service.SetUserIsActive(ctx, &input)
	if err != nil {
		if resp, err := s.handlePostUsersSetIsActiveError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	userPayload := toGeneratedUser(user)
	return generated.PostUsersSetIsActive200JSONResponse{
		User: &userPayload,
	}, nil
}

// handlePostPullRequestCreateError handles the error for the PostPullRequestCreate endpoint.
func (s *Server) handlePostPullRequestCreateError(err error) (generated.PostPullRequestCreateResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.PostPullRequestCreate404JSONResponse(payload.body), nil
	case http.StatusConflict:
		return generated.PostPullRequestCreate409JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

// handlePostPullRequestMergeError handles the error for the PostPullRequestMerge endpoint.
func (s *Server) handlePostPullRequestMergeError(err error) (generated.PostPullRequestMergeResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.PostPullRequestMerge404JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

// handlePostPullRequestReassignError handles the error for the PostPullRequestReassign endpoint.
func (s *Server) handlePostPullRequestReassignError(err error) (generated.PostPullRequestReassignResponseObject, error) { //nolint:lll
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.PostPullRequestReassign404JSONResponse(payload.body), nil
	case http.StatusConflict:
		return generated.PostPullRequestReassign409JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

func (s *Server) handlePostTeamAddError(err error) (generated.PostTeamAddResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusBadRequest:
		return generated.PostTeamAdd400JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

func (s *Server) handleGetTeamGetError(err error) (generated.GetTeamGetResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.GetTeamGet404JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

func (s *Server) handleGetUsersGetReviewError(err error) error {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return err
	}

	payload := makeErrorPayload(appErr)
	return echo.NewHTTPError(payload.status, payload.body)
}

func (s *Server) handlePostUsersSetIsActiveError(err error) (generated.PostUsersSetIsActiveResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.PostUsersSetIsActive404JSONResponse(payload.body), nil
	case http.StatusUnauthorized:
		return generated.PostUsersSetIsActive401JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

func unwrapDomainError(err error) (*domain.Error, error) {
	if err == nil {
		return nil, nil
	}

	if errors.Is(err, domain.ErrNotImplemented) {
		return nil, echo.NewHTTPError(http.StatusNotImplemented, domain.ErrNotImplemented.Error())
	}

	var appErr *domain.Error
	if errors.As(err, &appErr) {
		return appErr, nil
	}

	return nil, err
}

func toGeneratedPullRequest(pr *domain.PullRequest) generated.PullRequest {
	reviewers := append([]string{}, pr.AssignedReviewers...)

	return generated.PullRequest{
		AssignedReviewers: reviewers,
		AuthorId:          pr.AuthorID,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Name,
		Status:            generated.PullRequestStatus(pr.Status),
	}
}

func toGeneratedPullRequestShort(pr domain.PullRequestShort) generated.PullRequestShort {
	return generated.PullRequestShort{
		AuthorId:        pr.AuthorID,
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		Status:          generated.PullRequestShortStatus(pr.Status),
	}
}

func toGeneratedPullRequestShortList(list []domain.PullRequestShort) []generated.PullRequestShort {
	result := make([]generated.PullRequestShort, 0, len(list))
	for _, item := range list {
		result = append(result, toGeneratedPullRequestShort(item))
	}
	return result
}

func toGeneratedTeam(team *domain.Team) generated.Team {
	members := make([]generated.TeamMember, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, generated.TeamMember{
			UserId:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return generated.Team{
		TeamName: team.Name,
		Members:  members,
	}
}

func fromGeneratedTeam(team generated.Team) domain.Team {
	members := make([]domain.TeamMember, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, domain.TeamMember{
			UserID:   member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return domain.Team{
		Name:    team.TeamName,
		Members: members,
	}
}

func toGeneratedUser(user *domain.User) generated.User {
	return generated.User{
		UserId:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

// PostTeamsTeamNameDeactivate deactivates all users in a team.
func (s *Server) PostTeamsTeamNameDeactivate(ctx context.Context, request generated.PostTeamsTeamNameDeactivateRequestObject) (
	generated.PostTeamsTeamNameDeactivateResponseObject, error) {
	count, err := s.service.DeactivateTeam(ctx, request.TeamName)
	if err != nil {
		if resp, err := s.handlePostTeamsTeamNameDeactivateError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	return generated.PostTeamsTeamNameDeactivate200JSONResponse{
		DeactivatedCount: &count,
	}, nil
}

// PostPullRequestSafeReassign safely reassigns inactive reviewers on an open PR.
func (s *Server) PostPullRequestSafeReassign(ctx context.Context, request generated.PostPullRequestSafeReassignRequestObject) (
	generated.PostPullRequestSafeReassignResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	pr, err := s.service.SafeReassignPR(ctx, request.Body.PullRequestId)
	if err != nil {
		if resp, err := s.handlePostPullRequestSafeReassignError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	prPayload := toGeneratedPullRequest(pr)
	return generated.PostPullRequestSafeReassign200JSONResponse{
		Pr: &prPayload,
	}, nil
}

// GetStats returns statistics.
func (s *Server) GetStats(ctx context.Context, request generated.GetStatsRequestObject) (
	generated.GetStatsResponseObject, error) {
	stats, err := s.service.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return toGeneratedStats(stats), nil
}

func toGeneratedStats(stats *domain.Stats) generated.GetStats200JSONResponse {
	response := generated.GetStats200JSONResponse{}

	if stats.Users.ByUser != nil {
		byUser := make([]struct {
			AssignedReviewsTotal *int    `json:"assignedReviewsTotal,omitempty"`
			IsActive             *bool   `json:"isActive,omitempty"`
			MergedReviews        *int    `json:"mergedReviews,omitempty"`
			OpenReviews          *int    `json:"openReviews,omitempty"`
			Team                 *string `json:"team,omitempty"`
			UserId               *string `json:"userId,omitempty"`
			UserName             *string `json:"userName,omitempty"`
		}, 0, len(stats.Users.ByUser))

		for _, u := range stats.Users.ByUser {
			entry := struct {
				AssignedReviewsTotal *int    `json:"assignedReviewsTotal,omitempty"`
				IsActive             *bool   `json:"isActive,omitempty"`
				MergedReviews        *int    `json:"mergedReviews,omitempty"`
				OpenReviews          *int    `json:"openReviews,omitempty"`
				Team                 *string `json:"team,omitempty"`
				UserId               *string `json:"userId,omitempty"`
				UserName             *string `json:"userName,omitempty"`
			}{
				AssignedReviewsTotal: &u.AssignedReviewsTotal,
				IsActive:             &u.IsActive,
				MergedReviews:        &u.MergedReviews,
				OpenReviews:          &u.OpenReviews,
				Team:                 &u.Team,
				UserId:               &u.UserID,
				UserName:             &u.UserName,
			}
			byUser = append(byUser, entry)
		}
		response.Users = &struct {
			Active *int `json:"active,omitempty"`
			ByUser *[]struct {
				AssignedReviewsTotal *int    `json:"assignedReviewsTotal,omitempty"`
				IsActive             *bool   `json:"isActive,omitempty"`
				MergedReviews        *int    `json:"mergedReviews,omitempty"`
				OpenReviews          *int    `json:"openReviews,omitempty"`
				Team                 *string `json:"team,omitempty"`
				UserId               *string `json:"userId,omitempty"`
				UserName             *string `json:"userName,omitempty"`
			} `json:"byUser,omitempty"`
			Inactive *int `json:"inactive,omitempty"`
			Total    *int `json:"total,omitempty"`
		}{
			Active:   &stats.Users.Active,
			ByUser:   &byUser,
			Inactive: &stats.Users.Inactive,
			Total:    &stats.Users.Total,
		}
	}

	if stats.Teams.ByTeam != nil {
		byTeam := make([]struct {
			MembersActive   *int    `json:"membersActive,omitempty"`
			MembersTotal    *int    `json:"membersTotal,omitempty"`
			PrsCreatedTotal *int    `json:"prsCreatedTotal,omitempty"`
			PrsOpen         *int    `json:"prsOpen,omitempty"`
			TeamName        *string `json:"teamName,omitempty"`
		}, 0, len(stats.Teams.ByTeam))

		for _, t := range stats.Teams.ByTeam {
			entry := struct {
				MembersActive   *int    `json:"membersActive,omitempty"`
				MembersTotal    *int    `json:"membersTotal,omitempty"`
				PrsCreatedTotal *int    `json:"prsCreatedTotal,omitempty"`
				PrsOpen         *int    `json:"prsOpen,omitempty"`
				TeamName        *string `json:"teamName,omitempty"`
			}{
				MembersActive:   &t.MembersActive,
				MembersTotal:    &t.MembersTotal,
				PrsCreatedTotal: &t.PRsCreatedTotal,
				PrsOpen:         &t.PRsOpen,
				TeamName:        &t.TeamName,
			}
			byTeam = append(byTeam, entry)
		}
		response.Teams = &struct {
			ByTeam *[]struct {
				MembersActive   *int    `json:"membersActive,omitempty"`
				MembersTotal    *int    `json:"membersTotal,omitempty"`
				PrsCreatedTotal *int    `json:"prsCreatedTotal,omitempty"`
				PrsOpen         *int    `json:"prsOpen,omitempty"`
				TeamName        *string `json:"teamName,omitempty"`
			} `json:"byTeam,omitempty"`
			Total *int `json:"total,omitempty"`
		}{
			ByTeam: &byTeam,
			Total:  &stats.Teams.Total,
		}
	}

	response.Prs = &struct {
		Merged         *int `json:"merged,omitempty"`
		Open           *int `json:"open,omitempty"`
		Total          *int `json:"total,omitempty"`
		With0Reviewers *int `json:"with0Reviewers,omitempty"`
		With1Reviewer  *int `json:"with1Reviewer,omitempty"`
		With2Reviewers *int `json:"with2Reviewers,omitempty"`
	}{
		Merged:         &stats.PRs.Merged,
		Open:           &stats.PRs.Open,
		Total:          &stats.PRs.Total,
		With0Reviewers: &stats.PRs.With0Reviewers,
		With1Reviewer:  &stats.PRs.With1Reviewer,
		With2Reviewers: &stats.PRs.With2Reviewers,
	}

	return response
}

func (s *Server) handlePostTeamsTeamNameDeactivateError(err error) (generated.PostTeamsTeamNameDeactivateResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.PostTeamsTeamNameDeactivate404JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}

func (s *Server) handlePostPullRequestSafeReassignError(err error) (generated.PostPullRequestSafeReassignResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)

	switch payload.status {
	case http.StatusNotFound:
		return generated.PostPullRequestSafeReassign404JSONResponse(payload.body), nil
	default:
		return nil, echo.NewHTTPError(payload.status, payload.body.Error.Message)
	}
}
