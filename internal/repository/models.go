package repository

import (
	"movie-db/internal"
	"slices"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
)

type pgPerson struct {
	Id          int32         `json:"id"`
	Name        string        `json:"name"`
	RoleDetails zeronull.Text `json:"role_details"`
}

func (p *pgPerson) toModel() *internal.Person {
	return &internal.Person{
		Id:          p.Id,
		Name:        p.Name,
		RoleDetails: string(p.RoleDetails),
	}
}

type pgCrewMembers struct {
	Count   int32       `json:"count"`
	Members []*pgPerson `json:"members"`
}

func (cm *pgCrewMembers) toModel() *internal.CrewMembers {
	if cm == nil {
		return nil
	}

	modelMembers := make([]*internal.Person, len(cm.Members))
	for i := range modelMembers {
		modelMembers[i] = cm.Members[i].toModel()
	}

	return &internal.CrewMembers{
		Count:   cm.Count,
		Members: modelMembers,
	}
}

type pgCrew struct {
	Actors              *pgCrewMembers `json:"actor"`
	CastingDirectors    *pgCrewMembers `json:"casting_director"`
	Cinematographers    *pgCrewMembers `json:"cinematographer"`
	Composers           *pgCrewMembers `json:"composer"`
	Directors           *pgCrewMembers `json:"director"`
	Editors             *pgCrewMembers `json:"editor"`
	Producers           *pgCrewMembers `json:"producer"`
	ProductionDesigners *pgCrewMembers `json:"production_designer"`
	Writers             *pgCrewMembers `json:"writer"`
}

func (c *pgCrew) toModel() *internal.Crew {
	if c == nil {
		return nil
	}

	return &internal.Crew{
		Actors:              c.Actors.toModel(),
		CastingDirectors:    c.CastingDirectors.toModel(),
		Cinematographers:    c.Cinematographers.toModel(),
		Composers:           c.Composers.toModel(),
		Directors:           c.Directors.toModel(),
		Editors:             c.Editors.toModel(),
		Producers:           c.Producers.toModel(),
		ProductionDesigners: c.ProductionDesigners.toModel(),
		Writers:             c.Writers.toModel(),
	}
}

type pgMovieBase struct {
	Id               int32         `db:"id"`
	Type             string        `db:"type"`
	Status           string        `db:"status"`
	Title            string        `db:"title"`
	Description      zeronull.Text `db:"description"`
	ShortDescription zeronull.Text `db:"short_description"`
	Tagline          zeronull.Text `db:"tagline"`
	AgeRating        zeronull.Text `db:"age_rating"`
	Runtime          zeronull.Int4 `db:"runtime"`
	YearFrom         zeronull.Int4 `db:"year_from"`
	ReleaseDate      pgtype.Date   `db:"release_date"`
	LatestAirDate    pgtype.Date   `db:"latest_air_date"`
	GenreIds         []int32       `db:"genre_ids"`
	KpId             zeronull.Int4 `db:"kp_id"`
	ImdbId           zeronull.Int4 `db:"imdb_id"`

	KpVotes   zeronull.Int8 `db:"kp_votes"`
	KpValue   pgtype.Float4 `db:"kp_value"`
	ImdbVotes zeronull.Int8 `db:"imdb_votes"`
	ImdbValue pgtype.Float4 `db:"imdb_value"`
}

func (ms *pgMovieBase) toModel() *internal.MovieBase {
	return &internal.MovieBase{
		Id:               ms.Id,
		Type:             ms.Type,
		Status:           ms.Status,
		Title:            ms.Title,
		Description:      string(ms.Description),
		ShortDescription: string(ms.ShortDescription),
		Tagline:          string(ms.Tagline),
		AgeRating:        string(ms.AgeRating),
		Runtime:          int32(ms.Runtime),
		YearFrom:         int32(ms.YearFrom),
		ReleaseDate:      internal.Date{Time: ms.ReleaseDate.Time},
		LatestAirDate:    internal.Date{Time: ms.LatestAirDate.Time},
		GenreIds:         slices.Clone(ms.GenreIds),
		KpId:             int32(ms.KpId),
		ImdbId:           int32(ms.ImdbId),
		Ratings: &internal.Ratings{
			Kp: internal.Rating{
				Value: ms.KpValue.Float32,
				Votes: int64(ms.KpVotes),
			},
			Imdb: internal.Rating{
				Value: ms.ImdbValue.Float32,
				Votes: int64(ms.ImdbVotes),
			},
		},
	}
}

type pgMovie struct {
	pgMovieBase

	Crew *pgCrew `db:"crew"`
}

func (m *pgMovie) toModel() *internal.Movie {
	return &internal.Movie{
		MovieBase: m.pgMovieBase.toModel(),
		Crew:      m.Crew.toModel(),
	}
}
