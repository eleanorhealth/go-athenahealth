package athenahealth

import (
	"fmt"
	"net/url"
)

type Subscription struct {
	Status        string               `json:"status"`
	Subscriptions []*SubscriptionEvent `json:"subscriptions"`
}

// GetSubscription - Handles managing subscriptions for changed appointment slots.
// GET /v1/{practiceid}/appointments/changed/subscription
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-7
func (h *HTTPClient) GetSubscription(feedType string) (*Subscription, error) {
	out := &Subscription{}

	_, err := h.Get(fmt.Sprintf("/%s/changed/subscription", feedType), nil, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type SubscriptionEvent struct {
	EventName string `json:"eventname"`
}

type subscriptionEventsResponse struct {
	Subscriptions []*SubscriptionEvent `json:"subscriptions"`
}

// ListSubscriptionEvents - Returns the list of events you can subscribe to for changed appointment slots.
// GET /v1/{practiceid}/appointments/changed/subscription/events.
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-8
func (h *HTTPClient) ListSubscriptionEvents(feedType string) ([]*SubscriptionEvent, error) {
	out := &subscriptionEventsResponse{}

	_, err := h.Get(fmt.Sprintf("/%s/changed/subscription/events", feedType), nil, &out)
	if err != nil {
		return nil, err
	}

	return out.Subscriptions, nil
}

// Subscribe - Handles subscriptions for changed appointment slots.
// POST /v1/{practiceid}/appointments/changed/subscription
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-6
func (h *HTTPClient) Subscribe(feedType, eventName string) error {
	var form url.Values

	if eventName != "" {
		form = url.Values{}
		form.Add("eventname", eventName)
	}

	_, err := h.PostForm(fmt.Sprintf("/%s/changed/subscription", feedType), form, nil)
	if err != nil {
		return err
	}

	return nil
}

// Unsubscribe - Handles subscriptions for changed appointment slots.
// POST /v1/{practiceid}/appointments/changed/subscription
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-6
func (h *HTTPClient) Unsubscribe(feedType, eventName string) error {
	var form url.Values

	if eventName != "" {
		form = url.Values{}
		form.Add("eventname", eventName)
	}

	_, err := h.DeleteForm(fmt.Sprintf("/%s/changed/subscription", feedType), form, nil)
	if err != nil {
		return err
	}

	return nil
}
