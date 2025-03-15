package query

import (
	"testing"

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

func TestRenderQuery(t *testing.T) {
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
