package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	chlib "github.com/alridev/ch-lib"
)

func TestSetterMany(t *testing.T) {
	tests := []struct {
		name           string
		events         []chlib.SetterEvent
		responseStatus int
		responseBody   chlib.SetterManyResponseBody
		expectedErrors map[int32]string
		expectedErr    bool
	}{
		{
			name: "successful request",
			events: []chlib.SetterEvent{
				{
					TableName: "test_table",
					Data:      map[string]interface{}{"field": "value"},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: chlib.SetterManyResponseBody{
				Errors: map[int32]string{},
			},
			expectedErrors: map[int32]string{},
			expectedErr:    false,
		},
		{
			name: "with errors",
			events: []chlib.SetterEvent{
				{
					TableName: "test_table",
					Data:      map[string]interface{}{"field": "value"},
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: chlib.SetterManyResponseBody{
				Errors: map[int32]string{
					0: "invalid data",
				},
			},
			expectedErrors: map[int32]string{
				0: "invalid data",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/setter/many", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			client := chlib.NewChClient(server.URL, "test-token", "getter-token", nil)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			errors, err := client.SetterMany(ctx, tt.events)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErrors, errors)
			}
		})
	}
}

func TestSetterOne(t *testing.T) {
	tests := []struct {
		name           string
		event          chlib.SetterEvent
		responseStatus int
		responseBody   chlib.SetterOneResponseBody
		expectedErr    bool
	}{
		{
			name: "successful request",
			event: chlib.SetterEvent{
				TableName: "test_table",
				Data:      map[string]interface{}{"field": "value"},
			},
			responseStatus: http.StatusOK,
			responseBody:   chlib.SetterOneResponseBody{},
			expectedErr:    false,
		},
		{
			name: "with error",
			event: chlib.SetterEvent{
				TableName: "test_table",
				Data:      map[string]interface{}{"field": "value"},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: chlib.SetterOneResponseBody{
				Error: "invalid data",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/setter", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			client := chlib.NewChClient(server.URL, "test-token", "getter-token", nil)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := client.SetterOne(ctx, tt.event)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetterByTable(t *testing.T) {
	tests := []struct {
		name           string
		tableName      string
		data           interface{}
		responseStatus int
		responseBody   chlib.SetterByTableResponseBody
		expectedErr    bool
	}{
		{
			name:           "successful request",
			tableName:      "test_table",
			data:           map[string]interface{}{"field": "value"},
			responseStatus: http.StatusOK,
			responseBody:   chlib.SetterByTableResponseBody{},
			expectedErr:    false,
		},
		{
			name:           "with error",
			tableName:      "test_table",
			data:           map[string]interface{}{"field": "value"},
			responseStatus: http.StatusBadRequest,
			responseBody: chlib.SetterByTableResponseBody{
				Error: "invalid data",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/setter/"+tt.tableName, r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			client := chlib.NewChClient(server.URL, "test-token", "getter-token", nil)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := client.SetterByTable(ctx, tt.tableName, tt.data)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
