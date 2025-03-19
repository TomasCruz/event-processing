package eventlistener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://api.freecurrencyapi.com/v1/latest?apikey=fca_live_cg50JWFU9ahwY8G3JKfYTppj3vWd8E9jWihmiN4Z&currencies=EUR&base_currency=NZD
// {"data":{"EUR":0.5281759606}}

// curl --request GET 'https://api.apilayer.com/exchangerates_data/latest?base=BTC&symbols=EUR' --header 'apikey: mGlwwvQCovT2h9oLloIbVLc7T2c7rESq'
// {
//     "success": true,
//     "timestamp": 1742044563,
//     "base": "BTC",
//     "date": "2025-03-15",
//     "rates": {
//         "EUR": 77070.53064
//     }
// }

type APILayerResp struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
}

type FreecurrencyAPIResp struct {
	Data map[string]float64 `json:"data"`
}

func (eSvc enricherSvc) composeAPILayerReq(baseCurrency string) (*http.Request, error) {
	urlTemplate := "https://api.apilayer.com/exchangerates_data/latest?base=%s&symbols=EUR"
	url := fmt.Sprintf(urlTemplate, baseCurrency)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", eSvc.config.ApiLayerAPIKey)
	return req, nil
}

func (eSvc enricherSvc) composeFreecurrencyAPIReq(baseCurrency string) (*http.Request, error) {
	urlTemplate := "https://api.freecurrencyapi.com/v1/latest?apikey=%s&currencies=EUR&base_currency=%s"
	url := fmt.Sprintf(urlTemplate, eSvc.config.FreeCurrencyAPIKey, baseCurrency)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (eSvc enricherSvc) doAPILayerReq(baseCurrency string) (float64, error) {
	req, err := eSvc.composeAPILayerReq(baseCurrency)
	if err != nil {
		return 0, err
	}

	b, err := respBodyBytes(req)
	if err != nil {
		return 0, err
	}

	return extractRateFromAPILayer(b)
}

func (eSvc enricherSvc) doFreecurrencyAPIReq(baseCurrency string) (float64, error) {
	req, err := eSvc.composeFreecurrencyAPIReq(baseCurrency)
	if err != nil {
		return 0, err
	}

	b, err := respBodyBytes(req)
	if err != nil {
		return 0, err
	}

	return extractRateFromFreecurrencyAPI(b)
}

func respBodyBytes(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func extractRateFromAPILayer(respBytes []byte) (float64, error) {
	var respStruct APILayerResp
	err := json.Unmarshal(respBytes, &respStruct)
	if err != nil {
		return 0, err
	}

	if !respStruct.Success {
		return 0, fmt.Errorf("extractRateFromAPILayer !Success")
	}

	rate, ok := respStruct.Rates["EUR"]
	if !ok {
		return 0, fmt.Errorf("extractRateFromAPILayer no EUR rate")
	}

	return rate, nil
}

func extractRateFromFreecurrencyAPI(respBytes []byte) (float64, error) {
	var respStruct FreecurrencyAPIResp
	err := json.Unmarshal(respBytes, &respStruct)
	if err != nil {
		return 0, err
	}

	rate, ok := respStruct.Data["EUR"]
	if !ok {
		return 0, fmt.Errorf("extractRateFromFreecurrencyAPI no EUR rate")
	}

	return rate, nil
}
