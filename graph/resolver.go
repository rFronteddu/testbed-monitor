package graph

import (
	"github.com/hashicorp/go-memdb"
	"log"
	"testbed-monitor/graph/generated"
	"testbed-monitor/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

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
	log.Println("Host report committed to DB.")
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
					"Tower": {
						Name:    "Tower",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Tower"},
					},
					"BoardReached": {
						Name:    "BoardReached",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "BoardReached"},
					},
					"TowerReached": {
						Name:    "TowerReached",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "TowerReached"},
					},
					"BootTime": {
						Name:    "BootTime",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "BootTime"},
					},
					"Reboots": {
						Name:    "Reboots",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Reboots"},
					},
					"UsedRAM": {
						Name:    "UsedRAM",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "UsedRAM"},
					},
					"UsedDisk": {
						Name:    "UsedDisk",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "UsedDisk"},
					},
					"CPU": {
						Name:    "CPU",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "CPU"},
					},
					"Reachable": {
						Name:    "Reachable",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Reachable"},
					},
					"Temperature": {
						Name:    "Temperature",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Temperature"},
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
