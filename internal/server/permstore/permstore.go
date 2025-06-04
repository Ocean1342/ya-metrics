package permstore

import (
	"fmt"
	"ya-metrics/config"
)

type PermStore struct {
}

type PermanentStorable interface {
	//метод достаёт данные из перманентного хранилища и помещает их в стораджи
	ExtractFromPermStore() error
	PutDataToPermStore() error
}

func NewPermStore(_ *config.Config) PermanentStorable {
	return &PermStore{}
}

func (ps *PermStore) ExtractFromPermStore() error {
	fmt.Println("ExtractFromPermStore")
	return nil
}

func (ps *PermStore) PutDataToPermStore() error {
	fmt.Println("PutDataToPermStore")
	return nil
}
