export interface Movimentacao {
  id: number;
  item_id: number;
  natureza: 'ENTRADA' | 'SAIDA';
  quantidade: number;
  data: string;
}
