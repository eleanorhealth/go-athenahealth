package athenahealth

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_GetSubscription(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/GetSubscription.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	subscription, err := athenaClient.GetSubscription("appointments")

	assert.NotNil(subscription)
	assert.Nil(err)
}

func TestHTTPClient_ListSubscriptionEvents(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./resources/ListSubscriptionEvents.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	events, err := athenaClient.ListSubscriptionEvents("appointments")

	assert.Len(events, 11)
	assert.Nil(err)
}

func TestHTTPClient_Subscribe(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "eventname=UpdateAppointment")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &SubscribeOpts{
		EventName: "UpdateAppointment",
	}
	err := athenaClient.Subscribe("appointments", opts)

	assert.Nil(err)
	assert.True(called)
}

func TestHTTPClient_Unsubscribe(t *testing.T) {
	assert := assert.New(t)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "eventname=UpdateAppointment")

		called = true
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	opts := &UnsubscribeOpts{
		EventName: "UpdateAppointment",
	}
	err := athenaClient.Unsubscribe("appointments", opts)

	assert.Nil(err)
	assert.True(called)
}
