package db

import (
	"experiment_lwc/commons"
	"experiment_lwc/config"
	"github.com/boltdb/bolt"
)

type NetworkData struct {
	Name               string
	Address            string
	AttachedContainers map[string]bool
}

func (db *DB) SaveNetworkData(tx *bolt.Tx, name string, address string, containers map[string]bool) {
	b := tx.Bucket([]byte("networks"))
	v, err := marshal(NetworkData{
		Name:               name,
		Address:            address,
		AttachedContainers: containers,
	})
	commons.Must(err)
	commons.Must(b.Put([]byte(name), v))
}

func (db *DB) DeleteNetworkData(tx *bolt.Tx, name string) {
	b := tx.Bucket([]byte("networks"))
	commons.Must(b.Delete([]byte(name)))
}

func (db *DB) QueryAllNetworkData(tx *bolt.Tx) []NetworkData {
	result := make([]NetworkData, 0)
	b := tx.Bucket([]byte("networks"))
	commons.Must(b.ForEach(func(k, v []byte) error {
		var c NetworkData
		commons.Must(unmarshal(v, &c))
		result = append(result, c)
		return nil
	}))
	return result
}

func (db *DB) QueryNetworkData(tx *bolt.Tx, name string) *NetworkData {
	var result NetworkData
	b := tx.Bucket([]byte("networks"))
	v := b.Get([]byte(name))
	if v == nil || len(v) == 0 {
		return nil
	}
	commons.Must(unmarshal(v, result))
	return &result
}

func (db *DB) AddContainerToNetwork(tx *bolt.Tx, name string, address string, containerId string, veth string, cidr string) {
	b := tx.Bucket([]byte("networks"))
	var networkData NetworkData
	v := b.Get([]byte(name))
	commons.Must(unmarshal(v, &networkData))
	networkData.AttachedContainers[containerId] = true
	v, err := marshal(networkData)
	commons.Must(err)
	commons.Must(b.Put([]byte(name), v))
	b = tx.Bucket([]byte("containers"))
	var containerData ContainerData
	v = b.Get([]byte(containerId))
	commons.Must(unmarshal(v, &containerData))
	containerData.Configuration.Networks[name] = config.Network{
		Address: address,
		CIDR:    cidr,
	}
	containerData.VethMap[name] = Veth{
		Name: veth,
		CIDR: cidr,
	}
	v, err = marshal(containerData)
	commons.Must(err)
	commons.Must(b.Put([]byte(containerId), v))
}

func (db *DB) RemoveContainerFromNetwork(tx *bolt.Tx, name string, containerId string) {
	b := tx.Bucket([]byte("networks"))
	var data NetworkData
	v := b.Get([]byte(name))
	commons.Must(unmarshal(v, &data))
	delete(data.AttachedContainers, containerId)
	v, err := marshal(data)
	commons.Must(err)
	commons.Must(b.Put([]byte(name), v))
	b = tx.Bucket([]byte("containers"))
	var containerData ContainerData
	v = b.Get([]byte(containerId))
	commons.Must(unmarshal(v, &containerData))
	delete(containerData.Configuration.Networks, name)
	delete(containerData.VethMap, name)
	v, err = marshal(containerData)
	commons.Must(err)
	commons.Must(b.Put([]byte(containerId), v))
}
