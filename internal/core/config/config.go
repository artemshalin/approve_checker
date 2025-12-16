package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Approve Approve `yaml:"approve"`
	GitLab  GitLabConfig
}

type Approve struct {
	MinApprovalRole  int      `env:"APPROVE_MIN_APPROVAL_ROLE" env-default:"40"`
	ApprovalAuthors  []string `env:"APPROVE_APPROVAL_AUTHORS"`
	MinApprovalCount int      `env:"APPROVE_MIN_APPROVAL_COUNT" env-default:"1"`
}

type GitLabConfig struct {
	Token           string `env:"CI_JOB_TOKEN"`
	Host            string `env:"CI_SERVER_URL"`
	ProjectID       string `env:"CI_PROJECT_ID"`
	MergeRequestIID int64  `env:"CI_MERGE_REQUEST_IID"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{
		Approve: Approve{},
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("parse config was failed, err: %w", err)
	}

	slog.Info("current", "config", cfg.String())

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation was failed, err: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) String() string {
	presentation := make([]string, 0)
	presentation = append(presentation, fmt.Sprintf("APPROVE_MIN_APPROVAL_ROLE : %d", cfg.Approve.MinApprovalRole))
	presentation = append(presentation, fmt.Sprintf("APPROVE_APPROVAL_AUTHORS : %s", strings.Join(cfg.Approve.ApprovalAuthors, ", ")))
	presentation = append(presentation, fmt.Sprintf("APPROVE_MIN_APPROVAL_COUNT : %d", cfg.Approve.MinApprovalCount))
	presentation = append(presentation, fmt.Sprintf("CI_JOB_TOKEN : %s", cfg.GitLab.Token))
	presentation = append(presentation, fmt.Sprintf("CI_SERVER_URL : %s", cfg.GitLab.Host))
	presentation = append(presentation, fmt.Sprintf("CI_PROJECT_ID : %s", cfg.GitLab.ProjectID))
	presentation = append(presentation, fmt.Sprintf("CI_MERGE_REQUEST_IID : %d", cfg.GitLab.MergeRequestIID))

	return strings.Join(presentation, "\n")
}

func (cfg *Config) validate() error {
	if cfg.Approve.MinApprovalCount < 1 {
		return errors.New("the minimum number of approvals should be more than 1")
	}

	if len(cfg.Approve.ApprovalAuthors) == 0 && cfg.Approve.MinApprovalRole < 0 {
		return errors.New("should set approval authors or the minimum role that can approve an MR")
	}

	if cfg.GitLab.Token == "" {
		return errors.New("environment variables CI_JOB_TOKEN is required")
	}

	if cfg.GitLab.Host == "" {
		return errors.New("environment variables CI_SERVER_URL is required")
	}

	if cfg.GitLab.ProjectID == "" {
		return errors.New("environment variables CI_PROJECT_ID is required")
	}

	if cfg.GitLab.MergeRequestIID == 0 {
		return errors.New("environment variables CI_MERGE_REQUEST_IID is required")
	}

	return nil
}
