DROP INDEX IF EXISTS uq_itens_nome;

ALTER TABLE itens
DROP CONSTRAINT chk_item_estoque_minimo;

ALTER TABLE itens
DROP COLUMN estoque_minimo;
