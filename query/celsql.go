package query

import (
	"fmt"
	"time"

	"darkport.net/protoapi/model"
	"github.com/google/cel-go/cel"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SqlQueryBuilder struct {
	Table     string
	Prototype protoreflect.Message
}

func (s *SqlQueryBuilder) BuildQuery(searchRequest *model.SearchRequest) (string, []any, error) {
	env, err := cel.NewEnv(
		cel.TypeDescs(s.Prototype.Descriptor().ParentFile()),
		cel.Variable("record", cel.ObjectType(string(s.Prototype.Descriptor().FullName()))),
	)
	if err != nil {
		return "", nil, err
	}
	args := make([]any, 0)
	queryBase := fmt.Sprintf("SELECT * FROM %s", s.Table)
	var where string
	if searchRequest.Query != "" {
		ast, iss := env.Compile(searchRequest.Query)
		if iss != nil && iss.Err() != nil {
			return "", nil, iss.Err()
		}
		expression, err := cel.AstToCheckedExpr(ast)
		if err != nil {
			return "", nil, err
		}
		where, args, err = s.RenderExpr(expression.Expr, args)
		if err != nil {
			return "", args, err
		}
	}
	if where != "" {
		where = " WHERE " + where
	}
	sort, err := s.RenderSort(searchRequest)
	if err != nil {
		return "", args, err
	}
	limit := s.RenderLimitAndOffset(searchRequest)
	return queryBase + where + sort + limit, args, nil
}

func (s *SqlQueryBuilder) RenderSort(searchRequest *model.SearchRequest) (string, error) {
	if searchRequest.Sort == "" {
		return "", nil
	}
	validName := s.Prototype.Descriptor().Fields().ByName(protoreflect.Name(searchRequest.Sort))
	if validName == nil {
		return "", fmt.Errorf("invalid sort field: %s", searchRequest.Sort)
	}
	order := "ASC"
	switch searchRequest.Order {
	case model.SortOrder_ASC:
		order = "ASC"
	case model.SortOrder_DESC:
		order = "DESC"
	default:
		return "", fmt.Errorf("invalid sort order: %v", searchRequest.Order)
	}
	return fmt.Sprintf(" ORDER BY %s %s", validName.Name(), order), nil
}

func (s *SqlQueryBuilder) RenderLimitAndOffset(searchRequest *model.SearchRequest) string {
	limit := searchRequest.Size
	if searchRequest.Size == 0 {
		limit = 10
	}
	offset := ""
	if searchRequest.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", searchRequest.Offset)
	}
	return fmt.Sprintf(" LIMIT %d%s", limit, offset)
}

func (s *SqlQueryBuilder) RenderExpr(expression *exprpb.Expr, args []any) (string, []any, error) {
	call := expression.GetCallExpr()
	if call != nil {
		return s.RenderCall(call, args)
	}
	constant := expression.GetConstExpr()
	if constant != nil {
		switch constant.ConstantKind.(type) {
		case *exprpb.Constant_StringValue:
			args = append(args, constant.GetStringValue())
			return "?", args, nil
		case *exprpb.Constant_Int64Value:
			args = append(args, constant.GetInt64Value())
			return "?", args, nil
		case *exprpb.Constant_Uint64Value:
			args = append(args, constant.GetUint64Value())
			return "?", args, nil
		case *exprpb.Constant_DoubleValue:
			args = append(args, constant.GetDoubleValue())
			return "?", args, nil
		case *exprpb.Constant_BoolValue:
			args = append(args, constant.GetBoolValue())
			return "?", args, nil
		case *exprpb.Constant_NullValue:
			return "NULL", args, nil
		case *exprpb.Constant_BytesValue:
			args = append(args, constant.GetBytesValue())
			return "?", args, nil
		default:
			return "", args, fmt.Errorf("unsupported constant type: %v", constant.ConstantKind)
		}
	}
	ident := expression.GetIdentExpr()
	if ident != nil {
		return ident.GetName(), args, nil
	}
	select_expr := expression.GetSelectExpr()
	if select_expr != nil {
		return select_expr.Field, args, nil
	}
	return "", args, fmt.Errorf("unsupported expr type: %v", expression.ExprKind)
}

func (s *SqlQueryBuilder) RenderCall(call *exprpb.Expr_Call, args []any) (string, []any, error) {
	// handle special case of timestamp function
	if call.Function == "timestamp" {
		arg := call.Args[0]
		time, err := time.Parse(time.RFC3339, arg.GetConstExpr().GetStringValue())
		if err != nil {
			return "", args, err
		}
		args = append(args, time)
		return "?", args, nil
	}
	// handle other functions with 2 args
	lhs := call.Args[0]
	lhs_string, args, err := s.RenderExpr(lhs, args)
	if err != nil {
		return "", args, err
	}
	rhs := call.Args[1]
	rhs_string, args, err := s.RenderExpr(rhs, args)
	if err != nil {
		return "", args, err
	}
	opr := ""
	switch call.Function {
	case "_==_":
		opr = " = "
	case "_&&_":
		opr = " AND "
	case "_||_":
		opr = " OR "
	case "_>_":
		opr = " > "
	case "_>=_":
		opr = " >= "
	case "_<_":
		opr = " < "
	case "_<=_":
		opr = " <= "
	default:
		return "", args, fmt.Errorf("unsupported function: %s", call.Function)
	}
	return lhs_string + opr + rhs_string, args, nil

}
