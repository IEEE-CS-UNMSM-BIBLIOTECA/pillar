package types

type Document struct {
	Id                    int32    `json:"id"`
	Publisher_id          int32    `json:"publisher_id"`
	Mean_rating           *float64 `json:"mean_rating"`
	Language_id           int32    `json:"language_id"`
	Format_id             int32    `json:"format_id"`
	Publication_year      int32    `json:"publication_year"`
	Acquisition_date      string   `json:"acquisition_date"`
	Edition               int32    `json:"edition"`
	Total_pages           int32    `json:"total_pages"`
	External_lend_allowed bool     `json:"external_lend_allowed"`
	Base_price            float64  `json:"base_price"`
	Total_copies          int32    `json:"total_copies"`
	Available_copies      int32    `json:"available_copies"`
	Title                 string   `json:"title"`
	Isbn                  string   `json:"isbn"`
	Description           string   `json:"description"`
	Cover_url             *string  `json:"cover_url"`
}

type Author struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type AuthorDashboard struct {
	Id        int32   `json:"id"`
	Name      string  `json:"name"`
	BirthDate string  `json:"birthdate"`
	DeathDate *string `json:"deathdate"`
	Bio       string  `json:"bio"`
	GenderID  int32   `json:"gender_id"`
	CountryID int32   `json:"country_id"`
	ImageUrl  *string `json:"image_url"`
}

type Authors struct {
	Authors []Author `json:"authors"`
}

type Tag struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Tags struct {
	Tags []Tag `json:"tags"`
}

type Language struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Country struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Publisher struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Format struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Gender struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type Review struct {
	Id          int32          `json:"id"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	Rating      int32          `json:"rating"`
	Total_likes int32          `json:"total_likes"`
	Liked       bool           `json:"liked"`
	User        UserFromReview `json:"user"`
	Spoiler     bool           `json:"spoiler"`
}

type ReviewRequest struct {
	DocumentID int    `json:"document_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Rating     int    `json:"rating"`
	Spoiler    bool   `json:"spoiler"`
}

type UserFromReview struct {
	Id                  int32   `json:"id"`
	Name                string  `json:"user_name"`
	Profile_picture_url *string `json:"profile_picture_url"`
}

type SignUpRequest struct {
	Username            string  `json:"username"`
	Email               string  `json:"email"`
	Password            string  `json:"password"`
	Name                string  `json:"name"`
	BirthDate           string  `json:"birth_date"` // Expecting a date string
	Address             string  `json:"address"`
	MobilePhone         string  `json:"mobile_phone"`
	RoleID              int32   `json:"role_id"`
	GenderID            int32   `json:"gender_id"`
	Bio                 *string `json:"bio"`
	Profile_picture_url *string `json:"profile_picture_url"`
}

type PopularBook struct {
	BookID     int    `json:"book_id"`
	Title      string `json:"title"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	CoverURL   string `json:"cover_url"`
}

type RegisterLend struct {
	BookID int `json:"document_id"`
}

type ListAddDocument struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	HasDocument bool   `json:"has_document"`
}

type AddDocList struct {
	DocumentID int `json:"document_id"`
}

type RenameList struct {
	Title string `json:"title"`
}
