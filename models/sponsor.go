package models

import "github.com/satori/go.uuid"


// CREATE TABLE sponsors (id UUID PRIMARY KEY, name string, client_id uuid, conference_id uuid, is_active bool, sort_order int, sponsor_level_id uuid, 
// deleted_at timestamptz, created_at timestamptz, updated_at timestamptz, deleted bool, created_by uuid, updated_by uuid, deleted_by uuid);
	
// insert into sponsors(id, name, client_id, conference_id, is_active, sort_order, sponsor_level_id) 
// values(gen_random_uuid(), 'Moftak Solutions', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
//  (gen_random_uuid(), 'LaunchPad7', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
//  (gen_random_uuid(), 'Vitrucare', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
//  (gen_random_uuid(), 'MAVEN LOGIX', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
   
//    (gen_random_uuid(), 'University of Haripur', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'de443293-88d0-4fb0-a383-018546134411'),
//    (gen_random_uuid(), 'PCHF', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2, 'de443293-88d0-4fb0-a383-018546134411'),
//    (gen_random_uuid(), 'SYSTEK', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3, 'de443293-88d0-4fb0-a383-018546134411'),
//    (gen_random_uuid(), 'OPEN Islamabad', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4, 'de443293-88d0-4fb0-a383-018546134411'),
   
//    (gen_random_uuid(), 'rowthli', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'ea1cc919-f472-4692-998f-3710de28986e'),
//    (gen_random_uuid(), 'LaunchPad7', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2, 'ea1cc919-f472-4692-998f-3710de28986e'),
//    (gen_random_uuid(), 'Vitrucare', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3, 'ea1cc919-f472-4692-998f-3710de28986e'),
//    (gen_random_uuid(), 'MAVEN LOGIX', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4, 'ea1cc919-f472-4692-998f-3710de28986e'),
   
//    (gen_random_uuid(), 'WECREATEPAKISTAN', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'f27c01cb-6d9a-40a1-ac30-23b61e1d6f4c'),
//    (gen_random_uuid(), 'DICE ANALYTICS', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2, 'f27c01cb-6d9a-40a1-ac30-23b61e1d6f4c'),
//    (gen_random_uuid(), 'GEN PAKISTAN', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3, 'f27c01cb-6d9a-40a1-ac30-23b61e1d6f4c'),
//    (gen_random_uuid(), '11Values', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4, 'f27c01cb-6d9a-40a1-ac30-23b61e1d6f4c')
   
   //latest
   
// insert into sponsors(id, name, client_id, conference_id, is_active, sort_order, sponsor_level_id) 
// values(gen_random_uuid(), 'Bradford', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
// (gen_random_uuid(), 'DLA', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
// (gen_random_uuid(), 'Government of Pakistan', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
// (gen_random_uuid(), 'Ilaan', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4, 'b564ed3d-9118-4547-8d30-d8f3de36fbf8'),
   
// (gen_random_uuid(), 'PITB', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'de443293-88d0-4fb0-a383-018546134411'),
// (gen_random_uuid(), 'PSEB', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 2, 'de443293-88d0-4fb0-a383-018546134411'),
// (gen_random_uuid(), 'PTV', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 3, 'de443293-88d0-4fb0-a383-018546134411'),
// (gen_random_uuid(), 'UMT', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 4, 'de443293-88d0-4fb0-a383-018546134411'),
   
// (gen_random_uuid(), 'UOL', '8c6e1b9e-3ebb-4ca0-9a2c-100d4ca0c95e', 'b10f4c67-930e-4e76-93fd-ddd52c532b79', true, 1, 'ea1cc919-f472-4692-998f-3710de28986e')

type Sponsor struct {
	Base
	Name             string
	IsActive         bool
	SponsorLevelID   uuid.UUID
	ClientID         uuid.UUID
	ConferenceID     uuid.UUID
	SortOrder int
	Description string 
	SocialMedia  SocialMedia `gorm:"embedded"`
}