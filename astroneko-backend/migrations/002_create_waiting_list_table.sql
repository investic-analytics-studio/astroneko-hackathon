-- Migration: Create waiting list table
-- Created: 2025-01-10
-- Description: Creates the astroneko_waiting_list_users table for waiting list functionality

CREATE TABLE IF NOT EXISTS astroneko_waiting_list_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_astroneko_waiting_list_users_email ON astroneko_waiting_list_users(email);
