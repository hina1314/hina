package db

import (
	"bufio"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const dataFile = "snapshot.gob"

type DB struct {
	ctx      context.Context
	String   *String
	HashMap  *HashMap
	dataLock sync.RWMutex
}

func NewDB(ctx context.Context) *DB {
	db := &DB{
		ctx:     ctx,
		String:  NewStrings(),
		HashMap: NewHashMap(),
	}
	gob.Register(map[string]map[string]string{})
	gob.Register(map[string]string{})
	err := db.LoadSnapshot()
	if err != nil {
		log.Fatal(err)
	}
	go db.autoPersist()
	return db
}

func (d *DB) SaveSnapshot() error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	defer bufWriter.Flush()

	encoder := gob.NewEncoder(bufWriter)

	d.dataLock.RLock()
	defer d.dataLock.RUnlock()

	copies := map[string]interface{}{
		"string": d.String.Data,
		"hash":   d.HashMap.data,
	}

	return encoder.Encode(copies)
}

func (d *DB) LoadSnapshot() error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	copies := make(map[string]interface{})
	decoder := gob.NewDecoder(file)
	decoder.Decode(&copies)
	if copies["string"] == nil {
		return nil
	}
	d.dataLock.Lock()
	d.String.Data = copies["string"].(map[string]string)
	d.HashMap.data = copies["hash"].(map[string]map[string]string)
	d.dataLock.Unlock()
	return nil
}

func (d *DB) autoPersist() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("saving data...")
			if err := d.SaveSnapshot(); err != nil {
				log.Printf("Auto-save error: %v", err)
			}
		case <-d.ctx.Done():
			return
		}
	}
}
