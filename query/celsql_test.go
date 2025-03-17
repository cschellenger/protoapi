package query

import (
	"testing"
	"time"

	"darkport.net/protoapi/model"
)

func TestRenderEmpty(t *testing.T) {
	s := &SqlQueryBuilder{
		Table:     "la_crime",
		Prototype: (&model.CrimeData{}).ProtoReflect(),
	}
	result, args, err := s.BuildQuery(&model.SearchRequest{})
	if err != nil {
		panic(err)
	}
	resultExpected := "SELECT * FROM la_crime LIMIT 10"
	if result != resultExpected {
		t.Errorf("Expected %s, got %s", resultExpected, result)
	}
	if len(args) != 0 {
		t.Errorf("Expected no args, got %v", args)
	}
}

func TestRenderCompoundQuery(t *testing.T) {
	s := &SqlQueryBuilder{
		Table:     "la_crime",
		Prototype: (&model.CrimeData{}).ProtoReflect(),
	}
	result, args, err := s.BuildQuery(&model.SearchRequest{Query: "record.status == 'AA' && record.victim_age > 18"})
	if err != nil {
		panic(err)
	}
	resultExpected := "SELECT * FROM la_crime WHERE status = ? AND victim_age > ? LIMIT 10"
	if result != resultExpected {
		t.Errorf("Expected %s, got %s", resultExpected, result)
	}
	argsExpected := []any{"AA", int64(18)}
	if len(args) != len(argsExpected) {
		t.Errorf("Expected %v, got %v", argsExpected, args)
	}
	for i := range args {
		if args[i] != argsExpected[i] {
			t.Errorf("Expected %v, got %v", argsExpected[i], args[i])
		}
	}
}

func TestTimeComparison(t *testing.T) {
	s := &SqlQueryBuilder{
		Table:     "la_crime",
		Prototype: (&model.CrimeData{}).ProtoReflect(),
	}
	result, args, err := s.BuildQuery(&model.SearchRequest{Query: "record.date_reported > timestamp('2021-01-01T00:00:00Z')"})
	if err != nil {
		panic(err)
	}
	resultExpected := "SELECT * FROM la_crime WHERE date_reported > ? LIMIT 10"
	if result != resultExpected {
		t.Errorf("Expected %s, got %s", resultExpected, result)
	}
	argsExpected := []any{time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)}
	if len(args) != len(argsExpected) {
		t.Errorf("Expected %v, got %v", argsExpected, args)
	}
	for i := range args {
		if args[i] != argsExpected[i] {
			t.Errorf("Expected %v, got %v", argsExpected[i], args[i])
		}
	}
}

func TestSize(t *testing.T) {
	s := &SqlQueryBuilder{
		Table:     "la_crime",
		Prototype: (&model.CrimeData{}).ProtoReflect(),
	}
	result, args, err := s.BuildQuery(&model.SearchRequest{Query: "record.status == 'AA'", Size: 15})
	if err != nil {
		panic(err)
	}
	resultExpected := "SELECT * FROM la_crime WHERE status = ? LIMIT 15"
	if result != resultExpected {
		t.Errorf("Expected %s, got %s", resultExpected, result)
	}
	argsExpected := []any{"AA"}
	if len(args) != len(argsExpected) {
		t.Errorf("Expected %v, got %v", argsExpected, args)
	}
	for i := range args {
		if args[i] != argsExpected[i] {
			t.Errorf("Expected %v, got %v", argsExpected[i], args[i])
		}
	}
}

func TestSizeAndOffset(t *testing.T) {
	s := &SqlQueryBuilder{
		Table:     "la_crime",
		Prototype: (&model.CrimeData{}).ProtoReflect(),
	}
	result, args, err := s.BuildQuery(&model.SearchRequest{Query: "record.status == 'AA'", Size: 15, Offset: 5})
	if err != nil {
		panic(err)
	}
	resultExpected := "SELECT * FROM la_crime WHERE status = ? LIMIT 15 OFFSET 5"
	if result != resultExpected {
		t.Errorf("Expected %s, got %s", resultExpected, result)
	}
	argsExpected := []any{"AA"}
	if len(args) != len(argsExpected) {
		t.Errorf("Expected %v, got %v", argsExpected, args)
	}
	for i := range args {
		if args[i] != argsExpected[i] {
			t.Errorf("Expected %v, got %v", argsExpected[i], args[i])
		}
	}
}

func TestSortOrder(t *testing.T) {
	s := &SqlQueryBuilder{
		Table:     "la_crime",
		Prototype: (&model.CrimeData{}).ProtoReflect(),
	}
	result, args, err := s.BuildQuery(&model.SearchRequest{Query: "record.status == 'AA'", Order: model.SortOrder_DESC, Sort: "date_reported"})
	if err != nil {
		panic(err)
	}
	resultExpected := "SELECT * FROM la_crime WHERE status = ? ORDER BY date_reported DESC LIMIT 10"
	if result != resultExpected {
		t.Errorf("Expected %s, got %s", resultExpected, result)
	}
	argsExpected := []any{"AA"}
	if len(args) != len(argsExpected) {
		t.Errorf("Expected %v, got %v", argsExpected, args)
	}
	for i := range args {
		if args[i] != argsExpected[i] {
			t.Errorf("Expected %v, got %v", argsExpected[i], args[i])
		}
	}
}
