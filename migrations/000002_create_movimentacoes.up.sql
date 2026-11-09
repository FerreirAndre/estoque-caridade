CREATE TABLE movimentacoes (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    item_id BIGINT NOT NULL,
    natureza VARCHAR(10) NOT NULL,
    quantidade BIGINT NOT NULL,
    data TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_movimentacao_item
        FOREIGN KEY (item_id)
        REFERENCES itens(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_movimentacao_natureza
        CHECK (natureza IN ('ENTRADA', 'SAIDA')),

    CONSTRAINT chk_movimentacao_quantidade
        CHECK (quantidade > 0)
);

CREATE INDEX idx_movimentacoes_item_id
    ON movimentacoes(item_id);

CREATE INDEX idx_movimentacoes_data
    ON movimentacoes(data DESC);
