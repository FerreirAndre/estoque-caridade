import { Routes } from '@angular/router';
import { Layout } from './layout/layout';

export const routes: Routes = [
  {
    path: '',
    component: Layout,
    children: [
      {
        path: '',
        redirectTo: 'dashboard',
        pathMatch: 'full'
      },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./features/dashboard/dashboard')
            .then(m => m.Dashboard)
      },
      {
        path: 'itens',
        loadComponent: () =>
          import('./features/itens/itens')
            .then(m => m.Itens)
      },
      {
        path: 'movimentacoes',
        loadComponent: () =>
          import('./features/movimentacoes/movimentacoes')
            .then(m => m.Movimentacoes)
      }
    ]
  }
];
