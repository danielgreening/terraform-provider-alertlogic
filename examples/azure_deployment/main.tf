resource "random_password" "client_secret" {
  length = 32
  special = true
  override_special = "/@\" "
}

provider "azuread" {
  version = "0.5.1"
}

provider "azurerm" {
  version = "1.28"
}

data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "current" {}

resource "azuread_application" "alertlogic_app" {
  name = "alertlogic_siemless_app"
  available_to_other_tenants = false
}

resource azurerm_role_definition "alertlogic_custom_role" {
  name = "Alert Logic Resource Explorer"
  description = "Grants minimal set of 'read' permissions to enable discovery of resources by Alert Logic"
  scope = data.azurerm_subscription.primary.id
  permissions {
    actions = [
      "Microsoft.Authorization/permissions/read",
      "Microsoft.Compute/virtualMachines/read",
      "Microsoft.Compute/virtualMachineScaleSets/read",
      "Microsoft.Compute/virtualMachineScaleSets/virtualMachines/*/read",
      "Microsoft.Network/networkInterfaces/read",
      "Microsoft.Network/publicIPAddresses/read",
      "Microsoft.Network/virtualNetworks/read",
      "Microsoft.Network/virtualNetworks/subnets/read",
      "Microsoft.Network/virtualNetworks/subnets/virtualMachines/read",
      "Microsoft.Network/virtualNetworks/virtualMachines/read",
      "Microsoft.Resources/subscriptions/locations/read",
      "Microsoft.Resources/subscriptions/resourceGroups/read"
    ]
    not_actions = []
  }
  assignable_scopes = [data.azurerm_subscription.primary.id]
}


resource "azurerm_role_assignment" "alertlogic_role_assignment" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.alertlogic_custom_role.id
  principal_id       = azuread_service_principal.alertlogic_service_principal.object_id
}

resource "azuread_service_principal" "alertlogic_service_principal" {
  application_id = azuread_application.alertlogic_app.application_id
}

resource "azuread_application_password" "alertlogic_app" {
  application_object_id = azuread_application.alertlogic_app.id
  end_date_relative = "72000h00m" // 3000 days, aka a long time
  value = random_password.client_secret.result
}

provider "alertlogic" {
  endpoint = "https://api.product.dev.alertlogic.com"
}

resource "alertlogic_credential" "azure_credential" {
  account_id           = "134235891"
  name                 = "dgreening-terraform-credential"
  secret_type          = "azure_ad_client"
  secret_ad_id         = data.azurerm_client_config.current.tenant_id
  secret_client_id     = azuread_application.alertlogic_app.application_id
  secret_client_secret = azuread_application_password.alertlogic_app.value
}

resource "alertlogic_deployment" "azure_deployment" {
  account_id    = "134235891"
  name          = "dgreening-azure-terraform-1"
  platform_type = "azure"
  platform_id   = data.azurerm_subscription.primary.subscription_id
  mode          = "manual"
  credential {
      id = alertlogic_credential.azure_credential.id
      purpose = "discover"
  }
  depends_on = ["azurerm_role_assignment.alertlogic_role_assignment"]
}

