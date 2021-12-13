package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "initial_19101747",
			Migrate: func(db *gorm.DB) error {
				query := "create table if not exists users (" +
					"uuid text not null constraint users_pk primary key, " +
					"last_name text, " +
					"first_name text, " +
					"username text not null unique, " +
					"created_at timestamp default current_timestamp, " +
					"email text not null unique, " +
					"date_of_birth timestamp);"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "feat_topics_init",
			Migrate: func(db *gorm.DB) error {
				query := "create table if not exists topics (" +
					"id bigserial not null constraint topics_pk primary key, " +
					"name text not null, " +
					"created_at timestamp default current_timestamp, " +
					"updated_at timestamp default current_timestamp, " +
					"deleted_at timestamp, " +
					"general bool);"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "feat_words_init",
			Migrate: func(db *gorm.DB) error {
				query := "create table if not exists words (" +
					"id bigserial not null constraint words_pk primary key, " +
					"english_spelling text not null, " +
					"french_spelling text not null, " +
					"created_at timestamp default current_timestamp, " +
					"updated_at timestamp default current_timestamp, " +
					"deleted_at timestamp);"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "feat_topics_words_init",
			Migrate: func(db *gorm.DB) error {
				query := "create table if not exists topics_words (" +
					"id bigserial not null constraint topics_words_pk primary key, " +
					"topic_id bigint not null\n\t\tconstraint topics_words_topic_id_fk\n\t\t\treferences topics, " +
					"word_id bigint not null\n\t\tconstraint topics_words_word_id_fk\n\t\t\treferences words, " +
					"created_at timestamp default current_timestamp, " +
					"updated_at timestamp default current_timestamp, " +
					"deleted_at timestamp);"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "feat_topics_users_init",
			Migrate: func(db *gorm.DB) error {
				query := "create table if not exists topics_users (" +
					"id bigserial not null constraint topics_users_pk primary key, " +
					"topic_id bigint not null\n\t\tconstraint topics_users_topic_id_fk\n\t\t\treferences topics, " +
					"user_id text not null\n\t\tconstraint topics_users_user_id_fk\n\t\t\treferences users, " +
					"created_at timestamp default current_timestamp, " +
					"updated_at timestamp default current_timestamp, " +
					"deleted_at timestamp);"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "add_columns_user_141220210044",
			Migrate: func(db *gorm.DB) error {
				query := "alter table users add column avatar text;" +
					"alter table users add column points integer default 0"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
	}
}
