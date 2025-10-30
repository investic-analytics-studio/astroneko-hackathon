-- Update referral_code column to referral_code_id with proper foreign key
ALTER TABLE astroneko_referral_logs
DROP COLUMN referral_code;

ALTER TABLE astroneko_referral_logs
ADD COLUMN referral_code_id uuid;

-- Add foreign key constraint to reference astroneko_user_referral_codes table
ALTER TABLE astroneko_referral_logs
ADD CONSTRAINT astroneko_referral_logs_referral_code_id_fkey
FOREIGN KEY (referral_code_id) REFERENCES astroneko_user_referral_codes(id) ON DELETE SET NULL ON UPDATE CASCADE;