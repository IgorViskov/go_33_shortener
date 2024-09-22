package migrator

var (
	CreatePartialIndex = "CREATE UNIQUE INDEX IF NOT EXISTS idx_urls_is_deleted ON urls (\"Value\") WHERE \"IsDeleted\" = 0"
)
