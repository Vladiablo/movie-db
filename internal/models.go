package internal

import "time"

type Date struct {
	time.Time
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format("\"2006-01-02\"")), nil
}

type Person struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	RoleDetails string `json:"role_details,omitempty"`
}

type CrewMembers struct {
	Count   int32     `json:"count"`
	Members []*Person `json:"members"`
}

type Crew struct {
	Actors              *CrewMembers `json:"actor,omitempty"`
	CastingDirectors    *CrewMembers `json:"casting_director,omitempty"`
	Cinematographers    *CrewMembers `json:"cinematographer,omitempty"`
	Composers           *CrewMembers `json:"composer,omitempty"`
	Directors           *CrewMembers `json:"director,omitempty"`
	Editors             *CrewMembers `json:"editor,omitempty"`
	Producers           *CrewMembers `json:"producer,omitempty"`
	ProductionDesigners *CrewMembers `json:"production_designer,omitempty"`
	Writers             *CrewMembers `json:"writer,omitempty"`
}

type Movie struct {
	Id               int32   `json:"id"`
	Type             string  `json:"type"`
	Status           string  `json:"status"`
	Title            string  `json:"title"`
	Description      string  `json:"description,omitempty"`
	ShortDescription string  `json:"short_description,omitempty"`
	Tagline          string  `json:"tagline,omitempty"`
	AgeRating        string  `json:"age_rating,omitempty"`
	Runtime          *int32  `json:"runtime,omitempty"`
	YearFrom         *int32  `json:"year_from,omitempty"`
	ReleaseDate      Date    `json:"release_date,omitzero"`
	LatestAirDate    Date    `json:"latest_air_date,omitzero"`
	GenreIds         []int32 `json:"genre_ids,omitempty"`
	KpId             int32   `json:"kp_id,omitempty"`
	ImdbId           int32   `json:"imdb_id,omitempty"`
	Crew             *Crew   `json:"crew,omitempty"`
}
