package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformGcpHelloWorldExample(t *testing.T) {
	t.Parallel()

	workingDir := test_structure.CopyTerraformFolderToTemp(t, "../", "terraform")

	// website::tag::2:: Give the example instance a unique name
	inputClusterName := fmt.Sprintf("test-%s", strings.ToLower(random.UniqueId()))

	// website::tag::6:: Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// website::tag::3:: The path to where our Terraform code is located
		TerraformDir: workingDir,

		// website::tag::4:: Variables to pass to our Terraform code using -var options
		VarFiles: []string{"varfile.tfvars"},
		Vars: map[string]interface{}{
			"gke-name": inputClusterName,
		},
		Lock: true,
	})

	// website::tag::8:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::7:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Get the created cluster's name region and project to set up the access config
	outputClusterName := terraform.Output(t, terraformOptions, "cluster_name")
	outputRegion := terraform.Output(t, terraformOptions, "region")
	outputProject := terraform.Output(t, terraformOptions, "project")

	// define kubeconfig path
	configPath := "../terraform/kubeconfig"

	// prepare command to update kubeconfig context with gcloud beta command
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"beta", "container", "clusters", "get-credentials", outputClusterName, "--region", outputRegion, "--project", outputProject},
		Env: map[string]string{
			"KUBECONFIG": configPath,
		},
	}

	// execute gcloud command
	shell.RunCommand(t, cmd)

	// pattern of kubeconfig context name
	contextName := fmt.Sprintf("gke_%s_%s_%s", outputProject, outputRegion, outputClusterName)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("ns-%s", strings.ToLower(random.UniqueId()))

	// website::tag::2::Setup the kubectl config and context.
	// Here we choose to use the defaults, which is:
	// - Current context of the kubectl config file
	// - HOME/.kube/config for the kubectl config file
	// - Random namespace
	kubectlOptions := k8s.NewKubectlOptions(contextName, configPath, namespaceName)

	// Make sure the nodes are ready and the cluster is in an operational state
	verifyGkeNodesAreReady(t, kubectlOptions)

	// website::tag::1::Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../terraform/hello-world-deployment.yml")
	require.NoError(t, err)

	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	// website::tag::5::Make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	// website::tag::6::At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDelete(t, kubectlOptions, kubeResourcePath)

	// website::tag::3::Apply kubectl with 'kubectl apply -f RESOURCE_CONFIG' command.
	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, kubectlOptions, kubeResourcePath)

	// website::tag::4:: Verify the service is available and get the URL for it.
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "hello-world-service", 10, 5*time.Second)
	service := k8s.GetService(t, kubectlOptions, "hello-world-service")

	url := fmt.Sprintf("http://%s", k8s.GetServiceEndpoint(t, kubectlOptions, service, 5000))

	// Verify the service type
	assert.Equal(t, "LoadBalancer", string(service.Spec.Type))

	// website::tag::5:: Make an HTTP request to the URL and make sure it returns a 200 OK with the body "Hello, World".
	http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello world!", 10, 5*time.Second)

}

// func TestOutput(t *testing.T) {
// 	t.Parallel()

// 	kubectlOptions := k8s.NewKubectlOptions("", "", "default")

// 	service := k8s.GetService(t, kubectlOptions, "hello-world-service")

// 	// log.Printf("##### service [%T]: %+s\n", string(service.Spec.Type), service.Spec.Type)
// 	assert.Equal(t, "LoadBalancer", string(service.Spec.Type))

// }
