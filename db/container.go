package db

import (
	"experiment_lwc/commons"
	"experiment_lwc/config"
	"github.com/boltdb/bolt"
)

type Veth struct {
	Name string
	CIDR string
}

type ContainerData struct {
	Configuration config.Configuration
	VethMap       map[string]Veth
	PID           int
}

func (db *DB) SaveContainerData(tx *bolt.Tx, id string, config config.Configuration, pid int, vethMap map[string]Veth) {
	b := tx.Bucket([]byte("containers"))
	v, err := marshal(ContainerData{
		Configuration: config,
		VethMap:       vethMap,
		PID:           pid,
	})
	commons.Must(err)
	commons.Must(b.Put([]byte(id), v))
}

func (db *DB) DeleteContainerData(tx *bolt.Tx, id string) {
	b := tx.Bucket([]byte("containers"))
	commons.Must(b.Delete([]byte(id)))
}

func (db *DB) QueryAllContainerData(tx *bolt.Tx) []ContainerData {
	result := make([]ContainerData, 0)
	b := tx.Bucket([]byte("containers"))
	commons.Must(b.ForEach(func(k, v []byte) error {
		var c ContainerData
		commons.Must(unmarshal(v, &c))
		result = append(result, c)
		return nil
	}))
	return result
}

func (db *DB) QueryContainerData(tx *bolt.Tx, id string) *ContainerData {
	var result ContainerData
	b := tx.Bucket([]byte("containers"))
	v := b.Get([]byte(id))
	if v == nil || len(v) == 0 {
		return nil
	}
	commons.Must(unmarshal(v, &result))
	return &result
}
