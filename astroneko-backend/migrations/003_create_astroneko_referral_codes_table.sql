CREATE TABLE astroneko_general_referral_codes (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	referral_code varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT astroneko_general_referral_codes_pkey PRIMARY KEY (id)
);

-- Insert the existing hardcoded referral codes
INSERT INTO astroneko_general_referral_codes (referral_code) VALUES
('ASTRONEKO!'),
('Fortunade'),
('SolanaAstro');
