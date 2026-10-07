-- +goose Up
CREATE TABLE chirps(
    id UUID PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id UUID,
    constraint fk_user_id
    foreign key (user_id)
    references users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE chirps;
