CREATE TABLE itens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nome VARCHAR(150) NOT NULL,
    quantidade BIGINT NOT NULL DEFAULT 0,
    unidade_medida VARCHAR(20) NOT NULL,

    CONSTRAINT chk_item_quantidade
        CHECK (quantidade >= 0),

    CONSTRAINT chk_item_unidade
        CHECK (unidade_medida IN ('KG', 'PACOTE', 'LITRO'))
);
