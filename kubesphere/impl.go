package kubesphere

import (
	"fmt"
	"time"
)

const (
	defaultPassword = "DefPwd123@"
)

// CreateUserAccount creates a new user account with the provided details.
func CreateUserAccount(userName string) error {
	email := userName + "@titan.com"

	body := map[string]interface{}{
		"apiVersion": "iam.kubesphere.io/v1beta1",
		"kind":       "User",
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"iam.kubesphere.io/uninitialized": "true", "iam.kubesphere.io/globalrole": "platform-regular", "kubesphere.io/creator": "admin",
			}, "name": userName,
		},
		"spec": map[string]interface{}{"email": email, "password": defaultPassword},
	}

	_, err := doRequest("POST", "/kapis/iam.kubesphere.io/v1beta1/users", body)
	if err != nil {
		log.Errorf("CreateUserAccount err:%s", err.Error())
		return err
	}

	// log.Infoln("CreateUserAccount rsp-----")
	// log.Infoln(string(rsp))

	return nil
}

// CreateSpaceAndResourceQuotas creates a space and resource quotas for a user.
func CreateSpaceAndResourceQuotas(order, userName string, cpu, ram, storage int) error {
	err := createUserSpace(order, userName)
	if err != nil {
		log.Errorf("CreateUserSpace: %s", err.Error())
		return err
	}

	time.Sleep(1 * time.Second)
	err = changeWorkspaceMembers(order, userName)
	if err != nil {
		log.Errorf("changeWorkspaceMembers: %s", err.Error())
		return err
	}

	err = createUserResourceQuotas(order, cpu, ram, storage)
	if err != nil {
		log.Errorf("CreateUserResourceQuotas: %s", err.Error())
	}

	return err
}

// createUserSpace creates a user space for the given order and user.
func createUserSpace(order, userName string) error {
	body := map[string]interface{}{
		"apiVersion": "iam.kubesphere.io/v1beta1",
		"kind":       "WorkspaceTemplate",
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"kubesphere.io/creator": "admin",
			}, "name": order,
		},
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"manager": userName,
				},
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"kubesphere.io/creator": "admin",
					},
				},
			},
			"placement": map[string]interface{}{
				"clusters": []map[string]interface{}{
					{
						"name": cluster,
					},
				},
			},
		},
	}

	_, err := doRequest("POST", "/kapis/tenant.kubesphere.io/v1beta1/workspacetemplates", body)
	if err != nil {
		log.Errorf("CreateUserSpace err:%s", err.Error())
		return err
	}

	// log.Infoln("CreateUserSpace rsp-----")
	// log.Infoln(string(rsp))

	return nil
}

func changeWorkspaceMembers(order, userName string) error {
	body := map[string]interface{}{
		"roleRef":  fmt.Sprintf("%s-self-provisioner", order),
		"username": userName,
	}

	path := fmt.Sprintf("/kapis/iam.kubesphere.io/v1beta1/workspaces/%s/workspacemembers/%s", order, userName)
	_, err := doRequest("PUT", path, body)
	if err != nil {
		log.Errorf("changeWorkspaceMembers err:%s", err.Error())
		return err
	}

	// log.Infoln("changeWorkspaceMembers rsp-----")
	// log.Infoln(string(rsp))

	return nil
}

// createUserResourceQuotas creates resource quotas for a user.
// It takes an order string and resource limits for CPU, ram, and storage.
func createUserResourceQuotas(order string, cpu, ram, storage int) error {
	body := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"kubesphere.io/workspace": order,
			}, "name": order,
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"kubesphere.io/workspace": order,
			},
			"quota": map[string]interface{}{
				"hard": map[string]interface{}{
					"limits.cpu":             fmt.Sprintf("%d", cpu),
					"limits.memory":          fmt.Sprintf("%dGi", ram),
					"requests.cpu":           fmt.Sprintf("%d", cpu),
					"requests.memory":        fmt.Sprintf("%dGi", ram),
					"requests.storage":       fmt.Sprintf("%dGi", storage),
					"persistentvolumeclaims": fmt.Sprintf("%d", storage),
				},
			},
		},
	}

	path := fmt.Sprintf("/clusters/%s/kapis/tenant.kubesphere.io/v1beta1/workspaces/%s/resourcequotas", cluster, order)
	_, err := doRequest("POST", path, body)
	if err != nil {
		log.Errorf("CreateUserResourceQuotas err:%s", err.Error())
		return err
	}

	// log.Infoln("CreateUserResourceQuotas rsp-----")
	// log.Infoln(string(rsp))

	return nil
}
