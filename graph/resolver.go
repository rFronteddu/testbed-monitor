package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/hashicorp/go-memdb"
	"testbed-monitor/graph/generated"
	"testbed-monitor/graph/model"
)

type Resolver struct {
	DB *memdb.MemDB
}

const (
	HostStatusTable = "host_table"
)

func (r *Resolver) GetHosts() ([]*model.HostStatus, error) {
	txn := r.DB.Txn(false)
	defer txn.Abort()
	it, err := txn.Get(HostStatusTable, "id")
	if err != nil {
		panic(err)
	}
	var list []*model.HostStatus
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*model.HostStatus)
		list = append(list, p)
	}
	return list, nil
}

//GetHost nil if id was not found
func (r *Resolver) GetHost(id string) interface{} {
	// Create read-only transaction
	txn := r.DB.Txn(false)
	defer txn.Abort()
	raw, err := txn.First(HostStatusTable, "id", id)
	if err != nil {
		panic(err)
	}
	return raw
}

func (r *Resolver) CommitHost(hostStatus *model.HostStatus) {
	txn := r.DB.Txn(true)
	if err := txn.Insert(HostStatusTable, hostStatus); err != nil {
		log.Printf("Fatal error: %s while inserting struct: %v in database\n", err, hostStatus)
		return
	}
	// Commit the transaction
	txn.Commit()
}

func NewResolver() (generated.Config, *Resolver) {
	r := Resolver{}
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			HostStatusTable: {
				Name: HostStatusTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"Time": {
						Name:    "Time",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Time"},
					},
					"HostID": {
						Name:    "HostID",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "HostID"},
					},
					"HostName": {
						Name:    "HostName",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "HostName"},
					},
					"Os": {
						Name:    "Os",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Os"},
					},
					"Platform": {
						Name:    "Platform",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Platform"},
					},
					"Kernel": {
						Name:    "Kernel",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Kernel"},
					},
					"BootTime": {
						Name:    "BootTime",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "BootTime"},
					},
					"reachable": {
						Name:    "reachable",
						Unique:  false,
						Indexer: &memdb.BoolFieldIndex{Field: "reachable"},
					},
					"VirtualMemoryUsagePercent": {
						Name:    "VirtualMemoryUsagePercent",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "VirtualMemoryUsagePercent"},
					},
					"DiskUsagePercent": {
						Name:    "DiskUsagePercent",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "DiskUsagePercent"},
					},
					"DiskFree": {
						Name:    "DiskFree",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "DiskFree"},
					},
					"cpu": {
						Name:    "cpu",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "cpu"},
					},
					"Load1": {
						Name:    "Load1",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Load1"},
					},
					"VirtualMemoryFree": {
						Name:    "VirtualMemoryFree",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "VirtualMemoryFree"},
					},
					"Load5": {
						Name:    "Load5",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Load5"},
					},
					"Load15": {
						Name:    "Load15",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Load15"},
					},
				},
			},
		},
	}
	var err error
	r.DB, err = memdb.NewMemDB(schema)
	if err != nil {

		panic(err)
	}
	return generated.Config{
		Resolvers: &r,
	}, &r
}
