package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	repoPath     = "/root/projects/golang_training/"    // local dir where repo is expected
	filePath     = "checklist_server/rutina-ninos.html" // subdir + file to check
	remoteRef    = "origin/main"
	buildScript  = "/root/projects/golang_training/checklist_server/build_rpi2.sh"
	interval     = 60 * time.Minute
	gitTimeout   = 2 * time.Minute
	buildTimeout = 10 * time.Minute
)

// true when a pull succeeded but the build has not yet succeeded
var buildPending bool

func git(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", repoPath}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0") // fail instead of waiting for credentials
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %v: %w: %s", args, err, stderr.String())
	}
	return out, nil
}

func runBuild(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, buildTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, buildScript)
	cmd.Dir = repoPath // the script is relative to the repo, not to the service's cwd
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", buildScript, err, out)
	}
	log.Printf("build output:\n%s", out)
	return nil
}

func check(ctx context.Context) error {
	// update remote refs
	if _, err := git(ctx, "fetch", "origin"); err != nil {
		return err
	}

	// md5 of local file
	localData, err := os.ReadFile(filepath.Join(repoPath, filePath))
	if err != nil {
		return err
	}
	localSum := md5.Sum(localData)

	// md5 of remote file
	remoteData, err := git(ctx, "show", remoteRef+":"+filePath)
	if err != nil {
		return err
	}
	remoteSum := md5.Sum(remoteData)

	// pull if md5 changed
	if localSum != remoteSum {
		log.Printf("changed: local %x, remote %x, pulling", localSum, remoteSum)
		out, err := git(ctx, "pull", "--ff-only")
		if err != nil {
			return err
		}
		log.Printf("pull ok:\n%s", out)
		buildPending = true
	}

	if !buildPending {
		log.Printf("no changes (md5 %x)", localSum)
		return nil
	}

	if err := runBuild(ctx); err != nil {
		return err // buildPending stays true, retried next tick
	}
	buildPending = false
	log.Println("build ok")
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("started, checking every %s", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := check(ctx); err != nil {
			log.Printf("check failed: %v", err)
		}
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-ticker.C:
		}
	}
}
