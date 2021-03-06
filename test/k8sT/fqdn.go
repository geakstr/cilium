// Copyright 2019 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sTest

import (
	"context"
	"fmt"

	. "github.com/cilium/cilium/test/ginkgo-ext"
	"github.com/cilium/cilium/test/helpers"

	. "github.com/onsi/gomega"
)

var _ = Describe("K8sFQDNTest", func() {
	var (
		kubectl          *helpers.Kubectl
		microscopeErr    error
		microscopeCancel                    = func() error { return nil }
		backgroundCancel context.CancelFunc = func() { return }
		backgroundError  error

		bindManifest = helpers.ManifestGet("bind_deployment.yaml")
		demoManifest = helpers.ManifestGet("demo.yaml")

		apps    = []string{helpers.App2, helpers.App3}
		appPods map[string]string

		worldTarget          = "http://world1.cilium.test"
		worldTargetIP        = "192.168.9.10"
		worldInvalidTarget   = "http://world2.cilium.test"
		worldInvalidTargetIP = "192.168.9.11"
	)

	BeforeAll(func() {
		kubectl = helpers.CreateKubectl(helpers.K8s1VMName(), logger)
		ProvisionInfraPods(kubectl)

		bindManifest = helpers.ManifestGet("bind_deployment.yaml")

		res := kubectl.Apply(bindManifest)
		res.ExpectSuccess("Bind config cannot be deployed")

		res = kubectl.Apply(demoManifest)
		res.ExpectSuccess("Demo config cannot be deployed")

		err := kubectl.WaitforPods(helpers.DefaultNamespace, "-l zgroup=testapp", helpers.HelperTimeout)
		Expect(err).Should(BeNil(), "Testapp is not ready after timeout")

		appPods = helpers.GetAppPods(apps, helpers.DefaultNamespace, kubectl, "id")

		err = kubectl.WaitforPods(helpers.DefaultNamespace, "-l zgroup=bind", helpers.HelperTimeout)
		Expect(err).Should(BeNil(), "Bind app is not ready after timeout")

	})

	AfterFailed(func() {
		kubectl.CiliumReport(helpers.KubeSystemNamespace,
			"cilium service list",
			"cilium endpoint list")
	})

	AfterAll(func() {
		_ = kubectl.Delete(bindManifest)
		_ = kubectl.Delete(demoManifest)
		ExpectAllPodsTerminated(kubectl)
	})

	JustBeforeEach(func() {
		microscopeErr, microscopeCancel = kubectl.MicroscopeStart()
		Expect(microscopeErr).To(BeNil(), "Microscope cannot be started")
		backgroundCancel, backgroundError = kubectl.BackgroundReport("uptime")
		Expect(backgroundError).To(BeNil(), "Cannot start background report process")
	})

	JustAfterEach(func() {
		kubectl.ValidateNoErrorsInLogs(CurrentGinkgoTestDescription().Duration)
		Expect(microscopeCancel()).To(BeNil(), "cannot stop microscope")
		backgroundCancel()
	})

	AfterEach(func() {
		_ = kubectl.Exec(fmt.Sprintf("%s delete --all cnp", helpers.KubectlCmd))
	})

	It("Restart Cilium validate that FQDN is still working", func() {
		// Test functionality:
		// - When Cilium is running) Connectivity from App2 application can
		// connect to DNS because dns-proxy filter the DNS request. If the
		// connection is made correctly the IP is whitelisted by the FQDN rule
		// until the DNS TTL expires.
		// When Cilium is not running) The DNS-proxy is not working, so the IP
		// connectivity to an existing IP that was queried before will work,
		// meanwhile connections using new DNS request will fail.
		// On restart) Cilium will restore the IPS that were white-listted in
		// the FQDN and connection will work as normal.

		connectivityTest := func() {

			By("Testing that connection from %q to %q should work",
				appPods[helpers.App2], worldTarget)
			res := kubectl.ExecPodCmd(
				helpers.DefaultNamespace, appPods[helpers.App2],
				helpers.CurlFail(worldTarget))
			ExpectWithOffset(1, res).To(helpers.CMDSuccess(), "%q cannot curl to %q",
				appPods[helpers.App2], worldTarget)

			By("Testing that connection from %q to %q shouldn't work",
				appPods[helpers.App2], worldTarget)
			res = kubectl.ExecPodCmd(
				helpers.DefaultNamespace, appPods[helpers.App2],
				helpers.CurlFail(worldInvalidTarget))
			ExpectWithOffset(1, res).ShouldNot(helpers.CMDSuccess(),
				"%q can curl to %q when it should fail", appPods[helpers.App2], worldInvalidTarget)

			By("Testing that connection from %q to %q works",
				appPods[helpers.App2], worldInvalidTarget)
			res = kubectl.ExecPodCmd(
				helpers.DefaultNamespace, appPods[helpers.App2],
				helpers.CurlFail(worldTargetIP))
			res.ExpectSuccess("%q cannot curl to %q during restart", helpers.App2, worldTargetIP)

			By("Testing that connection from %q to %q should not work",
				appPods[helpers.App2], worldInvalidTargetIP)
			res = kubectl.ExecPodCmd(
				helpers.DefaultNamespace, appPods[helpers.App2],
				helpers.CurlFail(worldInvalidTargetIP))
			res.ExpectFail("%q can  connect when it should not work", helpers.App2)
		}

		fqndProxyPolicy := helpers.ManifestGet("fqdn-proxy-policy.yaml")

		_, err := kubectl.CiliumPolicyAction(
			helpers.KubeSystemNamespace, fqndProxyPolicy,
			helpers.KubectlApply, helpers.HelperTimeout)
		Expect(err).To(BeNil(), "Cannot install fqdn proxy policy")

		connectivityTest()
		By("Deleting cilium pods")

		res := kubectl.Exec(fmt.Sprintf("%s -n %s delete pods -l k8s-app=cilium",
			helpers.KubectlCmd, helpers.KubeSystemNamespace))
		res.ExpectSuccess()

		By("Testing connectivity when cilium is restoring using IPS without DNS")
		res = kubectl.ExecPodCmd(
			helpers.DefaultNamespace, appPods[helpers.App2],
			helpers.CurlFail(worldTargetIP))
		res.ExpectSuccess("%q cannot curl to %q during restart", helpers.App2, worldTargetIP)

		res = kubectl.ExecPodCmd(
			helpers.DefaultNamespace, appPods[helpers.App2],
			helpers.CurlFail(worldInvalidTargetIP))
		res.ExpectFail("%q can  connect when it should not work", helpers.App2)

		ExpectAllPodsTerminated(kubectl)
		ExpectCiliumReady(kubectl)

		By("Testing connectivity when cilium is *restored* using IPS without DNS")
		res = kubectl.ExecPodCmd(
			helpers.DefaultNamespace, appPods[helpers.App2],
			helpers.CurlFail(worldTargetIP))
		res.ExpectSuccess("%q cannot curl to %q during restart", helpers.App2, worldTargetIP)

		res = kubectl.ExecPodCmd(
			helpers.DefaultNamespace, appPods[helpers.App2],
			helpers.CurlFail(worldInvalidTargetIP))
		res.ExpectFail("%q can  connect when it should not work", helpers.App2)

		err = kubectl.CiliumEndpointWaitReady()
		Expect(err).To(BeNil(), "Endpoints are not ready after Cilium restarts")
		By("Testing connectivity using DNS request when cilium is restored correctly")
		connectivityTest()
	})
})
