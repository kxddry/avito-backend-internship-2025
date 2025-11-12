package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kxddry/avito-backend-internship-2025/internal/api/generated"
	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
)

type Server struct {
	service domain.AssignmentService
}

func NewServer(service domain.AssignmentService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) PostPullRequestCreate(ctx context.Context, request generated.PostPullRequestCreateRequestObject) (generated.PostPullRequestCreateResponseObject, error) {
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

func (s *Server) PostPullRequestMerge(ctx context.Context, request generated.PostPullRequestMergeRequestObject) (generated.PostPullRequestMergeResponseObject, error) {
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

func (s *Server) PostPullRequestReassign(ctx context.Context, request generated.PostPullRequestReassignRequestObject) (generated.PostPullRequestReassignResponseObject, error) {
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

func (s *Server) PostTeamAdd(ctx context.Context, request generated.PostTeamAddRequestObject) (generated.PostTeamAddResponseObject, error) {
	if request.Body == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "request body is required")
	}

	team := fromGeneratedTeam(*request.Body)

	created, err := s.service.UpsertTeam(ctx, &team)
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

func (s *Server) GetTeamGet(ctx context.Context, request generated.GetTeamGetRequestObject) (generated.GetTeamGetResponseObject, error) {
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

func (s *Server) GetUsersGetReview(ctx context.Context, request generated.GetUsersGetReviewRequestObject) (generated.GetUsersGetReviewResponseObject, error) {
	result, err := s.service.GetReviewerAssignments(ctx, request.Params.UserId)
	if err != nil {
		if resp, err := s.handleGetUsersGetReviewError(err); err != nil || resp != nil {
			return resp, err
		}
		return nil, err
	}

	return generated.GetUsersGetReview200JSONResponse{
		UserId:       result.UserID,
		PullRequests: toGeneratedPullRequestShortList(result.PullRequests),
	}, nil
}

func (s *Server) PostUsersSetIsActive(ctx context.Context, request generated.PostUsersSetIsActiveRequestObject) (generated.PostUsersSetIsActiveResponseObject, error) {
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

func (s *Server) handlePostPullRequestReassignError(err error) (generated.PostPullRequestReassignResponseObject, error) {
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

func (s *Server) handleGetUsersGetReviewError(err error) (generated.GetUsersGetReviewResponseObject, error) {
	appErr, err := unwrapDomainError(err)
	if err != nil || appErr == nil {
		return nil, err
	}

	payload := makeErrorPayload(appErr)
	return nil, echo.NewHTTPError(payload.status, payload.body)
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
	reviewers := append([]string(nil), pr.AssignedReviewers...)

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
