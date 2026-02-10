package render

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/RobertoCCC/postcode-pt-cli/internal/api"
)

func ptrString(s string) *string { return &s }

func TestPostalCodesText(t *testing.T) {
	var buf bytes.Buffer
	opts := NewOptions(&buf)
	opts.NoColor = true

	entries := []api.PostalCodeEntry{
		{
			Code:        "1100-038",
			Designation: "LISBOA",
			Street:      api.Street{Type: ptrString("Rua"), Name: ptrString("do Arsenal")},
			Locality:    api.LocalityBrief{Code: "21696", Name: "Lisboa"},
			Municipality: api.MunicipalityBrief{Code: "1106", Name: "Lisboa"},
			District:    api.DistrictBrief{Code: "11", Name: "Lisboa"},
		},
	}
	if err := PostalCodes(opts, entries); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"1100-038", "LISBOA", "Rua do Arsenal", "Lisboa", "1106", "11"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n--- output ---\n%s", want, out)
		}
	}
}

func TestPostalCodesJSON(t *testing.T) {
	var buf bytes.Buffer
	opts := NewOptions(&buf)
	opts.JSON = true

	entries := []api.PostalCodeEntry{{Code: "1100-038", Designation: "LISBOA"}}
	if err := PostalCodes(opts, entries); err != nil {
		t.Fatal(err)
	}

	var got []api.PostalCodeEntry
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if len(got) != 1 || got[0].Code != "1100-038" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestDistrictsText(t *testing.T) {
	var buf bytes.Buffer
	opts := NewOptions(&buf)
	opts.NoColor = true

	if err := Districts(opts, []api.DistrictBrief{{Code: "11", Name: "Lisboa"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Lisboa") || !strings.Contains(buf.String(), "11") {
		t.Fatalf("unexpected: %s", buf.String())
	}
}

func TestMunicipalitiesText(t *testing.T) {
	var buf bytes.Buffer
	opts := NewOptions(&buf)
	opts.NoColor = true

	if err := Municipalities(opts, []api.MunicipalityWithDistrict{
		{Code: "1106", Name: "Lisboa", District: api.DistrictBrief{Code: "11", Name: "Lisboa"}},
	}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Lisboa") || !strings.Contains(buf.String(), "1106") {
		t.Fatalf("unexpected: %s", buf.String())
	}
}
