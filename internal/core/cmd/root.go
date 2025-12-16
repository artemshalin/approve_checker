package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"gitlab.levelgroup.ru/devops/approve-checker/internal/core/config"
	"gitlab.levelgroup.ru/devops/approve-checker/internal/services/gitlab"
)

var configPath = ""

var rootCmd = &cobra.Command{
	Use:   "approve_checker",
	Short: "Check mr approve in GitLab CI",
	Long:  `CLI app, that check count of approve votes in GitLab merge-request.`,
	Run: func(cmd *cobra.Command, _ []string) {
		cfg, err := config.GetConfig(configPath)
		if err != nil {
			slog.Error("get config failed", "path", configPath, "err", err)
			os.Exit(1)
		}

		c, err := gitlab.NewClient(cfg)
		if err != nil {
			slog.Error("make gitlab client was failed", "path", configPath, "err", err)
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

			slog.Error(t)
			os.Exit(1)
		}

		slog.Info("Merge request was approved! Great job!")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./approve-checker.yml", "path to yaml config")
}
