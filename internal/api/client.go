package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

const DefaultBaseURL = "https://postcode-pt.onrender.com/v1"

var (
	ErrNotFound  = errors.New("not found")
	ErrInvalidCP = errors.New("invalid postal code format")
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 60 * time.Second},
	}
}

var postalCodeRe = regexp.MustCompile(`^(\d{4})-?(\d{3})$`)

func NormalizePostalCode(raw string) (cp4, cp3 string, err error) {
	m := postalCodeRe.FindStringSubmatch(raw)
	if m == nil {
		return "", "", ErrInvalidCP
	}
	return m[1], m[2], nil
}

func (c *Client) Lookup(ctx context.Context, cp4, cp3 string) ([]PostalCodeEntry, error) {
	var out []PostalCodeEntry
	err := c.get(ctx, fmt.Sprintf("/postal-codes/%s-%s", cp4, cp3), &out)
	return out, err
}

func (c *Client) ListDistricts(ctx context.Context) ([]DistrictBrief, error) {
	var out []DistrictBrief
	err := c.get(ctx, "/districts", &out)
	return out, err
}

func (c *Client) ListMunicipalities(ctx context.Context, districtCode string) ([]MunicipalityWithDistrict, error) {
	var out []MunicipalityWithDistrict
	err := c.get(ctx, fmt.Sprintf("/districts/%s/municipalities", districtCode), &out)
	return out, err
}

func (c *Client) get(ctx context.Context, path string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case resp.StatusCode >= 400:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s: %s", resp.Status, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}
