package storage

type Message struct {
	Val []byte
}

type filesystem struct {
}

type backend struct {
	txn *transaction
	fs  *filesystem
}

type MessageReader struct {
	Message
}

type MessageWriter struct {
	Message
}

func Put(db *gorocksdb.DB, key, value []byte) error {
	wo := gorocksdb.NewDefaultWriteOptions()
	//TODO: Set options
	return db.Put(wo, []byte(key), []byte(value))
}

func Get(db *gorocksdb.DB, key []byte) (*gorocksdb.Slice, error) {
	ro := gorocksdb.NewDefaultReadOptions()
	return db.Get(ro, []byte(key)) //TODO: handle response here and convert to a storage layer types, not rocksdb specific.
}
