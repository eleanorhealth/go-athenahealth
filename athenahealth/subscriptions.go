package athenahealth

import (
	"context"
	"fmt"
	"net/url"
)

type Subscription struct {
	Status        string               `json:"status"`
	Subscriptions []*SubscriptionEvent `json:"subscriptions"`
}

// GetSubscription - List modified events appointment slots
//
// GET /v1/{practiceid}/appointments/changed/subscription
//
// https://docs.athenahealth.com/api/api-ref/appointment#Get-list-of-appointment-slot-change-subscription(s)
func (h *HTTPClient) GetSubscription(ctx context.Context, feedType string) (*Subscription, error) {
	out := &Subscription{}

	_, err := h.Get(ctx, fmt.Sprintf("/%s/changed/subscription", feedType), nil, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type SubscriptionEvent struct {
	EventName string `json:"eventname"`
}

type listSubscriptionEventsResponse struct {
	Subscriptions []*SubscriptionEvent `json:"subscriptions"`
}

// ListSubscriptionEvents - List events for appointments or appointment slots
//
// GET /v1/{practiceid}/appointments/changed/subscription/events
//
// https://docs.athenahealth.com/api/api-ref/appointment#Get-list-of-appointment-slot-change-events-to-which-you-can-subscribe
func (h *HTTPClient) ListSubscriptionEvents(ctx context.Context, feedType string) ([]*SubscriptionEvent, error) {
	out := &listSubscriptionEventsResponse{}

	_, err := h.Get(ctx, fmt.Sprintf("/%s/changed/subscription/events", feedType), nil, &out)
	if err != nil {
		return nil, err
	}

	return out.Subscriptions, nil
}

type SubscribeOptions struct {
	EventName string
}

// Subscribe - Subscribes for changed appointment slots events
//
// POST /v1/{practiceid}/appointments/changed/subscription
//
// https://docs.athenahealth.com/api/api-ref/appointment#Subscribe-to-all/specific-change-events-for-appointment-slots
func (h *HTTPClient) Subscribe(ctx context.Context, feedType string, opts *SubscribeOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}
		form.Add("eventname", opts.EventName)
	}

	_, err := h.PostForm(ctx, fmt.Sprintf("/%s/changed/subscription", feedType), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type UnsubscribeOptions struct {
	EventName string
}

// Unsubscribe - Unsubscribe to all/specific change events for appointment slots
//
// POST /v1/{practiceid}/appointments/changed/subscription
//
// https://docs.athenahealth.com/api/api-ref/appointment#Unsubscribe-to-all/specific-change-events-for-appointment-slots
func (h *HTTPClient) Unsubscribe(ctx context.Context, feedType string, opts *UnsubscribeOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}
		form.Add("eventname", opts.EventName)
	}

	_, err := h.DeleteForm(ctx, fmt.Sprintf("/%s/changed/subscription", feedType), form, nil)
	if err != nil {
		return err
	}

	return nil
}
