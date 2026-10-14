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

func (cm *pgCrewMembers) copyToModelOrNil() *internal.CrewMembers {
	if cm == nil {
		return nil
	}

	return cm.copyToModel()
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
		Actors:              c.Actors.copyToModelOrNil(),
		CastingDirectors:    c.CastingDirectors.copyToModelOrNil(),
		Cinematographers:    c.Cinematographers.copyToModelOrNil(),
		Composers:           c.Composers.copyToModelOrNil(),
		Directors:           c.Directors.copyToModelOrNil(),
		Editors:             c.Editors.copyToModelOrNil(),
		Producers:           c.Producers.copyToModelOrNil(),
		ProductionDesigners: c.ProductionDesigners.copyToModelOrNil(),
		Writers:             c.Writers.copyToModelOrNil(),
	}
}

func (c *pgCrew) copyToModelOrNil() *internal.Crew {
	if c == nil {
		return nil
	}

	return c.copyToModel()
}

type pgRating struct {
	Votes int64   `json:"votes"`
	Value float32 `json:"value"`
}

func (r *pgRating) copyToModel() internal.Rating {
	return internal.Rating{
		Votes: r.Votes,
		Value: r.Value,
	}
}

type pgRatings struct {
	Kp   pgRating `json:"kp"`
	Imdb pgRating `json:"imdb"`
}

func (r *pgRatings) copyToModel() *internal.Ratings {
	return &internal.Ratings{
		Kp:   r.Kp.copyToModel(),
		Imdb: r.Imdb.copyToModel(),
	}
}

func (r *pgRatings) copyToModelOrNil() *internal.Ratings {
	if r == nil {
		return nil
	}

	return r.copyToModel()
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
	Ratings          *pgRatings  `db:"ratings"`
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
		Ratings:          m.Ratings.copyToModelOrNil(),
		Crew:             m.Crew.copyToModelOrNil(),
	}
}
