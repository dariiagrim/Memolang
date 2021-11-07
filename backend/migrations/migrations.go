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
					"id bigserial not null constraint users_pk primary key, " +
					"last_name text not null, " +
					"first_name text not null, " +
					"username text not null unique, " +
					"password text not null, " +
					"created_at timestamp not null, " +
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
					"description text not null, " +
					"created_at timestamp not null, " +
					"updated_at timestamp not null, " +
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
					"created_at timestamp not null, " +
					"updated_at timestamp not null, " +
					"deleted_at timestamp, " +
					"definition text);"
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
					"created_at timestamp not null, " +
					"updated_at timestamp not null, " +
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
					"user_id bigint not null\n\t\tconstraint topics_users_user_id_fk\n\t\t\treferences users, " +
					"created_at timestamp not null, " +
					"updated_at timestamp not null, " +
					"deleted_at timestamp);"
				return db.Exec(query).Error

			},
			Rollback: func(db *gorm.DB) error {
				return nil
			},
		},
	}
}
