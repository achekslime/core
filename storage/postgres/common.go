package postgres

const (
	UserTableName           = "client"
	RoomTableName           = "room"
	AvailableRoomsTableName = "available_rooms"

	ErrUniqueConstraintDuplicate = "duplicate key value violates unique constraint"
	ErrSqlNoRows                 = "sql: no rows in result set"
)
