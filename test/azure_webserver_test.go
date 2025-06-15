package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAzureLinuxVMCreation(t *testing.T) {
	subscriptionId := "fc12c2a4-1b36-460b-9a72-d18c25fffb97"
	labelPrefix := "taran2025"

	terraformOptions := &terraform.Options{
		TerraformDir:    "../",
		TerraformBinary: "C:/Users/taran/cst8918-w25-A05-tarang-0508/test/terraform.exe",
		Vars: map[string]interface{}{
			"labelPrefix": labelPrefix,
		},
		EnvVars: map[string]string{
			"ARM_SUBSCRIPTION_ID": subscriptionId,
		},
	}

	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	resourceGroup := terraform.Output(t, terraformOptions, "resource_group_name")

	assert.NotEmpty(t, vmName)

	nicList, err := azure.GetVirtualMachineNicsE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.NotEmpty(t, nicList)
	assert.Contains(t, nicList, nicName)

	vmImage, err := azure.GetVirtualMachineImageE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.True(t,
		strings.Contains(vmImage.SKU, "18.04") || strings.Contains(vmImage.SKU, "22_04"),
		"VM image SKU should be either 18.04 or 22_04, got: %s", vmImage.SKU)
}
