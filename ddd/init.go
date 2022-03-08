package ddd

var _eventStorage EventStorage

func Init(eventStorage EventStorage) {
	_eventStorage = eventStorage
}
