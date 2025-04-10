package debank

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Protocol struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Chain      string  `json:"chain"`
	Portfolio  Portfolio `json:"portfolio_item_list"`
}

type Portfolio struct {
	Name     string  `json:"name"`
	NetUSD   float64 `json:"net_usd_value"`
}

func FetchProtocols(wallet string) ([]Protocol, error) {
	url := fmt.Sprintf("https://pro-openapi.debank.com/v1/user/complex_protocol_list?id=%s", wallet)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []Protocol
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
