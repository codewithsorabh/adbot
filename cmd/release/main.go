package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var semverPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$`)

func main() {
	var (
		flagVersion = flag.String("version", "", "Release version (e.g., 1.0.0)")
		flagMessage = flag.String("message", "", "Commit message")
		flagFiles   = flag.String("files", "", "Files to stage for git add (e.g., .)")
		flagAutoPush = flag.Bool("push", true, "Auto push to git remote")
		flagYes     = flag.Bool("yes", false, "Skip confirmation prompts")
	)
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("==============================================")
	fmt.Println("       ADBOT AUTOMATED RELEASE MANAGER        ")
	fmt.Println("==============================================")
	fmt.Println()

	// 1. Read current version from VERSION file if present
	currentVersion := readCurrentVersion()

	// 2. Version input
	var targetVersion string
	if *flagVersion != "" {
		targetVersion = cleanVersion(*flagVersion)
	} else if *flagYes {
		targetVersion = currentVersion
	} else {
		for {
			prompt := fmt.Sprintf("Enter release version [%s]: ", currentVersion)
			fmt.Print(prompt)
			input, err := reader.ReadString('\n')
			if err != nil {
				fatal("Failed to read input", err)
			}
			input = strings.TrimSpace(input)
			if input == "" {
				input = currentVersion
			}
			cleaned := cleanVersion(input)
			if !semverPattern.MatchString(cleaned) {
				fmt.Printf("[ERROR] '%s' is not a valid semantic version (expected X.Y.Z, e.g. 1.0.0)\n", cleaned)
				continue
			}
			targetVersion = cleaned
			break
		}
	}

	tag := fmt.Sprintf("v%s", targetVersion)
	fmt.Printf("[INFO] Release Tag: %s\n\n", tag)

	// 3. Git files to stage input
	var stageFiles string
	if *flagFiles != "" {
		stageFiles = *flagFiles
	} else if *flagYes {
		stageFiles = "."
	} else {
		fmt.Print("Files to stage with 'git add' [.]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fatal("Failed to read input", err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			stageFiles = "."
		} else {
			stageFiles = input
		}
	}

	// 4. Commit message input
	defaultCommitMsg := fmt.Sprintf("release: %s", tag)
	var commitMessage string
	if *flagMessage != "" {
		commitMessage = *flagMessage
	} else if *flagYes {
		commitMessage = defaultCommitMsg
	} else {
		fmt.Printf("Git commit message [%s]: ", defaultCommitMsg)
		input, err := reader.ReadString('\n')
		if err != nil {
			fatal("Failed to read input", err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			commitMessage = defaultCommitMsg
		} else {
			commitMessage = input
		}
	}

	// 5. Auto push confirmation
	autoPush := *flagAutoPush
	if !*flagYes && *flagVersion == "" {
		fmt.Print("Auto push commit and tag to GitHub? (Y/n) [Y]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fatal("Failed to read input", err)
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "n" || input == "no" {
			autoPush = false
		}
	}

	// Determine current git branch
	currentBranch := getCurrentGitBranch()
	fmt.Printf("\n[INFO] Target Branch: %s\n", currentBranch)
	fmt.Printf("[INFO] Release Tag:   %s\n", tag)
	fmt.Printf("[INFO] Commit Msg:    %s\n\n", commitMessage)

	// Step 1: Update VERSION file
	step(fmt.Sprintf("Updating VERSION file to %s...", targetVersion))
	if err := os.WriteFile("VERSION", []byte(targetVersion+"\n"), 0644); err != nil {
		fatal("Failed to write VERSION file", err)
	}

	// Step 2: Git Add
	step(fmt.Sprintf("Staging files: git add %s", stageFiles))
	addArgs := append([]string{"add"}, strings.Fields(stageFiles)...)
	runGitCommand(addArgs...)

	// Step 3: Git Commit
	step(fmt.Sprintf("Committing: git commit -m %q", commitMessage))
	// Check if there are changes to commit
	statusOutput := getGitOutput("status", "--porcelain")
	if strings.TrimSpace(statusOutput) != "" {
		runGitCommand("commit", "-m", commitMessage)
	} else {
		fmt.Println("[INFO] Working tree clean, no new changes to commit.")
	}

	// Step 4: Git Tag
	step(fmt.Sprintf("Creating git tag: %s", tag))
	// Check if tag already exists locally
	tagExists := checkGitTagExists(tag)
	if tagExists {
		fmt.Printf("[WARN] Tag %s already exists locally. Overwriting...\n", tag)
		runGitCommand("tag", "-f", tag)
	} else {
		runGitCommand("tag", tag)
	}

	// Step 5: Git Push
	if autoPush {
		step(fmt.Sprintf("Pushing branch '%s' to origin...", currentBranch))
		runGitCommand("push", "origin", currentBranch)

		step(fmt.Sprintf("Pushing tag '%s' to origin...", tag))
		if tagExists {
			runGitCommand("push", "origin", tag, "--force")
		} else {
			runGitCommand("push", "origin", tag)
		}
	} else {
		fmt.Println("\n[INFO] Auto-push skipped. To push manually, run:")
		fmt.Printf("  git push origin %s\n", currentBranch)
		fmt.Printf("  git push origin %s\n", tag)
	}

	// Done
	fmt.Println()
	fmt.Println("==============================================")
	fmt.Printf("      RELEASE %s COMPLETED SUCCESSFULLY!      \n", tag)
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Println("GitHub Actions will automatically build & publish:")
	fmt.Println("  https://github.com/codewithsorabh/adbot/actions")
	fmt.Println()
}

func readCurrentVersion() string {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "1.0.0"
	}
	v := strings.TrimSpace(string(data))
	if v == "" {
		return "1.0.0"
	}
	return cleanVersion(v)
}

func cleanVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return v
}

func getCurrentGitBranch() string {
	branch := getGitOutput("rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return "main"
	}
	return branch
}

func checkGitTagExists(tag string) bool {
	out := getGitOutput("tag", "-l", tag)
	return strings.TrimSpace(out) == tag
}

func runGitCommand(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal(fmt.Sprintf("git %s failed", strings.Join(args, " ")), err)
	}
}

func getGitOutput(args ...string) string {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func step(msg string) {
	fmt.Printf("==> %s\n", msg)
}

func fatal(msg string, err error) {
	fmt.Printf("[ERROR] %s\n", msg)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
	}
	os.Exit(1)
}
