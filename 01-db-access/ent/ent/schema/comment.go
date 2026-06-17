package schema

import (
	"time"

	"entgo.io/ent"
	// "entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Comment struct {
	ent.Schema
}

func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Text("body"),
		field.Time("created_at").Default(time.Now),
	}
}

func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).Ref("comments").Unique().Required(),
		edge.From("post", Post.Type).Ref("comments").Unique().Required(),
	}
}