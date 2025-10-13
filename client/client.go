package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

type clientImpl struct {
	url     string
	headers map[string]string
}

func (c *clientImpl) Login(request *models.LoginRequest) (*models.LoginResponse, error) {
	action := fmt.Sprintf("%s/login", c.url)
	var response models.LoginResponse
	_, err := webRequest(http.MethodPost, action, nil, request, &response)
	c.headers["Authorization"] = fmt.Sprintf("Bearer %s", response.Token)
	return &response, err
}

func (c *clientImpl) GetRoom() (*models.RoomResponse, error) {
	action := fmt.Sprintf("%s/room", c.url)
	var response models.RoomResponse
	_, err := webRequest(string(http.MethodGet), action, c.headers, (*struct{})(nil), &response)
	return &response, err
}

func (c *clientImpl) GetOpenRooms(page int, pageSize int) (*models.RoomListResponse, error) {
	action := fmt.Sprintf("%s/rooms?page=%d&pageSize=%d", c.url, page, pageSize)
	var response models.RoomListResponse
	_, err := webRequest(string(http.MethodGet), action, c.headers, (*struct{})(nil), &response)
	return &response, err
}

func (c *clientImpl) CreateRoom(request *models.CreateRoomRequest) (*models.CreateRoomResponse, error) {
	action := fmt.Sprintf("%s/rooms", c.url)
	var response models.CreateRoomResponse
	_, err := webRequest(string(http.MethodPost), action, c.headers, request, &response)
	return &response, err
}

func (c *clientImpl) JoinRoom(id uuid.UUID) error {
	action := fmt.Sprintf("%s/rooms/%s/player", c.url, id)
	_, err := webRequest(string(http.MethodPost), action, c.headers, (*struct{})(nil), (*struct{})(nil))
	return err
}

func (c *clientImpl) LeaveRoom(id uuid.UUID) error {
	action := fmt.Sprintf("%s/rooms/%s/player", c.url, id)
	_, err := webRequest(string(http.MethodDelete), action, c.headers, (*struct{})(nil), (*struct{})(nil))
	return err
}

func (c *clientImpl) CreateGame(id uuid.UUID) (bool, error) {
	action := fmt.Sprintf("%s/rooms/%s/game", c.url, id)
	statusDode, err := webRequest(string(http.MethodPost), action, c.headers, (*struct{})(nil), (*struct{})(nil))
	if err != nil {
		return false, err
	}
	created := true
	if statusDode != http.StatusCreated {
		created = false
	}
	return created, nil
}

func (c *clientImpl) GetGame(id uuid.UUID) (*models.GameResponse, error) {
	action := fmt.Sprintf("%s/rooms/%s/game", c.url, id)
	var response models.GameResponse
	_, err := webRequest(string(http.MethodGet), action, c.headers, (*struct{})(nil), &response)
	return &response, err
}

func (c *clientImpl) MakeMove(id uuid.UUID, position int) error {
	action := fmt.Sprintf("%s/rooms/%s/game/board/%d", c.url, id, position)
	_, err := webRequest(string(http.MethodPost), action, c.headers, (*struct{})(nil), (*struct{})(nil))
	return err
}

func (c *clientImpl) GetRanking(page int, pageSize int) (*models.RankingResponse, error) {
	action := fmt.Sprintf("%s/ranking?page=%d&pageSize=%d", c.url, page, pageSize)
	var response models.RankingResponse
	_, err := webRequest(string(http.MethodGet), action, c.headers, (*struct{})(nil), &response)
	return &response, err
}

func webRequest[P, R any](method string, url string, headers map[string]string, payload *P, data *R) (int, error) {
	var requestBody io.Reader = nil
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return 0, NewClientError(ClientErrorCode, err.Error())
		}
		requestBody = bytes.NewBuffer(jsonData)
	}

	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return 0, NewClientError(ClientErrorCode, err.Error())
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, NewClientError(ClientErrorCode, err.Error())
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, NewClientError(ClientErrorCode, err.Error())
	}
	if len(responseBody) > 0 {
		switch response.StatusCode {
		case http.StatusOK, http.StatusCreated, http.StatusAccepted:
			err = json.Unmarshal(responseBody, data)
			if err != nil {
				return 0, NewClientError(ClientErrorCode, err.Error())
			}
		default:
			var errorResponse models.ErrorResponse
			err = json.Unmarshal(responseBody, &errorResponse)
			if err != nil {
				return 0, NewClientError(ClientErrorCode, err.Error())
			}
			return 0, NewClientErrorf(models.ErrorCode(errorResponse.Code), "%s: %s", errorResponse.Code, errorResponse.Message)
		}
	}
	return response.StatusCode, nil
}
