// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"com.enesuysal/go-chat/api/ent/message"
	"com.enesuysal/go-chat/api/ent/user"
	"entgo.io/ent/dialect/sql"
)

// Message is the model entity for the Message schema.
type Message struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// SenderUsername holds the value of the "senderUsername" field.
	SenderUsername string `json:"senderUsername,omitempty"`
	// ReceiverUsername holds the value of the "ReceiverUsername" field.
	ReceiverUsername string `json:"ReceiverUsername,omitempty"`
	// Message holds the value of the "message" field.
	Message string `json:"message,omitempty"`
	// SendTime holds the value of the "sendTime" field.
	SendTime time.Time `json:"sendTime,omitempty"`
	// Seen holds the value of the "seen" field.
	Seen int `json:"seen,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MessageQuery when eager-loading is set.
	Edges        MessageEdges `json:"edges"`
	user_message *int
}

// MessageEdges holds the relations/edges for other nodes in the graph.
type MessageEdges struct {
	// Owner holds the value of the owner edge.
	Owner *User `json:"owner,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageEdges) OwnerOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Owner == nil {
			// The edge owner was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Message) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case message.FieldID, message.FieldSeen:
			values[i] = new(sql.NullInt64)
		case message.FieldSenderUsername, message.FieldReceiverUsername, message.FieldMessage:
			values[i] = new(sql.NullString)
		case message.FieldSendTime:
			values[i] = new(sql.NullTime)
		case message.ForeignKeys[0]: // user_message
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Message", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Message fields.
func (m *Message) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case message.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			m.ID = int(value.Int64)
		case message.FieldSenderUsername:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field senderUsername", values[i])
			} else if value.Valid {
				m.SenderUsername = value.String
			}
		case message.FieldReceiverUsername:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ReceiverUsername", values[i])
			} else if value.Valid {
				m.ReceiverUsername = value.String
			}
		case message.FieldMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[i])
			} else if value.Valid {
				m.Message = value.String
			}
		case message.FieldSendTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field sendTime", values[i])
			} else if value.Valid {
				m.SendTime = value.Time
			}
		case message.FieldSeen:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field seen", values[i])
			} else if value.Valid {
				m.Seen = int(value.Int64)
			}
		case message.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_message", value)
			} else if value.Valid {
				m.user_message = new(int)
				*m.user_message = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryOwner queries the "owner" edge of the Message entity.
func (m *Message) QueryOwner() *UserQuery {
	return (&MessageClient{config: m.config}).QueryOwner(m)
}

// Update returns a builder for updating this Message.
// Note that you need to call Message.Unwrap() before calling this method if this Message
// was returned from a transaction, and the transaction was committed or rolled back.
func (m *Message) Update() *MessageUpdateOne {
	return (&MessageClient{config: m.config}).UpdateOne(m)
}

// Unwrap unwraps the Message entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (m *Message) Unwrap() *Message {
	tx, ok := m.config.driver.(*txDriver)
	if !ok {
		panic("ent: Message is not a transactional entity")
	}
	m.config.driver = tx.drv
	return m
}

// String implements the fmt.Stringer.
func (m *Message) String() string {
	var builder strings.Builder
	builder.WriteString("Message(")
	builder.WriteString(fmt.Sprintf("id=%v", m.ID))
	builder.WriteString(", senderUsername=")
	builder.WriteString(m.SenderUsername)
	builder.WriteString(", ReceiverUsername=")
	builder.WriteString(m.ReceiverUsername)
	builder.WriteString(", message=")
	builder.WriteString(m.Message)
	builder.WriteString(", sendTime=")
	builder.WriteString(m.SendTime.Format(time.ANSIC))
	builder.WriteString(", seen=")
	builder.WriteString(fmt.Sprintf("%v", m.Seen))
	builder.WriteByte(')')
	return builder.String()
}

// Messages is a parsable slice of Message.
type Messages []*Message

func (m Messages) config(cfg config) {
	for _i := range m {
		m[_i].config = cfg
	}
}
