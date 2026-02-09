package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizePostalCode(t *testing.T) {
	cases := []struct {
		in       string
		wantCP4  string
		wantCP3  string
		wantErr  bool
	}{
		{"1100-038", "1100", "038", false},
		{"1100038", "1100", "038", false},
		{"1100", "", "", true},
		{"abcd-038", "", "", true},
		{"11000-038", "", "", true},
		{"", "", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			cp4, cp3, err := NormalizePostalCode(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tc.wantErr)
			}
			if cp4 != tc.wantCP4 || cp3 != tc.wantCP3 {
				t.Fatalf("got %q-%q, want %q-%q", cp4, cp3, tc.wantCP4, tc.wantCP3)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/postal-codes/1100-038" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"code":"1100-038","designation":"LISBOA","street":{"type":"Rua","name":"do Arsenal"},"locality":{"code":"21696","name":"Lisboa"},"municipality":{"code":"1106","name":"Lisboa"},"district":{"code":"11","name":"Lisboa"}}]`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	entries, err := c.Lookup(context.Background(), "1100", "038")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Designation != "LISBOA" {
		t.Fatalf("unexpected entries: %+v", entries)
	}
}

func TestLookupNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.Lookup(context.Background(), "0000", "000")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListDistricts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"code":"01","name":"Aveiro"},{"code":"11","name":"Lisboa"}]`))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	got, err := c.ListDistricts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1].Name != "Lisboa" {
		t.Fatalf("unexpected: %+v", got)
	}
}
