# Honest Bank DevOps Technical Assessment

This is a hands-on "fill in the blank" style coding exercise. The "blank" in this
case is the `main.tf` file located in the [/terraform folder](./terraform).

The objective of this assessment is to demonstrate your ability to write and test
infrastructure code.

The specifications of the infrastructure (contents of `main.tf`) can be inferred
from the test code located in [/test/terraform_test.go](./test/terraform_test.go).
The tests are written using [Gruntwork's Terratest library](https://terratest.gruntwork.io),
which requires you to have [Go](https://golang.org/dl/) installed.

## Instructions

1. Examine the [/test/terraform_test.go](./test/terraform_test.go) file.
1. Create a `main.tf` file in [/terraform](./terraform) and populate it with code
   to:
   1. Build a [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine)
      cluster.
   1. Create a [Kubernetes](https://kubernetes.io) [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
      and [Service](https://kubernetes.io/docs/concepts/services-networking/service/)
      following the specs in [/terraform/hello-world-deployment.yaml](./terraform/hello-world-deployment.yml).
1. Ensure all tests pass (see instructions below for how to run tests).
1. Implement the missing feature outlined below.

## Running Tests

The tests are written using [Gruntwork's Terratest library](https://terratest.gruntwork.io),
which requites you to have [Go](https://golang.org/dl/) installed.

To run the tests, execute the following commands from the project root:

```shell 
$ cd test
$ go test -v -timeout 30m
```

## Missing Feature

### Overview

The `hello-world` Kubernetes Deployment has a hardcoded replica count of 1. This
is not recommended practice. For stateless workloads, a [HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
is the recommended way to address scaling.

### Specs

* Implement a [HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) for the `hello-world` Deployment using the
  [Kubernetes provider for Terraform](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs).
* The HPA should scale up from a minimum of 1 replica to a maximum of 10 replicas.
* Scale-up should be triggered by 50% CPU utilization.
