package errors

var (
	DatabaseError               = createDescription("DatabaseError")
	DataConsistencyError        = createDescription("DataConsistencyError")
	DataNotFoundError           = createDescription("DataNotFoundError")
	EncodeUncompletedBlockError = createDescription("EncodeUncompletedBlockError")
)
