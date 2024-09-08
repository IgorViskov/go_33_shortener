package queries

const (
	SelectRecord               = "SELECT * FROM main.\"ShortenUrls\" WHERE \"ID\" = $1"
	InsertRecord               = "INSERT INTO main.\"ShortenUrls\" VALUES ($1, $2, $3)"
	UpdateRecord               = ""
	DeleteRecord               = ""
	SelectSeveralRecentRecords = "SELECT * FROM \"main\".\"ShortenUrls\" ORDER BY \"Date\" DESC LIMIT $1"
)
