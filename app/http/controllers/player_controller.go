package controllers

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/http/requests"
	"goravel/app/models"
)

type PlayerController struct{}

func NewPlayerController() *PlayerController {
	return &PlayerController{}
}

// Index returns a club with its squad.
//
// Laravel:
//
//	$club->load(['players' => fn ($q) => $q->orderBy('squad_number')]);
//
// Goravel: Load takes the relation's struct field name ("Players", not
// "players") and an optional function to constrain the query.
func (r *PlayerController) Index(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}

	// Load fills club.Players with a second query. The extra argument
	// constrains it, like load(['players' => fn ($q) => $q->orderBy(...)])
	if err := facades.Orm().Query().Load(club, "Players", func(query orm.Query) orm.Query {
		return query.Order("squad_number")
	}); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": club})
}

// Store adds a player to the club's squad.
//
// Laravel:
//
//	Gate::authorize('update', $club);
//	$player = $club->players()->create($request->validated());
//	// squad_number unique per club: Rule::unique('players')->where('club_id', $club->id)
//
// Goravel: no create-through-relation, so ClubID is set on the struct. The
// unique rule has no where clause, so the per-club check is a query here.
func (r *PlayerController) Store(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}
	if _, response := authorise(ctx, "manage-players", map[string]any{"club": club}); response != nil {
		return response
	}

	var request requests.StorePlayerRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validationErrors != nil {
		return validationFailed(ctx, validationErrors)
	}

	// The shirt number must be free within this club. Goravel's unique rule
	// cannot add the extra club_id condition, so check it here.
	taken, err := facades.Orm().Query().Model(&models.Player{}).
		Where("club_id", club.ID).Where("squad_number", request.SquadNumber).Exists()
	if err != nil {
		return serverError(ctx, err)
	}
	if taken {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "The given data was invalid.",
			"errors":  http.Json{"squad_number": http.Json{"unique": "This squad number is already taken."}},
		})
	}

	// Setting ClubID is how you "associate" in Go: there is no
	// $club->players()->create(); you set the foreign key yourself
	player := models.Player{
		ClubID:      club.ID,
		Name:        request.Name,
		Position:    request.Position,
		SquadNumber: request.SquadNumber,
	}
	if err := facades.Orm().Query().Create(&player); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{"data": player})
}
