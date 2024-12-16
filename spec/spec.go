package spec

import (
	"fmt"
	"os"

	"log"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func loadSpec(specFile string) (*libopenapi.DocumentModel[v3.Document], error) {
	// Load config
	spec, err := os.ReadFile(specFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %e", err)
	}

	specDocument, err := libopenapi.NewDocument(spec)
	if err != nil {
		return nil, fmt.Errorf("failed creating spec from %s: %e", specFile, err)
	}

	docModel, errors := specDocument.BuildV3Model()

	if len(errors) > 0 {
		for i := range errors {
			log.Printf("error: %e\n", errors[i])
		}
		return nil, fmt.Errorf("cannot create openApi v3 model from document: %d errors reported", len(errors))
	}

	return docModel, nil
}

func PrintSpec(specFile string) error {
	spec, err := loadSpec(specFile)
	if err != nil {
		log.Printf("Error loading spec: %v", err) // Log the error
		return err
	}
	fmt.Printf("%s %s - %s\n\n", spec.Model.Info.Title, spec.Model.Info.Version, spec.Model.Info.Description)
	return nil
}
