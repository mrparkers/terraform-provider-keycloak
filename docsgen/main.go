package main

import (
	"fmt"
	"github.com/mrparkers/terraform-provider-keycloak/provider"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	OutputDir    = "docs"
	TemplateFile = "resource.tmpl"
)

type Attribute struct {
	Type          string
	Description   string
	Name          string
	ConflictsWith []string
}

type Resource struct {
	Name      string
	CanImport bool

	UserProvidedAttributes []Attribute
	ComputedAttributes     []Attribute
}

type ResourceMetadata struct {
	Description string `yaml:"description"`
	Example     string `yaml:"example"`
	Import      string `yaml:"import"`
}

type TemplateModel struct {
	Resource Resource
	Meta     ResourceMetadata
}

func readResourcesFromProvider() (resources []Resource, warnings []string) {
	terraformProviderResources := provider.KeycloakProvider().ResourcesMap

	for terraformResourceName, terraformResource := range terraformProviderResources {
		if !(terraformResourceName == "keycloak_openid_user_attribute_protocol_mapper" || terraformResourceName == "keycloak_group") {
			continue
		}
		resource := Resource{
			Name:      terraformResourceName,
			CanImport: terraformResource.Importer != nil,
		}

		for attributeName, attributeSchema := range terraformResource.Schema {
			attribute := Attribute{
				Name:          attributeName,
				ConflictsWith: attributeSchema.ConflictsWith,
				Description:   attributeSchema.Description,
			}

			if attribute.Description == "" {
				warnings = append(warnings, fmt.Sprintf("%s.%s is missing a description.\n", terraformResourceName, attributeName))
			}

			if attributeSchema.Required {
				attribute.Type = "Required"
			} else if attributeSchema.Optional {
				attribute.Type = "Optional"
			}

			if attributeSchema.Computed {
				resource.ComputedAttributes = append(resource.ComputedAttributes, attribute)
			} else {
				resource.UserProvidedAttributes = append(resource.UserProvidedAttributes, attribute)
			}
		}

		resources = append(resources, resource)
	}

	return resources, warnings
}

func generateDocs(resources []Resource) error {

	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		os.Mkdir(OutputDir, os.ModePerm)
	}

	for _, r := range resources {
		file, err := os.Create(fmt.Sprintf("%s/%s.md", OutputDir, r.Name))
		if err != nil {
			return err
		}

		metaFile, err := ioutil.ReadFile(fmt.Sprintf("resource-metadata/%s.meta.yml", r.Name))

		if err != nil {
			return err
		}

		var resourceMeta ResourceMetadata
		err = yaml.Unmarshal(metaFile, &resourceMeta)

		if err != nil {
			return err
		}

		t := template.Must(template.ParseFiles(TemplateFile))
		err = t.ExecuteTemplate(file, "base", TemplateModel{
			Resource: r,
			Meta:     resourceMeta,
		})

		if err != nil {
			return err
		}

		file.Close()
	}

	return nil
}

func main() {
	resources, warnings := readResourcesFromProvider()

	fmt.Printf("Discovered %d resources from provider \n", len(resources))

	if len(warnings) > 0 {
		fmt.Println("\nThe following issues should be addressed in order to generate high-quality documentation:")
	}

	for _, warning := range warnings {
		fmt.Printf("- %s", warning)
	}

	err := generateDocs(resources)

	if err != nil {
		fmt.Printf("Error occured while generating docs %v\n", err)
	}
}
