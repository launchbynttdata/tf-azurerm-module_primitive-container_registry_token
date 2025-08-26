package testimpl

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	if len(subscriptionId) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	clientFactory, err := armcontainerregistry.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client factory: %v", err)
	}

	tokensClient := clientFactory.NewTokensClient()

	rgName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	registryName := terraform.Output(t, ctx.TerratestTerraformOptions(), "container_registry_name")
	tokenName := terraform.Output(t, ctx.TerratestTerraformOptions(), "token_name")

	t.Run("TokenWasCreated", func(t *testing.T) {
		resp, err := tokensClient.Get(context.TODO(), rgName, registryName, tokenName, nil)

		if err != nil {
			log.Fatalf("failed to get token: %v", err)
		}
		assert.NotNilf(t, resp.Token, "Expected Token to be created")
	})
}
