package main_test

import (
	. "github.com/cloudcredo/dockpush"

	"github.com/cloudfoundry/cli/plugin/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cdock", func() {

	var (
		fakeCliConnection *fakes.FakeCliConnection
		pluginDockPush    *PluginDockPush
	)

	BeforeEach(func() {
		fakeCliConnection = new(fakes.FakeCliConnection)
		pluginDockPush = new(PluginDockPush)
		pluginDockPush.CliConnection = fakeCliConnection
	})

	Describe("pushing a Docker container into the currently selected space", func() {
		It("should ask the cli connection for the name of the currently selected space", func() {
			testTargetOutput := []string{"\n",
				" API endpoint:   https://api.10.244.0.34.xip.io (API version: 2.16.0)\n",
				" User:           admin\n",
				" Org:            diego\n",
				" Space:          veryuniquetestspace\n",
			}

			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(testTargetOutput, nil)

			pluginDockPush.GetSelectedSpace()

			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).To(Equal([]string{"target"}))

			Expect(pluginDockPush.Space).To(Equal("veryuniquetestspace"))
		})

		It("should find the GUID for the currently selected space", func() {
			testSpacesOutput := []string{"b6723fdf-fc2c-4753-b2eb-c972a0333519\n"}
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(testSpacesOutput, nil)
			pluginDockPush.Space = "testSpace"

			pluginDockPush.GetSelectedSpaceGUID()
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).To(Equal([]string{"space", "--guid", "testSpace"}))
			Expect(pluginDockPush.SpaceGUID).To(Equal("b6723fdf-fc2c-4753-b2eb-c972a0333519"))
		})

		It("should use CF curl to push a container to the currently selected space", func() {
			pluginDockPush.SpaceGUID = "b6723fdf-fc2c-4753-b2eb-c972a0333518"
			pluginDockPush.DockerApp = DockerApp{
				AppName:     "docker-test",
				Memory:      "1024",
				Instances:   "1",
				DiskQuota:   "1024",
				DockerImage: "cloudfoundry/inigodockertest:latest",
				Command:     "/dockerapp",
			}

			pluginDockPush.PushContainer()
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).To(Equal([]string{"curl", "/v2/apps", "-X", "POST", "-d", `{ "name": "` + pluginDockPush.AppName + `", "memory":` + pluginDockPush.Memory + `, "instances":` + pluginDockPush.Instances + `, "disk_quota":` + pluginDockPush.DiskQuota + `, "space_guid":"b6723fdf-fc2c-4753-b2eb-c972a0333518", "docker_image":"` + pluginDockPush.DockerImage + `", "command":"` + pluginDockPush.Command + `"}`}))
		})

		It("should set an environment variable to use Diego components for staging and runtime", func() {
			pluginDockPush.AppName = "testApp"
			pluginDockPush.SetDiegoEnvVars()
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).To(Equal([]string{"set-env", "testApp", "DIEGO_STAGE_BETA", "true"}))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)).To(Equal([]string{"set-env", "testApp", "DIEGO_RUN_BETA", "true"}))
		})

		It("should find the first shared domain", func() {
			testTargetOutput := []string{
				"Getting domains in org diego as admin...\n",
				" name                 status\n",
				" 10.244.0.34.xip.io   shared\n",
			}

			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(testTargetOutput, nil)

			pluginDockPush.GetSelectedDomain()

			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)).To(Equal([]string{"domains"}))

			Expect(pluginDockPush.Domain).To(Equal("10.244.0.34.xip.io"))
		})

		It("should map the default route to the app", func() {
			pluginDockPush.AppName = "testApp"
			pluginDockPush.Domain = "10.244.0.34.xip.io"
			pluginDockPush.MapDefaultRoute()
			Expect(fakeCliConnection.CliCommandCallCount()).To(Equal(1))
			Expect(fakeCliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"map-route", "testApp", "10.244.0.34.xip.io", "-n", "testApp"}))
		})

		It("should start the app", func() {
			pluginDockPush.AppName = "testApp"
			pluginDockPush.StartApp()
			Expect(fakeCliConnection.CliCommandCallCount()).To(Equal(1))
			Expect(fakeCliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"start", "testApp"}))
		})
	})
})
