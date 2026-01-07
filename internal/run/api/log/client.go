package log

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// Client defines the interface for Logging operations.
type Client interface {
	Entries(ctx context.Context, opts ...interface{}) EntryIterator
	Close() error
}

// EntryIterator defines the interface for iterating over log entries.
type EntryIterator interface {
	Next() (*logging.Entry, error)
}

// ClientFactory is a function that returns a Client.
type ClientFactory func(ctx context.Context, projectID string) (Client, error)

var clientFactory ClientFactory = NewGCPClient

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct {
	client *logadmin.Client
}

// NewGCPClient creates a new GCPClient.
func NewGCPClient(ctx context.Context, projectID string) (Client, error) {
	creds, err := google.FindDefaultCredentials(ctx, logging.ReadScope)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	c, err := logadmin.NewClient(ctx, projectID, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &GCPClient{client: c}, nil
}

func (c *GCPClient) Entries(ctx context.Context, opts ...interface{}) EntryIterator {
	var logOpts []logadmin.EntriesOption
	for _, o := range opts {
		if opt, ok := o.(logadmin.EntriesOption); ok {
			logOpts = append(logOpts, opt)
		}
	}
	return &GCPEntryIterator{it: c.client.Entries(ctx, logOpts...)}
}

func (c *GCPClient) Close() error {
	return c.client.Close()
}

// GCPEntryIterator wraps logadmin.EntryIterator.
type GCPEntryIterator struct {
	it *logadmin.EntryIterator
}

func (it *GCPEntryIterator) Next() (*logging.Entry, error) {
	return it.it.Next()
}
