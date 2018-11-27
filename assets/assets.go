package assets

import (
	"archive/zip"

	"github.com/GoogleCloudPlatform/gcp-service-broker/assets/tfsources"
	"github.com/GoogleCloudPlatform/gcp-service-broker/utils"
)

// Sources
//go:generate go run fetch.go --dest "terraform-sources/terraform-v0.11.9.zip" --url "https://github.com/hashicorp/terraform/archive/v0.11.9.zip"
//go:generate go run fetch.go --dest "terraform-sources/terraform-provider-google-v1.19.0.zip" --url "https://github.com/terraform-providers/terraform-provider-google/archive/v1.19.0.zip"

// Generators for Terraform
//go:generate go run fetch.go --dest "terraform-linux-amd64.zip" --url "https://releases.hashicorp.com/terraform/0.11.9/terraform_0.11.9_linux_amd64.zip"
//go:generate go run fetch.go --dest "terraform-linux-386.zip" --url "https://releases.hashicorp.com/terraform/0.11.9/terraform_0.11.9_linux_386.zip"
//go:generate go run unzip.go --src "terraform-linux-386.zip" --dest "terraform-linux-386/bin"
//go:generate go run unzip.go --src "terraform-linux-amd64.zip" --dest "terraform-linux-amd64/bin"

// Generators for the Google provider
//go:generate go run fetch.go --dest "terraform-google-beta-linux-amd64.zip" --url "https://releases.hashicorp.com/terraform-provider-google-beta/1.19.0/terraform-provider-google-beta_1.19.0_linux_amd64.zip"
//go:generate go run fetch.go --dest "terraform-google-beta-linux-386.zip" --url "https://releases.hashicorp.com/terraform-provider-google-beta/1.19.0/terraform-provider-google-beta_1.19.0_linux_386.zip"
//go:generate go run unzip.go --src "terraform-google-beta-linux-amd64.zip" --dest "terraform-linux-386/providers"
//go:generate go run unzip.go --src "terraform-google-beta-linux-386.zip" --dest "terraform-linux-amd64/providers"

//go:generate go run pack.go . tflinux386 terraform-linux-386
//go:generate go run pack.go . tflinuxamd64 terraform-linux-amd64
//go:generate go run pack.go . tfsources terraform-sources

// DumpSources extracts the Terraform sources zip to the given directory.
func DumpSources(outputDirectory string) error {
	return extractEmbedded(tfsources.NewZipReader, outputDirectory)
}

func extractEmbedded(readerBuilder func() (*zip.Reader, error), outputDirectory string) error {
	reader, err := readerBuilder()
	if err != nil {
		return err
	}

	return utils.Unzip(reader, outputDirectory)
}
