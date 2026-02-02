-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd


CREATE TABLE conversations (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid REFERENCES users(id) ON DELETE CASCADE ,
	meta_data JSON ,
    title VARCHAR(255),  
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() ,
	updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE graphs (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	conversation_id uuid REFERENCES conversations(id) ON DELETE CASCADE ,
	data JSONB ,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() ,
	updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE
);



CREATE TABLE messages (
	id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	conversation_id uuid REFERENCES conversations(id) ON DELETE CASCADE ,
    role VARCHAR(255),  
    content TEXT,
    tool_calls JSON, 
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() ,
	updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
