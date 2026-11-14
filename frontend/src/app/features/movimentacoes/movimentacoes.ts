import {
  Component,
  OnInit,
  computed,
  inject,
  signal
} from '@angular/core';

import {
  FormBuilder,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';

import { DatePipe } from '@angular/common';
import { forkJoin } from 'rxjs';

import { Item } from '../../core/models/item.model';
import { Movimentacao } from '../../core/models/movimentacao.model';

import { ItemService } from '../../core/services/item.service';
import { MovimentacaoService } from '../../core/services/movimentacao.service';

@Component({
  selector: 'app-movimentacoes',
  standalone: true,
  imports: [ReactiveFormsModule, DatePipe],
  templateUrl: './movimentacoes.html',
  styleUrl: './movimentacoes.css'
})
export class Movimentacoes implements OnInit {
  private fb = inject(FormBuilder);
  private itemService = inject(ItemService);
  private movimentacaoService = inject(MovimentacaoService);

  itens = signal<Item[]>([]);
  movimentacoes = signal<Movimentacao[]>([]);

  carregando = signal(false);
  salvando = signal(false);

  erro = signal('');
  mensagemSucesso = signal('');

  filtroItem = signal('TODOS');
  filtroNatureza = signal('TODOS');
  itemSelecionadoId = signal(0);

  itemSelecionado = computed(() =>
    this.itens().find(item => item.id === this.itemSelecionadoId())
  );

  form = this.fb.nonNullable.group({
    item_id: [0, [Validators.required, Validators.min(1)]],
    natureza: ['ENTRADA' as 'ENTRADA' | 'SAIDA', Validators.required],
    quantidade: [1, [Validators.required, Validators.min(1)]]
  });

  historicoFiltrado = computed(() => {
    const itens = this.itens();
    const movimentacoes = this.movimentacoes();

    return movimentacoes
      .filter(mov => {
        const correspondeItem =
          this.filtroItem() === 'TODOS' ||
          String(mov.item_id) === this.filtroItem();

        const correspondeNatureza =
          this.filtroNatureza() === 'TODOS' ||
          mov.natureza === this.filtroNatureza();

        return correspondeItem && correspondeNatureza;
      })
      .map(mov => ({
        ...mov,
        nomeItem: itens.find(item => item.id === mov.item_id)?.nome
          ?? `Item #${mov.item_id}`
      }))
      .sort((a, b) =>
        new Date(b.data).getTime() - new Date(a.data).getTime()
      );
  });

  ngOnInit(): void {
    this.carregarDados();
  }

  carregarDados(): void {
    this.carregando.set(true);
    this.erro.set('');

    forkJoin({
      itens: this.itemService.findAll(),
      movimentacoes: this.movimentacaoService.findAll()
    }).subscribe({
      next: (dados) => {
        this.itens.set(dados.itens ?? []);
        this.movimentacoes.set(dados.movimentacoes ?? []);
        this.carregando.set(false);
      },
      error: (error) => {
        console.error('Erro ao carregar movimentações:', error);
        this.erro.set('Não foi possível carregar os dados.');
        this.carregando.set(false);
      }
    });
  }

  registrarMovimentacao(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const { item_id, natureza, quantidade } = this.form.getRawValue();

    const item = this.itens().find(item => item.id === item_id);

    if (!item) {
      this.erro.set('Selecione um item válido.');
      return;
    }

    if (natureza === 'SAIDA' && quantidade > item.quantidade) {
      this.erro.set(
        `Estoque insuficiente. Disponível: ${item.quantidade} ${item.unidade_medida}.`
      );
      return;
    }

    this.salvando.set(true);
    this.erro.set('');
    this.mensagemSucesso.set('');

    this.movimentacaoService.create(this.form.getRawValue()).subscribe({
      next: () => {
        this.salvando.set(false);

        this.mensagemSucesso.set(
          'Movimentação registrada com sucesso!'
        );

        this.form.reset({
          item_id: 0,
          natureza: 'ENTRADA',
          quantidade: 1
        });
        this.itemSelecionadoId.set(0);

        this.carregarDados();
      },
      error: (error) => {
        console.error('Erro ao registrar movimentação:', error);

        this.erro.set(
          error.error?.erro ??
          'Não foi possível registrar a movimentação. Verifique o estoque disponível.'
        );

        this.salvando.set(false);
      }
    });
  }

  atualizarItemSelecionado(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.itemSelecionadoId.set(Number(select.value));
  }
}
