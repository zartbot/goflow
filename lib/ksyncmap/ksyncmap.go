//Package ksyncmap : a timeout based syncmap
package ksyncmap

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

//KSyncMap is the main struct of timeout based syncmap
//Data : data syncmap.
//Timeout : Timeout value
//UpdateTimestamp store last updatetime, and avoid datamap update freq.
type Map struct {
	Name       string
	Data       sync.Map
	Timeout    int64
	CheckFreq  int64
	ExpireTime sync.Map
	Verbose    bool
}

func NewMap(name string, timeout int64, checkfreq int64, verbose bool) *Map {
	return &Map{
		Name:      name,
		Timeout:   timeout,
		CheckFreq: checkfreq,
		Verbose:   verbose,
	}
}

func (k *Map) Load(key interface{}) (value interface{}, ok bool) {
	return k.Data.Load(key)
}

func (k *Map) Store(key interface{}, value interface{}, currentTime time.Time) {
	//Check ExpireTime Map.
	exp, ok := k.ExpireTime.Load(key)
	if !ok {
		expireTime := currentTime.Add(time.Duration(k.Timeout) * time.Second)
		k.ExpireTime.Store(key, expireTime)
	} else {
		elapsed := exp.(time.Time).Sub(currentTime)
		//elapsed time less than half of timeout, update ExpireTime Store.
		if elapsed < time.Duration(k.Timeout/2) {
			expireTime := currentTime.Add(time.Duration(k.Timeout) * time.Second)
			k.ExpireTime.Store(key, expireTime)
		}
	}
	k.Data.Store(key, value)
}

func (k *Map) UpdateTime(key interface{}, currentTime time.Time) {
	expireTime := currentTime.Add(time.Duration(k.Timeout) * time.Second)
	k.ExpireTime.Store(key, expireTime)
}

func (k *Map) Delete(key interface{}) {
	k.Data.Delete(key)
	k.ExpireTime.Delete(key)
}

func (k *Map) Run() {
	rand.Seed(time.Now().UnixNano())
	r := k.CheckFreq / 5
	for {
		currentTime := time.Now()
		k.ExpireTime.Range(func(key, v interface{}) bool {
			value := v.(time.Time)
			if value.Sub(currentTime) < 0 {
				//fmt.Println("DEBUG:::DELETE-KEY", reflect.ValueOf(key))
				k.Data.Delete(key)
				k.ExpireTime.Delete(key)
			}
			return true
		})
		if k.Verbose {
			k.ShowExpireTime()
			k.ShowData()
		}
		time.Sleep(time.Second * time.Duration(k.CheckFreq+rand.Int63n(r)))
	}
}

func (k *Map) ShowExpireTime() {
	fmt.Printf("%10s:--------------------Expire Time Table-------------------------------\n", k.Name)
	i := 1
	k.ExpireTime.Range(func(kt, v interface{}) bool {
		value := v.(time.Time)
		key := reflect.ValueOf(kt)
		fmt.Printf("%10s:[%6d]Key:%v | ExipreTime: %v \n", k.Name, i, key, value)
		i++
		return true
	})
	fmt.Printf("%10s:--------------------------------------------------------------------\n\n", k.Name)
}

func (k *Map) ShowData() {
	fmt.Printf("%10s:---------------------Data Table-------------------------------------\n", k.Name)
	i := 1
	k.Data.Range(func(kt, v interface{}) bool {
		value := reflect.ValueOf(v)
		key := reflect.ValueOf(kt)
		fmt.Printf("%10s:[%6d]Key:%v | Value: %v \n", k.Name, i, key, value)
		i++
		return true
	})
	fmt.Printf("%10s:--------------------------------------------------------------------\n\n", k.Name)
}
