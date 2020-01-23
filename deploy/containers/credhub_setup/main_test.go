package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func generateFakeMounts(deploymentName string, links map[string]interface{}) (string, error) {
	workDir, err := ioutil.TempDir("", "credhub-setup-test-quarks-link-")
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %w", err)
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(workDir)
		}
	}()

	podInfoDir := filepath.Join(workDir, "run/pod-info")
	err = os.MkdirAll(podInfoDir, 0755)
	if err != nil {
		return "", fmt.Errorf("could not create fake pod info directory: %w", err)
	}
	deploymentNameFile, err := os.Create(filepath.Join(podInfoDir, "deployment-name"))
	if err != nil {
		return "", fmt.Errorf("could not create fake deployment name file: %w", err)
	}
	_, err = deploymentNameFile.WriteString(deploymentName)
	if err != nil {
		return "", fmt.Errorf("could not write fake deployment name file: %w", err)
	}
	err = deploymentNameFile.Close()
	if err != nil {
		return "", fmt.Errorf("could not close deployment name file: %w", err)
	}

	for linkName, linkContents := range links {
		linkDir := filepath.Join(workDir, "quarks/link", deploymentName, linkName)
		err = os.MkdirAll(linkDir, 0755)
		if err != nil {
			return "", fmt.Errorf("could not create fake link directory %s: %w", linkName, err)
		}
		linkFile, err := os.Create(filepath.Join(linkDir, "link.yaml"))
		if err != nil {
			return "", fmt.Errorf("could not create fake link file for %s: %w", linkName, err)
		}
		err = json.NewEncoder(linkFile).Encode(linkContents)
		if err != nil {
			return "", fmt.Errorf("could not write fake link file for %s: %w", linkName, err)
		}
		err = linkFile.Close()
		if err != nil {
			return "", fmt.Errorf("could not close fake link file for %s: %w", linkName, err)
		}
		fmt.Printf("wrote link %s/%s\n", deploymentName, linkName)
	}

	return workDir, nil
}

func TestResolveLink(t *testing.T) {
	t.Parallel()

	const linkName = "linkName"
	const deploymentName = "deploymentName"
	var expected, actual struct {
		Field string `json:"field"`
	}
	expected.Field = "hello"
	workDir, err := generateFakeMounts(deploymentName, map[string]interface{}{
		linkName: expected,
	})
	require.NoError(t, err, "could not set up temporary mount directory")
	defer func() {
		require.NoError(t, os.RemoveAll(workDir), "could not remove temporary directory")
	}()

	ctx := context.WithValue(context.Background(), overrideMountRoot, workDir)
	err = resolveLink(ctx, linkName, &actual)
	require.Equal(t, expected, actual, "unexpected link result")
}
