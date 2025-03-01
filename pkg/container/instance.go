package container

import (
	"fmt"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/provisioner"
	"github.com/z0mbix/rolecule/pkg/verifier"
)

type InstanceConfig struct {
	Name  string   `mapstructure:"name"`
	Image string   `mapstructure:"image"`
	Arch  string   `mapstructure:"arch"`
	Args  []string `mapstructure:"args"`
}

type Instances []Instance

type Instance struct {
	Name    string
	Image   string
	Args    []string
	Arch    string
	WorkDir string
	Engine
	Provisioner provisioner.Provisioner
	Verifier    verifier.Verifier
}

func (i *Instance) Create() (string, error) {
	workDir := "/src"
	instanceArgs := []string{
		"run",
		"--privileged",
		"--rm",
		"--detach",
		"--tmpfs", "/tmp",
		"--tmpfs", "/run",
		"--tmpfs", "/run/lock",
		"--tmpfs", "/var/lib/docker",
		"--cgroupns", "host",
		"--volume", "/sys/fs/cgroup:/sys/fs/cgroup:rw",
		"--volume", fmt.Sprintf("%s:%s", i.WorkDir, workDir),
		"--workdir", workDir,
	}

	if i.Arch != "" {
		instanceArgs = append(instanceArgs, "--platform", fmt.Sprintf("linux/%s", i.Arch))
	}

	instanceArgs = append(instanceArgs, "--name", i.Name)

	args := append(instanceArgs, i.Args...)

	log.Debugf("%+v", args)
	output, err := i.Run(i.Image, args)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (i *Instance) Converge() (string, error) {
	env, cmd, args := i.Provisioner.GetCommand()
	return i.Exec(i.Name, env, cmd, args)
}

func (i *Instance) Verify() (string, error) {
	env, cmd, args := i.Verifier.GetCommand()
	return i.Exec(i.Name, env, cmd, args)
}

func (i *Instance) Shell() error {
	return i.Engine.Shell(i.Name)
}

func (i *Instance) Exists() bool {
	return i.Engine.Exists(i.Name)
}

func (i *Instance) Destroy() error {
	return i.Remove(i.Name)
}
