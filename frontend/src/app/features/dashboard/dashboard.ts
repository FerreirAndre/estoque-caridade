import { Component, OnInit, computed, inject, signal } from '@angular/core';

import { DatePipe } from '@angular/common';
import { finalize, forkJoin } from 'rxjs';

import { Item } from '../../core/models/item.model';
import { Movimentacao } from '../../core/models/movimentacao.model';

import { ItemService } from '../../core/services/item.service';
import { MovimentacaoService } from '../../core/services/movimentacao.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [DatePipe],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard implements OnInit {
  private itemService = inject(ItemService);
  private movimentacaoService = inject(MovimentacaoService);

  itens = signal<Item[]>([]);
  movimentacoes = signal<Movimentacao[]>([]);

  carregando = signal(false);
  erro = signal('');

  totalItens = computed(() => this.itens().length);

  itensSemEstoque = computed(() => this.itens().filter((item) => item.quantidade === 0));

  itensEstoqueBaixo = computed(() =>
    this.itens().filter(
      (item) =>
        item.estoque_minimo > 0 && item.quantidade > 0 && item.quantidade <= item.estoque_minimo,
    ),
  );

  movimentacoesDoMes = computed(() => {
    const agora = new Date();

    return this.movimentacoes().filter((mov) => {
      const data = new Date(mov.data);

      return data.getMonth() === agora.getMonth() && data.getFullYear() === agora.getFullYear();
    });
  });

  totalEntradasMes = computed(
    () => this.movimentacoesDoMes().filter((mov) => mov.natureza === 'ENTRADA').length,
  );

  totalSaidasMes = computed(
    () => this.movimentacoesDoMes().filter((mov) => mov.natureza === 'SAIDA').length,
  );

  movimentacoesRecentes = computed(() => {
    const itens = new Map(this.itens().map((item) => [item.id, item]));

    return [...this.movimentacoes()]
      .sort((a, b) => new Date(b.data).getTime() - new Date(a.data).getTime())
      .slice(0, 5)
      .map((mov) => {
        const item = itens.get(mov.item_id);

        return {
          ...mov,
          nomeItem: item?.nome ?? `Item #${mov.item_id}`,
          unidadeMedida: item?.unidade_medida ?? '',
        };
      });
  });

  ngOnInit(): void {
    this.carregarDashboard();
  }

  carregarDashboard(): void {
    this.carregando.set(true);
    this.erro.set('');

    forkJoin({
      itens: this.itemService.findAll(),
      movimentacoes: this.movimentacaoService.findAll(),
    })
      .pipe(finalize(() => this.carregando.set(false)))
      .subscribe({
      next: (dados) => {
        this.itens.set(dados.itens ?? []);
        this.movimentacoes.set(dados.movimentacoes ?? []);
      },
      error: (error) => {
        console.error('Erro ao carregar dashboard:', error);
        this.erro.set('Não foi possível carregar os dados do dashboard.');
      },
    });
  }
}
