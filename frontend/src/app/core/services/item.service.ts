import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { Item } from "../models/item.model";

export interface ItemRequest {
  nome: string;
  unidade_medida: 'KG' | 'PACOTE' | 'LITRO';
  estoque_minimo: number;
}

@Injectable({
  providedIn: 'root'
})
export class ItemService {
  private http = inject(HttpClient);
  private apiUrl = '/api/itens';

  findAll(): Observable<Item[]> {
    return this.http.get<Item[]>(this.apiUrl);
  }

  findById(id: number): Observable<Item> {
    return this.http.get<Item>(`${this.apiUrl}/${id}`);
  }

  create(item: ItemRequest): Observable<Item> {
    return this.http.post<Item>(this.apiUrl, item);
  }

  update(id: number, item: ItemRequest): Observable<void> {
    return this.http.put<void>(`${this.apiUrl}/${id}`, item);
  }

  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
