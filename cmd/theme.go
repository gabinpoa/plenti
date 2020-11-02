package cmd

import (
	"errors"
	"fmt"
	"os"
	"plenti/readers"
	"plenti/writers"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

// CommitFlag targets a specific commit hash when running the "git clone" operation.
var CommitFlag string

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme [url]",
	Short: "Downloads parent theme to inherit content, layouts, and assets from",
	Long: `Themes allow you to leverage an existing Plenti site as a starting point for your own site.

To use https://plenti.co as a theme for example, run: plenti new theme git@github.com:plentico/plenti.co
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a url argument")
		}
		if len(args) > 1 {
			return errors.New("urls cannot have spaces")
		}
		if len(args) == 1 {
			return nil
		}
		return fmt.Errorf("invalid url specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Get the repo URL passed via the CLI.
		url := args[0]

		// Get the last part of the URL to isolate the repository name.
		parts := strings.Split(url, "/")
		repoName := parts[len(parts)-1]

		themeDir := "themes/" + repoName

		// Run the "git clone" operation.
		repo, err := git.PlainClone(themeDir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Printf("Can't clone theme repository: %v\n", err)
		}

		// Get the latest commit hash from the repo.
		ref, _ := repo.Head()
		commitObj, _ := repo.CommitObject(ref.Hash())
		commitHash := commitObj.Hash.String()

		// Check if a --commit flag was used.
		if CommitFlag != "" {
			worktree, worktreeErr := repo.Worktree()
			if worktreeErr != nil {
				fmt.Printf("Can't get worktree: %v\n", worktreeErr)
			}
			// Resolve commit in case short hash is used instead of full hash.
			resolvedCommitHash, resolveErr := repo.ResolveRevision(plumbing.Revision(CommitFlag))
			if resolveErr != nil {
				fmt.Printf("Can't resolve commit hash: %v\n", resolveErr)
			}
			// Git checkout the commit hash that was sent via the flag.
			checkoutErr := worktree.Checkout(&git.CheckoutOptions{
				Hash: *resolvedCommitHash,
			})
			if checkoutErr != nil {
				fmt.Printf("Can't get commit: %v\n", checkoutErr)
			} else {
				// The --commit flag could be checkout out, so the hash is valid.
				commitHash = CommitFlag
			}
		}

		// Get the current site configuration file values.
		siteConfig, configPath := readers.GetSiteConfig(".")
		// Update the sitConfig struct with new values.
		themeOptions := new(readers.ThemeOptions)
		themeOptions.URL = url
		themeOptions.Commit = commitHash
		siteConfig.ThemeConfig = make(map[string]readers.ThemeOptions)
		siteConfig.ThemeConfig[repoName] = *themeOptions

		// Update the config file on the filesystem.
		writers.SetSiteConfig(siteConfig, configPath)

	},
}

func init() {
	newCmd.AddCommand(themeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// typeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// typeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	themeCmd.Flags().StringVarP(&CommitFlag, "commit", "c", "", "pull a specific commit hash for the theme")
}
