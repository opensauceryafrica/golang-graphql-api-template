package dal

import "sync"

// BlacheMutex is a mutex for the application
var BlacheMutex = sync.Mutex{}

// FactoryTableMutex is a mutex for the factory table
var FactoryTableMutex = sync.Mutex{}
