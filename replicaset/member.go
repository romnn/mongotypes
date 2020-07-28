package replicaset

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// See https://github.com/timvaillancourt/go-mongodb-replset

// MemberHealth ...
type MemberHealth int

// MemberState ...
type MemberState int

const (
	// MemberHealthDown ...
	MemberHealthDown MemberHealth = iota
	// MemberHealthUp ...
	MemberHealthUp
	// MemberStateStartup ...
	MemberStateStartup MemberState = 0
	// MemberStatePrimary ...
	MemberStatePrimary MemberState = 1
	// MemberStateSecondary ...
	MemberStateSecondary MemberState = 2
	// MemberStateRecovering ...
	MemberStateRecovering MemberState = 3
	// MemberStateStartup2 ...
	MemberStateStartup2 MemberState = 5
	// MemberStateUnknown ...
	MemberStateUnknown MemberState = 6
	// MemberStateArbiter ...
	MemberStateArbiter MemberState = 7
	// MemberStateDown ...
	MemberStateDown MemberState = 8
	// MemberStateRollback ...
	MemberStateRollback MemberState = 9
	// MemberStateRemoved ...
	MemberStateRemoved MemberState = 10
)

// MemberStateStrings ...
var MemberStateStrings = map[MemberState]string{
	MemberStateStartup:    "STARTUP",
	MemberStatePrimary:    "PRIMARY",
	MemberStateSecondary:  "SECONDARY",
	MemberStateRecovering: "RECOVERING",
	MemberStateStartup2:   "STARTUP2",
	MemberStateUnknown:    "UNKNOWN",
	MemberStateArbiter:    "ARBITER",
	MemberStateDown:       "DOWN",
	MemberStateRollback:   "ROLLBACK",
	MemberStateRemoved:    "REMOVED",
}

// String ...
func (ms MemberState) String() string {
	if str, ok := MemberStateStrings[ms]; ok {
		return str
	}
	return ""
}

// Member ...
type Member struct {
	ID                int                 `bson:"_id" json:"_id"`
	Name              string              `bson:"name" json:"name"`
	Health            MemberHealth        `bson:"health" json:"health"`
	State             MemberState         `bson:"state" json:"state"`
	StateStr          string              `bson:"stateStr" json:"stateStr"`
	Uptime            int64               `bson:"uptime" json:"uptime"`
	Optime            *Optime             `bson:"optime" json:"optime"`
	OptimeDate        time.Time           `bson:"optimeDate" json:"optimeDate"`
	ConfigVersion     int                 `bson:"configVersion" json:"configVersion"`
	ElectionTime      primitive.Timestamp `bson:"electionTime,omitempty" json:"electionTime,omitempty"`
	ElectionDate      time.Time           `bson:"electionDate,omitempty" json:"electionDate,omitempty"`
	InfoMessage       string              `bson:"infoMessage,omitempty" json:"infoMessage,omitempty"`
	OptimeDurable     *Optime             `bson:"optimeDurable,omitempty" json:"optimeDurable,omitempty"`
	OptimeDurableDate time.Time           `bson:"optimeDurableDate,omitempty" json:"optimeDurableDate,omitempty"`
	LastHeartbeat     time.Time           `bson:"lastHeartbeat,omitempty" json:"lastHeartbeat,omitempty"`
	LastHeartbeatRecv time.Time           `bson:"lastHeartbeatRecv,omitempty" json:"lastHeartbeatRecv,omitempty"`
	PingMs            int64               `bson:"pingMs,omitempty" json:"pingMs,omitempty"`
	Self              bool                `bson:"self,omitempty" json:"self,omitempty"`
	SyncingTo         string              `bson:"syncingTo,omitempty" json:"syncingTo,omitempty"`
}

// GetSelf ...
func (s *Status) GetSelf() *Member {
	for _, member := range s.Members {
		if member.Self == true {
			return member
		}
	}
	return nil
}

// GetMember ...
func (s *Status) GetMember(name string) *Member {
	for _, member := range s.Members {
		if member.Name == name {
			return member
		}
	}
	return nil
}

// HasMember ...
func (s *Status) HasMember(name string) bool {
	return s.GetMember(name) != nil
}

// GetMemberID ...
func (s *Status) GetMemberID(id int) *Member {
	for _, member := range s.Members {
		if member.ID == id {
			return member
		}
	}
	return nil
}

// GetMembersByState ...
func (s *Status) GetMembersByState(state MemberState, limit int) []*Member {
	members := make([]*Member, 0)
	for _, member := range s.Members {
		if member.State == state {
			members = append(members, member)
			if limit > 0 && len(members) == limit {
				return members
			}
		}
	}
	return members
}

// Primary ...
func (s *Status) Primary() *Member {
	primary := s.GetMembersByState(MemberStatePrimary, 1)
	if len(primary) == 1 {
		return primary[0]
	}
	return nil
}

// Secondaries ...
func (s *Status) Secondaries() []*Member {
	return s.GetMembersByState(MemberStateSecondary, 0)
}
