package internal

import "time"

type Date struct {
	time.Time
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format("\"2006-01-02\"")), nil
}

// Person model info
//
//	@Description	Information about crew member and their role
type Person struct {
	Id          int32  `json:"id" example:"15694579"`                  // Person ID
	Name        string `json:"name" example:"Elijah Wood"`             // Person name
	RoleDetails string `json:"role_details,omitempty" example:"Frodo"` // Additional role details
}

// CrewMembers model info
//
//	@Descriptions	Crew members of given role
type CrewMembers struct {
	Count   int32     `json:"count" example:"10"` // Overall count of crew members
	Members []*Person `json:"members"`            // Array with no more than five most important crew members
}

// Crew model info
//
//	@Description	Movie crew members separated by their role
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

// Rating model info
//
//	@Description	Movie rating at given film aggregator
type Rating struct {
	Votes int64   `json:"votes" example:"2240637"` // Total number of votes
	Value float32 `json:"value" example:"8.9"`     // Average rating
}

func (r *Rating) IsZero() bool {
	return r.Votes == 0 && r.Value == 0.0
}

// Ratings model info
//
//	@Description	Movie ratings at supported film aggregators
type Ratings struct {
	Kp   Rating `json:"kp,omitzero"`   // Movie rating at Kinopoisk
	Imdb Rating `json:"imdb,omitzero"` // Movie rating at IMDB
}

func (r *Ratings) IsZero() bool {
	return r.Kp.IsZero() && r.Imdb.IsZero()
}

// Movie model info
//
//	@Description	Complete information about movie including most important crew members
type Movie struct {
	Id               int32    `json:"id" example:"12163543"`                                                        // Movie ID
	Type             string   `json:"type" example:"movie"`                                                         // Movie type (movie or series)
	Status           string   `json:"status" example:"post_production"`                                             // Movie production status
	Title            string   `json:"title" example:"The Lord of the Rings: The Fellowship of the Ring"`            // Movie title
	Description      string   `json:"description,omitempty"`                                                        // Movie description
	ShortDescription string   `json:"short_description,omitempty"`                                                  // Movie short description
	Tagline          string   `json:"tagline,omitempty"`                                                            // Movie tagline
	AgeRating        string   `json:"age_rating,omitempty"`                                                         // Movie age rating
	Runtime          *int32   `json:"runtime,omitempty" example:"178"`                                              // Movie runtime in minutes
	YearFrom         *int32   `json:"year_from,omitempty" example:"2001"`                                           // Year the movie was first shown
	ReleaseDate      Date     `json:"release_date,omitzero" swaggertype:"primitive,string" example:"2001-12-19"`    // Movie release date
	LatestAirDate    Date     `json:"latest_air_date,omitzero" swaggertype:"primitive,string" example:"2001-12-19"` // Latest date on which series have been translated
	GenreIds         []int32  `json:"genre_ids,omitempty" example:"12,15,16"`                                       // IDs of genres the movie belongs to
	KpId             int32    `json:"kp_id,omitempty" example:"328"`                                                // Movie ID in Kinopoisk database
	ImdbId           int32    `json:"imdb_id,omitempty" example:"120737"`                                           // Movie ID in IMDB database
	Ratings          *Ratings `json:"ratings,omitzero"`                                                             // Movie ratings
	Crew             *Crew    `json:"crew,omitempty"`                                                               // Movie crew
}

// MovieSuggestion model info
//
//	@Description	Summary information about movie
type MovieSuggestion struct {
	Id               int32    `json:"id" example:"12163543"`                                                        // Movie ID
	Type             string   `json:"type" example:"movie"`                                                         // Movie type (movie or series)
	Status           string   `json:"status" example:"post_production"`                                             // Movie production status
	Title            string   `json:"title" example:"The Lord of the Rings: The Fellowship of the Ring"`            // Movie title
	Description      string   `json:"description,omitempty"`                                                        // Movie description
	ShortDescription string   `json:"short_description,omitempty"`                                                  // Movie short description
	Tagline          string   `json:"tagline,omitempty"`                                                            // Movie tagline
	AgeRating        string   `json:"age_rating,omitempty"`                                                         // Movie age rating
	Runtime          *int32   `json:"runtime,omitempty" example:"178"`                                              // Movie runtime in minutes
	YearFrom         *int32   `json:"year_from,omitempty" example:"2001"`                                           // Year the movie was first shown
	ReleaseDate      Date     `json:"release_date,omitzero" swaggertype:"primitive,string" example:"2001-12-19"`    // Movie release date
	LatestAirDate    Date     `json:"latest_air_date,omitzero" swaggertype:"primitive,string" example:"2001-12-19"` // Latest date on which series have been translated
	GenreIds         []int32  `json:"genre_ids,omitempty" example:"12,15,16"`                                       // IDs of genres the movie belongs to
	KpId             int32    `json:"kp_id,omitempty" example:"328"`                                                // Movie ID in Kinopoisk database
	ImdbId           int32    `json:"imdb_id,omitempty" example:"120737"`                                           // Movie ID in IMDB database
	Ratings          *Ratings `json:"ratings,omitzero"`                                                             // Movie ratings
}
