package log

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

// MockClient is a mock implementation of Client.
type MockClient struct {
	EntriesFunc func(ctx context.Context, opts ...interface{}) EntryIterator
	CloseFunc   func() error
}

func (m *MockClient) Entries(ctx context.Context, opts ...interface{}) EntryIterator {
	if m.EntriesFunc != nil {
		return m.EntriesFunc(ctx, opts...)
	}
	return &MockEntryIterator{}
}

func (m *MockClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockEntryIterator is a mock implementation of EntryIterator.
type MockEntryIterator struct {
	Items []*logging.Entry
	Index int
	Err   error
}

func (m *MockEntryIterator) Next() (*logging.Entry, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Index >= len(m.Items) {
		return nil, iterator.Done
	}
	item := m.Items[m.Index]
	m.Index++
	return item, nil
}

func TestFormatEntry(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:00Z")
	entry := &logging.Entry{
		Timestamp: ts,
		Payload:   "Log message",
	}

	result := formatEntry(entry)
	assert.Equal(t, "[10:00:00] Log message", result)
}

func TestStreamLogs(t *testing.T) {
	// Restore clientFactory after test
	origFactory := clientFactory
	defer func() { clientFactory = origFactory }()

	t.Run("Backlog and Polling", func(t *testing.T) {
		// Speed up polling
		origInterval := pollInterval
		pollInterval = 10 * time.Millisecond
		defer func() { pollInterval = origInterval }()

		// Prepare mock data
		ts1, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:00Z")
		ts2, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:01Z")
		ts3, _ := time.Parse(time.RFC3339, "2023-10-27T10:00:02Z")

		// Backlog: Newest First (so 2 then 1)
		backlogLogs := []*logging.Entry{
			{Timestamp: ts2, Payload: "Log 2"},
			{Timestamp: ts1, Payload: "Log 1"},
		}

		// New logs: Oldest First (so 3)
		newLogs := []*logging.Entry{
			{Timestamp: ts3, Payload: "Log 3"},
		}

		// Refined Mock Logic
		callCount := 0
		clientFactory = func(ctx context.Context, projectID string) (Client, error) {
			return &MockClient{
				EntriesFunc: func(ctx context.Context, opts ...interface{}) EntryIterator {
					callCount++
					if callCount == 1 {
						return &MockEntryIterator{Items: backlogLogs}
					}
					if callCount == 2 {
						return &MockEntryIterator{Items: newLogs}
					}
					return &MockEntryIterator{Items: nil} // No more logs
				},
			}, nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		logChan := make(chan string, 10)

		// Run StreamLogs in a goroutine
		errChan := make(chan error)
		go func() {
			errChan <- StreamLogs(ctx, "project", "filter", logChan)
		}()

		// Allow some time for processing
		time.Sleep(100 * time.Millisecond) // Wait for backlog

		// Check Backlog (Should be reversed: Log 1, Log 2)
		select {
		case msg := <-logChan:
			assert.Contains(t, msg, "Log 1")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Log 1")
		}

		select {
		case msg := <-logChan:
			assert.Contains(t, msg, "Log 2")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Log 2")
		}

		// Check Polling (Log 3)
		select {
		case msg := <-logChan:
			assert.Contains(t, msg, "Log 3")
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Log 3")
		}

		cancel() // Stop the loop
		<-errChan
	})

	t.Run("Client Creation Error", func(t *testing.T) {
		expectedErr := errors.New("client error")
		clientFactory = func(ctx context.Context, projectID string) (Client, error) {
			return nil, expectedErr
		}

		logChan := make(chan string)
		err := StreamLogs(context.Background(), "p", "f", logChan)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestStreamLogs_PollingError(t *testing.T) {
	// Restore clientFactory after test
	origFactory := clientFactory
	defer func() { clientFactory = origFactory }()

	// Speed up polling
	origInterval := pollInterval
	pollInterval = 10 * time.Millisecond
	defer func() { pollInterval = origInterval }()

	callCount := 0
	clientFactory = func(ctx context.Context, projectID string) (Client, error) {
		return &MockClient{
			EntriesFunc: func(ctx context.Context, opts ...interface{}) EntryIterator {
				callCount++
				if callCount == 1 {
					// Backlog: empty
					return &MockEntryIterator{Items: nil}
				}
				// Polling: Error
				return &MockEntryIterator{Err: errors.New("polling error")}
			},
		}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logChan := make(chan string)

	// StreamLogs should not return error, but just retry or stop?
	// The code breaks the inner loop on error, then waits for next tick.
	// So it should continue running.

	// We want to verify it DOES NOT exit StreamLogs.
	done := make(chan error)
	go func() {
		done <- StreamLogs(ctx, "p", "f", logChan)
	}()

	select {
	case <-done:
		t.Fatal("StreamLogs exited unexpectedly on polling error")
	case <-time.After(100 * time.Millisecond):
		// passed
	}
}
