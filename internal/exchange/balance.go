package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetTronBalance(address string) (float64, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", nileTestNetURL, address)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API вернуло ошибку: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	fmt.Println(string(body))

	var result AccountData

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("ошибка при обработке ответа: %v", err)
	}

	if len(result.Data) == 0 {
		return 0, fmt.Errorf("адрес не найден")
	}

	// Получение баланса и конвертация из Sun в TRX
	balance := float64(result.Data[0].Balance) / 1e6

	return balance, nil
}

type AccountData struct {
	Data []Data `json:"data"`
}

type Data struct {
	Balance int64 `json:"balance"`
}

const (
	mainNetURL     = "https://api.trongrid.io"
	nileTestNetURL = "https://nile.trongrid.io"
)
