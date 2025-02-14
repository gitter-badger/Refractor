/*
This file is part of Refractor.

Refractor is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package refractor

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/sniddunc/refractor/internal/params"
)

const (
	INFRACTION_TYPE_WARNING = "WARNING"
	INFRACTION_TYPE_MUTE    = "MUTE"
	INFRACTION_TYPE_KICK    = "KICK"
	INFRACTION_TYPE_BAN     = "BAN"
)

var InfractionTypes = []string{INFRACTION_TYPE_WARNING, INFRACTION_TYPE_MUTE, INFRACTION_TYPE_KICK, INFRACTION_TYPE_BAN}

type Infraction struct {
	InfractionID int64  `json:"id"`
	PlayerID     int64  `json:"playerId"`
	UserID       int64  `json:"userId"`
	ServerID     int64  `json:"serverId"`
	Type         string `json:"type"`
	Reason       string `json:"reason"`
	Duration     int    `json:"duration"`
	Timestamp    int64  `json:"timestamp"`
	SystemAction bool   `json:"systemAction"`
	StaffName    string `json:"staffName"`  // not a database field
	PlayerName   string `json:"playerName"` // not a database field
}

type DBInfraction struct {
	InfractionID int64
	PlayerID     int64
	UserID       int64
	ServerID     int64
	Type         string
	Reason       sql.NullString
	Duration     sql.NullInt32
	Timestamp    int64
	SystemAction bool
}

// Infraction builds a Infraction instance from the DBInstance it was called upon.
func (dbi *DBInfraction) Infraction() *Infraction {
	return &Infraction{
		InfractionID: dbi.InfractionID,
		PlayerID:     dbi.PlayerID,
		UserID:       dbi.UserID,
		ServerID:     dbi.ServerID,
		Reason:       dbi.Reason.String,
		Duration:     int(dbi.Duration.Int32),
		Type:         dbi.Type,
		Timestamp:    dbi.Timestamp,
		SystemAction: dbi.SystemAction,
	}
}

type InfractionRepository interface {
	Create(infraction *DBInfraction) (*Infraction, error)
	FindByID(id int64) (*Infraction, error)
	Exists(args FindArgs) (bool, error)
	FindOne(args FindArgs) (*Infraction, error)
	FindMany(args FindArgs) ([]*Infraction, error)
	FindManyByPlayerID(playerID int64) ([]*Infraction, error)
	FindAll() ([]*Infraction, error)
	Update(id int64, args UpdateArgs) (*Infraction, error)
	Delete(id int64) error
	Search(args FindArgs, limit int, offset int) (int, []*Infraction, error)
	GetRecent(count int) ([]*Infraction, error)
}

type InfractionService interface {
	CreateWarning(userID int64, body params.CreateWarningParams) (*Infraction, *ServiceResponse)
	CreateMute(userID int64, body params.CreateMuteParams) (*Infraction, *ServiceResponse)
	CreateKick(userID int64, body params.CreateKickParams) (*Infraction, *ServiceResponse)
	CreateBan(userID int64, body params.CreateBanParams) (*Infraction, *ServiceResponse)
	DeleteInfraction(id int64, user params.UserMeta) *ServiceResponse
	UpdateInfraction(id int64, body params.UpdateInfractionParams) (*Infraction, *ServiceResponse)
	GetPlayerInfractionsType(infractionType string, playerID int64) ([]*Infraction, *ServiceResponse)
	GetPlayerInfractions(playerID int64) ([]*Infraction, *ServiceResponse)
	GetRecentInfractions(count int) ([]*Infraction, *ServiceResponse)
}

type InfractionHandler interface {
	CreateWarning(c echo.Context) error
	CreateMute(c echo.Context) error
	CreateKick(c echo.Context) error
	CreateBan(c echo.Context) error
	DeleteInfraction(c echo.Context) error
	UpdateInfraction(c echo.Context) error
	GetPlayerInfractions(infractionType string) echo.HandlerFunc
	GetRecentInfractions(c echo.Context) error
}
