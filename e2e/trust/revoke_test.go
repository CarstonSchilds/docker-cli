package trust

import (
	"testing"

	"github.com/docker/cli/e2e/internal/fixtures"
	"github.com/docker/cli/internal/test/environment"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/fs"
	"gotest.tools/v3/icmd"
	"gotest.tools/v3/skip"
)

const (
	revokeImage = "registry:5000/revoke:v1"
	revokeRepo  = "registry:5000/revokerepo"
)

func TestRevokeImage(t *testing.T) {
	skip.If(t, environment.RemoteDaemon())

	dir := fixtures.SetupConfigFile(t)
	defer dir.Remove()
	setupTrustedImagesForRevoke(t, dir)
	result := icmd.RunCmd(
		icmd.Command("docker", "trust", "revoke", revokeImage),
		fixtures.WithPassphrase("root_password", "repo_password"),
		fixtures.WithNotary, fixtures.WithConfig(dir.Path()))
	result.Assert(t, icmd.Success)
	assert.Check(t, is.Contains(result.Stdout(), "Successfully deleted signature for registry:5000/revoke:v1"))
}

func TestRevokeRepo(t *testing.T) {
	skip.If(t, environment.RemoteDaemon())

	dir := fixtures.SetupConfigFile(t)
	defer dir.Remove()
	setupTrustedImagesForRevokeRepo(t, dir)
	result := icmd.RunCmd(
		icmd.Command("docker", "trust", "revoke", revokeRepo, "-y"),
		fixtures.WithPassphrase("root_password", "repo_password"),
		fixtures.WithNotary, fixtures.WithConfig(dir.Path()))
	result.Assert(t, icmd.Success)
	assert.Check(t, is.Contains(result.Stdout(), "Successfully deleted signature for registry:5000/revoke"))
}

func setupTrustedImagesForRevoke(t *testing.T, dir fs.Dir) {
	t.Helper()
	icmd.RunCmd(icmd.Command("docker", "pull", fixtures.AlpineImage)).Assert(t, icmd.Success)
	icmd.RunCommand("docker", "tag", fixtures.AlpineImage, revokeImage).Assert(t, icmd.Success)
	icmd.RunCmd(
		icmd.Command("docker", "-D", "trust", "sign", revokeImage),
		fixtures.WithPassphrase("root_password", "repo_password"),
		fixtures.WithConfig(dir.Path()), fixtures.WithNotary).Assert(t, icmd.Success)
}

func setupTrustedImagesForRevokeRepo(t *testing.T, dir fs.Dir) {
	t.Helper()
	icmd.RunCmd(icmd.Command("docker", "pull", fixtures.AlpineImage)).Assert(t, icmd.Success)
	icmd.RunCommand("docker", "tag", fixtures.AlpineImage, revokeRepo+":v1").Assert(t, icmd.Success)
	icmd.RunCmd(
		icmd.Command("docker", "-D", "trust", "sign", revokeRepo+":v1"),
		fixtures.WithPassphrase("root_password", "repo_password"),
		fixtures.WithConfig(dir.Path()), fixtures.WithNotary).Assert(t, icmd.Success)
	icmd.RunCmd(icmd.Command("docker", "pull", fixtures.BusyboxImage)).Assert(t, icmd.Success)
	icmd.RunCommand("docker", "tag", fixtures.BusyboxImage, revokeRepo+":v2").Assert(t, icmd.Success)
	icmd.RunCmd(
		icmd.Command("docker", "-D", "trust", "sign", revokeRepo+":v2"),
		fixtures.WithPassphrase("root_password", "repo_password"),
		fixtures.WithConfig(dir.Path()), fixtures.WithNotary).Assert(t, icmd.Success)
}
