package db_pkg

type DB interface {
	Get(name string)
	NewDao(opts *NewDaoOptions) Dao
}
