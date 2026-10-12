-- +goose Up
CREATE TABLE chirps(
    id UUID PRIMARY KEY not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id UUID not null,
    constraint fk_user_id
    foreign key (user_id)
    references users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE chirps;
