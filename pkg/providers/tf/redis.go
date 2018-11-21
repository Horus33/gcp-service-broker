// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tf

import (
	"log"

	accountmanagers "github.com/GoogleCloudPlatform/gcp-service-broker/brokerapi/brokers/account_managers"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

func init() {
	redis := TfServiceDefinitionV1{
		Version:          1,
		Name:             "google-redis",
		Id:               "0e86ad78-99b3-48b6-a986-b594e7995fd6",
		Description:      "Cloud Memorystore for Redis is a fully managed Redis service for the Google Cloud Platform.",
		DisplayName:      "Google Cloud Memorystore for Redis",
		ImageUrl:         "https://cloud.google.com/_static/images/cloud/products/logos/svg/cache.svg",
		DocumentationUrl: "https://cloud.google.com/memorystore/docs/redis/",
		SupportUrl:       "https://cloud.google.com/memorystore/docs/redis/getting-support",
		Tags:             []string{"preview", "gcp", "terraform", "redis", "memorystore"},
		Plans: []broker.ServicePlan{
			{
				ServicePlan: brokerapi.ServicePlan{
					ID:   "e1d11f65-da66-46ad-977c-6d56513baf43",
					Name: "standard",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Standard",
					},
					Description: "Standard storage class.",
				},
				ServiceProperties: map[string]string{
					"storage_class": "STANDARD",
				},
			},
		},
		ProvisionSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{
				{
					FieldName: "capacity_gb",
					Type:      broker.JsonTypeInteger,
					Details:   "The memory capacity of the database in gigabytes.",
					Required:  true,
				},
				{
					FieldName: "service_tier",
					Type:      broker.JsonTypeString,
					Details:   "The storage class of the bucket. See: https://cloud.google.com/storage/docs/storage-classes.",
					Required:  true,
					Enum: map[interface{}]string{
						"BASIC":       "A basic tier, data will be erased between restarts.",
						"STANDARD_HA": "A tier with a backup device in a different location.",
					},
				},
			},
			UserInputs: []broker.BrokerVariable{
				{
					FieldName: "instance_id",
					Type:      broker.JsonTypeString,
					Details:   "Permanent identifier for your instance",
					Default:   "pcf-sb-${counter.next()}-${time.nano()}",
					Constraints: validation.NewConstraintBuilder().
						Pattern("^[a-z][a-z0-9-]+$").
						MinLength(6).
						MaxLength(30).
						Build(),
				},
				{
					FieldName: "display_name",
					Type:      broker.JsonTypeString,
					Details:   "For display purposes only",
					Default:   "${instance_id}",
					Constraints: validation.NewConstraintBuilder().
						MaxLength(80).
						Build(),
				},
				{
					FieldName: "region",
					Type:      broker.JsonTypeString,
					Details:   "The region of the Redis instance.",
					Default:   "us-central1",
					Constraints: validation.NewConstraintBuilder().
						Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
						Examples("us-central1", "asia-northeast1").
						Build(),
				},
				{
					FieldName: "zone",
					Type:      broker.JsonTypeString,
					Details:   "The zone within the region or any.",
					Default:   "any",
					Constraints: validation.NewConstraintBuilder().
						Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
						Examples("us-central1", "asia-northeast1").
						Build(),
				},
				{
					FieldName: "authorized_network",
					Type:      broker.JsonTypeString,
					Details:   "The name of the authorized network this instance will be connected to.",
					Default:   "default",
					Constraints: validation.NewConstraintBuilder().
						Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
						Build(),
				},
			},
			Computed: []varcontext.DefaultVariable{
				{Name: "labels", Default: "${json.marshal(request.default_labels)}", Overwrite: true},
			},
			Template: `
	variable name {type = "string"}
	variable location {type = "string"}
	variable storage_class {type = "string"}

	resource "google_storage_bucket" "bucket" {
	  name     = "${var.name}"
	  location = "${var.location}"
	  storage_class = "${var.storage_class}"
	}

	output id {value = "${google_storage_bucket.bucket.id}"}
	output bucket_name {value = "${var.name}"}
	`,
			Outputs: []broker.BrokerVariable{
				{
					FieldName: "bucket_name",
					Type:      broker.JsonTypeString,
					Details:   "Name of the bucket this binding is for.",
					Required:  true,
					Constraints: validation.NewConstraintBuilder(). // https://cloud.google.com/storage/docs/naming
											Pattern("^[A-Za-z0-9_\\.]+$").
											MinLength(3).
											MaxLength(222).
											Build(),
				},
				{
					FieldName:   "id",
					Type:        broker.JsonTypeString,
					Details:     "The GCP ID of this bucket.",
					Required:    true,
					Constraints: validation.NewConstraintBuilder().Build(),
				},
			},
		},
		BindSettings: TfServiceDefinitionV1Action{
			PlanInputs: []broker.BrokerVariable{},
			UserInputs: []broker.BrokerVariable{},
			Computed:   []varcontext.DefaultVariable{},
			Template: `
	variable service_account_name {type = "string"}
	variable service_account_display_name {type = "string"}

	resource "google_service_account" "account" {
	  account_id = "${var.service_account_name}"
	  display_name = "${var.service_account_display_name}"
	}

	resource "google_service_account_key" "key" {
	  service_account_id = "${google_service_account.account.name}"
	}

	resource "google_storage_bucket_iam_member" "member" {
	  bucket = "${var.bucket}"
	  role   = "roles/${var.role}"
	  member = "serviceAccount:${google_service_account.account.email}"
	}

	output "Name" {value = "${google_service_account.account.display_name}"}
	output "Email" {value = "${google_service_account.account.email}"}
	output "UniqueId" {value = "${google_service_account.account.unique_id}"}
	output "PrivateKeyData" {value = "${google_service_account_key.key.private_key}"}
	output "ProjectId" {value = "${google_service_account.account.project}"}
	`,
			Outputs: accountmanagers.ServiceAccountBindOutputVariables(),
		},

		Examples: []broker.ServiceExample{
			{
				Name:            "Basic Configuration",
				Description:     "Create a tiny Redis instance for development.",
				PlanId:          "6ed44104-8777-4b57-8c03-826b3af7d0be",
				ProvisionParams: map[string]interface{}{},
				BindParams:      map[string]interface{}{},
			},
		},

		// Internal SHOULD be set to true for Google maintained services.
		Internal: true,
	}

	service, err := redis.ToService()
	if err != nil {
		log.Fatal(err)
	}
	broker.Register(service)
}
