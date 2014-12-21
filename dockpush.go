package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type PluginDockPush struct {
	CliConnection plugin.CliConnection
	help          *bool
	Space         string
	SpaceGUID     string
	Domain        string
	DockerApp
}

type DockerApp struct {
	AppName     string
	Memory      string
	Instances   string
	DiskQuota   string
	DockerImage string
	Command     string
}

//thanks Concourse geezahs
func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(new(PluginDockPush))
}

func (pluginDP *PluginDockPush) Run(cliConnection plugin.CliConnection, args []string) {
	pluginDP.CliConnection = cliConnection

	pluginDP.parseArgs(args)

	pluginDP.GetSelectedSpace()
	pluginDP.GetSelectedSpaceGUID()

	pluginDP.PushContainer()

	pluginDP.SetDiegoEnvVars()

	pluginDP.GetSelectedDomain()
	pluginDP.MapDefaultRoute()

	pluginDP.StartApp()
}

func (pluginDP *PluginDockPush) parseArgs(args []string) {
	dockerFlagSet := flag.NewFlagSet("cloudocker", flag.ExitOnError)

	pluginDP.help = dockerFlagSet.Bool("help", false, "passed to display help text")

	memory := dockerFlagSet.String("m", "1024", "Memory limit (in MB) for container")
	instances := dockerFlagSet.String("i", "1", "Number of instances")
	diskQuota := dockerFlagSet.String("d", "1024", "Disk space limit (in MB) for container")

	err := dockerFlagSet.Parse(args[1:])
	fatalIf(err)

	if *pluginDP.help {
		printHelp()
		os.Exit(0)
	}

	if len(dockerFlagSet.Args()) != 3 {
		printHelp()
		os.Exit(1)
	}

	pluginDP.Memory = *memory
	pluginDP.Instances = *instances
	pluginDP.DiskQuota = *diskQuota
	pluginDP.DockerImage = dockerFlagSet.Args()[0]
	pluginDP.Command = dockerFlagSet.Args()[1]
	pluginDP.AppName = dockerFlagSet.Args()[2]
}

func (pluginDP *PluginDockPush) GetSelectedSpace() {
	output, err := pluginDP.CliConnection.CliCommandWithoutTerminalOutput("target")
	fatalIf(err)

	spaceRegex, _ := regexp.Compile("Space:\\s+(.*)\\s")
	pluginDP.Space = strings.TrimSpace(spaceRegex.FindStringSubmatch(output[4])[1])
}

func (pluginDP *PluginDockPush) GetSelectedSpaceGUID() {
	output, err := pluginDP.CliConnection.CliCommandWithoutTerminalOutput("space", "--guid", pluginDP.Space)
	fatalIf(err)

	pluginDP.SpaceGUID = strings.TrimSpace(output[0])
}

func (pluginDP *PluginDockPush) PushContainer() {
	_, err := pluginDP.CliConnection.CliCommandWithoutTerminalOutput("curl", "/v2/apps", "-X", "POST", "-d", `{ "name": "`+pluginDP.AppName+`", "memory":`+pluginDP.Memory+`, "instances":`+pluginDP.Instances+`, "disk_quota":`+pluginDP.DiskQuota+`, "space_guid":"`+pluginDP.SpaceGUID+`", "docker_image":"`+pluginDP.DockerImage+`", "command":"`+pluginDP.Command+`"}`)
	fatalIf(err)
}

func (pluginDP *PluginDockPush) SetDiegoEnvVars() {
	_, err := pluginDP.CliConnection.CliCommandWithoutTerminalOutput("set-env", pluginDP.AppName, "DIEGO_STAGE_BETA", "true")
	fatalIf(err)
	_, err = pluginDP.CliConnection.CliCommandWithoutTerminalOutput("set-env", pluginDP.AppName, "DIEGO_RUN_BETA", "true")
	fatalIf(err)
}

func (pluginDP *PluginDockPush) GetSelectedDomain() {
	output, err := pluginDP.CliConnection.CliCommandWithoutTerminalOutput("domains")
	fatalIf(err)

	// Use the first shared domain - should be the apps domain from the manifest
	for _, out := range output[2:] {
		if strings.Fields(out)[1] == "shared" {
			pluginDP.Domain = strings.Fields(out)[0]
			break
		}
	}
}

func (pluginDP *PluginDockPush) MapDefaultRoute() {
	_, err := pluginDP.CliConnection.CliCommand("map-route", pluginDP.AppName, pluginDP.Domain, "-n", pluginDP.AppName)
	fatalIf(err)
}

func (pluginDP *PluginDockPush) StartApp() {
	_, err := pluginDP.CliConnection.CliCommand("start", pluginDP.AppName)
	fatalIf(err)
}

func (pluginDP *PluginDockPush) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DockPush",
		Commands: []plugin.Command{
			{
				Name:     "dockpush",
				Alias:    "dp",
				HelpText: "Push a docker container into a CF with Diego. To obtain more information use --help",
			},
		},
	}
}

func printHelp() {
	fmt.Println(`
cf dockpush docker-image run-command app-name
e.g.
cf dockpush cloudfoundry/inigodockertest:latest /dockerapp docker-test

OPTIONAL PARAMS:
-help: used to display this additional output.
-m: memory limit (in MB) for container, default 1024
-i: number of instances, default 1
-d: disk space limit (in MB) for container, default 1024

cf dp -m=512 -i=2 -d=1200 cloudfoundry/inigodockertest:latest /dockerapp docker-test

`)
}
