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

package infraction

import (
	"database/sql"
	"github.com/sniddunc/refractor/internal/mock"
	"github.com/sniddunc/refractor/internal/params"
	"github.com/sniddunc/refractor/internal/player"
	"github.com/sniddunc/refractor/internal/server"
	"github.com/sniddunc/refractor/internal/user"
	"github.com/sniddunc/refractor/pkg/config"
	"github.com/sniddunc/refractor/pkg/log"
	"github.com/sniddunc/refractor/pkg/perms"
	"github.com/sniddunc/refractor/refractor"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func Test_infractionService_CreateWarning(t *testing.T) {
	testLogger, _ := log.NewLogger(true, false)

	type fields struct {
		mockPlayers map[int64]*refractor.DBPlayer
		mockServers map[int64]*refractor.Server
	}
	type args struct {
		userID int64
		body   params.CreateWarningParams
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantInfraction *refractor.Infraction
		wantRes        *refractor.ServiceResponse
	}{
		{
			name: "infraction.createwarning.1",
			fields: fields{
				mockPlayers: map[int64]*refractor.DBPlayer{
					1: {
						PlayerID: 1,
					},
				},
				mockServers: map[int64]*refractor.Server{
					1: {
						ServerID: 1,
					},
				},
			},
			args: args{
				userID: 1,
				body: params.CreateWarningParams{
					PlayerID: 1,
					ServerID: 1,
					Reason:   "Test warning reason",
				},
			},
			wantInfraction: &refractor.Infraction{
				InfractionID: 1,
				PlayerID:     1,
				UserID:       1,
				ServerID:     1,
				Type:         refractor.INFRACTION_TYPE_WARNING,
				Reason:       "Test warning reason",
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction created",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlayerRepo := mock.NewMockPlayerRepository(tt.fields.mockPlayers)
			playerService := player.NewPlayerService(mockPlayerRepo, testLogger)
			mockServerRepo := mock.NewMockServerRepository(tt.fields.mockServers)
			serverService := server.NewServerService(mockServerRepo, nil, testLogger)
			mockInfractionRepo := mock.NewMockInfractionRepository(map[int64]*refractor.DBInfraction{})
			infractionService := NewInfractionService(mockInfractionRepo, playerService, serverService, nil, testLogger)

			warning, res := infractionService.CreateWarning(tt.args.userID, tt.args.body)

			assert.True(t, infractionsAreEqual(tt.wantInfraction, warning), "Infractions were not equal\nWant = %v\nGot  = %v", tt.wantInfraction, warning)
			assert.True(t, tt.wantRes.Equals(res), "tt.wantRes = %v and res = %v should be equal", tt.wantRes, res)
		})
	}
}

func Test_infractionService_DeleteInfraction(t *testing.T) {
	testLogger, _ := log.NewLogger(true, false)

	type fields struct {
		mockInfractions map[int64]*refractor.DBInfraction
	}
	type args struct {
		id   int64
		user params.UserMeta
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *refractor.ServiceResponse
	}{
		{
			name: "infraction.deleteinfraction.1",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id: 1,
				user: params.UserMeta{
					UserID:      1,
					Permissions: perms.DELETE_OWN_INFRACTIONS,
				},
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction deleted",
			},
		},
		{
			name: "infraction.deleteinfraction.2",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id: 1,
				user: params.UserMeta{
					UserID:      1,
					Permissions: perms.FULL_ACCESS,
				},
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction deleted",
			},
		},
		{
			name: "infraction.deleteinfraction.3",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id: 1,
				user: params.UserMeta{
					UserID:      2,
					Permissions: perms.DELETE_ANY_INFRACTION,
				},
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction deleted",
			},
		},
		{
			name: "infraction.deleteinfraction.4",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id: 1,
				user: params.UserMeta{
					UserID:      2,
					Permissions: perms.DELETE_OWN_INFRACTIONS,
				},
			},
			wantRes: &refractor.ServiceResponse{
				Success:    false,
				StatusCode: http.StatusBadRequest,
				Message:    config.MessageNoPermission,
			},
		},
		{
			name: "infraction.deleteinfraction.5",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id: 1,
				user: params.UserMeta{
					UserID:      2,
					Permissions: 0,
				},
			},
			wantRes: &refractor.ServiceResponse{
				Success:    false,
				StatusCode: http.StatusBadRequest,
				Message:    config.MessageNoPermission,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInfractionRepo := mock.NewMockInfractionRepository(tt.fields.mockInfractions)
			infractionService := NewInfractionService(mockInfractionRepo, nil, nil, nil, testLogger)

			res := infractionService.DeleteInfraction(tt.args.id, tt.args.user)

			assert.True(t, tt.wantRes.Equals(res), "tt.wantRes = %v and res = %v should be equal", tt.wantRes, res)
		})
	}
}

func Test_infractionService_UpdateInfraction(t *testing.T) {
	testLogger, _ := log.NewLogger(true, false)

	type fields struct {
		mockInfractions map[int64]*refractor.DBInfraction
	}
	type args struct {
		id       int64
		reason   string
		duration int
		userMeta *params.UserMeta
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantInfraction *refractor.Infraction
		wantRes        *refractor.ServiceResponse
	}{
		{
			name: "infraction.updateinfraction.1",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: false},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id:       1,
				reason:   "Updated reason test.1",
				duration: 1440,
				userMeta: &params.UserMeta{
					UserID:      1,
					Permissions: perms.SUPER_ADMIN,
				},
			},
			wantInfraction: &refractor.Infraction{
				InfractionID: 1,
				PlayerID:     1,
				UserID:       1,
				ServerID:     1,
				Type:         refractor.INFRACTION_TYPE_WARNING,
				Reason:       "Updated reason test.1",
				Duration:     0,
				Timestamp:    0,
				SystemAction: false,
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction updated",
			},
		},
		{
			name: "infraction.updateinfraction.2",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_KICK,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: false},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id:       1,
				reason:   "Updated reason test.2",
				duration: 1440,
				userMeta: &params.UserMeta{
					UserID:      1,
					Permissions: 0,
				},
			},
			wantInfraction: nil,
			wantRes: &refractor.ServiceResponse{
				Success:    false,
				StatusCode: http.StatusBadRequest,
				Message:    config.MessageNoPermission,
			},
		},
		{
			name: "infraction.updateinfraction.3",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_BAN,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 60, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id:       1,
				reason:   "Updated ban reason.3",
				duration: 1440,
				userMeta: &params.UserMeta{
					UserID:      1,
					Permissions: perms.EDIT_OWN_INFRACTIONS,
				},
			},
			wantInfraction: &refractor.Infraction{
				InfractionID: 1,
				PlayerID:     1,
				UserID:       1,
				ServerID:     1,
				Type:         refractor.INFRACTION_TYPE_BAN,
				Reason:       "Updated ban reason.3",
				Duration:     1440,
				Timestamp:    0,
				SystemAction: false,
				StaffName:    "",
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction updated",
			},
		},
		{
			name: "infraction.updateinfraction.4",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       2,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_BAN,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 60, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id:       1,
				reason:   "Updated ban reason.4",
				duration: 1440,
				userMeta: &params.UserMeta{
					UserID:      1,
					Permissions: perms.EDIT_OWN_INFRACTIONS,
				},
			},
			wantInfraction: nil,
			wantRes: &refractor.ServiceResponse{
				Success:    false,
				StatusCode: http.StatusBadRequest,
				Message:    config.MessageNoPermission,
			},
		},
		{
			name: "infraction.updateinfraction.5",
			fields: fields{
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       2,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_MUTE,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 60, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				id:       1,
				reason:   "Updated mute reason.5",
				duration: 120,
				userMeta: &params.UserMeta{
					UserID:      1,
					Permissions: perms.EDIT_ANY_INFRACTION,
				},
			},
			wantInfraction: &refractor.Infraction{
				InfractionID: 1,
				PlayerID:     1,
				UserID:       2,
				ServerID:     1,
				Type:         refractor.INFRACTION_TYPE_MUTE,
				Reason:       "Updated mute reason.5",
				Duration:     120,
				Timestamp:    0,
				SystemAction: false,
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Infraction updated",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInfractionRepo := mock.NewMockInfractionRepository(tt.fields.mockInfractions)
			infractionService := NewInfractionService(mockInfractionRepo, nil, nil, nil, testLogger)

			body := params.UpdateInfractionParams{
				Reason:   &tt.args.reason,
				Duration: &tt.args.duration,
				UserMeta: tt.args.userMeta,
			}

			updatedInfraction, res := infractionService.UpdateInfraction(tt.args.id, body)

			assert.True(t, tt.wantRes.Equals(res), "tt.wantRes = %v and res = %v should be equal", tt.wantRes, res)

			if tt.wantRes.Success {
				assert.True(t, infractionsAreEqual(tt.wantInfraction, updatedInfraction), "Infractions were not equal\nwant = %v\ngot  = %v", tt.wantInfraction, updatedInfraction)
			}
		})
	}
}

func Test_infractionService_GetPlayerInfractions(t *testing.T) {
	testLogger, _ := log.NewLogger(true, false)

	type fields struct {
		mockInfractions map[int64]*refractor.DBInfraction
		mockUsers       map[int64]*mock.MockUser
	}
	type args struct {
		infractionType string
		playerID       int64
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantInfractions []*refractor.Infraction
		wantRes         *refractor.ServiceResponse
	}{
		{
			name: "infraction.getplayerinfractions.1",
			fields: fields{
				mockUsers: map[int64]*mock.MockUser{
					1: {
						UnhashedPassword: "",
						User: &refractor.User{
							UserID:              1,
							Email:               "yudsgdus@ydwtgtsss.com",
							Username:            "infractionusername",
							Password:            "",
							Permissions:         0,
							Activated:           true,
							NeedsPasswordChange: false,
						},
					},
				},
				mockInfractions: map[int64]*refractor.DBInfraction{
					1: {
						InfractionID: 1,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("a", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: false},
						Timestamp:    0,
						SystemAction: false,
					},
					2: {
						InfractionID: 2,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("b", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: false},
						Timestamp:    0,
						SystemAction: false,
					},
					3: {
						InfractionID: 3,
						PlayerID:     2,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_WARNING,
						Reason:       sql.NullString{String: strings.Repeat("c", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 0, Valid: false},
						Timestamp:    0,
						SystemAction: false,
					},
					4: {
						InfractionID: 4,
						PlayerID:     1,
						UserID:       1,
						ServerID:     1,
						Type:         refractor.INFRACTION_TYPE_KICK,
						Reason:       sql.NullString{String: strings.Repeat("d", config.InfractionReasonMinLen), Valid: true},
						Duration:     sql.NullInt32{Int32: 1440, Valid: true},
						Timestamp:    0,
						SystemAction: false,
					},
				},
			},
			args: args{
				infractionType: refractor.INFRACTION_TYPE_WARNING,
				playerID:       1,
			},
			wantInfractions: []*refractor.Infraction{
				{
					InfractionID: 1,
					PlayerID:     1,
					UserID:       1,
					ServerID:     1,
					Type:         refractor.INFRACTION_TYPE_WARNING,
					Reason:       strings.Repeat("a", config.InfractionReasonMinLen),
					Duration:     0,
					Timestamp:    0,
					SystemAction: false,
					StaffName:    "infractionusername",
				},
				{
					InfractionID: 2,
					PlayerID:     1,
					UserID:       1,
					ServerID:     1,
					Type:         refractor.INFRACTION_TYPE_WARNING,
					Reason:       strings.Repeat("b", config.InfractionReasonMinLen),
					Duration:     0,
					Timestamp:    0,
					SystemAction: false,
					StaffName:    "infractionusername",
				},
			},
			wantRes: &refractor.ServiceResponse{
				Success:    true,
				StatusCode: http.StatusOK,
				Message:    "Fetched 2 infractions",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := mock.NewMockUserRepository(tt.fields.mockUsers)
			userService := user.NewUserService(mockUserRepo, testLogger)
			mockInfractionRepo := mock.NewMockInfractionRepository(tt.fields.mockInfractions)
			infractionService := NewInfractionService(mockInfractionRepo, nil, nil, userService, testLogger)

			foundInfractions, res := infractionService.GetPlayerInfractionsType(tt.args.infractionType, tt.args.playerID)

			assert.True(t, tt.wantRes.Equals(res), "tt.wantRes = %v and res = %v should be equal", tt.wantRes, res)

			if tt.wantRes.Success {
				assert.Equal(t, tt.wantInfractions, foundInfractions, "Infraction slices were not equal")
			}
		})
	}
}

// infractionsAreEqual compares the following fields to determine is two infractions are equal:
// InfractionID, PlayerID, ServerID, UserID, Type, Reason, SystemAction
func infractionsAreEqual(infraction1 *refractor.Infraction, infraction2 *refractor.Infraction) bool {
	if infraction1.InfractionID != infraction2.InfractionID {
		return false
	}

	if infraction1.PlayerID != infraction2.PlayerID {
		return false
	}

	if infraction1.ServerID != infraction2.ServerID {
		return false
	}

	if infraction1.UserID != infraction2.UserID {
		return false
	}

	if infraction1.Type != infraction2.Type {
		return false
	}

	if infraction1.Reason != infraction2.Reason {
		return false
	}

	if infraction1.Duration != infraction2.Duration {
		return false
	}

	if infraction1.SystemAction != infraction2.SystemAction {
		return false
	}

	return true
}
