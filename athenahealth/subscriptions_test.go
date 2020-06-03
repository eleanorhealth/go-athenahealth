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

	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "UpdateAppointment")

		b, _ := ioutil.ReadFile("./resources/Subscribe.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.Subscribe("appointments", "UpdateAppointment")

	assert.Nil(err)
}

func TestHTTPClient_Unsubscribe(t *testing.T) {
	assert := assert.New(t)

	h := func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		assert.Contains(string(reqBody), "UpdateAppointment")

		b, _ := ioutil.ReadFile("./resources/Unsubscribe.json")
		w.Write(b)
	}

	athenaClient, ts := testClient(h)
	defer ts.Close()

	err := athenaClient.Unsubscribe("appointments", "UpdateAppointment")

	assert.Nil(err)
}
