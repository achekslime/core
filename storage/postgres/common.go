package postgres

const (
	UserTableName           = "client"
	RoomTableName           = "room"
	AvailableRoomsTableName = "available_rooms"

	errUniqueConstraintDuplicate = "duplicate key value violates unique constraint"
	errSqlNoRows                 = "sql: no rows in result set"
)
