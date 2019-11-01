package alertlogic

import (
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/deployments"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenDeploymentCredentials(credentials []deployments.Credential) []map[string]string {
	var out = make([]map[string]string, len(credentials), len(credentials))
	for i, v := range credentials {
		m := make(map[string]string)
		m["id"] = *v.Id
		m["purpose"] = *v.Purpose

		out[i] = m
	}
	return out
}

func expandDeploymentCredentials(in *schema.Set) []deployments.Credential {
	var items = make([]deployments.Credential, len(in.List()))
	for i, v := range in.List() {
		m := deployments.Credential{
			Id:      alertlogic.String(v.(map[string]interface{})["id"].(string)),
			Purpose: alertlogic.String(v.(map[string]interface{})["purpose"].(string)),
		}

		items[i] = m
	}
	return items
}

func flattenDeploymentScope(scope []deployments.ScopeObject) []map[string]string {
	var out = make([]map[string]string, len(scope), len(scope))
	for i, v := range scope {
		m := make(map[string]string)
		m["type"] = *v.Type
		m["key"] = *v.Key

		if v.Policy != nil {
			m["policy_id"] = *v.Policy.Id
		}

		out[i] = m
	}
	return out
}

func expandDeploymentScope(in *schema.Set) []deployments.ScopeObject {
	var items = make([]deployments.ScopeObject, len(in.List()))
	for i, v := range in.List() {
		m := deployments.ScopeObject{
			Type: alertlogic.String(v.(map[string]interface{})["type"].(string)),
			Key:  alertlogic.String(v.(map[string]interface{})["key"].(string)),
		}

		if p := v.(map[string]interface{})["policy_id"].(string); p != "" {
			if p != "" {
				m.Policy = &deployments.Policy{
					Id: alertlogic.String(p),
				}
			}
		}

		items[i] = m
	}
	return items
}
