import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { Movimentacao } from '../models/movimentacao.model';

export interface MovimentacaoRequest {
  item_id: number;
  natureza: 'ENTRADA' | 'SAIDA';
  quantidade: number;
}

@Injectable({
  providedIn: 'root'
})
export class MovimentacaoService {
  private http = inject(HttpClient);
  private apiUrl = '/api/movimentacoes';

  findAll(): Observable<Movimentacao[]> {
    return this.http.get<Movimentacao[]>(this.apiUrl);
  }

  findByItemId(itemId: number): Observable<Movimentacao[]> {
    return this.http.get<Movimentacao[]>(
      `${this.apiUrl}/item/${itemId}`
    );
  }

  create(movimentacao: MovimentacaoRequest): Observable<Movimentacao> {
    return this.http.post<Movimentacao>(this.apiUrl, movimentacao);
  }
}
