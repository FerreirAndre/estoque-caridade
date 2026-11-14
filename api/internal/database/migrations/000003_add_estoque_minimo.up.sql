ALTER TABLE itens
ADD COLUMN estoque_minimo BIGINT NOT NULL DEFAULT 0;

ALTER TABLE itens
ADD CONSTRAINT chk_item_estoque_minimo
CHECK (estoque_minimo >= 0);

CREATE UNIQUE INDEX uq_itens_nome
ON itens (LOWER(TRIM(nome)));
