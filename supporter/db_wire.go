package supporter

import "gorm.io/gorm"

type DBWire func() (func(*gorm.DB), interface{})

type DBWires []DBWire