package storage

type VisitRecorder interface {
	RecordVisit(string) error
}

type VisitCounter interface {
	CountVisits(string) (int, error)
}

type VisitRecorderCounter interface {
	VisitRecorder
	VisitCounter
}

type PurchaseAllower interface {
	IsAllowed(cardNumber string, cents uint64) (bool, error)
}
