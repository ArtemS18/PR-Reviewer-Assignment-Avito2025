package handlers_test

import (
	"errors"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
)

type mockTeamRepo struct{}
type mockUserRepo struct{}
type mockPKRepo struct{}

func (m *mockTeamRepo) GetTeam(teamName string) (ds.Team, error) {
	if teamName == "fail" {
		return ds.Team{}, errors.New("fail")
	}
	return ds.Team{Name: teamName}, nil
}

func (m *mockTeamRepo) AddTeam(teamData dto.TeamDTO) (ds.Team, error) {
	if teamData.Name == "fail" {
		return ds.Team{}, errors.New("fail")
	}
	return ds.Team{Name: teamData.Name}, nil
}

func (m *mockUserRepo) SetUserFlag(userID string, isActive bool) (ds.User, error) {
	if userID == "fail" {
		return ds.User{}, errors.New("fail")
	}
	return ds.User{ID: userID, Name: "Test", IsActive: isActive}, nil
}

func (m *mockUserRepo) GetReview(userID string) (ds.User, error) {
	if userID == "fail" {
		return ds.User{}, errors.New("fail")
	}
	return ds.User{ID: userID, Name: "Test", IsActive: true}, nil
}

func (m *mockPKRepo) CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error) {
	if pkDTO.ID == "fail" {
		return ds.PullRequest{}, errors.New("fail")
	}
	return ds.PullRequest{ID: pkDTO.ID, Name: pkDTO.Name, AuthorID: pkDTO.AuthorID}, nil
}

func (m *mockPKRepo) ReassignReviewer(pkID, oldReviewerID string) (ds.PullRequest, error) {
	if pkID == "fail" {
		return ds.PullRequest{}, errors.New("fail")
	}
	return ds.PullRequest{ID: pkID}, nil
}

func (m *mockPKRepo) Merged(pkID string) (ds.PullRequest, error) {
	if pkID == "fail" {
		return ds.PullRequest{}, errors.New("fail")
	}
	return ds.PullRequest{ID: pkID}, nil
}
