-- Migration: Create astroneko_crm_users table
-- Description: Creates a table for CRM user authentication with username/password

CREATE TABLE astroneko_crm_users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username varchar(255) NOT NULL,
    "password" varchar(255) NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT astroneko_crm_users_pkey PRIMARY KEY (id),
    CONSTRAINT astroneko_crm_users_username_unique UNIQUE (username)
);