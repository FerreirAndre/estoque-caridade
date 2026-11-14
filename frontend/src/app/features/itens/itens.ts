import { Component, OnInit, inject, signal } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';

import { Item } from '../../core/models/item.model';
import { ItemService } from '../../core/services/item.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-itens',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './itens.html',
  styleUrl: './itens.css'
})
export class Itens implements OnInit {
  private itemService = inject(ItemService);
  private fb = inject(FormBuilder);

  itens = signal<Item[]>([]);
  carregando = signal(false);
  salvando = signal(false);
  excluindo = signal<number | null>(null);

  erro = signal('');
  mensagemSucesso = signal('');
  mostrarFormulario = signal(false);
  itemEditando = signal<number | null>(null);

  form = this.fb.nonNullable.group({
    nome: ['', [Validators.required, Validators.maxLength(150)]],
    unidade_medida: ['KG' as 'KG' | 'LITRO' | 'PACOTE', Validators.required],
    estoque_minimo: [0, [Validators.min(0), Validators.required]]
  });

  ngOnInit(): void {
    this.carregarItens();
  }

  carregarItens(): void {
    this.carregando.set(true);
    this.erro.set('');

    this.itemService.findAll().subscribe({
      next: (dados) => {
        this.itens.set(dados ?? []);
        this.carregando.set(false);
      },
      error: (error) => {
        console.error('Erro ao carregar itens:', error);
        this.erro.set('Não foi possível carregar os itens.');
        this.carregando.set(false);
      }
    });
  }

  abrirFormulario(): void {
    this.itemEditando.set(null);

    this.form.reset({
      nome: '',
      unidade_medida: 'KG',
      estoque_minimo: 0
    });

    this.erro.set('');
    this.mensagemSucesso.set('');
    this.mostrarFormulario.set(true);
  }

  abrirEdicao(item: Item): void {
    this.itemEditando.set(item.id);

    this.form.reset({
      nome: item.nome,
      unidade_medida: item.unidade_medida,
      estoque_minimo: item.estoque_minimo
    });

    this.erro.set('');
    this.mensagemSucesso.set('');
    this.mostrarFormulario.set(true);
  }

  cancelarCadastro(): void {
    this.mostrarFormulario.set(false);
    this.itemEditando.set(null);

    this.form.reset({
      nome: '',
      unidade_medida: 'KG',
      estoque_minimo: 0
    });
  }

  salvarItem(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.salvando.set(true);
    this.erro.set('');
    this.mensagemSucesso.set('');

    const dados = this.form.getRawValue();
    const id = this.itemEditando();

    const requisicao: Observable<unknown> = id === null
      ? this.itemService.create(dados)
      : this.itemService.update(id, dados);

    requisicao.subscribe({
      next: () => {
        this.salvando.set(false);

        this.mensagemSucesso.set(
          id === null
            ? 'Item cadastrado com sucesso!'
            : 'Item atualizado com sucesso!'
        );

        this.cancelarCadastro();
        this.carregarItens();
      },
      error: (error: any) => {
        console.error('Erro ao salvar item:', error);
        this.erro.set('Não foi possível salvar o item.');
        this.salvando.set(false);
      }
    });
  }

  excluirItem(item: Item): void {
    const confirmado = window.confirm(
      `Deseja realmente excluir o item "${item.nome}"?`
    );

    if (!confirmado) {
      return;
    }

    this.excluindo.set(item.id);
    this.erro.set('');
    this.mensagemSucesso.set('');

    this.itemService.delete(item.id).subscribe({
      next: () => {
        this.excluindo.set(null);
        this.mensagemSucesso.set('Item excluído com sucesso!');
        this.carregarItens();
      },
      error: (error) => {
        console.error('Erro ao excluir item:', error);
        this.erro.set(
          'Não foi possível excluir o item. Verifique se ele possui movimentações registradas.'
        );
        this.excluindo.set(null);
      }
    });
  }
}
