package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/eleanorhealth/go-athenahealth/athenahealth"
	"github.com/eleanorhealth/go-athenahealth/athenahealth/tokencacher"
)

func newAthenaClient() athenahealth.Client {
	athenaClient := athenahealth.NewHTTPClient(
		&http.Client{},
		os.Getenv("ATHENA_PRACTICE_ID"),
		os.Getenv("ATHENA_CLIENT_ID"),
		os.Getenv("ATHENA_SECRET"),
	)

	athenaClient.
		WithPreview(true).
		WithTokenCacher(tokencacher.NewFile("/tmp/token.json"))

	return athenaClient
}

func customFields(ctx context.Context, client athenahealth.Client, patientID string) {
	fields, err := client.ListCustomFields(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fields {
		log.Println("id %s name %s type %s", f.CustomFieldID, f.Name, f.Type)
	}

	cfields, err := client.GetPatientCustomFields(ctx, patientID, "36")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range cfields {
		log.Println("id %s values %s", f.CustomFieldID, f.CustomFieldValue)
	}

	// test if inactive reasons appear in output
}

func physicalExam(ctx context.Context, client athenahealth.Client, encounterID string) {
	exam, err := client.GetPhysicalExam(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(exam)
	
	// test if eh asam and loc appear in output
}

func social(ctx context.Context, client athenahealth.Client, patientID, depID string) {
	opts := &athenahealth.GetPatientSocialHistoryOptions{
		DepartmentID: depID,
	}
	res, err := client.GetPatientSocialHistory(ctx, patientID, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, q := range res.Questions {
		if q.QuestionID == "243" {
			log.Printf("legal needs answer: %s", q.Answer)
		}
		if q.QuestionID == "249" {
			log.Printf("living situation answer: %s", q.Answer)
		}

	}
}

func main() {
	client := newAthenaClient()
	ctx := context.Background()

	customFields(ctx, client, "85802")
	physicalExam(ctx, client, "4430925")
	social(ctx, client, "87476", "11")
}
