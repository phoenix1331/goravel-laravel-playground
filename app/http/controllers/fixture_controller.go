package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"

	"goravel/app/facades"
	"goravel/app/http/requests"
	"goravel/app/models"
	"goravel/app/results"
)

type FixtureController struct{}

func NewFixtureController() *FixtureController {
	return &FixtureController{}
}

// Index lists fixtures with both clubs eager loaded.
//
// Laravel:
//
//	Fixture::with(['homeClub', 'awayClub'])->orderBy('kick_off_at')->get();
//
// Goravel: one With per relation, named by struct field. A relation you did
// not load is nil, so there is no silent lazy loading (and no N+1).
func (r *FixtureController) Index(ctx http.Context) http.Response {
	var fixtures []models.Fixture
	err := facades.Orm().Query().
		With("HomeClub").
		With("AwayClub").
		Order("kick_off_at").
		Get(&fixtures)
	if err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": fixtures})
}

// RecordResult sets the final score. results.Record saves it and fires
// FixtureResultRecorded; the controller does not know who is listening.
//
// Laravel:
//
//	try {
//	    app(RecordResult::class)($fixture, $home, $away);
//	} catch (AlreadyRecorded $e) {
//	    return response()->json([...], 409);
//	}
//
// Goravel: no exceptions. Record returns an error, and errors.Is picks out
// the expected failure before anything else becomes a 500.
func (r *FixtureController) RecordResult(ctx http.Context) http.Response {
	var fixture models.Fixture
	err := facades.Orm().Query().FindOrFail(&fixture, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Fixture not found"})
	}
	if err != nil {
		return serverError(ctx, err)
	}
	if _, response := authorise(ctx, "record-result", map[string]any{"fixture": &fixture}); response != nil {
		return response
	}

	var request requests.RecordResultRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validationErrors != nil {
		return validationFailed(ctx, validationErrors)
	}

	err = results.Record(&fixture, request.HomeScore, request.AwayScore)
	if errors.Is(err, results.ErrAlreadyRecorded) {
		// An expected failure with its own status, like catching one
		// exception class and rendering it differently
		return ctx.Response().Json(http.StatusConflict, http.Json{"message": "This fixture already has a result."})
	}
	if err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": fixture})
}
