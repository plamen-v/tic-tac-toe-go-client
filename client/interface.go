package client

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/plamen-v/tic-tac-toe-models/models"
)

type Client interface {
	Login(*models.LoginRequest) (*models.LoginResponse, error)
	GetRoom() (*models.RoomResponse, error)
	GetOpenRooms(int, int) (*models.RoomListResponse, error)
	CreateRoom(*models.CreateRoomRequest) (*models.CreateRoomResponse, error)
	JoinRoom(uuid.UUID) error
	LeaveRoom(uuid.UUID) error
	CreateGame(uuid.UUID) (bool, error)
	GetGame(uuid.UUID) (*models.GameResponse, error)
	MakeMove(uuid.UUID, int) error
	GetRanking(int, int) (*models.RankingResponse, error)
}

func NewClient(url string, port int) Client {
	url, _ = strings.CutSuffix(url, "/")

	return &clientImpl{
		url: fmt.Sprintf("%s:%d/api", url, port),
		headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "",
		},
	}
}
