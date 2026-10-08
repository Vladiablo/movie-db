package repository

import (
	"movie-db/internal"

	"github.com/jackc/pgx/v5/pgtype"
)

type pgPerson struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	RoleDetails string `json:"role_details"`
}

func (p *pgPerson) copyToModel() *internal.Person {
	return &internal.Person{
		Id:          p.Id,
		Name:        p.Name,
		RoleDetails: p.RoleDetails,
	}
}

type pgCrewMembers struct {
	Count   int32       `json:"count"`
	Members []*pgPerson `json:"members"`
}

func (cm *pgCrewMembers) copyToModel() *internal.CrewMembers {
	modelMembers := make([]*internal.Person, len(cm.Members))
	for i := range modelMembers {
		modelMembers[i] = cm.Members[i].copyToModel()
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

func (c *pgCrew) copyToModel() *internal.Crew {
	return &internal.Crew{
		Actors:              c.Actors.copyToModel(),
		CastingDirectors:    c.CastingDirectors.copyToModel(),
		Cinematographers:    c.Cinematographers.copyToModel(),
		Composers:           c.Composers.copyToModel(),
		Directors:           c.Directors.copyToModel(),
		Editors:             c.Editors.copyToModel(),
		Producers:           c.Producers.copyToModel(),
		ProductionDesigners: c.ProductionDesigners.copyToModel(),
		Writers:             c.Writers.copyToModel(),
	}
}

type pgMovie struct {
	Id               int32       `db:"id"`
	Type             string      `db:"type"`
	Status           string      `db:"status"`
	Title            string      `db:"title"`
	Description      pgtype.Text `db:"description"`
	ShortDescription pgtype.Text `db:"short_description"`
	Tagline          pgtype.Text `db:"tagline"`
	AgeRating        pgtype.Text `db:"age_rating"`
	Runtime          pgtype.Int4 `db:"runtime"`
	YearFrom         pgtype.Int4 `db:"year_from"`
	ReleaseDate      pgtype.Date `db:"release_date"`
	LatestAirDate    pgtype.Date `db:"latest_air_date"`
	GenreIds         []int32     `db:"genre_ids"`
	KpId             pgtype.Int4 `db:"kp_id"`
	ImdbId           pgtype.Int4 `db:"imdb_id"`
	Crew             *pgCrew     `db:"crew"`
}

func (m *pgMovie) copyToModel() *internal.Movie {
	runtime := (*int32)(nil)
	if m.Runtime.Valid {
		runtime = new(m.Runtime.Int32)
	}

	yearFrom := (*int32)(nil)
	if m.YearFrom.Valid {
		yearFrom = new(m.YearFrom.Int32)
	}

	genreIds := make([]int32, len(m.GenreIds))
	for i := range genreIds {
		genreIds[i] = m.GenreIds[i]
	}

	return &internal.Movie{
		Id:               m.Id,
		Type:             m.Type,
		Status:           m.Status,
		Title:            m.Title,
		Description:      m.Description.String,
		ShortDescription: m.ShortDescription.String,
		Tagline:          m.Tagline.String,
		AgeRating:        m.AgeRating.String,
		Runtime:          runtime,
		YearFrom:         yearFrom,
		ReleaseDate:      internal.Date{Time: m.ReleaseDate.Time},
		LatestAirDate:    internal.Date{Time: m.LatestAirDate.Time},
		GenreIds:         genreIds,
		KpId:             m.KpId.Int32,
		ImdbId:           m.ImdbId.Int32,
		Crew:             m.Crew.copyToModel(),
	}
}
