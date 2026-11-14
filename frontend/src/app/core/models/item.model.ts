export interface Item {
  id: number;
  nome: string;
  quantidade: number;
  unidade_medida: 'KG' | 'PACOTE' | 'LITRO';
  estoque_minimo: number;
}

