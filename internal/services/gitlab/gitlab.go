package gitlab

import (
	"context"
	"fmt"
	"slices"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"gitlab.levelgroup.ru/devops/approve-checker/internal/core/config"
)

type Client struct {
	git *gitlab.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
	git, err := gitlab.NewClient(cfg.GitLab.Token, gitlab.WithBaseURL(cfg.GitLab.Host))
	if err != nil {
		return nil, fmt.Errorf("an error occurred while init gitlab client, err: %w", err)
	}

	return &Client{
		git: git,
	}, nil
}

func (c *Client) MergeRequestWasApproved(ctx context.Context, cfg *config.Config) (bool, error) {
	approvals, _, err := c.git.MergeRequests.GetMergeRequestApprovals(
		cfg.GitLab.ProjectID,
		cfg.GitLab.MergeRequestIID,
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("an error occurred while get merge request approvals, err: %w", err)
	}

	approvedBy := make(map[*gitlab.MergeRequestApproverUser]bool, cfg.Approve.MinApprovalCount)
	approvedCount := 0

	for i := 0; i < len(approvals.ApprovedBy); i++ {
		approved := approvals.ApprovedBy[i]

		// check personal approves by users
		if _, found := slices.BinarySearch(cfg.Approve.ApprovalAuthors, approved.User.Username); found {
			approvedBy[approved] = true
			approvedCount++
			continue
		}

		// if check personal approves by users was failed, than check by minimal role
		m, _, err := c.git.ProjectMembers.GetProjectMember(
			cfg.GitLab.ProjectID,
			approved.User.ID,
			gitlab.WithContext(ctx),
		)
		if err != nil {
			return false, fmt.Errorf("an error occurred while get project members, err: %w", err)
		}

		if m.AccessLevel >= gitlab.AccessLevelValue(cfg.Approve.MinApprovalRole) {
			approvedBy[approved] = true
			approvedCount++
		}
	}

	if approvedCount < cfg.Approve.MinApprovalCount {
		return false, nil
	}

	return true, nil
}

func AccessLevelString(level int) string {
	switch level {
	case 5:
		return "Minimal"
	case 10:
		return "Guest"
	case 15:
		return "Planner"
	case 20:
		return "Reporter"
	case 30:
		return "Developer"
	case 40:
		return "Maintainer"
	case 50:
		return "Owner"
	case 60:
		return "Admin"
	default:
		return "No permissions"
	}
}
