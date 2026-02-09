package api

type DistrictBrief struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type MunicipalityBrief struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type MunicipalityWithDistrict struct {
	Code     string        `json:"code"`
	Name     string        `json:"name"`
	District DistrictBrief `json:"district"`
}

type LocalityBrief struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Street struct {
	Type *string `json:"type"`
	Name *string `json:"name"`
}

type PostalCodeEntry struct {
	Code         string            `json:"code"`
	Designation  string            `json:"designation"`
	Street       Street            `json:"street"`
	Locality     LocalityBrief     `json:"locality"`
	Municipality MunicipalityBrief `json:"municipality"`
	District     DistrictBrief     `json:"district"`
}
