package athenahealth

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetProvider(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/GetProvider.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	provider, err := athenaClient.GetProvider(context.Background(), "1")

	assert.NotNil(provider)
	assert.NoError(err)
}

func TestHTTPClient_ListChangedProviders(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("leaveunprocessed"))
		assert.Equal("06/01/2020 15:30:45", r.URL.Query().Get("showprocessedstartdatetime"))
		assert.Equal("06/02/2020 12:30:45", r.URL.Query().Get("showprocessedenddatetime"))

		b, _ := os.ReadFile("./resources/ListChangedProviders.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListChangedProviderOptions{
		LeaveUnprocessed:           true,
		ShowProcessedStartDatetime: time.Date(2020, 6, 1, 15, 30, 45, 0, time.UTC),
		ShowProcessedEndDatetime:   time.Date(2020, 6, 2, 12, 30, 45, 0, time.UTC),
	}

	patients, err := athenaClient.ListChangedProviders(context.Background(), opts)

	assert.Len(patients, 1)
	assert.NoError(err)
}

func TestHTTPClient_ListProviders(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("true", r.URL.Query().Get("showallproviderids"))

		b, _ := os.ReadFile("./resources/ListProviders.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &ListProvidersOptions{
		ShowAllProviderIDs: true,
	}

	res, err := athenaClient.ListProviders(context.Background(), opts)

	assert.Len(res.Providers, 1)
	assert.Equal(res.Pagination.NextOffset, 30)
	assert.Equal(res.Pagination.PreviousOffset, 10)
	assert.Equal(res.Pagination.TotalCount, 1)
	assert.NoError(err)
}
