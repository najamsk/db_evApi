package models

import "github.com/satori/go.uuid"

// 	CREATE TABLE sponsor_levels (id UUID PRIMARY KEY, name string, client_id uuid, conference_id uuid, is_active bool, sort_order int,
// 	deleted_at timestamptz, created_at timestamptz, updated_at timestamptz, deleted bool, created_by uuid, updated_by uuid, deleted_by uuid);
	
//    insert into sponsor_levels(id, name, client_id, conference_id, is_active, sort_order) 
//    values(gen_random_uuid(), 'Diamond Sponsors', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1),
//    (gen_random_uuid(), 'Platinum Sponsors', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2),
//    (gen_random_uuid(), 'Gold Sponsors', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3),
//    (gen_random_uuid(), 'Silver Sponsors', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4)

type SponsorLevel struct {
	Base
	Name             string
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Sponsors         []Sponsor

	// CreatedAt        time.Time
	// UpdatedAt        time.Time
	// DeletedAt        time.Time
}