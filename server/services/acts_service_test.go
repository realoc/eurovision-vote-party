package services_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sipgate/eurovision-vote-party/server/models"
	"github.com/sipgate/eurovision-vote-party/server/services"
)

// createTestActsFile creates a temporary JSON file containing the given acts
// wrapped in an {"acts": [...]} structure. It returns the path to the file.
func createTestActsFile(t *testing.T, acts []models.Act) string {
	t.Helper()

	wrapper := struct {
		Acts []models.Act `json:"acts"`
	}{
		Acts: acts,
	}

	data, err := json.Marshal(wrapper)
	require.NoError(t, err)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "acts.json")
	err = os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)

	return filePath
}

func TestNewActsService_LoadsValidFile(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
		{
			ID:           "act-2",
			Country:      "Italy",
			Artist:       "Maneskin",
			Song:         "Zitti e buoni",
			RunningOrder: 2,
			EventType:    models.EventSemifinal1,
		},
	}

	filePath := createTestActsFile(t, acts)

	svc, err := services.NewActsService(filePath)

	require.NoError(t, err)
	require.NotNil(t, svc)

	result, err := svc.ListActs("")
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestNewActsService_ErrorsOnMissingFile(t *testing.T) {
	svc, err := services.NewActsService("/nonexistent/path/acts.json")

	assert.Error(t, err)
	assert.Nil(t, svc)
}

func TestNewActsService_ErrorsOnInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "acts.json")
	err := os.WriteFile(filePath, []byte("not valid json{{{"), 0644)
	require.NoError(t, err)

	svc, err := services.NewActsService(filePath)

	assert.Error(t, err)
	assert.Nil(t, svc)
}

func TestActsService_ListActs_ReturnsAllActsWithEmptyFilter(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
		{
			ID:           "act-2",
			Country:      "Italy",
			Artist:       "Maneskin",
			Song:         "Zitti e buoni",
			RunningOrder: 2,
			EventType:    models.EventSemifinal1,
		},
		{
			ID:           "act-3",
			Country:      "France",
			Artist:       "Barbara Pravi",
			Song:         "Voila",
			RunningOrder: 3,
			EventType:    models.EventSemifinal2,
		},
	}

	filePath := createTestActsFile(t, acts)
	svc, err := services.NewActsService(filePath)
	require.NoError(t, err)

	result, err := svc.ListActs("")

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, acts, result)
}

func TestActsService_ListActs_FiltersBySemifinal1(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
		{
			ID:           "act-2",
			Country:      "Italy",
			Artist:       "Maneskin",
			Song:         "Zitti e buoni",
			RunningOrder: 2,
			EventType:    models.EventSemifinal1,
		},
		{
			ID:           "act-3",
			Country:      "France",
			Artist:       "Barbara Pravi",
			Song:         "Voila",
			RunningOrder: 3,
			EventType:    models.EventSemifinal2,
		},
	}

	filePath := createTestActsFile(t, acts)
	svc, err := services.NewActsService(filePath)
	require.NoError(t, err)

	result, err := svc.ListActs("semifinal1")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "act-2", result[0].ID)
	assert.Equal(t, models.EventSemifinal1, result[0].EventType)
}

func TestActsService_ListActs_FiltersByGrandfinal(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
		{
			ID:           "act-2",
			Country:      "Italy",
			Artist:       "Maneskin",
			Song:         "Zitti e buoni",
			RunningOrder: 2,
			EventType:    models.EventSemifinal1,
		},
		{
			ID:           "act-3",
			Country:      "France",
			Artist:       "Barbara Pravi",
			Song:         "Voila",
			RunningOrder: 3,
			EventType:    models.EventSemifinal2,
		},
	}

	filePath := createTestActsFile(t, acts)
	svc, err := services.NewActsService(filePath)
	require.NoError(t, err)

	result, err := svc.ListActs("grandfinal")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "act-1", result[0].ID)
	assert.Equal(t, models.EventGrandFinal, result[0].EventType)
}

func TestActsService_ListActs_ReturnsErrorForInvalidEventType(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
	}

	filePath := createTestActsFile(t, acts)
	svc, err := services.NewActsService(filePath)
	require.NoError(t, err)

	result, err := svc.ListActs("invalid")

	assert.ErrorIs(t, err, services.ErrInvalidEventType)
	assert.Nil(t, result)
}

func TestActsService_ListActs_ReturnsEmptySliceForValidEventWithNoMatches(t *testing.T) {
	acts := []models.Act{
		{
			ID:           "act-1",
			Country:      "Sweden",
			Artist:       "ABBA",
			Song:         "Waterloo",
			RunningOrder: 1,
			EventType:    models.EventGrandFinal,
		},
		{
			ID:           "act-2",
			Country:      "Italy",
			Artist:       "Maneskin",
			Song:         "Zitti e buoni",
			RunningOrder: 2,
			EventType:    models.EventSemifinal1,
		},
	}

	filePath := createTestActsFile(t, acts)
	svc, err := services.NewActsService(filePath)
	require.NoError(t, err)

	result, err := svc.ListActs("semifinal2")

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
}
