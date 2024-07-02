// blockchain/utils.go

package blockchain

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func PostToServer(url string, data interface{}) (*http.Response, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    return resp, nil
}
