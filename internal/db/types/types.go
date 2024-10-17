package types

import (
	"time"
)

type Document struct {
    Publisher_id int32 `db:"publisher_id"`
    Avg_rating float64 `db:"avg_rating"`
    Language_id int32 `db:"language_id"`
    Format_id int32 `db:"format_id"`
    Id int32 `db:"id"`
    Publication_date time.Time `db:"publication_date"`
    Acquisition_date time.Time `db:"acquisition_date"`
    Edition int32 `db:"edition"`
    Total_pages int32 `db:"total_pages"`
    External_lend_allowed bool `db:"external_lend_allowed"`
    Base_price float64 `db:"base_price"`
    Total_copies int32 `db:"total_copies"`
    Available_copies int32 `db:"available_copies"`
    Title string `db:"title"`
    Isbn string `db:"isbn"`
    Description string `db:"description"`
    Cover_url string `db:"cover_url"`
}

