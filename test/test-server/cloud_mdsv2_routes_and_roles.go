package test_server

var v2RbacRoles = map[string]string{
	"CCloudRoleBindingAdmin": `{
			"name": "CCloudRoleBindingAdmin",
			"policy": {
				"bindingScope": "root",
				"bindWithResource": false,
				"allowedOperations": [
				{"resourceType":"SecurityMetadata","operations":["Describe","Alter"]},
				{"resourceType":"Organization","operations":["AlterAccess","DescribeAccess"]}]}}`,
	"CloudClusterAdmin": `{
			"name": "CloudClusterAdmin",
			"policies": [
			{
				"bindingScope": "cluster",
				"bindWithResource": false,
				"allowedOperations": [
				{"resourceType": "Topic","operations": ["All"]},
				{"resourceType": "KsqlCluster","operations": ["All"]},
				{"resourceType": "Subject","operations": ["All"]},
				{"resourceType": "Connector","operations": ["All"]},
				{"resourceType": "NetworkAccess","operations": ["All"]},
				{"resourceType": "ClusterMetric","operations": ["All"]},
				{"resourceType": "Cluster","operations": ["All"]},
				{"resourceType": "ClusterApiKey","operations": ["All"]},
				{"resourceType": "SecurityMetadata","operations": ["Describe", "Alter"]}]
			},
			{
				"bindingScope": "organization",
				"bindWithResource": false,
				"allowedOperations": [
				{"resourceType": "SupportPlan","operations": ["Describe"]},
				{"resourceType": "User","operations": ["Describe","Invite"]},
				{"resourceType": "ServiceAccount","operations": ["Describe"]}]
			}]}`,
	"EnvironmentAdmin": `{
			"name": "EnvironmentAdmin",
			"policies": [
			{
				"bindingScope": "ENVIRONMENT",
				"bindWithResource": false,
				"allowedOperations": [
				{"resourceType": "SecurityMetadata","operations": ["Describe", "Alter"]},
				{"resourceType": "ClusterApiKey","operations": ["All"]},
				{"resourceType": "Connector","operations": ["All"]},
				{"resourceType": "NetworkAccess","operations": ["All"]},
				{"resourceType": "KsqlCluster","operations": ["All"]},
				{"resourceType": "Environment","operations": ["Alter","Delete","AlterAccess","CreateKafkaCluster","DescribeAccess"]},
				{"resourceType": "Subject","operations": ["All"]},
				{"resourceType": "NetworkConfig","operations": ["All"]},
				{"resourceType": "ClusterMetric","operations": ["All"]},
				{"resourceType": "Cluster","operations": ["All"]},
				{"resourceType": "SchemaRegistry","operations": ["All"]},
				{"resourceType": "NetworkRegion","operations": ["All"]},
				{"resourceType": "Deployment","operations": ["All"]},
				{"resourceType": "Topic","operations": ["All"]}
				]
			},
			{
				"bindingScope": "organization",
				"bindWithResource": false,
				"allowedOperations": [
				{"resourceType": "User","operations": ["Describe","Invite"]},
				{"resourceType": "ServiceAccount","operations": ["Describe"]},
				{"resourceType": "SupportPlan","operations": ["Describe"]}
				]
			}]}`,
	"OrganizationAdmin": `{
			"name": "OrganizationAdmin",
			"policy": {
				"bindingScope": "organization",
				"bindWithResource": false,
				"allowedOperations": [
				{"resourceType": "Topic","operations": ["All"]},
				{"resourceType": "NetworkConfig","operations": ["All"]},
				{"resourceType": "SecurityMetadata","operations": ["Describe", "Alter"]},
				{"resourceType": "Billing","operations": ["All"]},
				{"resourceType": "ClusterApiKey","operations": ["All"]},
				{"resourceType": "Deployment","operations": ["All"]},
				{"resourceType": "SchemaRegistry","operations": ["All"]},
				{"resourceType": "KsqlCluster","operations": ["All"]},
				{"resourceType": "CloudApiKey","operations": ["All"]},
				{"resourceType": "NetworkAccess","operations": ["All"]},
				{"resourceType": "SecuritySSO","operations": ["All"]},
				{"resourceType": "SupportPlan","operations": ["All"]},
				{"resourceType": "Connector","operations": ["All"]},
				{"resourceType": "ClusterMetric","operations": ["All"]},
				{"resourceType": "ServiceAccount","operations": ["All"]},
				{"resourceType": "Subject","operations": ["All"]},
				{"resourceType": "Cluster","operations": ["All"]},
				{"resourceType": "Environment","operations": ["All"]},
				{"resourceType": "NetworkRegion","operations": ["All"]},
				{"resourceType": "Organization","operations": ["Alter","CreateEnvironment","AlterAccess","DescribeAccess"]},
				{"resourceType": "User","operations": ["All"]}
				]
			}
		}`,
}

var v2RoutesAndReplies = map[string]string{
	"/api/metadata/security/v2alpha1/principals/User:u-11aaa/roles/CloudClusterAdmin": `[]`,
	"/api/metadata/security/v2alpha1/roleNames": `[
			"CCloudRoleBindingAdmin",
			"CloudClusterAdmin",
			"EnvironmentAdmin",
			"OrganizationAdmin"
		]`,
	"/api/metadata/security/v2alpha1/lookup/rolebindings/principal/User:u-11aaa": `[
			{
				"scope": {
				  	"path": [
						"organization=1111aaaa-11aa-11aa-11aa-111111aaaaaac"
					],
					"clusters": {
					}
				},
				"rolebindings": {
					"User:u-11aaa": {
						"OrganizationAdmin": []
					}
				}
		  	},
		  	{
				"scope": {
				  	"path": [
						"organization=1234",
						"environment=a-595",
						"cloud-cluster=lkc-1111aaa"
					],
					"clusters": {
					}
				},
				"rolebindings": {
					"User:u-11aaa": {
						"CloudClusterAdmin": []
					}
				}
		  	}
		]`,
	"/api/metadata/security/v2alpha1/lookup/rolebindings/principal/User:u-22bbb": `[
		  	{
				"scope": {
				  	"path": [
						"organization=1234",
						"environment=a-595"
					],
					"clusters": {
					}
				},
				"rolebindings": {
					"User:u-22bbb": {
						"EnvironmentAdmin": []
					}
				}
		  	}
		]`,
	"/api/metadata/security/v2alpha1/lookup/rolebindings/principal/User:u-33ccc": `[
		  	{
				"scope": {
				  	"path": [
						"organization=1234",
						"environment=a-595",
						"cloud-cluster=lkc-1111aaa"
					],
					"clusters": {
					}
				},
				"rolebindings": {
					"User:u-33ccc": {
						"CloudClusterAdmin": []
					}
				}
		  	}
		]`,
	"/api/metadata/security/v2alpha1/lookup/rolebindings/principal/User:u-44ddd": `[
		  	{
				"scope": {
				  	"path": [
						"organization=1234",
						"environment=a-595",
						"cloud-cluster=lkc-1111aaa"
					],
					"clusters": {
					}
				},
				"rolebindings": {
					"User:u-44ddd": {
						"CloudClusterAdmin": []
					}
				}
		  	}
		]`,
	"/api/metadata/security/v2alpha1/lookup/role/OrganizationAdmin": `[
			"User:u-11aaa"
		]`,
	"/api/metadata/security/v2alpha1/lookup/role/EnvironmentAdmin": `[
			"User:u-22bbb"
		]`,
	"/api/metadata/security/v2alpha1/lookup/role/CloudClusterAdmin": `[
			"User:u-33ccc", "User:u-44ddd"
		]`,
}