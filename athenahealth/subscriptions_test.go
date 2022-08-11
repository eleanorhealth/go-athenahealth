package athenahealth

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetSubscription(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/GetSubscription.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	subscription, err := athenaClient.GetSubscription(context.Background(), "appointments")

	assert.NotNil(subscription)
	assert.NoError(err)
}

func TestHTTPClient_ListSubscriptionEvents(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := os.ReadFile("./resources/ListSubscriptionEvents.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	events, err := athenaClient.ListSubscriptionEvents(context.Background(), "appointments")

	assert.Len(events, 11)
	assert.NoError(err)
}

func TestHTTPClient_Subscribe(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "eventname=UpdateAppointment")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &SubscribeOptions{
		EventName: "UpdateAppointment",
	}
	err := athenaClient.Subscribe(context.Background(), "appointments", opts)

	assert.NoError(err)
	assert.True(called)
}

func TestHTTPClient_Unsubscribe(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "eventname=UpdateAppointment")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &UnsubscribeOptions{
		EventName: "UpdateAppointment",
	}
	err := athenaClient.Unsubscribe(context.Background(), "appointments", opts)

	assert.NoError(err)
	assert.True(called)
}
