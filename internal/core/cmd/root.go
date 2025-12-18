package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"gitlab.levelgroup.ru/devops/approve-checker/internal/core/config"
	"gitlab.levelgroup.ru/devops/approve-checker/internal/services/gitlab"
)

var rootCmd = &cobra.Command{
	Use:   "approve_checker",
	Short: "Check mr approve in GitLab CI",
	Long:  `CLI app, that check count of approve votes in GitLab merge-request.`,
	Run: func(cmd *cobra.Command, _ []string) {
		cfg, err := config.GetConfig()
		if err != nil {
			slog.Error("get config failed", "err", err)
			os.Exit(1)
		}

		color.Magenta("🚧 Config:")
		color.Magenta("\t⚙️ APPROVE_MIN_APPROVAL_ROLE : %d", cfg.Approve.MinApprovalRole)
		color.Magenta("\t⚙️ APPROVE_APPROVAL_AUTHORS: %s", strings.Join(cfg.Approve.ApprovalAuthors, ", "))
		color.Magenta("\t⚙️ APPROVE_MIN_APPROVAL_COUNT: %d", cfg.Approve.MinApprovalCount)
		if cfg.GitLab.Token == "" {
			color.Magenta("\t⚙️ GITLAB_TOKEN: [IS EMPTY]")
		} else {
			color.Magenta("\t⚙️ GITLAB_TOKEN: [CONFIGURED]")
		}

		color.Magenta("\t⚙️ CI_SERVER_URL: %s", cfg.GitLab.Host)
		color.Magenta("\t⚙️ CI_PROJECT_ID: %s", cfg.GitLab.ProjectID)
		color.Magenta("\t⚙️ CI_MERGE_REQUEST_IID: %d", cfg.GitLab.MergeRequestIID)

		c, err := gitlab.NewClient(cfg)
		if err != nil {
			slog.Error("make gitlab client was failed", "err", err)
			os.Exit(1)
		}

		approved, err := c.MergeRequestWasApproved(cmd.Context(), cfg)
		if err != nil {
			slog.Error("check merge request approve was failed", "err", err)
			os.Exit(1)
		}

		if !approved {
			t := fmt.Sprintf(`Merge request was not approved!

Please receive minimum %d approves.
From the next user(-s): (%s). Or from any project members with role greater or equal then "%s".`,
				cfg.Approve.MinApprovalCount,
				strings.Join(cfg.Approve.ApprovalAuthors, ","),
				gitlab.AccessLevelString(cfg.Approve.MinApprovalRole))

			color.Red("❌ %s", t)
			os.Exit(1)
		}

		color.Green("✅ Merge request was approved! Great job!")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
