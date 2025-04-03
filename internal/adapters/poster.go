package adapters

import (
	"bytes"
	"certificate/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

// Client — структура данных клиента
type posterClient struct {
	ClientName     string `json:"client_name"`
	ClientSex      int    `json:"client_sex"`
	ClientGroupsID int    `json:"client_groups_id_client"`
	Phone          string `json:"phone"`
	Birthday       string `json:"birthday"`
}

type Response struct {
	Response int `json:"response"`
}

type BonusUpdateRequest struct {
	ClientID int `json:"client_id"`
	Count    int `json:"count"`
}

type PosterAPI struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

// NewPosterAPI создает новый экземпляр PosterAPI
func NewPosterAPI(token string) *PosterAPI {
	return &PosterAPI{
		BaseURL: "https://joinposter.com/api/",
		Token:   token,
		Client:  &http.Client{},
	}
}

// ChangeClientBonus изменяет количество бонусов у клиента
func (p *PosterAPI) ChangeClientBonus(clientID int) error {
	url := fmt.Sprintf("%sclients.changeClientBonus?token=%s", p.BaseURL, p.Token)

	requestBody := BonusUpdateRequest{
		ClientID: clientID,
		Count:    1000,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		slog.Error("Ошибка маршалинга JSON", "error", err)
		return err
	}

	slog.Info("Отправка запроса на изменение бонусов", "clientID", clientID, "count", requestBody.Count, "url", url)

	resp, err := p.Client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("Ошибка при выполнении запроса", "error", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Ошибка чтения ответа", "error", err)
		return err
	}

	slog.Info("Ответ от Poster API", "response", string(body))
	return nil
}

// CreateClient создает нового клиента, если его нет в базе
func (p *PosterAPI) CreateClient(c domain.Client) (int, error) {
	existingClientID, err := p.findClientByPhone(c.Phone)
	if err != nil {
		slog.Error("Ошибка при поиске клиента", "phone", c.Phone, "error", err)
		return 0, fmt.Errorf("ошибка при поиске клиента: %w", err)
	}
	if existingClientID != 0 {
		slog.Info("Клиент уже зарегистрирован", "phone", c.Phone, "clientID", existingClientID)
		return existingClientID, nil
	}

	url := fmt.Sprintf("%sclients.createClient?token=%s", p.BaseURL, p.Token)
	posterC := toPosterClient(c)

	clientData, err := json.Marshal(posterC)
	if err != nil {
		slog.Error("Ошибка маршалинга JSON", "error", err)
		return 0, err
	}

	slog.Info("Отправка запроса на создание клиента", "client", posterC, "url", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(clientData))
	if err != nil {
		slog.Error("Ошибка создания HTTP-запроса", "error", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		slog.Error("Ошибка при выполнении запроса", "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Ошибка чтения ответа", "error", err)
		return 0, err
	}

	var response struct {
		Response int `json:"response"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Ошибка при разборе ответа API", "error", err)
		return 0, err
	}

	slog.Info("Клиент успешно создан", "clientID", response.Response)
	return response.Response, nil
}

// findClientByPhone ищет клиента в базе Poster по номеру телефона
func (p *PosterAPI) findClientByPhone(phone string) (int, error) {
	url := fmt.Sprintf("%sclients.getClients?token=%s", p.BaseURL, p.Token)

	slog.Info("Поиск клиента по номеру телефона", "phone", phone, "url", url)

	resp, err := p.Client.Get(url)
	if err != nil {
		slog.Error("Ошибка при выполнении запроса", "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Ошибка чтения ответа", "error", err)
		return 0, err
	}

	var response struct {
		Response []struct {
			ClientID string `json:"client_id"`
			Phone    string `json:"phone"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Ошибка при разборе ответа API", "error", err, "response_body", string(body))
		return 0, err
	}

	for _, client := range response.Response {
		if client.Phone == phone {
			clientID, err := strconv.Atoi(client.ClientID)
			if err != nil {
				slog.Error("Ошибка конвертации client_id в int", "client_id", client.ClientID, "error", err)
				return 0, err
			}

			slog.Info("Клиент найден", "phone", phone, "clientID", clientID)
			return clientID, nil
		}
	}

	slog.Warn("Клиент не найден", "phone", phone)
	return 0, nil // Клиент не найден
}

// toPosterClient преобразует доменную модель в API-структуру
func toPosterClient(c domain.Client) posterClient {
	return posterClient{
		ClientName:     c.Name,
		ClientSex:      c.Sex,
		ClientGroupsID: 2,
		Phone:          c.Phone,
		Birthday:       c.Birthday,
	}
}
