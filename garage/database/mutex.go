package database

import "sync"

// CenditMutex is a mutex for the application
var CenditMutex = sync.Mutex{}

// FactoryTableMutex is a mutex for the factory table
var FactoryTableMutex = sync.Mutex{}
