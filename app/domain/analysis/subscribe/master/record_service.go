package master

type recordYearDao struct {
}

type recordMonthDao struct {
}

type recordDayDao struct {
}

type RecordService struct {
	yearDao *recordYearDao
	mothDao *recordMonthDao
	dayDao  *recordDayDao
}

func NewRecordService() *RecordService {
	return &RecordService{
		yearDao: &recordYearDao{},
		mothDao: &recordMonthDao{},
		dayDao:  &recordDayDao{},
	}
}

func (s *RecordService) Create() {

}

func (s *RecordService) Update() {

}

func (s *RecordService) Delete() {

}
