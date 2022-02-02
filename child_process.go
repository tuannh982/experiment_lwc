package main

import (
	"experiment_lwc/commons"
	"experiment_lwc/db"
	"os"
	"syscall"
)

func getVethCidrMap() map[string]string {
	DB := db.Setup(dbPath)
	tx, err := DB.DB.Begin(true)
	commons.Must(err)
	defer func() {
		_ = tx.Rollback()
		commons.Must(DB.DB.Close())
	}()
	values := make([]db.Veth, 0)
	containerData := DB.QueryContainerData(tx, configuration.ID)
	for _, v := range containerData.VethMap {
		values = append(values, v)
	}
	vethCidrMap := make(map[string]string)
	for _, v := range values {
		vethCidrMap[v.Name] = v.CIDR
	}
	commons.Must(tx.Commit())
	return vethCidrMap
}

func HandleChildProcess() {
	ParseRunArgs()
	configuration.ID = os.Args[1]
	vethCidrMap := getVethCidrMap()
	SpawnContainer()
	SetupContainerNetworkInterface(vethCidrMap)
	commons.Must(syscall.Exec("/bin/sh", []string{"sh"}, os.Environ()))
}
