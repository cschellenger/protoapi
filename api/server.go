package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"darkport.net/protoapi/model"
	"darkport.net/protoapi/query"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Name         string
	Db           *sqlx.DB
	Prototype    protoreflect.Message
	QueryBuilder query.QueryBuilder
}

func NewSqlServer(db *sqlx.DB, prototype protoreflect.Message, table string) *Server {
	return &Server{
		Name:      "ProtoAPI",
		Db:        db,
		Prototype: prototype,
		QueryBuilder: &query.SqlQueryBuilder{
			Table:     table,
			Prototype: prototype,
		},
	}
}

func (s *Server) Close() {
	s.Db.Close()
}

func (s *Server) sayHello(c *gin.Context) {
	name := c.DefaultQuery("Name", "User")
	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello %s. My name is: %s", name, s.Name)})
}

func (s *Server) search(c *gin.Context) {
	var searchRequest model.SearchRequest
	if err := c.BindQuery(&searchRequest); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to bind: %v", err)})
		return
	}
	query, args, err := s.QueryBuilder.BuildQuery(&searchRequest)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to build query: %v", err)})
		return
	}
	rows, err := s.Db.Queryx(query, args...)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to query DB: %v", err)})
		return
	}
	defer rows.Close()
	var records []*protoreflect.ProtoMessage

	for rows.Next() {
		var recordMap = make(map[string]any)
		err = rows.MapScan(recordMap)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Failed to read rows as map: %v", err)})
			return
		}
		records = append(records, s.marshall(recordMap))

	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": records})
}

func (s *Server) marshall(recordMap map[string]any) *protoreflect.ProtoMessage {
	fieldDescriptors := s.Prototype.Descriptor().Fields()
	record := s.Prototype.New().Interface()
	recordReflect := record.ProtoReflect()

	for key, value := range recordMap {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		if value == nil {
			fmt.Printf("Value for %s is nil\n", key)
			continue
		}
		fd := fieldDescriptors.ByName(protoreflect.Name(key))
		if fd != nil {
			fdKind := fd.Kind()
			fmt.Printf("%s, kind=%v", key, fdKind)
			if fdKind == protoreflect.MessageKind {
				if fd.Message().FullName() == "google.protobuf.Timestamp" {
					// Convert time.Time to google.protobuf.Timestamp
					ts := timestamppb.New(value.(time.Time))
					recordReflect.Set(fd, protoreflect.ValueOfMessage(ts.ProtoReflect()))
				} else {
					fmt.Printf("Field %s is an unsupported message of type: %s\n", key, fd.Message().FullName())
				}
			} else {
				recordReflect.Set(fd, protoreflect.ValueOf(value))
			}
		}
	}
	return &record
}

func (s *Server) Serve() {
	router := gin.Default()
	router.GET("/hello", s.sayHello)
	router.GET("/search", s.search)

	router.Run("localhost:8080")
}
