package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validVote returns a Vote with all 10 valid Eurovision point values and unique act IDs.
func validVote() Vote {
	return Vote{
		ID:      "vote-1",
		GuestID: "guest-1",
		PartyID: "party-1",
		Votes: map[int]string{
			12: "act-1",
			10: "act-2",
			8:  "act-3",
			7:  "act-4",
			6:  "act-5",
			5:  "act-6",
			4:  "act-7",
			3:  "act-8",
			2:  "act-9",
			1:  "act-10",
		},
		CreatedAt: time.Now(),
	}
}

func TestVoteValidate(t *testing.T) {
	t.Run("valid vote", func(t *testing.T) {
		v := validVote()
		err := v.Validate()
		require.NoError(t, err)
	})

	t.Run("missing guest id", func(t *testing.T) {
		v := validVote()
		v.GuestID = ""
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "guest id is required")
	})

	t.Run("missing party id", func(t *testing.T) {
		v := validVote()
		v.PartyID = ""
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "party id is required")
	})

	t.Run("missing created at", func(t *testing.T) {
		v := validVote()
		v.CreatedAt = time.Time{}
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "created at timestamp is required")
	})

	t.Run("too few votes", func(t *testing.T) {
		v := validVote()
		v.Votes = map[int]string{
			12: "act-1",
			10: "act-2",
			8:  "act-3",
		}
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly 10 votes required, got 3")
	})

	t.Run("too many votes", func(t *testing.T) {
		v := validVote()
		v.Votes[9] = "act-11"
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly 10 votes required, got 11")
	})

	t.Run("invalid point value", func(t *testing.T) {
		v := validVote()
		delete(v.Votes, 8)
		v.Votes[9] = "act-3"
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing vote for point value 8")
	})

	t.Run("duplicate act ids", func(t *testing.T) {
		v := validVote()
		v.Votes[1] = "act-1" // same as 12-point entry
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate act id")
	})

	t.Run("empty act id for valid point value", func(t *testing.T) {
		v := validVote()
		v.Votes[7] = ""
		err := v.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "act id is required for points 7")
	})
}

func TestVoteResultValidate(t *testing.T) {
	base := VoteResult{
		ActID:       "act-1",
		Country:     "Sweden",
		Artist:      "Loreen",
		Song:        "Tattoo",
		TotalPoints: 583,
		Rank:        1,
	}

	t.Run("valid vote result", func(t *testing.T) {
		err := base.Validate()
		require.NoError(t, err)
	})

	tests := map[string]struct {
		result  VoteResult
		wantErr string
	}{
		"missing act id": {
			result: VoteResult{
				Country:     base.Country,
				Artist:      base.Artist,
				Song:        base.Song,
				TotalPoints: base.TotalPoints,
				Rank:        base.Rank,
			},
			wantErr: "act id is required",
		},
		"missing country": {
			result: VoteResult{
				ActID:       base.ActID,
				Artist:      base.Artist,
				Song:        base.Song,
				TotalPoints: base.TotalPoints,
				Rank:        base.Rank,
			},
			wantErr: "country is required",
		},
		"missing artist": {
			result: VoteResult{
				ActID:       base.ActID,
				Country:     base.Country,
				Song:        base.Song,
				TotalPoints: base.TotalPoints,
				Rank:        base.Rank,
			},
			wantErr: "artist is required",
		},
		"missing song": {
			result: VoteResult{
				ActID:       base.ActID,
				Country:     base.Country,
				Artist:      base.Artist,
				TotalPoints: base.TotalPoints,
				Rank:        base.Rank,
			},
			wantErr: "song is required",
		},
		"negative total points": {
			result: VoteResult{
				ActID:       base.ActID,
				Country:     base.Country,
				Artist:      base.Artist,
				Song:        base.Song,
				TotalPoints: -1,
				Rank:        base.Rank,
			},
			wantErr: "total points cannot be negative",
		},
		"negative rank": {
			result: VoteResult{
				ActID:       base.ActID,
				Country:     base.Country,
				Artist:      base.Artist,
				Song:        base.Song,
				TotalPoints: base.TotalPoints,
				Rank:        -1,
			},
			wantErr: "rank cannot be negative",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.result.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
