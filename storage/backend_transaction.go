package storage

import "sync"
import "fmt"
import (
	"github.com/tecbot/gorocksdb"
)

type transaction struct {
	wo gorocksdb.WriteOptions
	ro gorocksdb.ReadOptions
}

func count_buckets() int {

}

func (txn *transaction) init() {
	init_func := func() {

		for i := 0; i < server_config.Buckets; i++ {
			var buffer bytes.Buffer
			buffer.WriteString("/Users/jaykpatel/gowork/godb-")
			buffer.Write([]byte(strconv.Itoa(i)))
			buffer.WriteString(".db")
			db_handle, err := gorocksdb.OpenDb(opts, buffer.String())

			if err != nil {
				fmt.Println("Error opening rocksdb:", err.Error())
				//TODO: throw exception
			}
			dbs[i] = db_handle
		}
		return dbs
		wo := gorocksdb.NewDefaultWriteOptions()
		ro := gorocksdb.NewDefaultReadOptions()
	}

	var dbs []*gorocksdb.DB = make([]*gorocksdb.DB, server_config.Buckets, server_config.MaxBuckets)
}
