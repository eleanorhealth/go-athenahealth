package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetPhysicalExam(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/PhysicalExam.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	exam, err := athenaClient.GetPhysicalExam(context.Background(), "1", nil)
	assert.NoError(err)
	assert.NotNil(exam)
}
