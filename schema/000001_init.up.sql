CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               service_name TEXT NOT NULL,
                               price INTEGER NOT NULL,
                               user_id UUID NOT NULL,
                               start_date DATE NOT NULL,
                               end_date DATE,
                               created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
);